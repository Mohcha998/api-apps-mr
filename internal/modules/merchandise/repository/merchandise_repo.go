package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/db"
	"context"
	"fmt"
	"time"
)

type MerchandiseRepository interface {
	GetAll(ctx context.Context) (domain.MerchandiseAll, error)
	GetKategoriWithProducts(ctx context.Context, kategoriID int) ([]map[string]interface{}, error) // mz / primerry
	GetByTipe(ctx context.Context, tipeID int) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, merchID int) ([]domain.Merchandise, error)
	GetAllMerchandise(ctx context.Context) ([]domain.Merchandise, error)
}

type merchandiseRepo struct {
	db db.MysqlDBInterface
}

func NewMerchandiseRepository(db db.MysqlDBInterface) MerchandiseRepository {
	return &merchandiseRepo{db: db}
}

func (r *merchandiseRepo) GetAllMerchandise(ctx context.Context) ([]domain.Merchandise, error) {
	var merch []domain.Merchandise

	err := r.db.Conn().
		WithContext(ctx).
		Raw(`SELECT * FROM merchandise ORDER BY id DESC`).
		Scan(&merch).Error

	if err != nil {
		return nil, err
	}

	return merch, nil
}

func (r *merchandiseRepo) getCategories(ctx context.Context, tipeID int) ([]domain.MerchandiseKategori, error) {
	var kategori []domain.MerchandiseKategori
	err := r.db.Conn().Raw(
		`SELECT * FROM merchandise_kategori WHERE id_merchandise_tipe = ? AND status = 1`, tipeID,
	).Scan(&kategori).Error
	return kategori, err
}

func (r *merchandiseRepo) getProducts(ctx context.Context, kategoriID int) ([]domain.Merchandise, error) {
	var merch []domain.Merchandise
	err := r.db.Conn().Raw(
		`SELECT * FROM merchandise WHERE id_merchandise_kategori = ? AND status = 1`, kategoriID,
	).Scan(&merch).Error
	return merch, err
}

// -----------------------------
// MAIN FUNCTIONS
// -----------------------------

func (r *merchandiseRepo) GetAll(ctx context.Context) (domain.MerchandiseAll, error) {
	var result domain.MerchandiseAll

	// 1. Ambil semua tipe merchandise kecuali 5 & 6
	var tipe []domain.MerchandiseTipe
	if err := r.db.Conn().Raw(`SELECT * FROM merchandise_tipe WHERE status = 1 AND id NOT IN (5,6)`).Scan(&tipe).Error; err != nil {
		return result, err
	}
	result.MerchandiseTipe = tipe

	// 2. Mapping tipe -> kategori -> merchandise
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

	// 3. Ambil semua merchandise
	var merchAll []domain.Merchandise
	if err := r.db.Conn().Raw(`SELECT * FROM merchandise WHERE status = 1`).Scan(&merchAll).Error; err != nil {
		return result, err
	}
	result.Merchandise = merchAll

	return result, nil
}

// -----------------------------
// mz() & primerry()
// -----------------------------
func (r *merchandiseRepo) GetKategoriWithProducts(ctx context.Context, kategoriID int) ([]map[string]interface{}, error) {
	type kategoriRow struct {
		ID                  int
		IDMerchandiseTipe   int
		Code                string
		Name                string
		Status              int
		CreatedDate         time.Time
		NameMerchandiseTipe string
	}

	var row kategoriRow
	err := r.db.Conn().Raw(`
		SELECT mk.id, mk.id_merchandise_tipe, mk.code, mk.name, mk.status, mk.created_date,
		       mt.name as name_merchandise_tipe
		FROM merchandise_kategori mk
		JOIN merchandise_tipe mt ON mt.id = mk.id_merchandise_tipe
		WHERE mk.id = ? AND mk.status = 1
	`, kategoriID).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	// Ambil merchandise
	products, _ := r.getProducts(ctx, row.ID)

	// Mapping ke string agar sama seperti CI
	data := map[string]interface{}{
		"id":                    fmt.Sprintf("%d", row.ID),
		"id_merchandise_tipe":   fmt.Sprintf("%d", row.IDMerchandiseTipe),
		"code":                  row.Code,
		"name":                  row.Name,
		"status":                fmt.Sprintf("%d", row.Status),
		"created_date":          row.CreatedDate.Format("2006-01-02 15:04:05"),
		"name_merchandise_tipe": row.NameMerchandiseTipe,
		"merchandise":           products,
	}

	return []map[string]interface{}{data}, nil
}

// -----------------------------
// findByTipe()
// -----------------------------
func (r *merchandiseRepo) GetByTipe(ctx context.Context, tipeID int) ([]map[string]interface{}, error) {
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

	return list, nil
}

// -----------------------------
// findByID()
// -----------------------------
func (r *merchandiseRepo) GetByID(ctx context.Context, merchID int) ([]domain.Merchandise, error) {
	var merch []domain.Merchandise
	err := r.db.Conn().Raw(`SELECT * FROM merchandise WHERE id = ? AND status = 1`, merchID).Scan(&merch).Error
	return merch, err
}
