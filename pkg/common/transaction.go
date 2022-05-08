package common

import (
	"github.com/nrc-no/core/pkg/store"
	"gorm.io/gorm"
)

type TransactionStore interface {
	Begin() *gorm.DB
	Commit(tx *gorm.DB) *gorm.DB
	Rollback(tx *gorm.DB) *gorm.DB
}

type transactionStore struct {
	db store.Factory
}

func NewTransactionStore(db store.Factory) TransactionStore {
	return &transactionStore{db: db}
}

func (s *transactionStore) Begin() *gorm.DB {
	db, err := s.db.Get()

	if err != nil {
		panic(err)
	}

	return db.Begin()
}

func (s *transactionStore) Commit(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (s *transactionStore) Rollback(tx *gorm.DB) *gorm.DB {
	return tx.Rollback()
}
