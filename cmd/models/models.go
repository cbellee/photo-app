package models

import (
	"time"

	"github.com/google/uuid"
)

type PhotoState struct {
	Key   string `json:"key"`
	Value Photo  `json:"value"`
}

type Photo struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Album         string    `json:"album"`
	Collection    string    `json:"collection"`
	Type          string    `json:"type"`
	Url           string    `json:"url"`
	ThumbsUrl     string    `json:"thumbsUrl"`
	ContentLength int32     `json:"contentLength"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modified"`
}

/* `{
	"filter": {
		"EQ": { "album": "" }
	},
	"sort": [
		{
			"key": "person.id",
			"order": "DESC"
		}
	]
}` */

type filter struct {
	EQ string
}

type Query struct {
	Filter filter 
}

type Event struct {
	Topic           string
	Subject         string
	EventType       string
	Id              string
	DataVersion     string
	MetadataVersion string
	EventTime       string
	Data            struct {
		Api                string
		ClientRequestId    string
		RequestId          string
		ETag               string
		ContentType        string
		ContentLength      int32
		BlobType           string
		Url                string
		Sequencer          string
		StorageDiagnostics struct {
			BatchId string
		}
	}
}

/* {
    "topic": "/subscriptions/b2375b5f-8dab-4436-b87c-32bc7fdce5d0/resourceGroups/photo-app-rg/providers/Microsoft.Storage/storageAccounts/storjn45dot42kw3k",
    "subject": "/blobServices/default/containers/resize/blobs/scenario.png",
    "eventType": "Microsoft.Storage.BlobCreated",
    "id": "c2dcbc2f-001e-003e-718c-72aabe062a11",
    "data": {
        "api": "PutBlob",
        "clientRequestId": "0c893f6f-cbae-4a75-8c52-c72224d40038",
        "requestId": "c2dcbc2f-001e-003e-718c-72aabe000000",
        "eTag": "0x8DA40A3990F5C39",
        "contentType": "image/png",
        "contentLength": 165557,
        "blobType": "BlockBlob",
        "url": "https://storjn45dot42kw3k.blob.core.windows.net/resize/scenario.png",
        "sequencer": "000000000000000000000000000000C400000000000f2f1e",
        "storageDiagnostics": {
            "batchId": "185255d9-4006-0062-008c-72ffe6000000"
        }
    },
    "dataVersion": "",
    "metadataVersion": "1",
    "eventTime": "2022-05-28T12:14:21.81241Z"
} */
