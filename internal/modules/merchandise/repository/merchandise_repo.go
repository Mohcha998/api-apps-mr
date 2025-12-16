package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"context"
	"fmt"
	"time"
)

// -----------------------------
// Interface
// -----------------------------
type MerchandiseRepository interface {
	GetAll(ctx context.Context) (domain.MerchandiseAll, error)
	GetKategoriWithProducts(ctx context.Context, kategoriIDs []int) ([]map[string]interface{}, error)
	GetByTipe(ctx context.Context, tipeID int) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, merchID int) ([]domain.Merchandise, error)
	GetAllMerchandise(ctx context.Context) ([]domain.Merchandise, error)
}

// -----------------------------
// Struct
// -----------------------------
type merchandiseRepo struct {
	db    db.MysqlDBInterface
	cache *cache.Client
}

// -----------------------------
// Constructor
// -----------------------------
func NewMerchandiseRepository(
	db db.MysqlDBInterface,
	cache *cache.Client,
) MerchandiseRepository {
	return &merchandiseRepo{
		db:    db,
		cache: cache,
	}
}

// -----------------------------
// Cache Helpers
// -----------------------------
func merchKey(key string) string {
	return "merchandise:" + key
}

func merchTTL() time.Duration {
	return 24 * time.Hour
}

// -----------------------------
// GetAllMerchandise
// -----------------------------
func (r *merchandiseRepo) GetAllMerchandise(ctx context.Context) ([]domain.Merchandise, error) {
	cacheKey := merchKey("all")
	var cached []domain.Merchandise

	if err := r.cache.GetJSON(ctx, cacheKey, &cached); err == nil && len(cached) > 0 {
		return cached, nil
	}

	var merch []domain.Merchandise
	err := r.db.Conn().
		WithContext(ctx).
		Raw(`SELECT * FROM merchandise ORDER BY id DESC`).
		Scan(&merch).Error

	if err != nil {
		return nil, err
	}

	_ = r.cache.SetJSON(ctx, cacheKey, merch, merchTTL())
	return merch, nil
}

// -----------------------------
// Helpers (DB ONLY)
// -----------------------------
func (r *merchandiseRepo) getCategories(ctx context.Context, tipeID int) ([]domain.MerchandiseKategori, error) {
	var kategori []domain.MerchandiseKategori
	err := r.db.Conn().Raw(
		`SELECT * FROM merchandise_kategori WHERE id_merchandise_tipe = ? AND status = 1`,
		tipeID,
	).Scan(&kategori).Error
	return kategori, err
}

func (r *merchandiseRepo) getProducts(ctx context.Context, kategoriID int) ([]domain.Merchandise, error) {
	var merch []domain.Merchandise
	err := r.db.Conn().Raw(
		`SELECT * FROM merchandise WHERE id_merchandise_kategori = ? AND status = 1`,
		kategoriID,
	).Scan(&merch).Error
	return merch, err
}

// -----------------------------
// GetAll (Home)
// -----------------------------
func (r *merchandiseRepo) GetAll(ctx context.Context) (domain.MerchandiseAll, error) {
	var result domain.MerchandiseAll
	cacheKey := merchKey("home")

	if err := r.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		return result, nil
	}

	var tipe []domain.MerchandiseTipe
	if err := r.db.Conn().
		Raw(`SELECT * FROM merchandise_tipe WHERE status = 1 AND id NOT IN (5,6)`).
		Scan(&tipe).Error; err != nil {
		return result, err
	}
	result.MerchandiseTipe = tipe

	var list []map[string]interface{}
	for _, t := range tipe {
		item := map[string]interface{}{
			"id_merchandise_tipe":   t.ID,
			"name_merchandise_tipe": t.Name,
		}

		kategori, _ := r.getCategories(ctx, t.ID)
		var kategoriList []map[string]interface{}

		for _, k := range kategori {
			products, _ := r.getProducts(ctx, k.ID)
			kategoriList = append(kategoriList, map[string]interface{}{
				"id_merchandise_kategori":   k.ID,
				"name_merchandise_kategori": k.Name,
				"merchandise":               products,
			})
		}

		item["merchandise_kategori"] = kategoriList
		list = append(list, item)
	}

	result.MerchandiseAll = list

	var merchAll []domain.Merchandise
	_ = r.db.Conn().
		Raw(`SELECT * FROM merchandise WHERE status = 1`).
		Scan(&merchAll).Error

	result.Merchandise = merchAll

	_ = r.cache.SetJSON(ctx, cacheKey, result, merchTTL())
	return result, nil
}

