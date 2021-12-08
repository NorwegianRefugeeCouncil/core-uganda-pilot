package tests

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	tu "github.com/nrc-no/core/pkg/testutils"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestRecordCreateReadList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var db types.Database
	if err := s.createDatabase(ctx, &db); !assert.NoError(s.T(), err) {
		return
	}

	var form types.FormDefinition
	const (
		textFieldName          = "Text Field"
		monthFieldName         = "Month Field"
		dateFieldName          = "Date Field"
		weekFieldName          = "Week Field"
		multilineTextFieldName = "Multiline Text Field"
		quantityFieldName      = "Quantity Field"
	)

	if err := s.cli.CreateForm(ctx, &types.FormDefinition{
		Name:       "My Form",
		DatabaseID: db.ID,
		Fields: []*types.FieldDefinition{
			tu.ATextField(tu.FieldName(textFieldName)),
			tu.AMonthField(tu.FieldName(monthFieldName)),
			tu.ADateField(tu.FieldName(dateFieldName)),
			tu.AWeekField(tu.FieldName(weekFieldName)),
			tu.AMultilineTextField(tu.FieldName(multilineTextFieldName)),
			tu.AQuantityField(tu.FieldName(quantityFieldName)),
		},
	}, &form); !assert.NoError(s.T(), err) {
		return
	}

	var values types.FieldValues
	values, _ = values.SetValueForFieldName(&form, textFieldName, pointers.String("text value"))
	values, _ = values.SetValueForFieldName(&form, monthFieldName, pointers.String("2010-01"))
	values, _ = values.SetValueForFieldName(&form, dateFieldName, pointers.String("2010-12-31"))
	values, _ = values.SetValueForFieldName(&form, weekFieldName, pointers.String("2020-W10"))
	values, _ = values.SetValueForFieldName(&form, multilineTextFieldName, pointers.String("text\nvalue"))
	values, _ = values.SetValueForFieldName(&form, quantityFieldName, pointers.String("10"))

	var record types.Record
	if err := s.cli.CreateRecord(ctx, &types.Record{
		Values:     values,
		FormID:     form.ID,
		DatabaseID: form.DatabaseID,
	}, &record); !assert.NoError(s.T(), err) {
		return
	}

	var got types.Record
	if err := s.cli.GetRecord(ctx, types.RecordRef{
		ID:         record.ID,
		DatabaseID: record.DatabaseID,
		FormID:     record.FormID,
	}, &got); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), record, got)

	var list types.RecordList
	if err := s.cli.ListRecords(ctx, types.RecordListOptions{
		DatabaseID: form.DatabaseID,
		FormID:     form.ID,
	}, &list); !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, &got)

}
