package tgdelivery

import "time"

type Order struct {
	OrderID         int64            `json:"order_id,omitempty" validate:"required"`
	DeliveryDetails *DeliveryDetails `json:"delivery_details" validate:"required"`
	TotalCartPrice  int64            `json:"total_cart_price" validate:"required"`
	Pay             Pay              `json:"pay" validate:"required"`
	IsDeliveredAsap bool             `json:"is_delivered_asap"`
}

type DeliveryDetails struct {
	Address        string    `json:"address,omitempty"`
	FlatCall       int64     `json:"flat_call,omitempty"`
	EntranceNumber int64     `json:"entrance_number,omitempty"`
	Floor          int64     `json:"floor,omitempty"`
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

type CartProduct struct {
	ProductID int64  `json:"id,omitempty"`
	Name      string `json:"translate" validate:"required"`
	Price     int64  `json:"price" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
	Category  string `json:"category,omitempty"`
}

type OrderForCheck struct {
	OrderID         int64           `json:"order_id" validate:"required"`
	DeliveryDetails DeliveryDetails `json:"delivery_details,omitempty"`
	TotalCartPrice  int64           `json:"total_cart_price" validate:"required"`
	Pay             Pay             `json:"pay" validate:"required"`
	Cart            []CartProduct   `json:"cart" validate:"required"`
	IsDelivered     bool            `json:"is_delivered"`
}

type UserForCheck struct {
	Username    string `json:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type Pay string

var (
	Online   Pay = "online"
	OnPickup Pay = "onPickup"
)

type DeliveryStatus struct {
	OrderID int64 `json:"orderId"`
	Status  bool  `json:"status"`
}
