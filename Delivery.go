package tgdelivery

import "time"

type Delivery struct {
	DeliveryID int64     `json:"delivery_id,omitempty" db:"delivery_id"`
	OrderID    int64     `json:"order_id" db:"order_id"`
	RunnerID   int64     `json:"runner_id,omitempty" db:"runner_id"`
	ReservedAt time.Time `json:"reserved_at,omitempty" db:"reserved_at"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"created_at"`
	IsFree     bool      `json:"is_free" db:"is_free"`
	Pay        Pay       `json:"pay" db:"pay"`
}
