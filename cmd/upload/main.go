package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"models"
	"os"
	"utils"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var (
	serviceName = utils.GetEnvValue("SERVICE_NAME", "")
	servicePort = utils.GetEnvValue("SERVICE_PORT", "")

	storageConfig = models.StorageConfig{
		StorageAccount:   utils.GetEnvValue("STORAGE_ACCOUNT_NAME", ""),
		StorageContainer: utils.GetEnvValue("STORAGE_CONTAINER_NAME", ""),
	}

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

	port := fmt.Sprintf(":%s", servicePort)

	s, err := daprd.NewService(port)
	if err != nil {
		Error.Printf("%s: failed create the server: %v", serviceName, err)
	}

	// invocation handler
	if err := s.AddServiceInvocationHandler("/upload", UploadHandler); err != nil {
		Error.Printf("%s: error adding upload handler: %v", serviceName, err)
	}

	// start the service
	if err := s.Start(); err != nil {
		Error.Printf("%s: server failed to start: %v", serviceName, err)
	}
}

func UploadHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	Info.Printf("%s - ContentType:%s, Verb:%s, QueryString:%s, Data:%+v", serviceName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	blob := models.Blob{}
	err = json.Unmarshal(in.Data, blob)
	if err != nil {
		Error.Fatalf("Error unmarshalling blob JSON")
		return nil, err
	}
	
	o := common.Content{
		Data: nil,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	_, err = SaveBlob(ctx, storageConfig, blob.Prefix, blob.Name, []byte(blob.Data), blob.MetaData, blob.Tags)
	if err != nil {
		Error.Printf("error uploading blob: %s", err)
		return &o, err
	}


	return &o, nil
}

func SaveBlob(ctx context.Context, storageConfig models.StorageConfig, blobPrefix string, blobName string, blobBytes []byte, metadata map[string]string, tags map[string]string) (eTag *azcore.ETag, err error) {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", storageConfig.StorageAccount)
	ctx = context.Background()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	// upload blob
	blobPath := fmt.Sprintf("%s/%s", blobPrefix, blobName)
	opt := azblob.UploadBufferOptions{
		Metadata: metadata,
		Tags:     tags,
	}

	resp, err := client.UploadBuffer(ctx, storageConfig.StorageContainer, blobPath, blobBytes, &opt)
	if err != nil {
		log.Fatal("Blob upload failed with error: " + err.Error())
		return nil, err
	}

	Info.Printf("response: %s", resp)
	return resp.ETag, nil
}
