package documents

import (
	"fmt"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

const bucketCollName = "document_buckets"

func CreateBucket(
	mongoClientFn func() (*mongo.Client, error),
	databaseName string,
	generator utils.UIDGenerator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		var bucket Bucket
		if err := utils.BindJSON(req, &bucket); err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to decode bucket from request body: %v", err))
			return
		}
		bucket.ID = generator.GenUID()

		mongoCli, err := mongoClientFn()
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		collection := mongoCli.Database(databaseName).Collection(bucketCollName)

		if _, err := collection.InsertOne(ctx, bucket); err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to create bucket: %v", err))
			return
		}

		utils.JSONResponse(w, http.StatusOK, bucket)

	}
}
