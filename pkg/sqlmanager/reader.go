package sqlmanager

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"gorm.io/gorm"
)

const (
	errInvalidSQLDataType = "invalid sql data type"
)

func NewFormReader(db *gorm.DB) Reader {
	return &reader{
		db: db,
	}
}

type Reader interface {
	GetRecords(ctx context.Context, form types.FormInterface) (*types.RecordList, error)
	GetRecord(ctx context.Context, form types.FormInterface, recordRef types.RecordRef) (*types.Record, error)
}

type reader struct {
	db *gorm.DB
}

func (f reader) GetRecords(ctx context.Context, form types.FormInterface) (*types.RecordList, error) {
	return queryRecords(ctx, f.db, form, "")
}

func (f reader) GetRecord(ctx context.Context, form types.FormInterface, recordRef types.RecordRef) (*types.Record, error) {
	recs, err := queryRecords(ctx, f.db, form, "where id = ?", recordRef.ID)
	if err != nil {
		return nil, err
	}
	if len(recs.Items) == 0 {
		return nil, meta.NewNotFound(types.RecordGR, recordRef.ID)
	}
	if len(recs.Items) > 1 {
		return nil, meta.NewInternalServerError(fmt.Errorf("unexpected number of records"))
	}
	return recs.Items[0], nil
}

type sqlReader interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...interface{}) error
}

// queryRecords will iterate through a series of SQL Rows and return a list of populated records
func queryRecords(ctx context.Context, db *gorm.DB, form types.FormInterface, sqlQuery string, args ...interface{}) (*types.RecordList, error) {
	qry := db.Raw(fmt.Sprintf("select * from %s.%s %s",
		pq.QuoteIdentifier(form.GetDatabaseID()),
		pq.QuoteIdentifier(form.GetFormID()),
		sqlQuery,
	), args...)

	rows, err := qry.Rows()
	if err != nil {
		return nil, err
	}

	subRecords, err := querySubRecords(ctx, db, form, args...)
	if err != nil {
		return nil, err
	}

	records, err := readRecords(form, rows, subRecords)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// querySubRecords will get all sub records for each subform field in a form and return a nested map of recordId -> fieldId -> recordList
func querySubRecords(ctx context.Context, db *gorm.DB, form types.FormInterface, args ...interface{}) (map[string]map[string]*types.RecordList, error) {
	subFormRecordMap := map[string]map[string]*types.RecordList{}

	// Iterate over form fields
	for _, field := range form.GetFields() {
		fieldKind, err := field.FieldType.GetFieldKind()
		if err != nil {
			return nil, err
		}
		// If field is a subform, query sub records
		if fieldKind == types.FieldKindSubForm {
			subForm, err := form.FindSubForm(field.ID)
			if err != nil {
				return nil, err
			}

			query := ""
			if len(args) > 0 && args[0] != nil {
				query = fmt.Sprintf("where %s = ?", keyOwnerIdColumn)
			}

			subFormRecordList, err := queryRecords(ctx, db, subForm, query, args...)
			if err != nil {
				return nil, err
			}
			// Create or append to map entries
			for _, subRecord := range subFormRecordList.Items {
				if _, ok := subFormRecordMap[*subRecord.OwnerID]; !ok {
					subFormRecordMap[*subRecord.OwnerID] = map[string]*types.RecordList{}
				}
				if _, ok := subFormRecordMap[*subRecord.OwnerID][field.ID]; !ok {
					subFormRecordMap[*subRecord.OwnerID][field.ID] = &types.RecordList{
						Items: []*types.Record{},
					}
				}
				subFormRecordMap[*subRecord.OwnerID][field.ID].Items = append(subFormRecordMap[*subRecord.OwnerID][field.ID].Items, subRecord)
			}
		}
	}

	return subFormRecordMap, nil
}

// readRecords will iterate through a series of SQL Rows and return a list of populated records
func readRecords(form types.FormInterface, rows sqlReader, subRecords map[string]map[string]*types.RecordList) (*types.RecordList, error) {

	result := &types.RecordList{
		Items: []*types.Record{},
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(values))
		for i := range values {
			valuePointers[i] = &values[i]
		}
		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}
		record := &types.Record{
			FormID:     form.GetFormID(),
			DatabaseID: form.GetDatabaseID(),
		}
		if err := readInRecord(record, form, subRecords, columns, values); err != nil {
			return nil, err
		}
		result.Items = append(result.Items, record)
	}

	return result, nil
}

