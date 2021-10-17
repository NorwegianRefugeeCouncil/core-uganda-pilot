package documents

import (
	"context"
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/generic/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//getDocumentMediaTypeFromHTTPHeader retrieves and parses the media type from the given header.
func getDocumentMediaTypeFromHTTPHeader(header http.Header) (string, map[string]string, error) {
	mediaType, params, err := mime.ParseMediaType(header.Get(headerContentType))
	if err != nil {
		return "", nil, meta.NewBadRequest(fmt.Sprintf("failed to parse media type: %v", err))
	}
	return mediaType, params, err
}

// getBucketIDFromQueryParam will ensure that the paramBucketID query parameter
// is present, and will return it or an error
func getBucketIDFromQueryParam(values url.Values) (string, error) {
	b := values.Get(paramBucketID)
	if len(b) == 0 {
		return "", meta.NewBadRequest(fmt.Sprintf("request parameter '%s' is required", paramBucketID))
	}
	return b, nil
}

// findDocumentVersionFromQueryParam will retrieve the document ID from the
// given query parameters, will make sure that it is well-formed and will
// return it.
func findDocumentVersionFromQueryParam(values url.Values) (*int64, error) {
	v := values.Get(paramVersion)
	if len(v) > 0 {
		version, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, meta.NewBadRequest(fmt.Sprintf("failed to parse version: %s", v))
		}
		return &version, nil
	}
	return nil, nil
}

// prepareDocumentDataForStorage will encode the given document data in the right format.
// If the format it json, it will convert it to a json object.
// this is used to store in the database
func prepareDocumentDataForStorage(bodyBytes []byte, contentType string) (interface{}, error) {
	var dataIntf interface{}
	dataIntf = bodyBytes
	if isMediaType(contentType, mimeTypeApplicationJson) {
		dataMap := map[string]interface{}{}
		if err := json.Unmarshal(bodyBytes, &dataMap); err != nil {
			return nil, meta.NewBadRequest(fmt.Sprintf("failed to decode json: %v", err))
		}
		dataIntf = dataMap
	} else if isMediaType(contentType, mimeTypeTextPlain, mimeTypeTextHtml) {
		dataIntf = string(bodyBytes)
	} else {
		base64Data := base64.StdEncoding.EncodeToString(bodyBytes)
		dataIntf = base64Data
	}
	return dataIntf, nil
}

// transformDocumentDataFromStorage will convert the given data to bytes.
// this is used to retrieve data from the database
func transformDocumentDataFromStorage(data interface{}, contentType string) ([]byte, error) {
	if isMediaType(contentType, mimeTypeApplicationJson) {
		dataMap := data.(bson.D)
		bsonBytes, err := bson.Marshal(dataMap)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal bson data: %v", err)
		}
		jsonMap := map[string]interface{}{}
		err = bson.Unmarshal(bsonBytes, &jsonMap)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal json data: %v", err)
		}
		jsonBytes, err := json.Marshal(jsonMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json data: %v", err)
		}
		return jsonBytes, nil
	} else if isMediaType(contentType, mimeTypeTextPlain, mimeTypeTextHtml) {
		dataStr := data.(string)
		return []byte(dataStr), nil
	} else {
		dataStr := data.(string)
		base64Bytes, err := base64.StdEncoding.DecodeString(dataStr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 data: %v", err)
		}
		return base64Bytes, nil
	}
}

// isMediaType will check if the given content-type is one of the given media types
func isMediaType(contentType string, mimeTypes ...string) bool {
	for _, mimeType := range mimeTypes {
		if strings.HasPrefix(contentType, mimeType) {
			return true
		}
	}
	return false
}

// getContentLengthFromHTTPHeader will retrieve the headerContentLength header from the request
func getContentLengthFromHTTPHeader(header http.Header) (int32, error) {
	contentLengthStr := header.Get(headerContentLength)
	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 32)
	if err != nil {
		return 0, meta.NewBadRequest(fmt.Sprintf("failed to parse Content-Length: %v", err))
	}
	return int32(contentLength), err
}

