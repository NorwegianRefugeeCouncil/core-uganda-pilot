package documents

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/generic/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func Get(
	databaseName string,
	collectionName string,
	mongoClientFn func() (*mongo.Client, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		id := getObjectIDFromPath(req.URL.Path)
		if !strings.HasPrefix(id, "/") {
			id = fmt.Sprintf("/%s", id)
		}

		bucketId, err := getBucketIdFromHeader(req.URL.Query())
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		mongoCli, err := mongoClientFn()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		// ensure bucket exists
		_, err = getBucket(ctx, mongoCli, databaseName, bucketId)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		collection := mongoCli.Database(databaseName).Collection(collectionName)

		findOneResult := collection.FindOne(ctx, bson.M{
			"id":             id,
			"bucketId":       bucketId,
			"isDeleted":      false,
			"isLastRevision": true,
		})
		err = findOneResult.Err()
		if err != nil {
			if errors.Is(mongo.ErrNoDocuments, err) {
				writeError(w, http.StatusNotFound, fmt.Errorf("object not found"))
				return
			} else {
				writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to get object: %v", err))
				return
			}
		}

		doc := &StoredDocument{}
		if err := findOneResult.Decode(doc); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to unmarshal object: %v", err))
			return
		}

		w.Header().Set("Content-Type", doc.ContentType)
		w.Header().Set("Content-Length", strconv.Itoa(int(doc.ContentLength)))
		w.Header().Set("ETag", doc.MD5Checksum)
		w.Header().Set("Last-Modified", getLastModified(doc.CreatedAt))
		w.Header().Set("x-bucket-id", bucketId)
		w.Header().Set("Location", path.Join("https://%s/%s", req.Host, server.DocumentsEndpoint))
		w.WriteHeader(http.StatusOK)

		decoded, err := decodeData(doc.Data, doc.ContentType)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to decode data: %v", err))
			return
		}

		w.Write(decoded)

	}
}