// readInRecord will populate a record for a form from the given SQL columns and values
func readInRecord(record *types.Record, form types.FormInterface, subRecords map[string]map[string]*types.RecordList, columns []string, values []interface{}) error {

	formFields := form.GetFields()
	fieldValueMap := map[string]interface{}{}

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
			fieldValueMap[column] = columnValue
		}
	}

	for _, field := range formFields {
		fieldValue := fieldValueMap[field.ID]
		if err := readInRecordField(record, subRecords[record.ID][field.ID], field, fieldValue); err != nil {
			return err
		}
	}

	return nil

}

// readInRecordField will populate a record types.FieldValue for the given field and value
func readInRecordField(
	record *types.Record,
	subRecords *types.RecordList,
	field *types.FieldDefinition,
	value interface{},
) error {

	fieldKind, err := field.FieldType.GetFieldKind()
	if err != nil {
		return err
	}

	switch fieldKind {
	case types.FieldKindMonth, types.FieldKindDate, types.FieldKindWeek:
		return readInDateField(record, field, value, fieldKind)
	case types.FieldKindQuantity:
		return readInQuantityField(record, field, value)
	case types.FieldKindSingleSelect:
		return readInSingleSelectField(record, field, value)
	case types.FieldKindMultiSelect:
		return readInMultiSelectField(record, field, value)
	case types.FieldKindReference:
		return readInReferenceField(record, field, value)
	case types.FieldKindText, types.FieldKindMultilineText:
		return readInTextField(record, field, value)
	case types.FieldKindCheckbox:
		return readInBooleanField(record, field, value)
	case types.FieldKindSubForm:
		return readInSubFormField(record, subRecords, field)
	}

	return nil
}

// readInReferenceField will populate a record types.FieldValue for a types.FieldTypeDate or types.FieldTypeMonth field from
// an  SQL value
func readInDateField(record *types.Record, field *types.FieldDefinition, value interface{}, fieldKind types.FieldKind) error {
	var valueStr *string
	switch fieldKind {
	case types.FieldKindMonth:
		switch t := value.(type) {
		case time.Time:
			valueStr = pointers.String(t.Format(monthFieldFormat))
		case *time.Time:
			if t != nil {
				valueStr = pointers.String(t.Format(monthFieldFormat))
			}
		}
	case types.FieldKindWeek:
		switch t := value.(type) {
		case time.Time:
			year, week := t.ISOWeek()
			valueStr = pointers.String(fmt.Sprintf("%d-W%d", year, week))
		case *time.Time:
			if t != nil {
				year, week := t.ISOWeek()
				valueStr = pointers.String(fmt.Sprintf("%d-W%d", year, week))
			}
		}

	case types.FieldKindDate:
		switch t := value.(type) {
		case time.Time:
			valueStr = pointers.String(t.Format(dateFieldFormat))
		case *time.Time:
			if t != nil {
				valueStr = pointers.String(t.Format(dateFieldFormat))
			}
		}
	}

	if valueStr == nil {
		record.Values = append(record.Values, types.NewFieldNullValue(field.ID))
	} else {
		record.Values = append(record.Values, types.NewFieldStringValue(field.ID, *valueStr))
	}

	return nil
}

// readInReferenceField will populate a record types.FieldValue for a types.FieldTypeReference field from
// an  SQL value
func readInReferenceField(record *types.Record, field *types.FieldDefinition, value interface{}) error {
	fieldValue, err := getStringValue(value)
	if err != nil {
		return err
	}
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   fieldValue,
	})
	return nil
}

// readInSingleSelectField will populate a record types.FieldValue for a types.FieldTypeSingleSelect field from
// an  SQL value
func readInSingleSelectField(record *types.Record, field *types.FieldDefinition, value interface{}) error {
	fieldValue, err := getStringValue(value)
	if err != nil {
		return err
	}
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   fieldValue,
	})
	return nil
}

// readInMultiSelectField will populate a record types.FieldValue for a types.FieldTypeMultiSelect field from
// an  SQL value
func readInMultiSelectField(record *types.Record, field *types.FieldDefinition, value interface{}) error {
	fieldValue, err := getStringListValue(value)
	if err != nil {
		return err
	}
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   fieldValue,
	})
	return nil
}

