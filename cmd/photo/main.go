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
	"utils"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	serviceName   = os.Getenv("SERVICE_NAME")
	servicePort   = os.Getenv("SERVICE_PORT")
	cosmosBinding = os.Getenv("COSMOS_BINDING")

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var dbConfig = models.DbConfig{
	DbURL:       os.Getenv("COSMOS_URL"),
	DbName:      os.Getenv("COSMOS_DB"),
	DbKey:       os.Getenv("COSMOS_KEY"),
	DbContainer: os.Getenv("COSMOS_CONTAINER"),
}

var storageConfig = models.StorageConfig{
	StorgeURL: ,
}

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
	Info.Printf("dbContainer: %s", dbConfig.DbContainer)
	Info.Printf("dbKey: %s", dbConfig.DbKey)
	Info.Printf("dbName: %s", dbConfig.DbName)
	Info.Printf("dbURL: %s", dbConfig.DbURL)

	port := fmt.Sprintf(":%s", servicePort)

	s := daprd.NewService(port)

	// add invocation handler
	if err := s.AddServiceInvocationHandler("/photos", getPhotos); err != nil {
		Error.Printf("%s: error adding service invocation handler: %v", serviceName, err)
	}

	// add invocation handler
	if err := s.AddServiceInvocationHandler("/collections", getCollectionAlbumTags); err != nil {
		Error.Printf("%s: error adding service invocation handler: %v", serviceName, err)
	}

	// start the service
	if err := s.Start(); err != nil {
		Error.Printf("%s: server failed to start: %v", serviceName, err)
	}

	Info.Printf("%s listening on port %s", serviceName, port)
}

func getCollectionAlbumTags(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	Info.Printf("%s - ContentType:%s, Verb:%s, QueryString:%s, Data:%+v", serviceName, in.ContentType, in.Verb, in.QueryString, string(in.Data))
	
	utils.GetBlobMetadata(&azblob.Client{}, containerName)
}

func getCollections(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	Info.Printf("%s - ContentType:%s, Verb:%s, QueryString:%s, Data:%+v", serviceName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	query := "SELECT DISTINCT c['value'].collection AS name FROM c"

	queryResult, err := queryCollections(ctx, dbConfig.DbKey, dbConfig.DbURL, dbConfig.DbName, "photos", "/url", query, azcosmos.QueryOptions{})
	if err != nil {
		Error.Printf("query error: %v", err)
	}

	fmt.Print("getCollections - queryResult: ", queryResult)
	bytArr, _ := json.Marshal(queryResult)

	out = &common.Content{
		Data:        bytArr,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return out, nil
}

func getAlbums(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	Info.Printf("%s - ContentType:%s, Verb:%s, QueryString:%s, Data:%+v", serviceName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	query := "SELECT DISTINCT c['value'].album AS name FROM c"

	queryResult, err := queryAlbums(ctx, dbConfig.DbKey, dbConfig.DbURL, dbConfig.DbName, "photos", "/url", query, azcosmos.QueryOptions{})
	if err != nil {
		Error.Printf("query error: %v", err)
	}

	fmt.Print("getAlbums - queryResult: ", queryResult)
	bytArr, _ := json.Marshal(queryResult)

	out = &common.Content{
		Data:        bytArr,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return out, nil
}

func getPhotos(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	Info.Printf("%s - ContentType:%s, Verb:%s, QueryString:%s, Data:%+v", serviceName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	q, err := url.ParseQuery(in.QueryString)
	if err != nil {
		panic(err)
	}

	options := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@collection", Value: q.Get("collection")},
			{Name: "@album", Value: q.Get("album")},
		},
	}

	query := "SELECT * FROM p where p[\"value\"].collection = @collection AND p[\"value\"].album = @album"

	queryResult, err := queryPhotos(ctx, dbConfig.DbKey, dbConfig.DbURL, dbConfig.DbName, "photos", "/url", query, options)
	if err != nil {
		Error.Printf("query error: %v", err)
	}

	fmt.Print("getAlbumPhotos - queryResult: ", queryResult)
	bytArr, _ := json.Marshal(queryResult)

	out = &common.Content{
		Data:        bytArr,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return out, nil
}

func queryAlbums(ctx context.Context, accountKey string, accountEndpointURL string, databaseId string, containerId string, partitionKey string, query string, options azcosmos.QueryOptions) (response []models.Album, err error) {
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	cred, err := azcosmos.NewKeyCredential(accountKey)
	if err != nil {
		Error.Printf("failed to create cosmos credential: %s", err)
	}

	client, err := azcosmos.NewClientWithKey(accountEndpointURL, cred, nil)
	if err != nil {
		Error.Printf("failed to create cosmos client: %s", err)
	}

	container, err := client.NewContainer(databaseId, containerId)
	if err != nil {
		Error.Printf("failed to get container: %v with error: %s", containerId, err)
	}

	var resp []models.Album
	queryPager := container.NewQueryItemsPager(query, pk, &options)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			Error.Printf("failed to query cosmosdb: %s", err)
			return nil, err
		}

		for _, item := range queryResponse.Items {
			d := models.Album{}
			err := json.Unmarshal(item, &d)
			if err != nil {
				Error.Printf("failed to Unmarshal type to json: %s", err)
				return nil, err
			}
			resp = append(resp, d)
		}
	}

	return resp, nil
}

func queryCollections(ctx context.Context, accountKey string, accountEndpointURL string, databaseId string, containerId string, partitionKey string, query string, options azcosmos.QueryOptions) (response []models.Album, err error) {
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	cred, err := azcosmos.NewKeyCredential(accountKey)
	if err != nil {
		Error.Printf("failed to create cosmos credential: %s", err)
	}

	client, err := azcosmos.NewClientWithKey(accountEndpointURL, cred, nil)
	if err != nil {
		Error.Printf("failed to create cosmos client: %s", err)
	}

	container, err := client.NewContainer(databaseId, containerId)
	if err != nil {
		Error.Printf("failed to get container: %v with error: %s", containerId, err)
	}

	var resp []models.Album
	queryPager := container.NewQueryItemsPager(query, pk, &options)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			Error.Printf("failed to query cosmosdb: %s", err)
			return nil, err
		}

		for _, item := range queryResponse.Items {
			d := models.Album{}
			err := json.Unmarshal(item, &d)
			if err != nil {
				Error.Printf("failed to Unmarshal type to json: %s", err)
				return nil, err
			}
			resp = append(resp, d)
		}
	}

	return resp, nil
}

