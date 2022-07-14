package tgdelivery

import "time"

type Payload struct {
	Order *Order `json:"order"`
	User  *User  `json:"user"`
}

type Order struct {
	OrderID         int64            `json:"order_id"`
	DeliveryDetails *DeliveryDetails `json:"delivery_details"`
	TotalCartPrice  int64            `json:"total_cart_price"`
	Pay             Pay              `json:"pay"`
}

type DeliveryDetails struct {
	Address        string    `json:"address"`
	FlatCall       int64     `json:"flat_call"`
	EntranceNumber int64     `json:"entrance_number"`
	Floor          int64     `json:"floor"`
	DeliveredAt    time.Time `json:"delivered_at"`
}

type User struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Marks       []Mark `json:"marks,omitempty"`
}

type Mark struct {
	MarkID      int64     `json:"mark_id"`
	UserID      int64     `json:"user_id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	IsImportant bool      `json:"is_important"`
}

type Runner struct {
	RunnerID    int64
	Username    string
	PhoneNumber string
}

type Pay string

var (
	Paid     Pay = "paid"
	Cash     Pay = "cash"
	WithCard Pay = "withCard"
)

type BotConfig struct {
	Token   string
	Timeout int
	Debug   bool
}

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
