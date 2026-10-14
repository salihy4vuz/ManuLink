package repositories

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type gormUnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &gormUnitOfWork{db: db}
}

func (u *gormUnitOfWork) Do(ctx context.Context, fn func(tx *gorm.DB) error) error {
	// Enterprise Deep Dive: Use Serializable Isolation Level to prevent Dirty/Phantom reads during concurrent quotes
	opts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	return u.db.WithContext(ctx).Transaction(fn, opts)
}
