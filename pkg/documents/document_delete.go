package documents

import (
	"fmt"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

func Delete(
	mongoFn func() (*mongo.Client, error),
	databaseName string,
	timeTeller utils.TimeTeller,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		id := getObjectIDFromPath(req.URL.Path)

		bucketId, err := getBucketIdFromHeader(req.URL.Query())
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		versionIdStr := req.URL.Query().Get("version")

		mongoCli, err := mongoFn()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("could not connect to the database: %v", err))
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
			"isDeleted": false,
		}

		if len(versionIdStr) > 0 {
			versionId, err := strconv.Atoi(versionIdStr)
			if err != nil {
				writeError(w, http.StatusBadRequest, fmt.Errorf("could not parse version query parameter: %v", err))
				return
			}
			filter["revision"] = versionId
		} else {
			filter["isLastRevision"] = true
		}

		updateResult, err := collection.UpdateOne(ctx,
			filter, bson.M{
				"$set": bson.M{
					"isDeleted": true,
					"deletedAt": timeTeller.TellTime(),
				},
			})
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("could not delete the object: %v", err))
			return
		}
		if updateResult.ModifiedCount == 0 {
			writeError(w, http.StatusNotFound, fmt.Errorf("object not found"))
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