// readInReferenceField will populate a record types.FieldValue for a types.FieldTypeText or
// types.FieldTypeMultilineText field from an  SQL value
func readInTextField(record *types.Record, field *types.FieldDefinition, value interface{}) error {
	fieldValue, err := getStringValue(value)
	if err != nil {
		return err
	}
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   fieldValue,
	})
	return nil
}

// readInQuantityField will populate a record types.FieldValue for a types.FieldTypeQuantity from an  SQL value
func readInQuantityField(record *types.Record, field *types.FieldDefinition, value interface{}) error {
	fieldValue, err := getIntValue(value)
	if err != nil {
		return err
	}
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   fieldValue,
	})
	return nil
}

// readInBooleanField will populate a record types.FieldValue for a types.FieldTypeCheckbox from an SQL value
func readInBooleanField(record *types.Record, field *types.FieldDefinition, value interface{}) error {
	fieldValue, err := getBooleanValue(value)
	if err != nil {
		return err
	}
	record.Values = append(record.Values, types.FieldValue{
		FieldID: field.ID,
		Value:   fieldValue,
	})
	return nil
}

func readInSubFormField(record *types.Record, subRecords *types.RecordList, field *types.FieldDefinition) error {
	// When we create a record we don't create the sub records and the resulting GET will fail while getting the subrecords
	// We should create both together, but for now we just create the record and ignore the subrecords
	if subRecords == nil {
		return nil
	}

	value := make([]types.FieldValues, 0)

	for _, subRecord := range subRecords.Items {
		if *subRecord.OwnerID != record.ID {
			continue
		}

		subValue := make([]types.FieldValue, 0)
		for _, subFieldValue := range subRecord.Values {
			subValue = append(subValue, subFieldValue)
		}
		value = append(value, subValue)
	}

	record.Values = append(record.Values, types.NewFieldSubFormValue(
		field.ID,
		value,
	))

	return nil
}

// getStringValue will coerce a string or *string into a *string
func getStringValue(value interface{}) (types.StringOrArray, error) {
	switch t := value.(type) {
	case string:
		return types.NewStringValue(t), nil
	case *string:
		if t == nil {
			return types.NewNullValue(), nil
		} else {
			return types.NewStringValue(*t), nil
		}
	case nil:
		return types.NewNullValue(), nil
	default:
		return types.StringOrArray{}, fmt.Errorf("cannot convert type %T to types.StringOrArray", value)
	}
}

// getStringListValue will coerce an interface{} into a []string
func getStringListValue(value interface{}) (types.StringOrArray, error) {
	switch t := value.(type) {
	case string:
		return types.NewArrayValue(parseArrayStr(t)), nil
	case *string:
		if t == nil {
			return types.NewNullValue(), nil
		} else {
			return types.NewArrayValue(parseArrayStr(*t)), nil
		}
	case nil:
		return types.NewNullValue(), nil
	default:
		return types.StringOrArray{}, fmt.Errorf("cannot convert type %T to types.StringOrArray", value)
	}
}

func parseArrayStr(str string) []string {
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")
	return strings.Split(str, ",")
}

// getIntValue will convert an int or *int into a *string
func getIntValue(value interface{}) (types.StringOrArray, error) {
	switch t := value.(type) {
	case int:
		return types.NewStringValue(strconv.Itoa(t)), nil
	case *int:
		if t == nil {
			return types.NewNullValue(), nil
		} else {
			return types.NewStringValue(strconv.Itoa(*t)), nil
		}
	case int64:
		return types.NewStringValue(strconv.FormatInt(t, 10)), nil
	case *int64:
		if t == nil {
			return types.NewNullValue(), nil
		} else {
			return types.NewStringValue(strconv.FormatInt(*t, 10)), nil
		}
	case nil:
		return types.NewNullValue(), nil
	default:
		return types.StringOrArray{}, fmt.Errorf("cannot convert type %T to types.StringOrArray", value)
	}
}

func getBooleanValue(value interface{}) (types.StringOrArray, error) {
	switch t := value.(type) {
	case bool:
		return types.NewStringValue(strconv.FormatBool(t)), nil
	case *bool:
		if t == nil {
			return types.NewNullValue(), nil
		} else {
			return types.NewStringValue(strconv.FormatBool(*t)), nil
		}
	case nil:
		return types.NewNullValue(), nil
	default:
		return types.StringOrArray{}, fmt.Errorf("cannot convert type %T to types.StringOrArray", value)
	}
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
