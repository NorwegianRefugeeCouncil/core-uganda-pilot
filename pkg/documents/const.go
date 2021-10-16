package documents

const (
	DocumentsCollection = "documents"
	BucketsCollection   = "document_buckets"
)

var AllCollections = []string{
	DocumentsCollection,
	BucketsCollection,
}

const (
	headerContentType    = "Content-Type"
	headerContentLength  = "Content-Length"
	headerETag           = "ETag"
	headerLastModified   = "Last-Modified"
	headerBucketID       = "x-bucket-id"
	headerObjectVersion  = "x-object-version"
	headerObjectKey      = "x-object-key"
	headerSha512Checksum = "x-sha512-checksum"
	headerTags           = "x-tags"

	paramVersion  = "version"
	paramBucketID = "bucketId"

	keyID             = "id"
	keyBucketID       = "bucketId"
	keyIsLastRevision = "isLatestVersion"
	keyRevision       = "resourceVersion"
	keyIsDeleted      = "isDeleted"
	keyDeletedAt      = "deletedAt"

	mimeTypeApplicationJson = "application/json"
	mimeTypeTextPlain       = "text/plain"
	mimeTypeTextHtml        = "text/html"
)
