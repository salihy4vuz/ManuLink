package repositories

import (
	"context"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"gorm.io/gorm"
)

type QuoteRepository interface {
	Create(ctx context.Context, tx *gorm.DB, quote *models.Quote) error
}

type gormQuoteRepository struct{}

func NewQuoteRepository() QuoteRepository {
	return &gormQuoteRepository{}
}

func (r *gormQuoteRepository) Create(ctx context.Context, tx *gorm.DB, quote *models.Quote) error {
	return tx.WithContext(ctx).Create(quote).Error
}
