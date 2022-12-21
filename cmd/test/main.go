package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"utils"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

var (
	containerName      = "uploads"
	storageAccountName = "storaodiwusgqeyiu"
)

func main() {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)
	ctx := context.Background()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	// get container client
	containerClient := client.ServiceClient().NewContainerClient(containerName)

	// List top-level directories in the container
	root := "" // empty string denotes root folder in a container
	opt := container.ListBlobsHierarchyOptions{
		Prefix: &root,
	}

	blobs := utils.GetBlobDirectories(containerClient, ctx, opt, map[string][]string{})
	metadata := utils.GetBlobMetadata(client, containerName, ctx)
	tags := utils.GetBlobTags(client, containerName, ctx)

	jsonString, err := json.Marshal(tags)
	if err != nil {
		log.Fatal("error marshalling JSON: err", err)
	}
	fmt.Printf("\nTags: %s \n", jsonString)

	jsonString, err = json.Marshal(metadata)
	if err != nil {
		log.Fatal("error marshalling JSON: err", err)
	}
	fmt.Printf("\nMetadata: %s \n", jsonString)

	jsonString, err = json.Marshal(blobs)
	if err != nil {
		log.Fatal("error marshalling JSON: err", err)
	}
	fmt.Printf("\nBlobs: %s \n", jsonString)
}
