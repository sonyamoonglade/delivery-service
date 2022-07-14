package tgdelivery

import "time"

type Delivery struct {
	DeliveryID int64     `json:"delivery_id,omitempty"`
	OrderID    int64     `json:"order_id"`
	RunnerID   int64     `json:"runner_id,omitempty"`
	ReservedAt time.Time `json:"reserved_at,omitempty"`
	IsFree     bool      `json:"is_free"`
	Pay        *Pay      `json:"pay"`
}
