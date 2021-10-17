package documents

import (
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

		docRef, err := getDocumentRefFromReq(req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		if docRef.HasVersion() {
			writeError(w, http.StatusBadRequest, fmt.Errorf("query parameter '%s' is illegal for %s operation", paramVersion, http.MethodPut))
			return
		}

		mediaType, mediaTypeParams, err := getMediaType(req.Header)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("failed to get media type: %v", err))
			return
		}

		formattedMediaType := mime.FormatMediaType(mediaType, mediaTypeParams)

		contentLength, err := getContentLength(req)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("failed to get content-length: %v", err))
			return
		}

		metadata, err := getMetadata(req.Header)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("failed to get tags: %v", err))
			return
		}

		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to upload object: %v", err.Error()))
			return
		}

		sha512ChecksumStr := getSha512Checksum(bodyBytes)
		md5ChecksumStr := getMD5Checksum(bodyBytes)

		dataIntf, err := encodeData(bodyBytes, mediaType)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("failed to get data: %v", err.Error()))
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
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		// ensure bucket exists
		if err := ensureBucketExists(ctx, db, databaseName, docRef); err != nil {
			writeError(w, http.StatusInternalServerError, err)
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

			collection := db.Database(databaseName).Collection(collDocuments)
			// Update the previous version if it exists
			result := collection.FindOneAndUpdate(ctx, getDocumentFilter(docRef), bson.M{
				"$set": bson.M{
					keyIsLastRevision: false,
				},
			})

			if result.Err() != nil {
				if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
					return nil, docNotFound(docRef)
				}
			} else {
				oldDoc := StoredDocument{}
				if err := result.Decode(&oldDoc); err != nil {
					return nil, meta.NewInternalServerError(err)
				}
				doc.Revision = oldDoc.Revision + 1
			}

			// Insert new version
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
