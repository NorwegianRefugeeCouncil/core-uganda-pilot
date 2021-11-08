package test

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestRecordsApi() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli := s.PublicClient(ctx)

	db := &types.Database{}
	if err := cli.CreateDatabase(ctx, &types.Database{Name: "test-db"}, db); !assert.NoError(s.T(), err) {
		return
	}
	defer cli.DeleteDatabase(ctx, db.ID)

	form := &types.FormDefinition{}
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
	}, form); !assert.NoError(s.T(), err) {
		return
	}

	defer cli.DeleteForm(ctx, form.ID)

	created := &types.Record{}
	if err := cli.CreateRecord(ctx, &types.Record{
		DatabaseID: db.ID,
		FormID:     form.ID,
		Values: map[string]interface{}{
			form.Fields[0].ID: "somevalue",
		},
	}, created); !assert.NoError(s.T(), err) {
		return
	}

	list := &types.RecordList{}
	if err := cli.ListRecords(ctx, types.RecordListOptions{FormID: form.ID, DatabaseID: db.ID}, list); !assert.NoError(s.T(), err) {
		return
	}

	assert.Len(s.T(), list.Items, 1)

}
