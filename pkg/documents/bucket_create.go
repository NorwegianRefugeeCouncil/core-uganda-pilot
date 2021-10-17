package documents

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func CreateBucket(
	dbFactory storage.Factory,
	databaseName string,
	generator utils.UIDGenerator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		var bucket Bucket
		if err := utils.BindJSON(req, &bucket); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		bucket.ID = generator.GenUID()
		bucket.Revision = 1
		bucket.IsLastRevision = true
		bucket.IsDeleted = false

		if errs := validateBucket(&bucket); errs.HasAny() {
			utils.ErrorResponse(w, meta.NewInvalid(GroupVersion.WithResource("bucket"), "", errs))
			return
		}

		mongoCli, err := dbFactory.New()
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		collection := mongoCli.Database(databaseName).Collection(BucketsCollection)

		if _, err := collection.InsertOne(ctx, bucket); err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to create bucket: %v", err))
			return
		}

		utils.JSONResponse(w, http.StatusOK, bucket)

	}
}
