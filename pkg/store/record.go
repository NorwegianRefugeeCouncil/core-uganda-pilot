package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/types"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type RecordStore interface {
	Get(ctx context.Context, recordRef types.RecordRef) (*types.Record, error)
	List(ctx context.Context, options types.RecordListOptions) (*types.RecordList, error)
	Create(ctx context.Context, record *types.Record) (*types.Record, error)
	Update(ctx context.Context, record *types.Record) (*types.Record, error)
	Delete(ctx context.Context, recordRef types.RecordRef) error
}

type recordStore struct {
	db        Factory
	formStore FormStore
}

func NewRecordStore(db Factory, formStore FormStore) RecordStore {
	return &recordStore{db: db, formStore: formStore}
}

func (r recordStore) Get(ctx context.Context, recordRef types.RecordRef) (*types.Record, error) {

	db, err := r.db.Get()
	if err != nil {
		return nil, err
	}

	rootForm, err := r.formStore.Get(ctx, recordRef.FormID)
	if err != nil {
		return nil, err
	}

	form, err := rootForm.GetFormInterface(recordRef.FormID)
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	if rootForm.DatabaseID != recordRef.DatabaseID {
		return nil, meta.NewNotFound(meta.GroupResource{
			Group:    "nrc.no",
			Resource: "records",
		}, fmt.Sprintf("%s/%s/%s",
			recordRef.DatabaseID,
			recordRef.FormID,
			recordRef.ID))
	}

	query := strings.Builder{}
	query.WriteString(fmt.Sprintf("select * from %s.%s where id = $1",
		pq.QuoteIdentifier(rootForm.DatabaseID),
		pq.QuoteIdentifier(form.GetID()),
	))

	sqlDB, err := db.DB()
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	result, err := sqlDB.Query(query.String(), recordRef.ID)
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	recordList, err := mapRecordList(result)
	if err != nil {
		return nil, err
	}

	if len(recordList.Items) != 1 {
		return nil, meta.NewInternalServerError(fmt.Errorf("unexpected number of records"))
	}

	return recordList.Items[0], nil

}

func (r recordStore) List(ctx context.Context, options types.RecordListOptions) (*types.RecordList, error) {

	rootForm, err := r.formStore.Get(ctx, options.FormID)
	if err != nil {
		return nil, err
	}

	if rootForm.DatabaseID != options.DatabaseID {
		return nil, fmt.Errorf("wrong database id: %s", options.DatabaseID)
	}

	form, err := rootForm.GetFormInterface(options.FormID)
	if err != nil {
		return nil, fmt.Errorf("wrong database id: %s", options.DatabaseID)
	}

	query := strings.Builder{}
	query.WriteString(fmt.Sprintf("select * from %s.%s",
		pq.QuoteIdentifier(rootForm.DatabaseID),
		pq.QuoteIdentifier(form.GetID()),
	))

	db, err := r.db.Get()
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	result, err := sqlDB.Query(query.String())
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	recordList, err := mapRecordList(result)
	if err != nil {
		return nil, err
	}

	return recordList, nil
}

func mapRecordList(result *sql.Rows) (*types.RecordList, error) {

	cols, err := result.Columns()
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	var records []*types.Record
	for result.Next() {
		record, err := mapRecord(cols, result)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	recordList := &types.RecordList{
		Items: records,
	}
	return recordList, nil
}

func mapRecord(cols []string, result *sql.Rows) (*types.Record, error) {
	var recordValues = map[string]interface{}{}
	if err := mapRecordRow(cols, result, recordValues); err != nil {
		return nil, err
	}

	recordID := recordValues["id"].(string)
	recordDatabaseID := recordValues["database_id"].(string)
	recordFormID := recordValues["form_id"].(string)
	var parentId *string = nil
	if parentIdIntf, ok := recordValues["parent_id"]; ok {
		if parentIdStr, ok := parentIdIntf.(string); ok {
			parentId = &parentIdStr
		}
	}

	delete(recordValues, "id")
	delete(recordValues, "database_id")
	delete(recordValues, "form_id")
	delete(recordValues, "parent_id")

	record := &types.Record{
		ID:         recordID,
		DatabaseID: recordDatabaseID,
		FormID:     recordFormID,
		ParentID:   parentId,
		Values:     recordValues,
	}
	return record, nil
}

func (r recordStore) Create(ctx context.Context, record *types.Record) (*types.Record, error) {
	var keys []string
	var values []interface{}
	var params []string
	i := 1

	rootform, err := r.formStore.Get(ctx, record.FormID)
	if err != nil {
		return nil, err
	}
	databaseId := rootform.DatabaseID

	form, err := rootform.GetFormInterface(record.FormID)
	if err != nil {
		return nil, err
	}

	record.ID = uuid.NewV4().String()

	for _, field := range form.GetFields() {
		if fieldValue, ok := record.Values[field.ID]; ok {
			keys = append(keys, pq.QuoteIdentifier(field.ID))
			values = append(values, fieldValue)
			params = append(params, fmt.Sprintf("$%d", i))
			i++
		}
	}

	keys = append(keys, "id")
	values = append(values, record.ID)
	params = append(params, fmt.Sprintf("$%d", i))
	i++

	keys = append(keys, "database_id")
	values = append(values, record.DatabaseID)
	params = append(params, fmt.Sprintf("$%d", i))
	i++

	keys = append(keys, "form_id")
	values = append(values, record.FormID)
	params = append(params, fmt.Sprintf("$%d", i))
	i++

	if record.ParentID != nil {
		keys = append(keys, "parent_id")
		values = append(values, record.ParentID)
		params = append(params, fmt.Sprintf("$%d", i))
		i++
	}

	insertQry := strings.Builder{}
	insertQry.WriteString(fmt.Sprintf("insert into %s.%s (%s) values (%s) returning id",
		pq.QuoteIdentifier(databaseId),
		pq.QuoteIdentifier(form.GetID()),
		strings.Join(keys, ","),
		strings.Join(params, ","),
	))

	db, err := r.db.Get()
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	var lastInsertedId string
	insertSQLQuery := insertQry.String()
	if err := sqlDB.QueryRow(insertSQLQuery, values...).Scan(&lastInsertedId); err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	return r.Get(ctx, types.RecordRef{
		ID:         lastInsertedId,
		DatabaseID: record.DatabaseID,
		FormID:     record.FormID,
	})
}

func (r recordStore) Update(ctx context.Context, record *types.Record) (*types.Record, error) {
	panic("implement me")
}

func (r recordStore) Delete(ctx context.Context, recordRef types.RecordRef) error {
	panic("implement me")
}

var _ RecordStore = &recordStore{}

func mapRecordRow(cols []string, rows *sql.Rows, into map[string]interface{}) error {

	cols, err := rows.Columns()
	if err != nil {
		return meta.NewInternalServerError(err)
	}

	dest := make([]interface{}, len(cols))
	args := make([]interface{}, len(cols))

	for i := range cols {
		args[i] = &(dest[i])
	}

	if err := rows.Scan(args...); err != nil {
		return meta.NewInternalServerError(err)
	}

	for i, col := range cols {
		into[col] = dest[i]
	}

	return nil
}
