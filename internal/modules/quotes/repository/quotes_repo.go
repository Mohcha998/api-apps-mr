package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"context"
	"fmt"
	"time"
	"math/rand"
)

type QuotesRepository interface {
	CreatePool(ctx context.Context, q *domain.QuotePool) error
	UpdatePool(ctx context.Context, q *domain.QuotePool, changedBy int64) error
	ListPools(ctx context.Context, limit, offset int) ([]domain.QuotePool, error)
	GetPoolByID(ctx context.Context, id int64) (*domain.QuotePool, error)
	GetOrAssignToday(ctx context.Context, userID int) (*domain.QuotePool, error)
}

type quotesRepo struct {
	db    db.MysqlDBInterface
	cache *cache.Client
}

func NewQuotesRepository(db db.MysqlDBInterface, cache *cache.Client) QuotesRepository {
	return &quotesRepo{
		db:    db,
		cache: cache,
	}
}

// =====================================================

func (r *quotesRepo) CreatePool(ctx context.Context, q *domain.QuotePool) error {
	return r.db.Create(ctx, q)
}

func (r *quotesRepo) UpdatePool(ctx context.Context, q *domain.QuotePool, changedBy int64) error {
	// 1️⃣ Simpan ke history
	history := domain.QuotePoolHistory{
		QuotePoolID: q.ID,
		Text:        q.Text,
		Author:      q.Author,
		ChangedBy:   &changedBy,
		ChangedAt:   time.Now(),
	}
	if err := r.db.Create(ctx, &history); err != nil {
		return fmt.Errorf("failed to create history: %v", err)
	}

	// 2️⃣ Ambil data lama untuk pertahankan CreatedAt
	var old domain.QuotePool
	if err := r.db.FindOne(ctx, &old, db.WithQuery(db.NewQuery("id = ?", q.ID))); err != nil {
		return fmt.Errorf("failed to get existing quote: %v", err)
	}

	// 3️⃣ Pertahankan created_at lama
	q.CreatedAt = old.CreatedAt
	now := time.Now()
	q.UpdatedAt = &now

	// 4️⃣ Update ke DB tanpa ubah created_at
	if err := r.db.Update(ctx, q); err != nil {
		return fmt.Errorf("failed to update quote: %v", err)
	}

	return nil
}


func (r *quotesRepo) ListPools(ctx context.Context, limit, offset int) ([]domain.QuotePool, error) {
	var pools []domain.QuotePool
	err := r.db.Find(ctx, &pools, db.WithLimit(limit), db.WithOffset(offset))
	return pools, err
}

func (r *quotesRepo) GetPoolByID(ctx context.Context, id int64) (*domain.QuotePool, error) {
	cacheKey := fmt.Sprintf("quotes:pool:%d", id)
	var cached domain.QuotePool

	if err := r.cache.GetJSON(ctx, cacheKey, &cached); err == nil && cached.ID != 0 {
		return &cached, nil
	}

	var pool domain.QuotePool
	if err := r.db.FindOne(ctx, &pool, db.WithQuery(db.NewQuery("id = ?", id))); err != nil {
		return nil, err
	}

	_ = r.cache.SetJSON(ctx, cacheKey, pool, 24*time.Hour)
	return &pool, nil
}

func (r *quotesRepo) GetOrAssignToday(ctx context.Context, userID int) (*domain.QuotePool, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")

	cacheKey := fmt.Sprintf("quotes:user:%d:%s", userID, today)

	// 1️⃣ Cek dari cache
	var cached domain.QuotePool
	if err := r.cache.GetJSON(ctx, cacheKey, &cached); err == nil && cached.ID != 0 {
		return &cached, nil
	}

	// 2️⃣ Ambil assignment terakhir user (pakai query manual)
	var assigned domain.QuoteAssigned
	query := db.NewQuery("user_id = ? ORDER BY created_at DESC LIMIT 1", userID)

	err := r.db.FindOne(ctx, &assigned, db.WithQuery(query))
	if err == nil && assigned.ID != 0 {
		// Cek tanggal assign
		if assigned.AssignDate == today {
			// Sudah ada quote hari ini
			quote, err := r.GetPoolByID(ctx, assigned.QuotePoolID)
			if err == nil {
				_ = r.cache.SetJSON(ctx, cacheKey, quote, ttlUntilMidnight())
				return quote, nil
			}
		}
	}

	// 3️⃣ Ambil semua quote aktif
	var quotes []domain.QuotePool
	if err := r.db.Find(ctx, &quotes, db.WithQuery(db.NewQuery("is_active = ?", true))); err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return nil, fmt.Errorf("no active quotes found")
	}

	// 4️⃣ Pilih quote acak (hindari yang sama)
	rand.Seed(time.Now().UnixNano())
	randomQuote := quotes[rand.Intn(len(quotes))]
	if assigned.QuotePoolID == randomQuote.ID && len(quotes) > 1 {
		for randomQuote.ID == assigned.QuotePoolID {
			randomQuote = quotes[rand.Intn(len(quotes))]
		}
	}

	// 5️⃣ Simpan assignment baru
	newAssign := domain.QuoteAssigned{
		UserID:      userID,
		QuotePoolID: randomQuote.ID,
		AssignDate:  today,
		CreatedAt:   now,
	}
	_ = r.db.Create(ctx, &newAssign)

	// 6️⃣ Cache sampai tengah malam
	_ = r.cache.SetJSON(ctx, cacheKey, randomQuote, ttlUntilMidnight())

	return &randomQuote, nil
}

// Fungsi bantu untuk hitung durasi cache sampai jam 23:59:59
func ttlUntilMidnight() time.Duration {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	return time.Until(endOfDay)
}
