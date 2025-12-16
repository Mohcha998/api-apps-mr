package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"context"
	"time"
)

type GalleryRepository interface {
	GetAll(ctx context.Context) ([]domain.QuoteGallery, error)
}

type galleryRepo struct {
	db    db.MysqlDBInterface
	cache *cache.Client
}

func NewGalleryRepository(
	db db.MysqlDBInterface,
	cache *cache.Client,
) GalleryRepository {
	return &galleryRepo{
		db:    db,
		cache: cache,
	}
}

func galleryTTL() time.Duration {
	return 24 * time.Hour
}

func (r *galleryRepo) GetAll(ctx context.Context) ([]domain.QuoteGallery, error) {
	cacheKey := "gallery:all"

	// 1️⃣ Try cache
	var list []domain.QuoteGallery
	if err := r.cache.GetJSON(ctx, cacheKey, &list); err == nil {
		return list, nil
	}

	// 2️⃣ Query DB
	err := r.db.Conn().
		WithContext(ctx).
		Raw(`SELECT * FROM quote ORDER BY created_at DESC`).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	// 3️⃣ Save to cache
	_ = r.cache.SetJSON(ctx, cacheKey, list, galleryTTL())

	return list, nil
}
