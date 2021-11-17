package test

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestFormsApi() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli := s.PublicClient(ctx)

	db := &types.Database{}
	if err := cli.CreateDatabase(ctx, &types.Database{Name: "test-db"}, db); !assert.NoError(s.T(), err) {
		return
	}
	defer cli.DeleteDatabase(ctx, db.ID)

	created := &types.FormDefinition{}
	if err := cli.CreateForm(ctx, &types.FormDefinition{
		DatabaseID: db.ID,
		Name:       "my-form",
		Fields: []*types.FieldDefinition{
			{
				Name:        "field-1",
				Description: "some field description",
				Key:         true,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
		},
	}, created); !assert.NoError(s.T(), err) {
		return
	}

	if err := cli.DeleteForm(ctx, created.ID); !assert.NoError(s.T(), err) {
		return
	}

}
