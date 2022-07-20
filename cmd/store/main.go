package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"models"
	"net/url"
	"os"
	"strings"
	"time"
	"utils"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/google/uuid"
)

var (
	serviceName    = os.Getenv("SERVICE_NAME")
	servicePort    = os.Getenv("SERVICE_PORT")
	serviceVersion = os.Getenv("SERVICE_VERSION")

	cosmosBinding       = os.Getenv("COSMOS_BINDING")
	imagesQueueBinding = os.Getenv("IMAGES_QUEUE_BINDING")

	imagesContainerBinding  = os.Getenv("IMAGES_CONTAINER_BINDING")
	thumbsContainerBinding  = os.Getenv("THUMBS_CONTAINER_BINDING")
	uploadsContainerBinding = os.Getenv("UPLOADS_CONTAINER_BINDING")

	partitionKey = "/url"

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
	Info.Printf("cosmosBinding: %s", cosmosBinding)
	Info.Printf("imagesQueueBinding: %s", imagesQueueBinding)
	Info.Printf("uploadsContainerBinding: %s", uploadsContainerBinding)

	port := fmt.Sprintf(":%s", servicePort)

	s, err := daprd.NewService(port)
	if err != nil {
		Error.Printf("%s: failed create the server: %v", serviceName, err)
	}

	// add images queue input binding invocation handler
	if err := s.AddBindingInvocationHandler(imagesQueueBinding, ImagesHandler); err != nil {
		Error.Printf("%s: error adding storage queue binding handler: %v", serviceName, err)
	}

	// start the service
	if err := s.Start(); err != nil {
		Error.Printf("%s: server failed to start: %v", serviceName, err)
	}
}

func ImagesHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	Info.Printf("%s: imagesHandler", serviceName)

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	// cast event grid event to struct
	evt, err := utils.ConvertToEvent(in)
	if err != nil {
		Error.Printf("%s: error converting BindingEvent to struct: %v", serviceName, err)
	}

	u, err := url.Parse(evt.Data.Url)
	if err != nil {
		Error.Printf("%s: error parsing string to url: %v", serviceName, err)
	}

	blobName := u.Path[strings.LastIndex(u.Path, "/")+1:]
	Info.Printf("blobName: %s", blobName)

	storageUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
	Info.Printf("storageUrl: %s", storageUrl)

	now := time.Now().UTC()

	Info.Printf("%s: getting blob: '%s'", serviceName, blobName)
	blob, err := getBlob(ctx, uploadsContainerBinding, blobName)
	if err != nil {
		Error.Printf("%s: error getting blob '%s': %v", serviceName, blobName, err)
		return nil, err
	}

	Info.Printf("found blob: %s", blobName)
	Info.Printf("blob Metadata: album=%s collection=%s", blob.Metadata["album"], blob.Metadata["collection"])

	p := models.Photo{
		Id:            uuid.New(),
		Name:          blobName,
		Album:         blob.Metadata["album"],
		Collection:    blob.Metadata["collection"],
		Type:          "application/json",
		Url:           evt.Data.Url,
		ThumbsUrl:     fmt.Sprintf("%s/thumbs/%s", storageUrl, blobName),
		ContentLength: evt.Data.ContentLength,
		Created:       now,
		Modified:      now,
	}

	photo, _ := json.Marshal(p)
	Info.Printf("photo JSON: %s", string(photo))

	Info.Printf("saving blob '%s' to database '%s'", blobName, cosmosBinding)

	meta := make(map[string]string)
	meta["partitionKey"] = partitionKey

	si := &dapr.SetStateItem{
		Key:  blobName,
		Value:  photo,
		Metadata: meta,
		Options:  &dapr.StateOptions{
			Concurrency:  dapr.StateConcurrencyUndefined,
			Consistency:  dapr.StateConsistencyUndefined,
		},
	}

	err = client.SaveBulkState(ctx, cosmosBinding, si)
	if err != nil {
		Error.Printf("%s: error saving state '%s': %v", serviceName, string(photo), err)
		return nil, err
	}

	return photo, nil
}

func getBlob(ctx context.Context, bindingName string, blobName string) (out *dapr.BindingEvent, err error) {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	in := &dapr.InvokeBindingRequest{
		Name:      bindingName,
		Operation: "get",
		Metadata: map[string]string{
			"blobName":        blobName,
			"includeMetadata": "true",
		},
	}

	out, err = client.InvokeBinding(ctx, in)
	if err != nil {
		Error.Printf("%s: error invoking binding '%s': %v", serviceName, bindingName, err)
		return nil, err
	}

	return out, nil
}
