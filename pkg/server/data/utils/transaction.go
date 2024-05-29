package utils

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nrc-no/core/pkg/server/data/api"
)

type Transaction struct {
	tx      *sqlx.Tx
	onQuery func(qry string, args []interface{})
}

func (t Transaction) Query(ctx context.Context, query string, args []interface{}) (api.ResultReader, error) {
	if t.onQuery != nil {
		t.onQuery(query, args)
	}
	res, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	r, err := NewSQLResultReader(res)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (t Transaction) Exec(ctx context.Context, query string, args []interface{}) (interface{}, error) {
	if t.onQuery != nil {
		t.onQuery(query, args)
	}
	return t.tx.ExecContext(ctx, query, args...)
}

func (t Transaction) Commit() error {
	return t.tx.Commit()
}

func (t Transaction) Rollback() error {
	return t.tx.Rollback()
}

var _ api.Transaction = &Transaction{}

func NewTransaction(ctx context.Context, db *sqlx.DB, onQuery func(qry string, args []interface{})) (api.Transaction, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		tx:      tx,
		onQuery: onQuery,
	}, nil
}
