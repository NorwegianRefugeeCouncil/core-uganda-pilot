package documents

import (
	"context"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
)

func Put(
	timeTeller utils.TimeTeller,
	dbFactory storage.Factory,
	databaseName string,
) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		docRef, err := getDocumentRefFromHTTPRequest(req)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		if docRef.HasVersion() {
			reason := fmt.Sprintf("query parameter '%s' is illegal for %s operation", paramVersion, http.MethodPut)
			utils.ErrorResponse(w, meta.NewBadRequest(reason))
			return
		}

		mediaType, mediaTypeParams, err := getDocumentMediaTypeFromHTTPHeader(req.Header)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		formattedMediaType := mime.FormatMediaType(mediaType, mediaTypeParams)

		contentLength, err := getContentLengthFromHTTPHeader(req.Header)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		metadata, err := getDocumentMetadataFromHTTPHeader(req.Header)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("failed to read body: %v", err)))
			return
		}

		sha512ChecksumStr := getSha512Checksum(bodyBytes)
		md5ChecksumStr := getMD5Checksum(bodyBytes)

		dataIntf, err := prepareDocumentDataForStorage(bodyBytes, mediaType)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		doc := &StoredDocument{
			ID:             docRef.GetKey(),
			BucketID:       docRef.GetBucketID(),
			CreatedAt:      timeTeller.TellTime(),
			DeletedAt:      nil,
			CreatedBy:      "",
			UpdatedBy:      "",
			DeletedBy:      "",
			ContentType:    formattedMediaType,
			ContentLength:  contentLength,
			SHA512Checksum: sha512ChecksumStr,
			MD5Checksum:    md5ChecksumStr,
			IsLastRevision: true,
			Metadata:       metadata,
			Data:           dataIntf,
			Revision:       1,
		}

		db, err := dbFactory.New()
		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("database unreachable: %v", err)))
			return
		}

		// ensure bucket exists
		if err := assertDocumentBucketExists(ctx, db, databaseName, docRef); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		sess, err := db.StartSession()
		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("failed to start db session: %v", err.Error())))
			return
		}
		defer sess.EndSession(ctx)

		wc := writeconcern.New(writeconcern.WMajority())
		rc := readconcern.Snapshot()
		txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

		_, err = sess.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {

			collection := db.Database(databaseName).Collection(DocumentsCollection)

			if err := preparePuttingDocument(ctx, collection, doc); err != nil {
				return nil, err
			}

			if _, err := collection.InsertOne(ctx, doc); err != nil {
				return nil, meta.NewInternalServerError(err)
			}

			return nil, nil

		}, txnOpts)

		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		w.Header().Set(headerETag, md5ChecksumStr)
		w.Header().Set(headerObjectKey, doc.ID)
		w.Header().Set(headerObjectVersion, strconv.Itoa(doc.Revision))
		w.Header().Set(headerBucketID, docRef.GetBucketID())

	}
}

// preparePuttingDocument will check if a previous document at the given key exists.
// If yes, it will update that document to indicate that this is not the current version.
// The preparePuttingDocument will return the Document.Version for the next document.
func preparePuttingDocument(ctx context.Context, collection *mongo.Collection, doc *StoredDocument) error {

	docRef := doc.DocumentRef().WithCurrentVersion()
	filter := getDocumentDBFilter(docRef)

	result := collection.FindOneAndUpdate(ctx, filter, bson.M{
		"$set": bson.M{
			keyIsLastRevision: false,
		},
	})

	if result.Err() != nil {
		if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return docNotFound(docRef)
		}
	} else {
		oldDoc := StoredDocument{}
		if err := result.Decode(&oldDoc); err != nil {
			return meta.NewInternalServerError(err)
		}
		doc.Revision = oldDoc.Revision + 1
	}
	return nil
}
