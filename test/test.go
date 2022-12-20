package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

var (
	containerName      = "uploads"
	storageAccountName = "storaodiwusgqeyiu"
)

type Album struct {
	Name string `json:"name"`
}

type Collection struct {
	Name   string  `json:"name"`
	Albums []Album `json:"albums"`
}

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

	//c := Collection{}
	fmt.Printf("getting colections & albums in the container: %s \n", containerName)
	result := listBlobDirectories(containerClient, ctx, opt, map[string][]string{})

	jsonString, err := json.Marshal(result)
	fmt.Printf("Collections data: %s \n", jsonString)
}

var cl = []Collection{}
/* 
func listBlobDirectories(containerClient *container.Client, ctx context.Context, opt container.ListBlobsHierarchyOptions) Collection {
	pager := containerClient.NewListBlobsHierarchyPager("/", &opt)
	var c = Collection{}
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		segment := resp.Segment
		
		if segment.BlobPrefixes != nil {
			
			for _, prefix := range segment.BlobPrefixes {
				
				str := strings.Split(strings.Trim(*prefix.Name, "/"), "/")
				c.Name = str[0]

				if len(str) > 1 {
					album := Album{
						Name: strings.Trim(str[1], "/"),
					}
					c.Albums = append(c.Albums, album)
				}

				opt := container.ListBlobsHierarchyOptions{
					Prefix: prefix.Name,
				}
				t := listBlobDirectories(containerClient, ctx, opt)
				cl = append(cl, t)
				//return c
			}
			//return c
		}
		//return c
		//cl = append(cl, c)
	}
	
	return c
}
 */

func listBlobDirectories(containerClient *container.Client, ctx context.Context, opt container.ListBlobsHierarchyOptions, m map[string][]string) map[string][]string {
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
				listBlobDirectories(containerClient, ctx, opt, m)
			}
		}
	}
	return m
}