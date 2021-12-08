package sqlmanager

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/testutils"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func Test_readRecords(t *testing.T) {

	const (
		formId       = "formId"
		databaseId   = "databaseId"
		keyMyFieldID = "myField"
		ownerId      = "ownerId"
		recordId     = "recordId"
	)

	formWithField := func(kind types.FieldKind, ) types.FormInterface {
		fldOpt := []testutils.FieldOption{
			testutils.FieldID(keyMyFieldID),
			testutils.FieldTypeKind(kind),
		}
		fldOpt = append(fldOpt)
		formOpts := []testutils.FormOption{
			testutils.FormID(formId),
			testutils.FormDatabaseID(databaseId),
			testutils.FormFields(
				testutils.AField(fldOpt...),
			),
		}
		formOpts = append(formOpts)
		return testutils.AForm(formOpts...)
	}

	recordWithValue := func(value *string, options ...testutils.RecordOption) types.Record {
		r := &types.Record{
			ID:         recordId,
			DatabaseID: databaseId,
			FormID:     formId,
			Values: []types.FieldValue{
				{
					FieldID: keyMyFieldID,
					Value:   value,
				},
			},
		}
		r = testutils.RecordOptions(options...)(r)
		return *r
	}

	columns := func(cols ...string) []string {
		return cols
	}

	singleRow := func(values ...interface{}) [][]interface{} {
		return [][]interface{}{
			values,
		}
	}

	singleRecord := func(record types.Record) *types.RecordList {
		return &types.RecordList{
			Items: []*types.Record{
				&record,
			},
		}
	}

	tests := []struct {
		name          string
		form          types.FormInterface
		columns       []string
		values        [][]interface{}
		want          *types.RecordList
		wantErr       bool
		scanThrows    bool
		columnsThrows bool
	}{
		{
			name:    "form with text field",
			form:    formWithField(types.FieldKindText),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, "abc"),
			want:    singleRecord(recordWithValue(pointers.String("abc"))),
		}, {
			name:    "form with optional text field",
			form:    formWithField(types.FieldKindText),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.String("abc")),
			want:    singleRecord(recordWithValue(pointers.String("abc"))),
		}, {
			name:    "form with multiline text field",
			form:    formWithField(types.FieldKindMultilineText),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, "abc"),
			want:    singleRecord(recordWithValue(pointers.String("abc"))),
		}, {
			name:    "form with optional multiline text field",
			form:    formWithField(types.FieldKindMultilineText),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.String("abc")),
			want:    singleRecord(recordWithValue(pointers.String("abc"))),
		}, {
			name:    "form with quantity field",
			form:    formWithField(types.FieldKindQuantity),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, 123),
			want:    singleRecord(recordWithValue(pointers.String("123"))),
		}, {
			name:    "form with nullable quantity field",
			form:    formWithField(types.FieldKindQuantity),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.Int(123)),
			want:    singleRecord(recordWithValue(pointers.String("123"))),
		}, {
			name:    "form with date field",
			form:    formWithField(types.FieldKindDate),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, time.Date(2020, 6, 26, 10, 20, 30, 0, time.UTC)),
			want:    singleRecord(recordWithValue(pointers.String("2020-06-26"))),
		}, {
			name:    "form with nullable date field",
			form:    formWithField(types.FieldKindDate),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.Time(time.Date(2020, 6, 26, 10, 20, 30, 0, time.UTC))),
			want:    singleRecord(recordWithValue(pointers.String("2020-06-26"))),
		}, {
			name:    "form with month field",
			form:    formWithField(types.FieldKindMonth),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, time.Date(2020, 6, 26, 10, 20, 30, 0, time.UTC)),
			want:    singleRecord(recordWithValue(pointers.String("2020-06"))),
		}, {
			name:    "form with nullable month field",
			form:    formWithField(types.FieldKindMonth),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.Time(time.Date(2020, 6, 26, 10, 20, 30, 0, time.UTC))),
			want:    singleRecord(recordWithValue(pointers.String("2020-06"))),
		}, {
			name:    "form with week field",
			form:    formWithField(types.FieldKindWeek),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, time.Date(2020, 6, 26, 10, 20, 30, 0, time.UTC)),
			want:    singleRecord(recordWithValue(pointers.String("2020-W26"))),
		}, {
			name:    "form with nullable week field",
			form:    formWithField(types.FieldKindWeek),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.Time(time.Date(2020, 6, 26, 10, 20, 30, 0, time.UTC))),
			want:    singleRecord(recordWithValue(pointers.String("2020-W26"))),
		}, {
			name:    "form with reference field",
			form:    formWithField(types.FieldKindReference),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, "otherRefId"),
			want:    singleRecord(recordWithValue(pointers.String("otherRefId"))),
		}, {
			name:    "form with nullable reference field",
			form:    formWithField(types.FieldKindReference),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.String("otherRefId")),
			want:    singleRecord(recordWithValue(pointers.String("otherRefId"))),
		}, {
			name:    "form with owner id",
			form:    formWithField(types.FieldKindText),
			columns: columns(keyIdColumn, keyOwnerIdColumn, keyMyFieldID),
			values:  singleRow(recordId, ownerId, "abc"),
			want:    singleRecord(recordWithValue(pointers.String("abc"), testutils.RecordOwnerID(pointers.String(ownerId)))),
		}, {
			name:    "form with nullable owner id",
			form:    formWithField(types.FieldKindText),
			columns: columns(keyIdColumn, keyOwnerIdColumn, keyMyFieldID),
			values:  singleRow(recordId, pointers.String(ownerId), "abc"),
			want:    singleRecord(recordWithValue(pointers.String("abc"), testutils.RecordOwnerID(pointers.String(ownerId)))),
		}, {
			name:    "record with bad id",
			form:    formWithField(types.FieldKindText),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(123, "abc"),
			wantErr: true,
		}, {
			name:    "record with bad owner id",
			form:    formWithField(types.FieldKindText),
			columns: columns(keyIdColumn, keyOwnerIdColumn, keyMyFieldID),
			values:  singleRow(recordId, 123, "abc"),
			wantErr: true,
		}, {
			name:    "record with bad created_at",
			form:    formWithField(types.FieldKindText),
			columns: columns(keyIdColumn, keyCreatedAtColumn, keyMyFieldID),
			values:  singleRow(recordId, 123, "abc"),
			wantErr: true,
		}, {
			name:    "record with unknown field type",
			form:    formWithField(types.FieldKindUnknown),
			columns: columns(keyIdColumn, keyMyFieldID),
			values:  singleRow(recordId, "abc"),
			wantErr: true,
		}, {
			name:          "columns throws",
			form:          formWithField(types.FieldKindUnknown),
			columns:       columns(),
			values:        singleRow(),
			columnsThrows: true,
			wantErr:       true,
		}, {
			name:       "scan throws",
			form:       formWithField(types.FieldKindUnknown),
			columns:    columns(),
			values:     singleRow(),
			scanThrows: true,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSql := newMockSQLReader(tt.columns, tt.values)
			mockSql.columnsThrows = tt.columnsThrows
			mockSql.scanThrows = tt.scanThrows
			reader := NewFormReader(tt.form, mockSql)
			got, err := reader.GetRecords()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

type mockSqlReader struct {
	i             int
	columns       []string
	values        [][]interface{}
	columnsThrows bool
	scanThrows    bool
}

func newMockSQLReader(columns []string, values [][]interface{}) *mockSqlReader {
	return &mockSqlReader{
		i:       -1,
		columns: columns,
		values:  values,
	}
}

func (m *mockSqlReader) Columns() ([]string, error) {
	if m.columnsThrows {
		return nil, fmt.Errorf("mock error")
	}
	return m.columns, nil
}

func (m *mockSqlReader) Next() bool {
	m.i = m.i + 1
	return m.i < len(m.values)
}

func (m *mockSqlReader) Scan(intf ...interface{}) error {
	if m.scanThrows {
		return fmt.Errorf("mock error")
	}

	for i, a := range intf {
		firstValue := reflect.ValueOf(a)
		firstValue.Elem().Set(reflect.ValueOf(m.values[m.i][i]))
	}

	return nil
}

var _ sqlReader = &mockSqlReader{}
