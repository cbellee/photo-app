package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"models"
	"net"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"golang.org/x/image/draw"
	"google.golang.org/grpc"
)

func ResizeImage(imgBytes []byte, imageFormat string, blobName string, maxHeight int, maxWidth int) (img []byte, err error) {

	var dst *image.RGBA
	var buf = new(bytes.Buffer)

	log.Printf("ResizeImage: imgBytes length='%d'", len(imgBytes))

	src, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return buf.Bytes(), fmt.Errorf("Error ResizeImage: %w", err)
	}

	height := src.Bounds().Dy()
	width := src.Bounds().Dx()

	// if height > width, then the image is portrait so resize height to maxHeight
	if height > width {
		newWidth := maxHeight * width / height
		log.Printf("ResizeImage: resizing blob '%s' with height: %s width: %s to newHeight: %s newWidth: %s", blobName, fmt.Sprint(height), fmt.Sprint(width), fmt.Sprint(maxHeight), fmt.Sprint(newWidth))
		dst = image.NewRGBA((image.Rect(0, 0, newWidth, maxHeight)))
		// if height <= width, then the image is landscape or square so resize width to maxWidth
	} else {
		newHeight := maxWidth * height / width
		log.Printf("ResizeImage: resizing blob '%s' with height: %s width: %s to newHeight: %s newWidth: %s", blobName, fmt.Sprint(height), fmt.Sprint(width), fmt.Sprint(newHeight), fmt.Sprint(maxWidth))
		dst = image.NewRGBA((image.Rect(0, 0, maxWidth, newHeight)))
	}

	// detect image type from 'imageFormat' value
	switch imageFormat {
	case "image/jpeg":
		log.Printf("ResizeImage: encoding jpeg '%s' to type '%s'", blobName, imageFormat)
		src, _ = jpeg.Decode(bytes.NewReader(imgBytes))

		draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
		err := jpeg.Encode(buf, dst, nil)
		if err != nil {
			return nil, fmt.Errorf("ResizeImage: %w", err)
		}
	case "image/png":
		log.Printf("ResizeImage: encoding png '%s' to type '%s'", blobName, imageFormat)
		src, _ = png.Decode(bytes.NewReader(imgBytes))

		draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
		err := png.Encode(buf, dst)
		if err != nil {
			return nil, fmt.Errorf("ResizeImage: %w", err)
		}
	case "image/gif":
		log.Printf("ResizeImage: encoding gif '%s' to type '%s'", blobName, imageFormat)
		src, _ = gif.Decode(bytes.NewReader(imgBytes))

		draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
		err := gif.Encode(buf, dst, nil)
		if err != nil {
			return nil, fmt.Errorf("ResizeImage: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func ConvertToEvent(b *common.BindingEvent) (models.Event, error) {
	var evt models.Event

	byt := make([]byte, base64.StdEncoding.DecodedLen(len(b.Data)))
	l, err := base64.StdEncoding.Decode(byt, b.Data)
	if err != nil {
		return evt, fmt.Errorf("ConvertToEvent: %w", err)
	}

	err = json.Unmarshal(byt[:l], &evt)
	if err != nil {
		return evt, fmt.Errorf("ConvertToEvent: %w", err)
	}
	return evt, nil
}

func GetEnvValue(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}



func GetBlob(ctx context.Context, bindingName string, blobName string, maxRequestBodySize int) (out *dapr.BindingEvent, err error) {
	headerBuffer := 1
	var opts []grpc.CallOption
	opts = append(opts, grpc.MaxCallRecvMsgSize((maxRequestBodySize+headerBuffer)*1024*1024))
	conn, err := grpc.Dial(net.JoinHostPort("127.0.0.1",
		GetEnvValue("DAPR_GRPC_PORT", "50001")),
		grpc.WithDefaultCallOptions(opts...), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := dapr.NewClientWithConnection(conn)
	if err != nil {
		panic(err)
	}
	//defer client.Close()

	in := &dapr.InvokeBindingRequest{
		Name:      bindingName,
		Operation: "get",
		Metadata: map[string]string{
			"blobName":        blobName,
			"includeMetadata": "true",
		},
	}

	fmt.Printf("getting blob '%s'", blobName)

	out, err = client.InvokeBinding(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("getLargeBlob: %w", err)
	}

	fmt.Printf("blob size in utils.go: %s", fmt.Sprint((out.Data)) )
	return out, nil
}

func GetBlobDirectories(containerClient *container.Client, ctx context.Context, opt container.ListBlobsHierarchyOptions, m map[string][]string) map[string][]string {
	pager := containerClient.NewListBlobsHierarchyPager("/", &opt)

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		segment := resp.Segment
		if segment.BlobPrefixes != nil {
			for _, prefix := range segment.BlobPrefixes {
				str := strings.Split(strings.Trim(*prefix.Name, "/"), "/")
				if len(str) > 1 {
					m[str[0]] = append(m[str[0]], strings.Trim(str[1], "/"))
				}

				opt := container.ListBlobsHierarchyOptions{
					Prefix: prefix.Name,
				}
				GetBlobDirectories(containerClient, ctx, opt, m)
			}
		}
	}
	return m
}

func GetBlobTags(blobClient *azblob.Client, containerName string, ctx context.Context) map[string][]string {
	pager := blobClient.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{
			Deleted:  false,
			Versions: false,
			Metadata: true,
			Tags:     true,
		},
	})

	blobTagMap := make(map[string][]string)

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("error occurred during blob listing: %s", err)
			break
		}

		for _, _blob := range resp.Segment.BlobItems {
			tags := *_blob.BlobTags
			album := ""
			collection := ""

			for _, t := range tags.BlobTagSet {
				if *t.Key == "collection" {
					collection = *t.Value
				}

				if *t.Key == "album" {
					album = *t.Value
				}
			}

			if !Contains(blobTagMap[collection], album) {
				blobTagMap[collection] = append(blobTagMap[collection], album)
			}
		}
	}
	return blobTagMap
}

func GetBlobMetadata(blobClient *azblob.Client, containerName string, ctx context.Context) map[string][]string {
	pager := blobClient.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{
			Deleted:  false,
			Versions: false,
			Metadata: true,
			Tags:     true,
		},
	})

	metadataMap := make(map[string][]string)

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("error occurred during blob listing: %s", err)
			break
		}

		for _, _blob := range resp.Segment.BlobItems {
			if !Contains(metadataMap[*_blob.Metadata["collection"]], *_blob.Metadata["album"]) {
				metadataMap[*_blob.Metadata["collection"]] = append(metadataMap[*_blob.Metadata["collection"]], *_blob.Metadata["album"])
			}
		}
	}
	return metadataMap
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
