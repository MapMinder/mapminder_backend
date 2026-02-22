package infrastructure

import (
	"context"

	"github.com/MapMinder/mapminder_backend/shared/tx"
	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "gorm_tx"

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) tx.Manager {
	return &TransactionManager{
		db: db,
	}
}

func (m *TransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(gormTx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, gormTx)
		return fn(txCtx)
	})
}

func ExtractTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if ok {
		return tx
	}
	return tx
}
