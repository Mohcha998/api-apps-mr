package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/cache"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type youtubeRepository struct {
	apiKey     string
	channelID  string
	httpClient *http.Client
	cache      *cache.Client
}

// ✅ huruf besar di "New" supaya bisa diakses dari luar package
func NewYouTubeRepository(apiKey, channelID string, httpTimeout time.Duration, cache *cache.Client) domain.YoutubeRepository {
	return &youtubeRepository{
		apiKey:     apiKey,
		channelID:  channelID,
		httpClient: &http.Client{Timeout: httpTimeout},
		cache:      cache,
	}
}

// =====================================================
// Implementasi YoutubeRepository interface
// =====================================================

func (r *youtubeRepository) BuildURL(path string, params map[string]string) string {
	base := "https://www.googleapis.com/youtube/v3/" + path
	v := url.Values{}
	v.Set("key", r.apiKey)
	if _, ok := params["channelId"]; !ok && r.channelID != "" {
		v.Set("channelId", r.channelID)
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return base + "?" + v.Encode()
}

func (r *youtubeRepository) FetchOrCache(ctx context.Context, cacheKey, reqURL string, ttl time.Duration) (*domain.Youtube, error) {
	// 🔹 Coba ambil dari cache dulu
	if r.cache != nil {
		var cached domain.Youtube
		if err := r.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
			return &cached, nil
		}
	}

	// 🔹 Fetch dari API
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("youtube api error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data domain.Youtube
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	// 🔹 Simpan ke Redis
	if r.cache != nil {
		_ = r.cache.SetJSON(ctx, cacheKey, data, ttl)
	}

	return &data, nil
}
