package response

type SuccessResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func NewSuccessResponse(code int, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Code:    code,
		Data:    data,
	}
}
