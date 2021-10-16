package documents

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
)

func Put(
	timeTeller utils.TimeTeller,
	mongoFn func() (*mongo.Client, error),
	databaseName string,
) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		id := getObjectIDFromPath(req.URL.Path)

		if err := validateObjectId(id); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid object key: %v", err))
			return
		}

		bucketId, err := getBucketIdFromHeader(req.URL.Query())
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("could not get bucketId: %v", err))
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
			ID:             id,
			BucketID:       bucketId,
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
		}

		mongoClient, err := mongoFn()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to connect to database: %v", err))
			return
		}

		// ensure bucket exists
		_, err = getBucket(ctx, mongoClient, databaseName, bucketId)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		session, err := mongoClient.StartSession()
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to start database session: %v", err))
			return
		}
		defer session.EndSession(ctx)

		_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

			collection := mongoClient.Database(databaseName).Collection(collDocuments)

			// Update the previous version if it exists
			result := collection.FindOneAndUpdate(sessCtx, bson.M{
				"id":             id,
				"isLastRevision": true,
			}, bson.M{
				"$set": bson.M{
					"isLastRevision": false,
				},
			})
			if result.Err() != nil {
				if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
					return nil, fmt.Errorf("failed to update previous document version: %v", err)
				}
			} else {
				oldDoc := StoredDocument{}
				if err := result.Decode(&oldDoc); err != nil {
					return nil, fmt.Errorf("failed to decode previous document: %v", err)
				}
				doc.Revision = oldDoc.Revision + 1
			}

			// Insert new version
			if _, err := collection.InsertOne(sessCtx, doc); err != nil {
				return nil, err
			}

			return nil, nil
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to save document: %v", err))
			return
		}

		w.Header().Set("ETag", md5ChecksumStr)
		w.Header().Set("x-object-key", doc.ID)
		w.Header().Set("x-object-version", strconv.Itoa(doc.Revision))
		w.Header().Set("x-object-bucket", bucketId)

	}
}
