package tgdelivery

import "time"

type Order struct {
	OrderID         int64            `json:"order_id,omitempty" validate:"required"`
	DeliveryDetails *DeliveryDetails `json:"delivery_details" validate:"required"`
	TotalCartPrice  int64            `json:"total_cart_price" validate:"required"`
	Pay             Pay              `json:"pay" validate:"required"`
	IsPaid          bool             `json:"is_paid"`
	IsDeliveredAsap bool             `json:"is_delivered_asap"`
}

type DeliveryDetails struct {
	Address        string    `json:"address" validate:"required"`
	FlatCall       int64     `json:"flat_call" validate:"required"`
	EntranceNumber int64     `json:"entrance_number" validate:"required"`
	Floor          int64     `json:"floor" validate:"required"`
	DeliveredAt    time.Time `json:"delivered_at,omitempty"`
	Comment        string    `json:"comment,omitempty"`
}

type User struct {
	UserID      int64  `json:"user_id" validate:"required"`
	Username    string `json:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Marks       []Mark `json:"marks,omitempty" validate:"required"`
}

type Mark struct {
	MarkID      int64     `json:"mark_id" validate:"required"`
	UserID      int64     `json:"user_id" validate:"required"`
	Content     string    `json:"content" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	IsImportant bool      `json:"is_important"`
}

type Pay string

var (
	Cash           Pay = "cash"
	WithCardRunner Pay = "withCardRunner"
	WithCard       Pay = "withCard"
)

var (
	BOT_TOKEN = "BOT_TOKEN"
)

type DeliveryStatus struct {
	OrderID int64 `json:"orderId"`
	Status  bool  `json:"status"`
}
