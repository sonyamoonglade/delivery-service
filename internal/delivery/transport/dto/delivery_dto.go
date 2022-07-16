package dto

import (
	tgdelivery "github.com/sonyamoonglade/delivery-service"
)

type CreateDeliveryDto struct {
	OrderID int64          `json:"order_id"`
	Pay     tgdelivery.Pay `json:"pay"`
}
