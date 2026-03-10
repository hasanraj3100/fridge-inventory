package dto

type PaginationRequest struct {
	Page     int `json:"page" validate:"omitempty,min=1"`
	PageSize int `json:"page_size" validate:"omitempty,min=1,max=100"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedItemUsageResponse struct {
	Data       []ItemUsageWithDetails `json:"data"`
	Pagination PaginationResponse     `json:"pagination"`
}

type ItemUsageWithDetails struct {
	ID           int64   `json:"id"`
	ItemID       int64   `json:"item_id"`
	ItemName     string  `json:"item_name"`
	QuantityUsed float32 `json:"quantity_used"`
	Unit         string  `json:"unit"`
	Reason       string  `json:"reason"`
	UsedAt       string  `json:"used_at"`
}
