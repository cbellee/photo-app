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

type StorageConfig struct {
	StorgeURL        string
	StorageAccount   string
	StorageKey       string
	StorageContainer string
}

type Blob struct {
	Name     string
	Prefix   string
	Data     string
	Tags     map[string]string
	MetaData map[string]string
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
	Id            uuid.UUID `json:"Id"`
	Name          string    `json:"Name"`
	Album         string    `json:"Album"`
	Collection    string    `json:"Collection"`
	Type          string    `json:"Type"`
	Url           string    `json:"Url"`
	ThumbsUrl     string    `json:"ThumbsUrl"`
	ContentLength int32     `json:"ContentLength"`
	Created       time.Time `json:"Created"`
	Modified      time.Time `json:"Modified"`
}

type Album struct {
	Name string `json:"Name"`
}

type Collection struct {
	Name string `json:"Name"`
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
