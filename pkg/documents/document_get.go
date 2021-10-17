package documents

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

func Get(
	databaseName string,
	storageFactory storage.Factory,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		docRef, err := getDocumentRefFromReq(req)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		db, err := storageFactory.New()
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		if err := ensureBucketExists(ctx, db, databaseName, docRef); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		collection := db.Database(databaseName).Collection(collDocuments)
		filter := getDocumentFilter(docRef)
		findOneResult := collection.FindOne(ctx, filter)
		if findOneResult.Err() != nil {
			if errors.Is(findOneResult.Err(), mongo.ErrNoDocuments) {
				utils.ErrorResponse(w, docNotFound(docRef))
				return
			} else {
				utils.ErrorResponse(w, meta.NewInternalServerError(findOneResult.Err()))
				return
			}
		}

		doc := &StoredDocument{}
		if err := findOneResult.Decode(doc); err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("failed to decode document: %v", err)))
			return
		}

		w.Header().Set(headerContentType, doc.ContentType)
		w.Header().Set(headerContentLength, strconv.Itoa(int(doc.ContentLength)))
		w.Header().Set(headerETag, doc.MD5Checksum)
		w.Header().Set(headerLastModified, getLastModified(doc.CreatedAt))
		w.Header().Set(headerBucketID, docRef.GetBucketID())
		w.WriteHeader(http.StatusOK)

		decoded, err := decodeData(doc.Data, doc.ContentType)
		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("failed to decode document data: %v", err)))
			return
		}

		w.Write(decoded)

	}
}
