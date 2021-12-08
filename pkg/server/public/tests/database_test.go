package tests

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestCanCreateDatabase() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var out types.Database
	if err := s.cli.CreateDatabase(ctx, &types.Database{Name: "My Database"}, &out); !assert.NoError(s.T(), err) {
		return
	}
}
