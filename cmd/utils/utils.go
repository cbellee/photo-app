package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"models"

	"github.com/dapr/go-sdk/service/common"
	"golang.org/x/image/draw"
)

func ResizeImage(imgBytes []byte, imageFormat string, blobName string, maxHeight int, maxWidth int) (img []byte, err error) {

	var dst *image.RGBA
	var buf = new(bytes.Buffer)

	log.Printf("ResizeImage: imgBytes length='%d'", len(imgBytes))

	src, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return buf.Bytes(), fmt.Errorf("ResizeImage: %w", err)
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
