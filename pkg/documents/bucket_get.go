package documents

import (
	"context"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

func GetBucket(
	dbFactory storage.Factory,
	databaseName string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		id := req.URL.Path
		if strings.HasPrefix(id, server.BucketsEndpoint) {
			id = strings.TrimPrefix(id, server.BucketsEndpoint)
		}
		if strings.HasPrefix(id, "/") {
			id = strings.TrimPrefix(id, "/")
		}

		db, err := dbFactory.New()
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		bucket, err := getBucket(ctx, db, databaseName, id)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, &bucket)

	}
}

func getBucket(ctx context.Context, db *mongo.Client, databaseName string, id string) (*Bucket, error) {
	collection := db.Database(databaseName).Collection(collBuckets)
	result := collection.FindOne(ctx, bson.M{"id": id})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, bucketNotFound(id)
		} else {
			return nil, meta.NewInternalServerError(result.Err())
		}
	}
	var bucket Bucket
	if err := result.Decode(&bucket); err != nil {
		return nil, meta.NewInternalServerError(fmt.Errorf("failed to decode bucket: %v", err))
	}
	return &bucket, nil
}
