package documents

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

func GetBucket(
	mongoClientFn func() (*mongo.Client, error),
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

		mongoCli, err := mongoClientFn()
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		bucket, err := getBucket(ctx, mongoCli, databaseName, id)
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to get bucket: %v", err))
			return
		}

		utils.JSONResponse(w, http.StatusOK, &bucket)

	}
}

func getBucket(ctx context.Context, mongoCli *mongo.Client, databaseName string, id string) (*Bucket, error) {
	collection := mongoCli.Database(databaseName).Collection(collBuckets)
	result := collection.FindOne(ctx, bson.M{
		keyID: id,
	})
	if result.Err() != nil {
		return nil, fmt.Errorf("failed to get bucket: %v", result.Err())
	}
	var bucket Bucket
	if err := result.Decode(&bucket); err != nil {
		return nil, fmt.Errorf("failed to decode bucket: %v", err)
	}
	return &bucket, nil
}
