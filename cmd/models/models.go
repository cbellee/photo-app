package models

import (
	"time"

	"github.com/google/uuid"
)

type PhotoState struct {
	Key   string `json:"key"`
	Value Photo  `json:"value"`
}

type DbConfig struct {
	DbURL       string
	DbName      string
	DbKey       string
	DbContainer string
}

type CosmosDocument struct {
	Attachments  string `json:"_attachments"`
	ETag         string `json:"_etag"`
	Rid          string `json:"_rid"`
	Self         string `json:"_self"`
	TimeStamp    int64  `json:"_ts"`
	Id           string `json:"id"`
	IsBinary     bool   `json:"isBinary"`
	PartitionKey string `json:"partitionKey"`
	Value        Photo  `json:"value"`
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

type Album struct {
	Name string `json:"name"`
}

type Collection struct {
	Name string `json:"name"`
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
