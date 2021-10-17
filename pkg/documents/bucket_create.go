package documents

import (
	"fmt"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
	"regexp"
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

		mongoCli, err := dbFactory.New()
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		collection := mongoCli.Database(databaseName).Collection(collBuckets)

		if _, err := collection.InsertOne(ctx, bucket); err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to create bucket: %v", err))
			return
		}

		utils.JSONResponse(w, http.StatusOK, bucket)

	}
}

func validateBucket(fieldPath validation.Path, bucket *Bucket) validation.ErrorList {
	result := validation.ErrorList{}
	result = append(result, validateBucketName(fieldPath.Key("name"), bucket.Name)...)
	return result
}

var bucketNameRegex = regexp.MustCompile("^[a-zA-Z0-9]+([_-][a-zA-Z0-9]+)*$")
var bucketNameMinLength = 3
var bucketNameMaxLength = 128

func validateBucketName(fieldPath *validation.Path, name string) validation.ErrorList {
	result := validation.ErrorList{}

	if len(name) == 0 {
		result = append(result, validation.Required(fieldPath, name))
	}
	if len(name) < bucketNameMinLength {
		result = append(result, validation.TooShort(fieldPath, name, bucketNameMinLength))
	}
	if len(name) > bucketNameMaxLength {
		result = append(result, validation.TooLong(fieldPath, name, bucketNameMaxLength))
	}
	if !bucketNameRegex.MatchString(name) {
		result = append(result, validation.Invalid(fieldPath, name, "invalid character sequence"))
	}

	return result
}