// -----------------------------
// GetKategoriWithProducts
// -----------------------------
func (r *merchandiseRepo) GetKategoriWithProducts(
	ctx context.Context,
	kategoriIDs []int,
) ([]map[string]interface{}, error) {

	cacheKey := merchKey(fmt.Sprintf("kategori:%v", kategoriIDs))
	var cached []map[string]interface{}

	if err := r.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	type kategoriRow struct {
		ID                  int
		IDMerchandiseTipe   int
		Code                string
		Name                string
		Status              int
		CreatedDate         time.Time
		NameMerchandiseTipe string
	}

	var rows []kategoriRow

	err := r.db.Conn().Raw(`
		SELECT mk.id,
		       mk.id_merchandise_tipe,
		       mk.code,
		       mk.name,
		       mk.status,
		       mk.created_date,
		       mt.name AS name_merchandise_tipe
		FROM merchandise_kategori mk
		JOIN merchandise_tipe mt ON mt.id = mk.id_merchandise_tipe
		WHERE mk.id IN ?
		  AND mk.status = 1
		ORDER BY mk.id ASC
	`, kategoriIDs).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, row := range rows {
		products, _ := r.getProducts(ctx, row.ID)

		result = append(result, map[string]interface{}{
			"id":                    fmt.Sprintf("%d", row.ID),
			"id_merchandise_tipe":   fmt.Sprintf("%d", row.IDMerchandiseTipe),
			"code":                  row.Code,
			"name":                  row.Name,
			"status":                fmt.Sprintf("%d", row.Status),
			"created_date":          row.CreatedDate.Format("2006-01-02 15:04:05"),
			"name_merchandise_tipe": row.NameMerchandiseTipe,
			"merchandise":           products,
		})
	}

	_ = r.cache.SetJSON(ctx, cacheKey, result, merchTTL())
	return result, nil
}

// -----------------------------
// GetByTipe
// -----------------------------
func (r *merchandiseRepo) GetByTipe(ctx context.Context, tipeID int) ([]map[string]interface{}, error) {
	cacheKey := merchKey(fmt.Sprintf("tipe:%d", tipeID))
	var cached []map[string]interface{}

	if err := r.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	kategori, _ := r.getCategories(ctx, tipeID)
	var list []map[string]interface{}

	for _, k := range kategori {
		products, _ := r.getProducts(ctx, k.ID)
		list = append(list, map[string]interface{}{
			"id_merchandise_kategori":   k.ID,
			"name_merchandise_kategori": k.Name,
			"merchandise":               products,
		})
	}

	_ = r.cache.SetJSON(ctx, cacheKey, list, merchTTL())
	return list, nil
}

// -----------------------------
// GetByID
// -----------------------------
func (r *merchandiseRepo) GetByID(ctx context.Context, merchID int) ([]domain.Merchandise, error) {
	cacheKey := merchKey(fmt.Sprintf("id:%d", merchID))
	var cached []domain.Merchandise

	if err := r.cache.GetJSON(ctx, cacheKey, &cached); err == nil && len(cached) > 0 {
		return cached, nil
	}

	var merch []domain.Merchandise
	err := r.db.Conn().
		Raw(`SELECT * FROM merchandise WHERE id = ? AND status = 1`, merchID).
		Scan(&merch).Error

	if err != nil {
		return nil, err
	}

	_ = r.cache.SetJSON(ctx, cacheKey, merch, merchTTL())
	return merch, nil
}