// getDocumentMetadataFromHTTPHeader will extract the document metadata from the supplied headers
// document metadata is in the form of
// x-tags key1=value1,key2=value2,...
func getDocumentMetadataFromHTTPHeader(header http.Header) (map[string]string, error) {
	if header == nil {
		return map[string]string{}, nil
	}
	tags := header.Get(headerTags)
	if len(tags) == 0 {
		return map[string]string{}, nil
	}
	tagParts := strings.Split(tags, ",")
	metadata := map[string]string{}
	for _, part := range tagParts {
		entryParts := strings.Split(part, "=")
		if len(entryParts) != 2 {
			reason := fmt.Sprintf("invalid tags: %s", tags)
			return nil, meta.NewBadRequest(reason)
		}
		metadata[entryParts[0]] = entryParts[1]
	}
	return metadata, nil
}

// formatHTTPLastModified formats the given time
// to the HTTP (time.RFC1123) format
func formatHTTPLastModified(t time.Time) string {
	return t.Format(http.TimeFormat)
}

// parseHTTPLastModified parses the given string
// to time.Time using (time.RFC1123) format
func parseHTTPLastModified(t string) (time.Time, error) {
	return time.Parse(http.TimeFormat, t)
}

// getMD5Checksum gets the md5 checksum for the given bytes
func getMD5Checksum(bodyBytes []byte) string {
	md5Checksum := md5.Sum(bodyBytes)
	md5ChecksumStr := hex.EncodeToString(md5Checksum[:])
	return md5ChecksumStr
}

// getSha512Checksum gets the SHA512 checksum from the given bytes
func getSha512Checksum(bodyBytes []byte) string {
	sha512Checksum := sha512.Sum512(bodyBytes)
	sha512ChecksumStr := hex.EncodeToString(sha512Checksum[:])
	return sha512ChecksumStr
}

//getDocumentIdFromHTTPRequestPath retrieves the Document.ID from the path
func getDocumentIdFromHTTPRequestPath(objectPath string) string {
	id := objectPath
	if strings.HasPrefix(id, path.Join(server.DocumentsEndpoint)) {
		id = strings.TrimPrefix(id, server.DocumentsEndpoint)
	}
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}
	return id
}

// documentIDRegex is the regex for valid document keys
var documentIDRegex = regexp.MustCompile("^([a-zA-Z0-9]+([/\\.\\-_][[a-zA-Z0-9]+)*)$")

// validateDocumentID validates that the given document ID is valid
func validateDocumentID(id string) error {

	if strings.HasPrefix(id, "/") {
		id = strings.TrimPrefix(id, "/")
	}

	parts := strings.Split(id, "/")
	for _, part := range parts {
		if len(part) == 0 {
			return meta.NewBadRequest(fmt.Sprintf("invalid object key: %s", id))
		}
	}

	if !documentIDRegex.MatchString(id) {
		return meta.NewBadRequest(fmt.Sprintf("invalid object key format: %s", id))
	}

	return nil

}

// getDocumentDBFilter gets the appropriate database filter for the given DocumentRef
func getDocumentDBFilter(docRef DocumentRef) interface{} {
	filter := bson.M{
		keyID:        docRef.GetKey(),
		keyIsDeleted: false,
		keyBucketID:  docRef.GetBucketID(),
	}
	if docRef.HasVersion() {
		filter[keyRevision] = docRef.GetVersion()
	} else {
		filter[keyIsLastRevision] = true
	}
	return filter
}

// assertDocumentBucketExists checks that the given bucket exists
func assertDocumentBucketExists(ctx context.Context, db *mongo.Client, databaseName string, docRef DocumentRef) error {
	// ensure bucket exists
	_, err := getBucket(ctx, db, databaseName, docRef.GetBucketID())
	return err
}

// removeLeadingTrailingSlashes will normalize trailing/leading slashes
func removeLeadingTrailingSlashes(key string) string {
	if strings.HasPrefix(key, "/") {
		key = strings.TrimPrefix(key, "/")
	}
	if strings.HasSuffix(key, "/") {
		key = strings.TrimSuffix(key, "/")
	}
	return key
}

//getDocumentRefFromHTTPRequest returns the DocumentRef from the given request
// it will parse the bucketId, documentId and version if present
func getDocumentRefFromHTTPRequest(req *http.Request) (DocumentRef, error) {
	id := getDocumentIdFromHTTPRequestPath(req.URL.Path)

	if err := validateDocumentID(id); err != nil {
		return DocumentRef{}, err
	}

	bucketId, err := getBucketIDFromQueryParam(req.URL.Query())
	if err != nil {
		return DocumentRef{}, err
	}

	objectVersion, err := findDocumentVersionFromQueryParam(req.URL.Query())
	if err != nil {
		return DocumentRef{}, err
	}

	return NewDocumentRef(bucketId, id, objectVersion), nil
}
