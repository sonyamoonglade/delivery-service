package dto

import (
	tgdelivery "github.com/sonyamoonglade/delivery-service"
)

type CreateDelivery struct {
	Order *tgdelivery.Order `json:"order" validate:"required"`
	User  *tgdelivery.User  `json:"user" validate:"required"`
}

type CreateDeliveryDatabaseDto struct {
	OrderID int64          `json:"order_id"`
	Pay     tgdelivery.Pay `json:"pay"`
}

type ReserveDeliveryDto struct {
	RunnerID   int64
	DeliveryID int64
}

type StatusOfDeliveryDto struct {
	OrderIDs []int64 `json:"order_ids" validate:"required"`
}

type CheckDto struct {
	Order tgdelivery.OrderForCheck `json:"order" validate:"required"`
	User  tgdelivery.UserForCheck  `json:"user" validate:"required"`
}

type CheckDtoForCli struct {
	Data CheckDto `json:"data"`
}
