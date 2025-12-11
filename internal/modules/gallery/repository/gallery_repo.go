package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/db"
	"context"
)

type GalleryRepository interface {
	GetAll(ctx context.Context) ([]domain.QuoteGallery, error)
}

type galleryRepo struct {
	db db.MysqlDBInterface
}

func NewGalleryRepository(db db.MysqlDBInterface) GalleryRepository {
	return &galleryRepo{db: db}
}

func (r *galleryRepo) GetAll(ctx context.Context) ([]domain.QuoteGallery, error) {
	var list []domain.QuoteGallery

	err := r.db.Conn().
		WithContext(ctx).
		Raw(`SELECT * FROM quote ORDER BY created_at DESC`).
		Scan(&list).Error

	return list, err
}
