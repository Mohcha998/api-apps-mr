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

	// kategori khusus
	GetMZ(ctx context.Context) ([]map[string]interface{}, error)
	GetMrs(ctx context.Context) ([]map[string]interface{}, error)
	GetPrimerry(ctx context.Context) ([]map[string]interface{}, error)

	// generic
	GetKategoriWithProducts(ctx context.Context, kategoriIDs []int) ([]map[string]interface{}, error)

	// lainnya
	FindByTipe(ctx context.Context, id int) ([]map[string]interface{}, error)
	FindByID(ctx context.Context, id int) ([]domain.Merchandise, error)
	GetAllMerchandise(ctx context.Context) ([]domain.Merchandise, error)
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
// Usecase Methods
// -----------------------------
func (u *merchandiseUC) GetAllMerchandise(ctx context.Context) ([]domain.Merchandise, error) {
	return u.repo.GetAllMerchandise(ctx)
}

func (u *merchandiseUC) GetAll(ctx context.Context) (*domain.MerchandiseAll, error) {
	result, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Mapping agar konsisten seperti CI lama
	for i, t := range result.MerchandiseTipe {
		var kategoriList []map[string]interface{}

		for _, k := range result.MerchandiseAll {
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

	return &result, nil
}

func (u *merchandiseUC) GetMZ(ctx context.Context) ([]map[string]interface{}, error) {
	return u.repo.GetKategoriWithProducts(ctx, []int{5})
}

func (u *merchandiseUC) GetMrs(ctx context.Context) ([]map[string]interface{}, error) {
	return u.repo.GetKategoriWithProducts(ctx, []int{1, 2, 3})
}

func (u *merchandiseUC) GetPrimerry(ctx context.Context) ([]map[string]interface{}, error) {
	return u.repo.GetKategoriWithProducts(ctx, []int{6})
}

func (u *merchandiseUC) GetKategoriWithProducts(
	ctx context.Context,
	kategoriIDs []int,
) ([]map[string]interface{}, error) {
	return u.repo.GetKategoriWithProducts(ctx, kategoriIDs)
}

func (u *merchandiseUC) FindByTipe(ctx context.Context, id int) ([]map[string]interface{}, error) {
	return u.repo.GetByTipe(ctx, id)
}

func (u *merchandiseUC) FindByID(ctx context.Context, id int) ([]domain.Merchandise, error) {
	return u.repo.GetByID(ctx, id)
}
