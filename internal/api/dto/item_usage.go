package dto

type ItemUsageCreateRequest struct {
	ItemID       int64   `json:"item_id" validate:"required"`
	QuantityUsed float32 `json:"quantity_used" validate:"required,gt=0"`
	Reason       string  `json:"reason" validate:"required,oneof=CONSUMED EXPIRED DISCARDED"`
}
