package documents

import (
	"fmt"
	"github.com/nrc-no/core/pkg/pointers"
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
	IsLastRevision bool              `bson:"isLatestVersion"`
	Revision       int               `bson:"resourceVersion"`
	Metadata       map[string]string `bson:"metadata"`
	Data           interface{}       `bson:"data"`
}

func (d *StoredDocument) DocumentRef() DocumentRef {
	return NewDocumentRef(d.BucketID, d.ID, pointers.Int64(int64(d.Revision)))
}

type DocumentRef struct {
	key      string
	bucketID string
	version  *int64
}

func (d DocumentRef) GetKey() string {
	return d.key
}

func (d DocumentRef) GetBucketID() string {
	return d.bucketID
}

func (d DocumentRef) HasVersion() bool {
	return d.version != nil
}

func (d DocumentRef) GetVersion() int64 {
	if d.version != nil {
		return *d.version
	}
	return 0
}

func (d DocumentRef) WithVersion(v int64) DocumentRef {
	d.version = &v
	return d
}

func (d DocumentRef) WithCurrentVersion() DocumentRef {
	d.version = nil
	return d
}

func (d DocumentRef) String() string {
	if d.HasVersion() {
		return fmt.Sprintf("%s, version=%d, bucket=%s", d.key, d.version, d.bucketID)
	} else {
		return fmt.Sprintf("%s, bucket=%s", d.key, d.bucketID)
	}

}

func NewDocumentRef(bucketID, key string, version *int64) DocumentRef {
	return DocumentRef{
		key:      key,
		bucketID: bucketID,
		version:  version,
	}
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
	ID             string `bson:"id"`
	Name           string `bson:"name"`
	IsLastRevision bool   `bson:"isLatestVersion"`
	Revision       int    `bson:"resourceVersion"`
	IsDeleted      bool   `bson:"isDeleted"`
}
