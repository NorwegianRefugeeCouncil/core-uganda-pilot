package sqlmanager

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"strconv"
	"time"
)

const (
	errInvalidSQLDataType = "invalid sql data type"
	errUnknownFieldF      = "unknown field '%s'"
)

func NewFormReader(form types.FormInterface, rows sqlReader) FormReader {
	return &formReader{
		form:      form,
		sqlReader: rows,
	}
}

type FormReader interface {
	GetRecords() (*types.RecordList, error)
}

type formReader struct {
	form      types.FormInterface
	sqlReader sqlReader
}

func (f formReader) GetRecords() (*types.RecordList, error) {
	return readRecords(f.form, f.sqlReader)
}

type sqlReader interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...interface{}) error
}

// readRecords will iterate through a series of SQL Rows and return a list of populated records
func readRecords(form types.FormInterface, rows sqlReader) (*types.RecordList, error) {

	result := &types.RecordList{
		Items: []*types.Record{},
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		if err := rows.Scan(&values); err != nil {
			return nil, err
		}
		record := &types.Record{
			FormID:     form.GetFormID(),
			DatabaseID: form.GetDatabaseID(),
		}
		if err := readInRecord(record, form, columns, values); err != nil {
			return nil, err
		}
		result.Items = append(result.Items, record)
	}

	return result, nil
}

// readInRecord will populate a record for a form from the given SQL columns and values
func readInRecord(record *types.Record, form types.FormInterface, columns []string, values []interface{}) error {

	formFields := form.GetFields()
	formFieldMap := map[string]*types.FieldDefinition{}
	for _, field := range formFields {
		formFieldMap[field.ID] = field
	}

	for i, column := range columns {
		columnValue := values[i]

		switch column {
		case keyIdColumn:
			recordId, err := mapStringValue(columnValue)
			if err != nil {
				return err
			}
			record.ID = recordId
			continue
		case keyOwnerIdColumn:
			ownerId, err := mapStringPointerValue(columnValue)
			if err != nil {
				return err
			}
			record.OwnerID = ownerId
			continue
		case keyCreatedAtColumn:
			_, err := mapTimeValue(columnValue)
			if err != nil {
				return err
			}
			// todo: record.CreatedAt
			continue
		default:
			formField, ok := formFieldMap[column]
			if !ok {
				return fmt.Errorf(errUnknownFieldF, column)
			}
			if err := readInRecordField(record, formField, columnValue); err != nil {
				return err
			}
		}

	}

	return nil

}

// readInRecordField will populate a record types.FieldValue for the given field and value
func readInRecordField(
	record *types.Record,
	field *types.FieldDefinition,
	value interface{},
) error {

	fieldKind, err := field.FieldType.GetFieldKind()
	if err != nil {
		return err
	}

	switch fieldKind {
	case types.FieldKindMonth, types.FieldKindDate:
		readInDateField(record, field, value, fieldKind)
	case types.FieldKindQuantity:
		readInQuantityField(record, field, value)
	case types.FieldKindSingleSelect:
		readInSingleSelectField(record, field, value)
	case types.FieldKindReference:
		readInReferenceField(record, field, value)
	case types.FieldKindText, types.FieldKindMultilineText:
		readInTextField(record, field, value)
	}

	return nil
}

// readInReferenceField will populate a record types.FieldValue for a types.FieldTypeDate or types.FieldTypeMonth field from
// an  SQL value
func readInDateField(record *types.Record, field *types.FieldDefinition, value interface{}, fieldKind types.FieldKind) {
	var timeFormat string
	switch fieldKind {
	case types.FieldKindMonth:
		timeFormat = monthFieldFormat
	case types.FieldKindDate:
		timeFormat = dateFieldFormat
	}
	var monthStr *string
	switch t := value.(type) {
	case time.Time:
		monthStr = pointers.String(t.Format(timeFormat))
	case *time.Time:
		if t != nil {
			monthStr = pointers.String(t.Format(timeFormat))
		}
	}
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   monthStr,
	})
}

// readInReferenceField will populate a record types.FieldValue for a types.FieldTypeReference field from
// an  SQL value
func readInReferenceField(record *types.Record, field *types.FieldDefinition, value interface{}) {
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   getStringColumnStringPtr(value),
	})
}

// readInSingleSelectField will populate a record types.FieldValue for a types.FieldTypeSingleSelect field from
// an  SQL value
func readInSingleSelectField(record *types.Record, field *types.FieldDefinition, value interface{}) {
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   getStringColumnStringPtr(value),
	})
}

// readInReferenceField will populate a record types.FieldValue for a types.FieldTypeText or
// types.FieldTypeMultilineText field from an  SQL value
func readInTextField(record *types.Record, field *types.FieldDefinition, value interface{}) {
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   getStringColumnStringPtr(value),
	})
}

// readInQuantityField will populate a record types.FieldValue for a types.FieldTypeQuantity from an  SQL value
func readInQuantityField(record *types.Record, field *types.FieldDefinition, value interface{}) {
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   getIntColumnStringPtr(value),
	})
}

// getStringColumnStringPtr will coerce a string or *string into a *string
func getStringColumnStringPtr(value interface{}) *string {
	var textStr *string
	switch t := value.(type) {
	case string:
		textStr = pointers.String(t)
	case *string:
		textStr = t
	}
	return textStr
}

// getIntColumnStringPtr will convert an int or *int into a *string
func getIntColumnStringPtr(value interface{}) *string {
	var intStr *string
	switch t := value.(type) {
	case int:
		intStr = pointers.String(strconv.Itoa(t))
	case *int:
		if t != nil {
			intStr = pointers.String(strconv.Itoa(*t))
		}
	}
	return intStr
}

// mapStringValue will return a string value from an interface
// useful when we know for sure that a columns contains a non-nullable string field
func mapStringValue(value interface{}) (string, error) {
	recordId, ok := value.(string)
	if !ok {
		return "", errors.New(errInvalidSQLDataType)
	}
	return recordId, nil
}

// mapStringPointerValue will return a string pointer value from an interface
func mapStringPointerValue(value interface{}) (*string, error) {
	var val *string
	switch v := value.(type) {
	case *string:
		val = v
	case string:
		val = &v
	default:
		return nil, errors.New(errInvalidSQLDataType)
	}
	return val, nil
}

// mapTimeValue is useful when we know for sure that a columns contains a non-nullable time field
func mapTimeValue(value interface{}) (time.Time, error) {
	timeValue, ok := value.(time.Time)
	if !ok {
		return time.Time{}, errors.New(errInvalidSQLDataType)
	}
	return timeValue, nil
}
