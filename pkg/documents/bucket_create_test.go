package documents

import (
	"context"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestCreateBucket() {
	ctx := context.Background()
	out, err := s.client.Buckets().Create(ctx, &Bucket{
		Name: "test-create-bucket",
	})

	if !assert.NoError(s.T(), err) {
		return
	}
	assert.NotEmpty(s.T(), out.ID)
	assert.Equal(s.T(), "test-create-bucket", out.Name)
}
