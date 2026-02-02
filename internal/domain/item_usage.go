package domain

import "time"

type ItemUsageReason string

const (
	ReasonConsumed  ItemUsageReason = "CONSUMED"
	ReasonExpired   ItemUsageReason = "EXPIRED"
	ReasonDiscarded ItemUsageReason = "DISCARDED"
)

type ItemUsage struct {
	ID           int64           `json:"id" db:"id"`
	ItemID       int64           `json:"item_id" db:"item_id"`
	QuantityUsed float32         `json:"quantity_used" db:"quantity_used"`
	Reason       ItemUsageReason `json:"reason" db:"reason"`
	UsedAt       time.Time       `json:"used_at" db:"used_at"`
}