func queryPhotos(ctx context.Context, accountKey string, accountEndpointURL string, databaseId string, containerId string, partitionKey string, query string, options azcosmos.QueryOptions) (response []models.CosmosDocument, err error) {
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	cred, err := azcosmos.NewKeyCredential(accountKey)
	if err != nil {
		Error.Printf("failed to create cosmos credential: %s", err)
	}

	client, err := azcosmos.NewClientWithKey(accountEndpointURL, cred, nil)
	if err != nil {
		Error.Printf("failed to create cosmos client: %s", err)
	}

	container, err := client.NewContainer(databaseId, containerId)
	if err != nil {
		Error.Printf("failed to get container: %v with error: %s", containerId, err)
	}

	var resp []models.CosmosDocument
	queryPager := container.NewQueryItemsPager(query, pk, &options)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			Error.Printf("failed to query cosmosdb: %s", err)
			return nil, err
		}

		for _, item := range queryResponse.Items {
			d := models.CosmosDocument{}
			err := json.Unmarshal(item, &d)
			if err != nil {
				Error.Printf("failed to Unmarshal type to json: %s", err)
				return nil, err
			}
			resp = append(resp, d)
		}
	}

	return resp, nil
}

/* func queryCosmosDb(ctx context.Context, accountKey string, accountEndpointURL string, databaseId string, containerId string, partitionKey string, queryBody models.Query) (response []models.CosmosDocument, err error) {
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	cred, err := azcosmos.NewKeyCredential(accountKey)
	if err != nil {
		Error.Printf("failed to create cosmos credential: %s", err)
	}

	client, err := azcosmos.NewClientWithKey(accountEndpointURL, cred, nil)
	if err != nil {
		Error.Printf("failed to create cosmos client: %s", err)
	}

	container, err := client.NewContainer(databaseId, containerId)
	if err != nil {
		Error.Printf("failed to get container: %v with error: %s", containerId, err)
	}

	options := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@collection", Value: queryBody.Collection},
			{Name: "@album", Value: queryBody.Album},
		},
	}

	q := "SELECT * FROM p where p[\"value\"].collection = @collection AND p[\"value\"].album = @album"

	var resp []models.CosmosDocument
	queryPager := container.NewQueryItemsPager(q, pk, &options)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			Error.Printf("failed to query cosmosdb: %s", err)
			return nil, err
		}

		for _, item := range queryResponse.Items {
			d := models.CosmosDocument{}
			err := json.Unmarshal(item, &d)
			if err != nil {
				Error.Printf("failed to Unmarshal type to json: %s", err)
				return nil, err
			}
			resp = append(resp, d)
		}
	}

	return resp, nil
} */

/* func queryStateDataAlpha(ctx context.Context, query string) (r *dapr.QueryResponse, err error) {
	client, err := dapr.NewClient()
	if err != nil {
		Error.Printf("%s: failed to create client: %v", serviceName, err)
		return nil, err
	}

	Info.Printf("query: %s", query)

	r, err = client.QueryStateAlpha1(ctx, cosmosBinding, query, nil)
	if err != nil {
		Error.Printf("%s: failed to query state: %v", serviceName, err)
		return nil, err
	}

	return r, nil
} */
