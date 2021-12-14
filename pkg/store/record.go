package store

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sqlmanager"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	l.Debug("finding records")
	formReader := sqlmanager.NewFormReader(db)
	records, err := formReader.GetRecords(ctx, form)
	if err != nil {
		l.Error("failed to list records", zap.Error(err))
		return nil, err
	}

	if len(records.Items) != 1 {
		err := meta.NewInternalServerError(fmt.Errorf("unexpected number of records"))
		l.Error("should only have 1 record in result", zap.Error(err))
		return nil, err
	}

	return records.Items[0], nil

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

	formReader := sqlmanager.NewFormReader(db)
	records, err := formReader.GetRecords(ctx, form)
	if err != nil {
		l.Error("failed to list records", zap.Error(err))
		return nil, err
	}

	return records, nil
}

func (r recordStore) Create(ctx context.Context, record *types.Record) (*types.Record, error) {
	ctx, db, l, done, err := actionContext(ctx, r.db, "record", "create", zap.String("form_id", record.FormID))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("getting record form")
	rootForm, err := r.formStore.Get(ctx, record.FormID)
	if err != nil {
		l.Error("failed to get record form")
		return nil, err
	}

	l.Debug("getting form interface")
	form, err := rootForm.GetFormOrSubForm(record.FormID)
	if err != nil || form == nil {
		l.Error("failed to get form interface")
		return nil, meta.NewInternalServerError(err)
	}

	formWriter := sqlmanager.New()
	formWriter, err = formWriter.PutRecords(form, &types.RecordList{
		Items: []*types.Record{record},
	})
	if err != nil {
		l.Error("failed to put records", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, ddlItem := range formWriter.GetStatements() {
			if err := tx.Exec(ddlItem.Query, ddlItem.Args...).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, types.RecordRef{
		ID:         record.ID,
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
