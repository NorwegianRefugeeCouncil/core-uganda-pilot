package documents

import (
	"fmt"
	"net/http"
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

type documentRef struct {
	key      string
	bucketID string
	version  *int64
}

func (d documentRef) GetKey() string {
	return d.key
}

func (d documentRef) GetBucketID() string {
	return d.bucketID
}

func (d documentRef) HasVersion() bool {
	return d.version != nil
}

func (d documentRef) GetVersion() int64 {
	if d.version != nil {
		return *d.version
	}
	return 0
}

func (d documentRef) String() string {
	if d.HasVersion() {
		return fmt.Sprintf("%s, version=%d, bucket=%s", d.key, d.version, d.bucketID)
	} else {
		return fmt.Sprintf("%s, bucket=%s", d.key, d.bucketID)
	}

}

type DocumentRef interface {
	fmt.Stringer
	GetKey() string
	GetBucketID() string
	HasVersion() bool
	GetVersion() int64
}

func NewDocumentVersionRef(bucketID, key string, version *int64) DocumentRef {
	return &documentRef{
		key:      key,
		bucketID: bucketID,
		version:  version,
	}
}

func getDocumentRefFromReq(req *http.Request) (DocumentRef, error) {
	id := getObjectIDFromPath(req.URL.Path)

	if err := validateObjectId(id); err != nil {
		return nil, err
	}

	bucketId, err := requireBucketIDFromQueryParam(req.URL.Query())
	if err != nil {
		return nil, err
	}

	objectVersion, err := findObjectVersionFromQueryParam(req.URL.Query())
	if err != nil {
		return nil, err
	}

	return NewDocumentVersionRef(bucketId, id, objectVersion), nil
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
