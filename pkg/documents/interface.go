package documents

import (
	"time"
)

type StoredDocument struct {
	ID             string            `bson:"id"`
	BucketID       string            `bson:"bucketId"`
	CreatedAt      time.Time         `bson:"createdAt"`
	DeletedAt      *time.Time        `bson:"deletedAt"`
	CreatedBy      string            `bson:"createdBy"`
	UpdatedBy      string            `bson:"updatedBy"`
	DeletedBy      string            `bson:"deletedBy"`
	ContentType    string            `bson:"contentType"`
	ContentLength  int32             `bson:"contentLength"`
	SHA512Checksum string            `bson:"sha512Checksum"`
	MD5Checksum    string            `bson:"md5Checksum"`
	IsDeleted      bool              `bson:"isDeleted"`
	IsLastRevision bool              `bson:"isLastRevision"`
	Revision       int               `bson:"revision"`
	Metadata       map[string]string `bson:"metadata"`
	Data           interface{}       `bson:"data"`
}

type Document struct {
	ID             string
	BucketId       string
	CreatedAt      time.Time
	DeletedAt      *time.Time
	CreatedBy      string
	UpdatedBy      string
	DeletedBy      string
	ContentType    string
	ContentLength  int32
	SHA512Checksum string
	MD5Checksum    string
	IsLastRevision bool
	Revision       int
	Metadata       map[string]string
	Data           []byte
}

type Bucket struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}
