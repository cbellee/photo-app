package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"utils"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

var (
	serviceName = utils.GetEnvValue("SERVICE_NAME", "some service")
	servicePort = utils.GetEnvValue("SERVICE_PORT", "8080")

	uploadsQueueBinding     = utils.GetEnvValue("UPLOADS_QUEUE_BINDING", "")
	imagesContainerBinding  = utils.GetEnvValue("IMAGES_CONTAINER_BINDING", "")
	thumbsContainerBinding  = utils.GetEnvValue("THUMBS_CONTAINER_BINDING", "")
	uploadsContainerBinding = utils.GetEnvValue("UPLOADS_CONTAINER_BINDING", "")
	maxRequestBodySizeMb    = utils.GetEnvValue("GRPC_MAX_REQUEST_BODY_SIZE_MB", "30")

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

	// get image
	blobName := evt.Data.Url[strings.LastIndex(evt.Data.Url, "/")+1:]

	// maxRequestBodySize := 30 //
	maxRequestBodySize, err := strconv.Atoi(maxRequestBodySizeMb)
	if err != nil {
		Error.Printf("%s: error converting string to integer", err)
		return nil, err
	}
	
	blob, err := utils.GetLargeBlob(ctx, uploadsContainerBinding, blobName, maxRequestBodySize)
	if err != nil {
		Error.Printf("%s: error getting blob '%s': %v", serviceName, blobName, err)
		return nil, err
	}

	// resize image as thumbnail
	thumbBytes, err := utils.ResizeImage(blob.Data, evt.Data.ContentType, blobName, mth, mtw)
	if err != nil {
		Error.Printf("%s: error resizing image: '%s': %v", serviceName, blobName, err)
		return nil, err
	}

	// write thumbnail to blob storage
	_, err = setBlob(ctx, thumbsContainerBinding, blobName, thumbBytes, blob.Metadata["collection"], blob.Metadata["album"])
	if err != nil {
		Error.Printf("%s: error saving blob: '%s': %v", serviceName, blobName, err)
		return nil, err
	}

	// resize main image
	imgBytes, err := utils.ResizeImage(blob.Data, evt.Data.ContentType, blobName, mih, miw)
	if err != nil {
		Error.Printf("%s: error resizing image: '%s': %v", serviceName, blobName, err)
		return nil, err
	}

	// write main image to blob storage
	_, err = setBlob(ctx, imagesContainerBinding, blobName, imgBytes, blob.Metadata["collection"], blob.Metadata["album"])
	if err != nil {
		Error.Printf("%s: error saving blob '%s': %v", serviceName, blobName, err)
		return nil, err
	}

	return nil, nil
}

func setBlob(ctx context.Context, bindingName string, blobName string, blob []byte, collection string, album string) (out *dapr.BindingEvent, err error) {
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
		},
	}

	Info.Printf("%s: saving image '%s' to container '%s'", serviceName, blobName, bindingName)

	out, err = client.InvokeBinding(ctx, in)
	if err != nil {
		return out, fmt.Errorf("setBlob: %w", err)
	}

	return out, nil
}
