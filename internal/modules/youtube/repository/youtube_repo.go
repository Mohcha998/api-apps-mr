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

func NewYouTubeRepository(apiKey, channelID string, httpTimeout time.Duration, redis *redis.Client) domain.YoutubeRepository {
	return &youtubeRepository{
		apiKey:     apiKey,
		channelID:  channelID,
		httpClient: &http.Client{Timeout: httpTimeout},
		redis:      redis,
	}
}

// BuildURL membuat URL API YouTube
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

// FetchOrCache mengambil data dari Redis jika ada, jika tidak fetch API dan simpan ke Redis
func (r *youtubeRepository) FetchOrCache(ctx context.Context, cacheKey, reqURL string, ttl time.Duration) (*domain.Youtube, error) {
	if r.redis != nil {
		if s, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
			var cached domain.Youtube
			if err := json.Unmarshal([]byte(s), &cached); err == nil {
				return &cached, nil
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

	var data domain.Youtube
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if r.redis != nil {
		_ = r.redis.Set(ctx, cacheKey, string(body), ttl).Err()
	}

	return &data, nil
}

// GetActivity
func (r *youtubeRepository) GetActivity(ctx context.Context) (*domain.Youtube, error) {
	url := r.BuildURL("activities", map[string]string{
		"part":       "snippet,contentDetails",
		"maxResults": "1",
	})
	return r.FetchOrCache(ctx, "youtube_activity", url, 60*time.Second)
}

// GetLatest
func (r *youtubeRepository) GetLatest(ctx context.Context) (*domain.Youtube, error) {
	url := r.BuildURL("search", map[string]string{
		"part":       "snippet",
		"order":      "date",
		"type":       "video",
		"maxResults": "5",
	})
	return r.FetchOrCache(ctx, "youtube_latest", url, 60*time.Second)
}

// GetRecent
func (r *youtubeRepository) GetRecent(ctx context.Context) (*domain.Youtube, error) {
	url := r.BuildURL("activities", map[string]string{
		"part":       "snippet,contentDetails",
		"maxResults": "5",
	})
	return r.FetchOrCache(ctx, "youtube_recent", url, 60*time.Second)
}

// GetPlaylists
func (r *youtubeRepository) GetPlaylists(ctx context.Context) (*domain.Youtube, error) {
	url := r.BuildURL("playlists", map[string]string{
		"part":       "snippet,id,status,contentDetails",
		"maxResults": "20",
	})
	return r.FetchOrCache(ctx, "youtube_playlists", url, 60*time.Second)
}

// GetPlaylistItems
func (r *youtubeRepository) GetPlaylistItems(ctx context.Context, playlistId string) (*domain.Youtube, error) {
	cacheKey := "youtube_playlist_items_" + playlistId
	url := r.BuildURL("playlistItems", map[string]string{
		"part":       "snippet,status,id",
		"maxResults": "50",
		"playlistId": playlistId,
	})
	return r.FetchOrCache(ctx, cacheKey, url, 60*time.Second)
}
