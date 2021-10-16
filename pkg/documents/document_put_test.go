package documents

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func (s *Suite) TestPutDocumentWithNoBucketShouldThrow() {
	_, err := s.client.Documents().Put(context.Background(), &Document{
		ID:          "some-id",
		BucketId:    "",
		Data:        []byte(`{"a": "b"}`),
		ContentType: "application/json",
	}, PutDocumentOptions{})
	assert.Error(s.T(), err)
}

func (s *Suite) TestPutDocumentNonExistingBucketShouldThrow() {
	_, err := s.client.Documents().Put(context.Background(), &Document{
		ID:          "some-id",
		BucketId:    "non-existing",
		Data:        []byte(`{"a": "b"}`),
		ContentType: "application/json",
	}, PutDocumentOptions{})
	assert.Error(s.T(), err)
}

func (s *Suite) TestPutDocument() {

	ctx := context.Background()

	type args struct {
		name             string
		key              string
		contentType      string
		tags             string
		content          []byte
		expectStatusCode int
		expectRevision   int
	}

	bucket, err := s.client.Buckets().Create(ctx, &Bucket{Name: "test-document-get"}, CreateBucketOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	existingObj := &Document{
		ID:          "/put-existing",
		ContentType: "application/json",
		BucketId:    bucket.ID,
		Data:        []byte(`{"a":"b"}`),
	}

	if _, err := s.client.Documents().Put(context.Background(), existingObj, PutDocumentOptions{}); !assert.NoError(s.T(), err) {
		return
	}

	tcs := []args{
		{
			name:             "putJsonDocument",
			key:              "/json-document",
			contentType:      "application/json",
			content:          []byte(`{"a":"b"}`),
			expectStatusCode: 200,
			expectRevision:   0,
		}, {
			name:             "putTextDocument",
			key:              "/text-document",
			contentType:      "text/plain",
			content:          []byte(`hello`),
			expectStatusCode: 200,
			expectRevision:   0,
		}, {
			name:             "putHtmlDocument",
			key:              "/html-document",
			contentType:      "text/html",
			content:          []byte(`<html></html>`),
			expectStatusCode: 200,
			expectRevision:   0,
		}, {
			name:             "putBadJson",
			key:              "/bad-json",
			contentType:      "application/json",
			content:          []byte(`abc`),
			expectStatusCode: 400,
		}, {
			name:             "putExisting",
			key:              existingObj.ID,
			contentType:      "application/json",
			content:          []byte(`{"d":"a"}`),
			expectStatusCode: 200,
			expectRevision:   1,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {

			body := bytes.NewReader(tc.content)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/apis/documents/v1/documents%s?bucketId=%s", s.server.GetAddress(), tc.key, bucket.ID), body)
			if !assert.NoError(t, err) {
				return
			}

			if len(tc.contentType) != 0 {
				req.Header.Set(headerContentType, tc.contentType)
			}

			if len(tc.tags) != 0 {
				req.Header.Set(headerTags, tc.tags)
			}

			res, err := http.DefaultClient.Do(req)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tc.expectStatusCode, res.StatusCode)
			responseBytes, err := ioutil.ReadAll(res.Body)
			if !assert.NoError(t, err) {
				return
			}

			// ignore next assertions if error is expected
			if tc.expectStatusCode > 399 {
				return
			}

			assert.Equal(t, getMD5Checksum(tc.content), res.Header.Get(headerETag))
			assert.Equal(t, strconv.Itoa(tc.expectRevision), res.Header.Get(headerObjectVersion))

			t.Log(string(responseBytes))

		})
	}

}
