package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
	tu "github.com/nrc-no/core/pkg/testutils"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_formWriter_WriteRecords(t *testing.T) {

	const (
		formId     = "formId"
		databaseId = "databaseId"
	)

	formWithFields := func(options ...tu.FormOption) types.FormInterface {
		opts := []tu.FormOption{
			tu.FormID(formId),
			tu.FormDatabaseID(databaseId),
		}
		f := tu.AForm(append(opts, options...)...)
		return f
	}

	aRecord := func(options ...tu.RecordOption) *types.Record {
		opts := []tu.RecordOption{
			tu.RecordID("recordId"),
			tu.RecordFormID(formId),
			tu.RecordDatabaseID(databaseId),
		}
		opts = append(opts, options...)
		return tu.RecordOptions(opts...)(&types.Record{})
	}

	aSingleRecord := func(options ...tu.RecordOption) *types.RecordList {
		return &types.RecordList{
			Items: []*types.Record{
				aRecord(options...),
			},
		}
	}

	mustParseTime := func(layout string, str string) time.Time {
		tm, err := time.Parse(layout, str)
		if err != nil {
			panic(err)
		}
		return tm
	}

	tests := []struct {
		name    string
		form    types.FormInterface
		records *types.RecordList
		want    []schema.DDL
		wantErr bool
	}{
		{
			name:    "record with text value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeText()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("myValue"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", "myValue"},
				},
			},
		}, {
			name:    "record with null text value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeText()))),
			records: aSingleRecord(tu.RecordValue("fieldId", nil)),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with multiline text value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeMultilineText()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("myValue"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", "myValue"},
				},
			},
		}, {
			name:    "record with nil multiline text value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeMultilineText()))),
			records: aSingleRecord(tu.RecordValue("fieldId", nil)),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with month value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeMonth()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("2020-01"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", mustParseTime(monthFieldFormat, "2020-01")},
				},
			},
		}, {
			name:    "record with nil month value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeMonth()))),
			records: aSingleRecord(tu.RecordValue("fieldId", nil)),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with bad month value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeMonth()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("abc"))),
			wantErr: true,
		}, {
			name:    "record with date value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeDate()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("2020-01-01"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", mustParseTime(dateFieldFormat, "2020-01-01")},
				},
			},
		}, {
			name:    "record with nil date value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeDate()))),
			records: aSingleRecord(tu.RecordValue("fieldId", nil)),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with bad date value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeDate()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("abc"))),
			wantErr: true,
		}, {
			name:    "record with reference value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeReference()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("refId"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", "refId"},
				},
			},
		}, {
			name:    "record with nil reference value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeReference()))),
			records: aSingleRecord(tu.RecordValue("fieldId", nil)),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with quantity value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeQuantity()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("3"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", 3},
				},
			},
		}, {
			name:    "record with nil quantity value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeQuantity()))),
			records: aSingleRecord(tu.RecordValue("fieldId", nil)),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,fieldId) values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with bad quantity value",
			form:    formWithFields(tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeQuantity()))),
			records: aSingleRecord(tu.RecordValue("fieldId", pointers.String("abc"))),
			wantErr: true,
		}, {
			name: "record with owner",
			form: formWithFields(
				tu.FormHasOwner(true),
				tu.FormField(tu.AField(tu.FieldID("fieldId"), tu.FieldTypeQuantity())),
			),
			records: aSingleRecord(
				tu.RecordOwnerID(pointers.String("ownerId")),
				tu.RecordValue("fieldId", pointers.String("123")),
			),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" (id,owner_id,fieldId) values ($1,$2,$3);`,
					Args:  []interface{}{"recordId", "ownerId", 123},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &formWriter{
				form: tt.form,
			}
			got, err := f.WriteRecords(tt.records)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			t.Log(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
