package documents

import (
	"context"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestDeleteBucket() {
	ctx := context.Background()
	out, err := s.client.Buckets().Create(ctx, &Bucket{
		Name: "test-create-bucket",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.client.Buckets().Delete(ctx, out.ID)
	if !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.client.Buckets().Get(ctx, out.ID, GetBucketOptions{})
	if !assert.Error(s.T(), err) {
		return
	}
}
