package sqlmanager

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
	"github.com/nrc-no/core/pkg/utils/dates"
	"strconv"
	"time"
)

type sqlState struct {
	Tables schema.SQLTables
}

type Writer interface {
	PutForms(formDefinitions *types.FormDefinitionList) (Writer, error)
	PutRecords(form types.FormInterface, records *types.RecordList) (Writer, error)
	GetStatements() []schema.DDL
}

type writer struct {
	State      sqlState
	Statements []schema.DDL
}

func New() Writer {
	return writer{}
}

func (s writer) GetStatements() []schema.DDL {
	return s.Statements
}

func (s writer) PutForms(formDefinitions *types.FormDefinitionList) (Writer, error) {
	actions := sqlActions{}
	for _, item := range formDefinitions.Items {
		formActions, err := getSQLActionsForForm(item)
		if err != nil {
			return nil, err
		}
		actions = append(actions, formActions...)
	}
	return s.handleActions(actions)
}

func (s writer) handleActions(actions sqlActions) (writer, error) {
	walk := s
	var err error
	for _, action := range actions {
		walk, err = walk.handleAction(action)
		if err != nil {
			return writer{}, err
		}
	}
	return walk, nil
}

func (s writer) handleAction(action sqlAction) (writer, error) {
	if action.createColumn != nil {
		return s.handleCreateColumn(*action.createColumn)
	}
	if action.createTable != nil {
		return s.handleCreateTable(*action.createTable)
	}
	if action.createUniqueConstraint != nil {
		return s.handleCreateConstraint(*action.createUniqueConstraint)
	}
	if action.insertRow != nil {
		return s.handleInsertRow(*action.insertRow)
	}
	return writer{}, errors.New("could not handle action")
}

func (s writer) PutRecords(form types.FormInterface, records *types.RecordList) (Writer, error) {
	var actions sqlActions
	for _, item := range records.Items {
		recordActions, err := s.writeRecord(form, item)
		if err != nil {
			return nil, err
		}
		actions = append(actions, recordActions...)
	}
	return s.handleActions(actions)
}

type sqlArg struct {
	columnName string
	value      interface{}
}
type sqlArgs []sqlArg

func (s writer) writeRecord(form types.FormInterface, record *types.Record) (sqlActions, error) {

	var err error

	sqlParams := writeIdColumn(record)

	if _, ok := form.(types.SubFormInterface); ok {
		if sqlParams, err = writeOwnerIdColumn(record, sqlParams); err != nil {
			return nil, err
		}
	}

	for _, field := range form.GetFields() {

		fieldKind, err := field.FieldType.GetFieldKind()
		if err != nil {
			return nil, err
		}

		fieldValue, ok := record.Values.Find(field.ID)
		if !ok {
			continue
		}

		if fieldValue.Value.Kind == types.NullValue {
			sqlParams = append(sqlParams, sqlArg{
				columnName: field.ID,
				value:      nil,
			})
			continue
		}

		switch fieldKind {
		case types.FieldKindText:
			sqlParams, err = prepareTextFieldColumn(fieldValue, sqlParams)
		case types.FieldKindMultilineText:
			sqlParams, err = prepareMultilineTextFieldColumn(fieldValue, sqlParams)
		case types.FieldKindReference:
			sqlParams, err = prepareReferenceFieldColumn(fieldValue, sqlParams)
		case types.FieldKindSingleSelect:
			sqlParams, err = prepareSingleSelectFieldColumn(fieldValue, sqlParams)
		case types.FieldKindDate:
			sqlParams, err = prepareDateFieldColumn(fieldValue, sqlParams)
		case types.FieldKindMonth:
			sqlParams, err = prepareMonthFieldColumn(fieldValue, sqlParams)
		case types.FieldKindWeek:
			sqlParams, err = prepareWeekFieldColumn(fieldValue, sqlParams)
		case types.FieldKindQuantity:
			sqlParams, err = prepareQuantityFieldColumn(fieldValue, sqlParams)
		case types.FieldKindMultiSelect:
			continue
		default:
			err = fmt.Errorf("unhandled field kind %v", fieldKind)
		}
		if err != nil {
			return nil, err
		}
	}

	columnNames := make([]string, len(sqlParams))
	args := make([]interface{}, len(sqlParams))
	for i := 0; i < len(sqlParams); i++ {
		columnNames[i] = sqlParams[i].columnName
		args[i] = sqlParams[i].value
	}

	var result sqlActions

	result = append(result, sqlAction{
		insertRow: &sqlActionInsertRow{
			schemaName: form.GetDatabaseID(),
			tableName:  form.GetFormID(),
			columns:    columnNames,
			values:     args,
		},
	})

	return result, nil

}

