package documents

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core/pkg/generic/server"
	"go.mongodb.org/mongo-driver/bson"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(err.Error()))
	return
}

func getMediaType(header http.Header) (string, map[string]string, error) {
	mediaType, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse media type: %v", err)
	}
	return mediaType, params, err
}

func getBucketIdFromHeader(header url.Values) (string, error) {
	b := header.Get("bucketId")
	if len(b) == 0 {
		return "", fmt.Errorf("bucketId request parameter is required")
	}
	return b, nil
}

func encodeData(bodyBytes []byte, contentType string) (interface{}, error) {
	var dataIntf interface{}
	dataIntf = bodyBytes
	if isMediaType(contentType, "application/json") {
		dataMap := map[string]interface{}{}
		if err := json.Unmarshal(bodyBytes, &dataMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json data: %v", err)
		}
		dataIntf = dataMap
	} else if isMediaType(contentType, "text/plain", "text/html") {
		dataIntf = string(bodyBytes)
	} else {
		base64Data := base64.StdEncoding.EncodeToString(bodyBytes)
		dataIntf = base64Data
	}
	return dataIntf, nil
}

func decodeData(data interface{}, contentType string) ([]byte, error) {
	if isMediaType(contentType, "application/json") {
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
	} else if isMediaType(contentType, "text/plain", "text/html") {
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

func isMediaType(contentType string, mimeTypes ...string) bool {
	for _, mimeType := range mimeTypes {
		if strings.HasPrefix(contentType, mimeType) {
			return true
		}
	}
	return false
}

func getContentLength(req *http.Request) (int32, error) {
	contentLengthStr := req.Header.Get("Content-Length")
	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Content-Length: %v", err)
	}
	return int32(contentLength), err
}

func getMetadata(header http.Header) (map[string]string, error) {
	if header == nil {
		return map[string]string{}, nil
	}
	tags := header.Get("x-tags")
	if len(tags) == 0 {
		return map[string]string{}, nil
	}
	tagParts := strings.Split(tags, ",")
	metadata := map[string]string{}
	for _, part := range tagParts {
		entryParts := strings.Split(part, "=")
		if len(entryParts) != 2 {
			return nil, fmt.Errorf("invalid tags")
		}
		metadata[entryParts[0]] = entryParts[1]
	}
	return metadata, nil
}

func getLastModified(t time.Time) string {
	return t.Format(http.TimeFormat)
}

func parseLastModified(t string) (time.Time, error) {
	return time.Parse(http.TimeFormat, t)
}

func getMD5Checksum(bodyBytes []byte) string {
	md5Checksum := md5.Sum(bodyBytes)
	md5ChecksumStr := hex.EncodeToString(md5Checksum[:])
	return md5ChecksumStr
}

func getSha512Checksum(bodyBytes []byte) string {
	sha512Checksum := sha512.Sum512(bodyBytes)
	sha512ChecksumStr := hex.EncodeToString(sha512Checksum[:])
	return sha512ChecksumStr
}

func getObjectIDFromPath(objectPath string) string {
	id := objectPath
	if strings.HasPrefix(id, path.Join(server.DocumentsEndpoint)) {
		id = strings.TrimPrefix(id, server.DocumentsEndpoint)
	}
	return id
}
