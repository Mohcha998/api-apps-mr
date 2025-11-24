package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/db"
	"context"
)

type AdsRepository interface {
	GetActivePopup(ctx context.Context) (*domain.AdsPopup, error)
}

type adsRepository struct {
	db db.MysqlDBInterface
}

func NewAdsRepository(db db.MysqlDBInterface) AdsRepository {
	return &adsRepository{db: db}
}

func (r *adsRepository) GetActivePopup(ctx context.Context) (*domain.AdsPopup, error) {
	var ads domain.AdsPopup

	query := db.NewQuery("status = ?", 1)

	err := r.db.FindOne(
		ctx,
		&ads,
		db.WithQuery(query),
		db.WithOrder("id ASC"),
		db.WithLimit(1),
	)

	if err != nil {
		return nil, err
	}

	return &ads, nil
}
