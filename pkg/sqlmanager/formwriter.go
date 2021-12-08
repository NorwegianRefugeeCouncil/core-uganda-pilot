package sqlmanager

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
	"github.com/nrc-no/core/pkg/utils/dates"
	"strconv"
	"strings"
	"time"
)

func NewFormWriter(form types.FormInterface) FormWriter {
	return &formWriter{
		form: form,
	}
}

type FormWriter interface {
	WriteRecords(records *types.RecordList) ([]schema.DDL, error)
}

type formWriter struct {
	form types.FormInterface
}

func (f *formWriter) WriteRecords(records *types.RecordList) ([]schema.DDL, error) {
	var result []schema.DDL
	for _, item := range records.Items {
		recordDDL, err := f.writeRecord(item)
		if err != nil {
			return nil, err
		}
		result = append(result, recordDDL)
	}
	return result, nil
}

func (f *formWriter) writeRecord(record *types.Record) (schema.DDL, error) {

	ddl := fmt.Sprintf(`insert into %s.%s (`,
		pq.QuoteIdentifier(f.form.GetDatabaseID()),
		pq.QuoteIdentifier(f.form.GetFormID()),
	)

	paramCount := 1
	columnNames := []string{keyIdColumn}
	args := []interface{}{record.ID}

	if f.form.IsSubForm() {
		if record.OwnerID == nil {
			return schema.DDL{}, errors.New("no owner id on record")
		}
		paramCount++
		args = append(args, *record.OwnerID)
		columnNames = append(columnNames, keyOwnerIdColumn)
	}

	for _, field := range f.form.GetFields() {
		fieldKind, err := field.FieldType.GetFieldKind()
		if err != nil {
			return schema.DDL{}, err
		}
		switch fieldKind {
		case types.FieldKindText, types.FieldKindMultilineText, types.FieldKindReference, types.FieldKindSingleSelect:
			if fieldValue, ok := record.Values.Find(field.ID); ok {
				paramCount++
				columnNames = append(columnNames, field.ID)
				if fieldValue.Value != nil {
					args = append(args, *fieldValue.Value)
				} else {
					args = append(args, nil)
				}
			}
		case types.FieldKindDate, types.FieldKindMonth:
			if fieldValue, ok := record.Values.Find(field.ID); ok {
				paramCount++
				columnNames = append(columnNames, field.ID)
				if fieldValue.Value != nil {
					var dateFormat string
					switch fieldKind {
					case types.FieldKindDate:
						dateFormat = dateFieldFormat
					case types.FieldKindMonth:
						dateFormat = monthFieldFormat
					}
					sqlDate, err := time.Parse(dateFormat, *fieldValue.Value)
					if err != nil {
						return schema.DDL{}, err
					}
					args = append(args, sqlDate)
				} else {
					args = append(args, nil)
				}
			}
		case types.FieldKindWeek:
			if fieldValue, ok := record.Values.Find(field.ID); ok {
				paramCount++
				columnNames = append(columnNames, field.ID)
				if fieldValue.Value != nil {
					t, err := dates.ParseIsoWeekTime(*fieldValue.Value)
					if err != nil {
						return schema.DDL{}, err
					}
					args = append(args, t)
				} else {
					args = append(args, nil)
				}
			}
		case types.FieldKindQuantity:
			if fieldValue, ok := record.Values.Find(field.ID); ok {
				paramCount++
				columnNames = append(columnNames, field.ID)
				if fieldValue.Value != nil {
					intValue, err := strconv.Atoi(*fieldValue.Value)
					if err != nil {
						return schema.DDL{}, err
					}
					args = append(args, intValue)
				} else {
					args = append(args, nil)
				}
			}
		}
	}

	for i := 0; i < len(columnNames); i++ {
		columnNames[i] = pq.QuoteIdentifier(columnNames[i])
	}

	ddl = ddl + strings.Join(columnNames, ",")
	ddl = ddl + ") values ("

	for i := 1; i <= paramCount; i++ {
		ddl = ddl + "$" + strconv.Itoa(i)
		if i < paramCount {
			ddl = ddl + ","
		}
	}

	ddl = ddl + ");"

	return schema.NewDDL(ddl, args...), nil

}
