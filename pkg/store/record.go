package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
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
	ctx, db, l, done, err := actionContext(ctx, r.db, "record", "get", zap.String("record_id", recordRef.ID))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("getting record form")
	rootForm, err := r.formStore.Get(ctx, recordRef.FormID)
	if err != nil {
		l.Error("failed to get record form")
		return nil, err
	}

	l.Debug("getting form interface")
	form, err := rootForm.GetFormOrSubForm(recordRef.FormID)
	if err != nil || form == nil {
		l.Error("failed to get form interface")
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("ensuring record database id match")
	if rootForm.DatabaseID != recordRef.DatabaseID {
		l.Error("failed to verify record database id match")
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
		pq.QuoteIdentifier(form.GetFormID()),
	))

	l.Debug("getting raw sql database")
	sqlDB, err := db.DB()
	if err != nil {
		l.Error("failed to get raw sql database", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("finding records")
	result, err := sqlDB.Query(query.String(), recordRef.ID)
	if err != nil {
		l.Error("failed to find records", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("mapping records")
	recordList, err := mapRecordList(form.GetDatabaseID(), form.GetFormID(), result)
	if err != nil {
		l.Error("failed to map records", zap.Error(err))
		return nil, err
	}

	if len(recordList.Items) != 1 {
		err := meta.NewInternalServerError(fmt.Errorf("unexpected number of records"))
		l.Error("should only have 1 record in result", zap.Error(err))
		return nil, err
	}

	return recordList.Items[0], nil

}

func (r recordStore) List(ctx context.Context, options types.RecordListOptions) (*types.RecordList, error) {
	ctx, db, l, done, err := actionContext(ctx, r.db, "record", "list",
		zap.String("database_id", options.DatabaseID),
		zap.String("form_Id", options.FormID))

	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("getting record form")
	rootForm, err := r.formStore.Get(ctx, options.FormID)
	if err != nil {
		l.Error("failed to get record form", zap.Error(err))
		return nil, err
	}

	l.Debug("verifying form database id match")
	if rootForm.DatabaseID != options.DatabaseID {
		l.Error("failed to verify form database id match")
		return nil, fmt.Errorf("wrong database id: %s", options.DatabaseID)
	}

	l.Debug("getting form interface")
	form, err := rootForm.GetFormOrSubForm(options.FormID)
	if err != nil || form == nil {
		l.Error("failed to get form interface")
		return nil, meta.NewInternalServerError(err)
	}

	query := strings.Builder{}
	query.WriteString(fmt.Sprintf("select * from %s.%s",
		pq.QuoteIdentifier(rootForm.DatabaseID),
		pq.QuoteIdentifier(form.GetFormID()),
	))

	l.Debug("getting raw sql database")
	sqlDB, err := db.DB()
	if err != nil {
		l.Error("failed to get raw sql database", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("listing records")
	result, err := sqlDB.Query(query.String())
	if err != nil {
		l.Error("failed to list records", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("mapping records")
	recordList, err := mapRecordList(form.GetDatabaseID(), form.GetFormID(), result)
	if err != nil {
		l.Error("failed to map records", zap.Error(err))
		return nil, err
	}

	return recordList, nil
}

func mapRecordList(databaseId, formId string, result *sql.Rows) (*types.RecordList, error) {

	cols, err := result.Columns()
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	var records []*types.Record
	for result.Next() {
		record, err := mapRecord(databaseId, formId, cols, result)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if records == nil {
		records = []*types.Record{}
	}

	recordList := &types.RecordList{
		Items: records,
	}
	return recordList, nil
}

func mapRecord(databaseId, formId string, cols []string, result *sql.Rows) (*types.Record, error) {
	var recordValues = map[string]interface{}{}
	if err := mapRecordRow(cols, result, recordValues); err != nil {
		return nil, err
	}

	recordID := recordValues["id"].(string)
	recordSeq := recordValues["seq"].(int64)
	var ownerId *string = nil
	if ownerIdIntf, ok := recordValues["owner_id"]; ok {
		if ownerIdStr, ok := ownerIdIntf.(string); ok {
			ownerId = &ownerIdStr
		}
	}

	delete(recordValues, "id")
	delete(recordValues, "seq")
	delete(recordValues, "database_id")
	delete(recordValues, "form_id")
	delete(recordValues, "owner_id")
	delete(recordValues, "created_at")

	record := &types.Record{
		ID:         recordID,
		Seq:        recordSeq,
		DatabaseID: databaseId,
		FormID:     formId,
		OwnerID:    ownerId,
		Values:     recordValues,
	}
	return record, nil
}

func (r recordStore) Create(ctx context.Context, record *types.Record) (*types.Record, error) {
	ctx, db, l, done, err := actionContext(ctx, r.db, "record", "create", zap.String("form_id", record.FormID))
	if err != nil {
		return nil, err
	}
	defer done()

	var keys []string
	var values []interface{}
	var params []string
	i := 1

	l.Debug("getting root form")
	rootForm, err := r.formStore.Get(ctx, record.FormID)
	if err != nil {
		l.Error("failed to get record root form", zap.Error(err))
		return nil, err
	}
	databaseId := rootForm.DatabaseID

	l.Debug("getting form interface")
	form, err := rootForm.GetFormOrSubForm(record.FormID)
	if err != nil || form == nil {
		l.Error("failed to get form interface")
		return nil, meta.NewInternalServerError(err)
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

	if record.OwnerID != nil {
		keys = append(keys, "owner_id")
		values = append(values, record.OwnerID)
		params = append(params, fmt.Sprintf("$%d", i))
		i++
	}

	insertQry := strings.Builder{}
	insertQry.WriteString(fmt.Sprintf("insert into %s.%s (%s) values (%s) returning id",
		pq.QuoteIdentifier(databaseId),
		pq.QuoteIdentifier(form.GetFormID()),
		strings.Join(keys, ","),
		strings.Join(params, ","),
	))

	l.Debug("getting raw sql database")
	sqlDB, err := db.DB()
	if err != nil {
		l.Error("failed to get raw sql database", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("inserting record")
	var lastInsertedId string
	insertSQLQuery := insertQry.String()
	if err := sqlDB.QueryRow(insertSQLQuery, values...).Scan(&lastInsertedId); err != nil {
		l.Error("failed to insert record", zap.Error(err))
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
