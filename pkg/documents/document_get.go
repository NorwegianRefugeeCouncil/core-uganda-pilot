package documents

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

func Get(
	databaseName string,
	mongoClientFn func() (*mongo.Client, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		id := getObjectIDFromPath(req.URL.Path)

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

		collection := mongoCli.Database(databaseName).Collection(collDocuments)

		filter := bson.M{
			"id":        id,
			"bucketId":  bucketId,
			"isDeleted": false,
		}

		versionStr := req.URL.Query().Get("version")
		if len(versionStr) > 0 {
			version, err := strconv.Atoi(versionStr)
			if err != nil {
				writeError(w, http.StatusNotFound, fmt.Errorf("could not parse version: %v", err))
				return
			}
			filter["revision"] = version
		} else {
			filter["isLastRevision"] = true
		}

		findOneResult := collection.FindOne(ctx, filter)
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
		w.WriteHeader(http.StatusOK)

		decoded, err := decodeData(doc.Data, doc.ContentType)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to decode data: %v", err))
			return
		}

		w.Write(decoded)

	}
}
