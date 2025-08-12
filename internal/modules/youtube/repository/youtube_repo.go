package repository

import (
	"apps/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

type youtubeRepository struct {
	apiKey     string
	channelID  string
	httpClient *http.Client
	redis      *redis.Client
}

func NewYouTubeRepository(apiKey, channelID string, httpTimeout time.Duration, redis *redis.Client) *youtubeRepository {
	return &youtubeRepository{
		apiKey:     apiKey,
		channelID:  channelID,
		httpClient: &http.Client{Timeout: httpTimeout},
		redis:      redis,
	}
}

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
	if r.redis != nil {
		if s, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
			var cached *domain.Youtube
			if err := json.Unmarshal([]byte(s), &cached); err == nil {
				return cached, nil
			}
		}
	}

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

	var data *domain.Youtube
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if r.redis != nil {
		_ = r.redis.Set(ctx, cacheKey, string(body), ttl).Err()
	}

	return data, nil
}