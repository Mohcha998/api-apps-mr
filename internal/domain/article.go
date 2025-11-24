package domain

import "context"

type ArticleRepository interface {
    GetLatest(ctx context.Context) (LatestArticles, error)
}

type Article struct {
    Title     string `json:"title"`
    Image     string `json:"image"`
    CreatedAt string `json:"created_at"`
}

type LatestArticles struct {
    Status bool     `json:"status"`
    Data   []Article `json:"data"`
}
