package req

type CreateQuoteReq struct {
	Text      string `json:"text" validate:"required"`
	Author    string `json:"author"`
	CreatedBy int64  `json:"created_by"`
}

type UpdateQuoteReq struct {
	Text     string `json:"text" validate:"required"`
	Author   string `json:"author"`
	IsActive bool   `json:"is_active"`
}
