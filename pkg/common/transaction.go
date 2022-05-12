package common

import (
	"github.com/nrc-no/core/pkg/store"
	"gorm.io/gorm"
)

type TransactionManager interface {
	Begin() *gorm.DB
	Commit(tx *gorm.DB) *gorm.DB
	Rollback(tx *gorm.DB) *gorm.DB
}

type transactionManager struct {
	db store.Factory
}

func NewTransactionManager(db store.Factory) TransactionManager {
	return &transactionManager{db: db}
}

func (s *transactionManager) Begin() *gorm.DB {
	db, err := s.db.Get()

	if err != nil {
		panic(err)
	}

	return db.Begin()
}

func (s *transactionManager) Commit(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (s *transactionManager) Rollback(tx *gorm.DB) *gorm.DB {
	return tx.Rollback()
}