func writeOwnerIdColumn(record *types.Record, sqlParams sqlArgs) (sqlArgs, error) {
	if record.OwnerID == nil {
		return nil, errors.New("no owner id on record")
	}
	sqlParams = append(sqlParams, sqlArg{
		columnName: keyOwnerIdColumn,
		value:      *record.OwnerID,
	})
	return sqlParams, nil
}

func writeIdColumn(record *types.Record) sqlArgs {
	var sqlParams = sqlArgs{
		{
			columnName: keyIdColumn,
			value:      record.ID,
		},
	}
	return sqlParams
}

func prepareTextFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	return writeGenericTextValue("text", fieldValue, sqlParams)
}

func prepareMultilineTextFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	return writeGenericTextValue("multiline text", fieldValue, sqlParams)
}

func prepareReferenceFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	return writeGenericTextValue("reference", fieldValue, sqlParams)
}

func prepareSingleSelectFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	return writeGenericTextValue("single select", fieldValue, sqlParams)
}

func writeGenericTextValue(fieldType string, fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	if err := assertStringValueType(fieldType, fieldValue); err != nil {
		return nil, err
	}
	return append(sqlParams, sqlArg{
		columnName: fieldValue.FieldID,
		value:      fieldValue.Value.StringValue,
	}), nil
}

func prepareMonthFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	if err := assertStringValueType("month", fieldValue); err != nil {
		return nil, err
	}
	sqlDate, err := time.Parse(monthFieldFormat, fieldValue.Value.StringValue)
	if err != nil {
		return nil, err
	}
	return append(sqlParams, sqlArg{
		columnName: fieldValue.FieldID,
		value:      sqlDate,
	}), nil
}

func prepareDateFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	if err := assertStringValueType("date", fieldValue); err != nil {
		return nil, err
	}
	sqlDate, err := time.Parse(dateFieldFormat, fieldValue.Value.StringValue)
	if err != nil {
		return nil, err
	}
	return append(sqlParams, sqlArg{
		columnName: fieldValue.FieldID,
		value:      sqlDate,
	}), nil
}

func prepareWeekFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	if err := assertStringValueType("week", fieldValue); err != nil {
		return nil, err
	}
	isoWeekTime, err := dates.ParseIsoWeekTime(fieldValue.Value.StringValue)
	if err != nil {
		return nil, err
	}
	return append(sqlParams, sqlArg{
		columnName: fieldValue.FieldID,
		value:      isoWeekTime,
	}), nil
}

func prepareQuantityFieldColumn(fieldValue types.FieldValue, sqlParams sqlArgs) (sqlArgs, error) {
	if err := assertStringValueType("quantity", fieldValue); err != nil {
		return nil, err
	}
	intValue, err := strconv.Atoi(fieldValue.Value.StringValue)
	if err != nil {
		return nil, err
	}
	return append(sqlParams, sqlArg{
		columnName: fieldValue.FieldID,
		value:      intValue,
	}), nil
}

func assertStringValueType(fieldType string, value types.FieldValue) error {
	if value.Value.Kind != types.StringValue {
		return fmt.Errorf("unsupported value kind %v for %s column", value.Value.Kind, fieldType)
	}
	return nil
}
