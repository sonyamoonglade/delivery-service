package dto

import "time"

type PersonalReserveReplyDto struct {
	DeliveryID int64
	ReservedAt time.Time
}

type GroupReserveReplyDto struct {
	DeliveryID     int64
	OrderID        string
	Username       string
	TotalCartPrice int64
	ReservedAt     time.Time
	RunnerUsername string
}

type PersonalCompleteReplyDto struct {
	DeliveryID     int64
	OrderID        string
	Username       string
	TotalCartPrice int64
}
