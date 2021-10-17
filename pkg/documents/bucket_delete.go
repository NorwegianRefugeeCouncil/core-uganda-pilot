package documents

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
)

func DeleteBucket(
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

		mongoCli, err := dbFactory.New()
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		collection := mongoCli.Database(databaseName).Collection(collBuckets)

		deleteRes, err := collection.DeleteOne(ctx, bson.M{"id": id})
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to delete bucket: %v", err))
			return
		}
		if deleteRes.DeletedCount == 0 {
			utils.ErrorResponse(w, bucketNotFound(id))
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}

func bucketNotFound(id string) *meta.StatusError {
	return meta.NewNotFound(GroupVersion.WithResource("buckets"), id)
}
