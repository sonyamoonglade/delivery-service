package tgdelivery

import "time"

type Payload struct {
	Order *Order `json:"order" validate:"required"`
	User  *User  `json:"user" validate:"required"`
}

type Order struct {
	OrderID         int64            `json:"order_id,omitempty" validate:"required"`
	DeliveryDetails *DeliveryDetails `json:"delivery_details" validate:"required"`
	TotalCartPrice  int64            `json:"total_cart_price" validate:"required"`
	Pay             Pay              `json:"pay" validate:"required"`
}

type DeliveryDetails struct {
	Address        string    `json:"address" validate:"required"`
	FlatCall       int64     `json:"flat_call" validate:"required"`
	EntranceNumber int64     `json:"entrance_number" validate:"required"`
	Floor          int64     `json:"floor" validate:"required"`
	DeliveredAt    time.Time `json:"delivered_at" validate:"required"`
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
	IsImportant bool      `json:"is_important" validate:"required"`
}

type Pay string

var (
	Paid     Pay = "paid"
	Cash     Pay = "cash"
	WithCard Pay = "withCard"
)

var (
	BOT_TOKEN = "BOT_TOKEN"
)

const MessageTemplate = "" +
	"Заказ #orderId\n\r" +
	"\n\r" +
	"Сумма | sum.0 ₽\n\r" +
	"Способ оплаты | pay\n\r" +
	"\n\r" +
	"Заказчик: username\n\r" +
	"Номер телефона  phoneNumber\n\r" +
	"marks" +
	"\n\r" +
	"Доставка: да\n\r" +
	"Адрес: ул. address\n\r" +
	"Подъезд ent | Этаж gr | Квартира fl\n\r" +
	"\n\r" +
	"Время доставки: к time"
