package documents

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Delete(
	dbFactory storage.Factory,
	databaseName string,
	timeTeller utils.TimeTeller,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		docRef, err := getDocumentRefFromHTTPRequest(req)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		db, err := dbFactory.New()
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		// ensure bucket exists
		if err := assertDocumentBucketExists(ctx, db, databaseName, docRef); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		collection := db.Database(databaseName).Collection(DocumentsCollection)
		updateResult, err := collection.UpdateOne(ctx,
			getDocumentDBFilter(docRef),
			bson.M{
				"$set": bson.M{
					keyIsDeleted: true,
					keyDeletedAt: timeTeller.TellTime(),
				},
			})

		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("could not delete bucket: %v", err)))
			return
		}

		if updateResult.ModifiedCount == 0 {
			utils.ErrorResponse(w, docNotFound(docRef))
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}

func docNotFound(docRef DocumentRef) *meta.StatusError {
	return meta.NewNotFound(GroupVersion.WithResource("documents"), docRef.GetKey())
}
