package repository

import (
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "time"

    "apps/internal/domain"
    "apps/internal/infrastructure/cache"
)

type articleRepository struct {
    redis *cache.Client
}

func NewArticleRepository(redis *cache.Client) domain.ArticleRepository {
    return &articleRepository{redis: redis}
}

func (r *articleRepository) GetLatest(ctx context.Context) (domain.LatestArticles, error) {
    var result domain.LatestArticles
    const cacheKey = "latest_articles"

    // 🔹 Cek Redis terlebih dahulu
    cached, err := r.redis.Get(ctx, cacheKey)
    if err == nil && cached != "" {
        _ = json.Unmarshal([]byte(cached), &result)
        return result, nil
    }

    // 🔹 Hit API eksternal
    resp, err := http.Get("https://article.merryriana.com/api/latest-posts")
    if err != nil {
        return result, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return result, errors.New("failed to fetch latest posts")
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return result, err
    }

    // 🔹 Decode JSON
    if err := json.Unmarshal(body, &result); err != nil {
        return result, err
    }

    // 🔹 Simpan ke Redis (10 menit)
    _ = r.redis.Set(ctx, cacheKey, string(body), 10*time.Minute)

    return result, nil
}
