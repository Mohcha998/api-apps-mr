package usecase

import (
	"apps/internal/domain"
	"apps/internal/modules/merchandise/repository"
	"context"
	"fmt"
)

// -----------------------------
// Interface
// -----------------------------
type MerchandiseUsecase interface {
	GetAll(ctx context.Context) (*domain.MerchandiseAll, error)
	GetMZ(ctx context.Context) ([]map[string]interface{}, error)
	GetPrimerry(ctx context.Context) ([]map[string]interface{}, error)
	FindByTipe(ctx context.Context, id int) ([]map[string]interface{}, error)
	FindByID(ctx context.Context, id int) ([]domain.Merchandise, error)
}

// -----------------------------
// Struct
// -----------------------------
type merchandiseUC struct {
	repo repository.MerchandiseRepository
}

// -----------------------------
// Constructor
// -----------------------------
func NewMerchandiseUsecase(repo repository.MerchandiseRepository) MerchandiseUsecase {
	return &merchandiseUC{repo: repo}
}

// -----------------------------
// GetAll
// -----------------------------
func (u *merchandiseUC) GetAll(ctx context.Context) (*domain.MerchandiseAll, error) {
	result, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Mapping semua kategori dan merchandise ke string agar sama seperti CI
	for i, t := range result.MerchandiseTipe {
		var kategoriList []map[string]interface{}
		for _, k := range result.MerchandiseAll {
			// Pastikan tipe sesuai
			if fmt.Sprintf("%v", k["id_merchandise_tipe"]) == fmt.Sprintf("%d", t.ID) {
				kategoriList = append(kategoriList, k)
			}
		}
		result.MerchandiseAll[i] = map[string]interface{}{
			"id_merchandise_tipe":   fmt.Sprintf("%d", t.ID),
			"name_merchandise_tipe": t.Name,
			"merchandise_kategori":  kategoriList,
		}
	}

	// Semua merchandise sudah diambil dari repository
	return &result, nil
}

// -----------------------------
// GetMZ (kategori id=5)
// -----------------------------
func (u *merchandiseUC) GetMZ(ctx context.Context) ([]map[string]interface{}, error) {
	return u.repo.GetKategoriWithProducts(ctx, 5)
}

// -----------------------------
// GetPrimerry (kategori id=6)
// -----------------------------
func (u *merchandiseUC) GetPrimerry(ctx context.Context) ([]map[string]interface{}, error) {
	return u.repo.GetKategoriWithProducts(ctx, 6)
}

// -----------------------------
// FindByTipe
// -----------------------------
func (u *merchandiseUC) FindByTipe(ctx context.Context, id int) ([]map[string]interface{}, error) {
	return u.repo.GetByTipe(ctx, id)
}

// -----------------------------
// FindByID
// -----------------------------
func (u *merchandiseUC) FindByID(ctx context.Context, id int) ([]domain.Merchandise, error) {
	return u.repo.GetByID(ctx, id)
}
