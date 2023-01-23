package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"models"
	"utils"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var (
	serviceName = utils.GetEnvValue("SERVICE_NAME", "some service")
	servicePort = utils.GetEnvValue("SERVICE_PORT", "8080")

	uploadsQueueBinding     = utils.GetEnvValue("UPLOADS_QUEUE_BINDING", "")
	imagesContainerBinding  = utils.GetEnvValue("IMAGES_CONTAINER_BINDING", "")
	thumbsContainerBinding  = utils.GetEnvValue("THUMBS_CONTAINER_BINDING", "")
	uploadsContainerBinding = utils.GetEnvValue("UPLOADS_CONTAINER_BINDING", "")
	maxRequestBodySizeMb    = utils.GetEnvValue("GRPC_MAX_REQUEST_BODY_SIZE_MB", "30")

	storageConfig = models.StorageConfig{
		StorageAccount:   utils.GetEnvValue("STORAGE_ACCOUNT_NAME", ""),
		StorageContainer: utils.GetEnvValue("STORAGE_CONTAINER_NAME", ""),
	}

	thumbsContainerName = "thumbs"
	imagesContainerName = "images"

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Lshortfile)
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	Info.Printf("serviceName: %s", serviceName)
	Info.Printf("servicePort: %s", servicePort)

	Info.Printf("storageQueueBinding: %s", uploadsQueueBinding)
	Info.Printf("imagesContainerBinding: %s", imagesContainerBinding)
	Info.Printf("thumbsContainerBinding: %s", thumbsContainerBinding)
	Info.Printf("uploadsContainerBinding: %s", uploadsContainerBinding)

	port := fmt.Sprintf(":%s", servicePort)

	s, err := daprd.NewService(port)
	if err != nil {
		Error.Printf("%s: failed create the server: %v", serviceName, err)
	}

	// add storage queue input binding invocation handler
	if err := s.AddBindingInvocationHandler(uploadsQueueBinding, ResizeHandler); err != nil {
		Error.Printf("%s: error adding storage queue binding handler: %v", serviceName, err)
	}

	// start the service
	if err := s.Start(); err != nil {
		Error.Printf("%s: server failed to start: %v", serviceName, err)
	}
}

func ResizeHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	Info.Printf("%s: resizeHandler", serviceName)

	// get env Vars
	mth, err := strconv.Atoi(utils.GetEnvValue("MAX_THUMB_HEIGHT", "96"))
	if err != nil {
		Error.Printf("%s: %v", serviceName, err)
	}
	mtw, err := strconv.Atoi(utils.GetEnvValue("MAX_THUMB_WIDTH", "128"))
	if err != nil {
		Error.Printf("%s: %v", serviceName, err)
	}
	mih, err := strconv.Atoi(utils.GetEnvValue("MAX_IMAGE_HEIGHT", "640"))
	if err != nil {
		Error.Printf("%s: %v", serviceName, err)
	}
	miw, err := strconv.Atoi(utils.GetEnvValue("MAX_IMAGE_WIDTH", "480"))
	if err != nil {
		Error.Printf("%s: %v", serviceName, err)
	}

	evt, err := utils.ConvertToEvent(in)
	if err != nil {
		Error.Printf("%s: error converting BindingEvent to struct: %v", serviceName, err)
	}

	Info.Printf("%s: input binding handler '%s': Url: '%s', EventTime: '%s' MetaData: '%v'", serviceName, uploadsQueueBinding, evt.Data.Url, evt.EventTime, in.Metadata)

	storageUrl := fmt.Sprintf("https://%s.blob.core.windows.net/", storageConfig.StorageAccount)

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	client, err := azblob.NewClient(storageUrl, credential, &azblob.ClientOptions{})
	if err != nil {
		Error.Printf("error creating blob client: %s", err)
	}

	u, err := url.Parse(evt.Data.Url)
	if err != nil {
		log.Fatal(err)
	}

	path := strings.Split(u.Path, "/")
	blobPath := fmt.Sprintf("%s/%s/%s", path[len(path)-3], path[len(path)-2], path[len(path)-1])
	fmt.Printf("blobPath: %s\n", blobPath)

	// maxRequestBodySize := 30 //
	maxRequestBodySize, err := strconv.Atoi(maxRequestBodySizeMb)
	if err != nil {
		Error.Printf("%s: error converting string to integer", err)
		return nil, err
	}

	blob, err := utils.GetBlob(ctx, uploadsContainerBinding, blobPath, maxRequestBodySize)
	if err != nil {
		Error.Printf("%s: error getting blob '%s': %v", serviceName, blobPath, err)
		return nil, err
	}

	// resize image as thumbnail
	Info.Printf("got blob size: %s", fmt.Sprint(len(blob.Data)))
	thumbBytes, err := utils.ResizeImage(blob.Data, evt.Data.ContentType, blobPath, mth, mtw)
	if err != nil {
		Error.Printf("%s: error resizing image: '%s': %v", serviceName, blobPath, err)
		return nil, err
	}

	Info.Printf("'thumbBytes' blob size: %s", fmt.Sprint(len(thumbBytes)))

	// write thumbnail to blob storage
	_, err = setBlob(ctx, thumbsContainerBinding, blobPath, thumbBytes, blob.Metadata["collection"], blob.Metadata["album"], "image/jpeg")
	if err != nil {
		Error.Printf("%s: error saving blob: '%s': %v", serviceName, blobPath, err)
		return nil, err
	}

	// set tags
	tags := make(map[string]string)
	tags["collection"] = blob.Metadata["collection"]
	tags["album"] = blob.Metadata["album"]
	tags["isThumb"] = "true"
	tags["url"] = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", storageConfig.StorageAccount, "thumbs", string(blobPath))
	tags["imgUrl"] = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", storageConfig.StorageAccount, "images", string(blobPath))
	tags["name"] = path[len(path)-1]
	tags["prefix"] = fmt.Sprintf("%s/%s/%s", path[len(path)-3], path[len(path)-2], path[len(path)-1])

	t, err  := json.Marshal(tags)
	if err != nil {
		Info.Fatal(err)
	}

	Info.Println("Tags \n -------- \n", string(t))

	Info.Printf("setting tags for blob: %s", blobPath)
	Info.Printf("setting tags: %s \n tags: %s", blobPath, tags)
	_, err = client.ServiceClient().NewContainerClient(thumbsContainerName).NewBlockBlobClient(blobPath).SetTags(ctx, tags, nil)
	if err != nil {
		Error.Printf("error setting tags: %s", err)
	}

	// resize main image
	imgBytes, err := utils.ResizeImage(blob.Data, evt.Data.ContentType, blobPath, mih, miw)
	if err != nil {
		Error.Printf("%s: error resizing image: '%s': %v", serviceName, blobPath, err)
		return nil, err
	}

	// write main image to blob storage
	Info.Printf("'imgBytes' blob size: %s", fmt.Sprint(len(imgBytes)))
	_, err = setBlob(ctx, imagesContainerBinding, blobPath, imgBytes, blob.Metadata["collection"], blob.Metadata["album"], "image/jpeg")
	if err != nil {
		Error.Printf("%s: error saving blob '%s': %v", serviceName, blobPath, err)
		return nil, err
	}

	tags["isThumb"] = "false"
	tags["url"] = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", storageConfig.StorageAccount, "images", path)

	_, err = client.ServiceClient().NewContainerClient(imagesContainerName).NewBlockBlobClient(blobPath).SetTags(ctx, tags, nil)
	if err != nil {
		Error.Printf("error setting tags: %s", err)
	}

	return nil, nil
}

func setBlob(ctx context.Context, bindingName string, blobName string, blob []byte, collection string, album string, contentType string) (out *dapr.BindingEvent, err error) {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	in := &dapr.InvokeBindingRequest{
		Name:      bindingName,
		Operation: "create",
		Data:      blob,
		Metadata: map[string]string{
			"blobName":        blobName,
			"includeMetadata": "true",
			"collection":      collection,
			"album":           album,
			"ContentType":     contentType,
		},
	}

	Info.Printf("%s: saving image '%s' to binding '%s'", serviceName, blobName, bindingName)

	out, err = client.InvokeBinding(ctx, in)
	if err != nil {
		return out, fmt.Errorf("error setBlob: %w", err)
	}

	return out, nil
}
