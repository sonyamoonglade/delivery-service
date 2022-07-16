package entity

import (
	"github.com/sonyamoonglade/delivery-service"
	"time"
)

type Delivery struct {
	DeliveryID int64          `json:"delivery_id,omitempty" db:"delivery_id"`
	OrderID    int64          `json:"order_id" db:"order_id"`
	RunnerID   int64          `json:"runner_id,omitempty" db:"runner_id"`
	ReservedAt time.Time      `json:"reserved_at,omitempty" db:"reserved_at"`
	CreatedAt  time.Time      `json:"created_at,omitempty" db:"created_at"`
	IsFree     bool           `json:"is_free" db:"is_free"`
	Pay        tgdelivery.Pay `json:"pay" db:"pay"`
}

type BaseDelivery struct {
	DeliveryID int64          `json:"delivery_id" db:"delivery_id"`
	OrderID    int64          `json:"order_id" db:"order_id"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	IsFree     bool           `json:"is_free" db:"is_free"`
	Pay        tgdelivery.Pay `json:"pay" db:"pay"`
}

type ReservedDelivery struct {
	DeliveryID int64     `json:"delivery_id" db:"delivery_id"`
	RunnerID   int64     `json:"runner_id" db:"runner_id"`
	ReservedAt time.Time `json:"reserved_at" db:"reserved_at"`
}

func CombineDelivery(b *BaseDelivery, r *ReservedDelivery) *Delivery {
	return &Delivery{
		DeliveryID: b.DeliveryID,
		OrderID:    b.OrderID,
		RunnerID:   r.RunnerID,
		ReservedAt: r.ReservedAt,
		CreatedAt:  b.CreatedAt,
		IsFree:     b.IsFree,
		Pay:        b.Pay,
	}
}
