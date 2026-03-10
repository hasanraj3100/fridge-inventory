package dto

type ItemUsageUpdateRequest struct {
	QuantityUsed float32 `json:"quantity_used" validate:"required,gt=0"`
	Reason       string  `json:"reason" validate:"required,oneof=CONSUMED EXPIRED DISCARDED"`
}
