package usecase

import (
	"context"
	"time"

	"apps/internal/domain"
	"apps/internal/modules/quotes/repository"
)

// QuotesUsecase adalah interface untuk semua logika bisnis (usecase) modul quotes.
// Layer ini berfungsi sebagai penghubung antara handler (delivery) dan repository.
type QuotesUsecase interface {
	CreatePool(ctx context.Context, q *domain.QuotePool) error
	UpdatePool(ctx context.Context, q *domain.QuotePool, changedBy int64) error
	ListPools(ctx context.Context, limit, offset int) ([]domain.QuotePool, error)
	GetOrAssignToday(ctx context.Context, userID int) (*domain.QuotePool, error)
	GetPoolByID(ctx context.Context, id int64) (*domain.QuotePool, error)
}

// quotesUc adalah implementasi konkret dari QuotesUsecase.
type quotesUc struct {
	repo repository.QuotesRepository
}

// NewQuotesUsecase membuat instance baru dari QuotesUsecase.
func NewQuotesUsecase(r repository.QuotesRepository) QuotesUsecase {
	return &quotesUc{repo: r}
}

// CreatePool membuat quote baru dan menyimpannya ke database.
func (u *quotesUc) CreatePool(ctx context.Context, q *domain.QuotePool) error {
	q.CreatedAt = time.Now()
	return u.repo.CreatePool(ctx, q)
}

// UpdatePool memperbarui quote yang sudah ada, serta menyimpan perubahan ke tabel history.
func (u *quotesUc) UpdatePool(ctx context.Context, q *domain.QuotePool, changedBy int64) error {
	now := time.Now()
	q.UpdatedAt = &now
	return u.repo.UpdatePool(ctx, q, changedBy)
}

// ListPools menampilkan daftar semua quotes, dengan limit dan offset (paging).
func (u *quotesUc) ListPools(ctx context.Context, limit, offset int) ([]domain.QuotePool, error) {
	return u.repo.ListPools(ctx, limit, offset)
}

// GetOrAssignToday akan mengembalikan quote untuk user tertentu di tanggal hari ini.
// Jika user belum memiliki quote hari ini, sistem akan otomatis mengassign satu quote acak dari pool.
func (u *quotesUc) GetOrAssignToday(ctx context.Context, userID int) (*domain.QuotePool, error) {
	pool, err := u.repo.GetOrAssignToday(ctx, userID)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// GetPoolByID mengambil data quote berdasarkan ID-nya.
func (u *quotesUc) GetPoolByID(ctx context.Context, id int64) (*domain.QuotePool, error) {
	return u.repo.GetPoolByID(ctx, id)
}
