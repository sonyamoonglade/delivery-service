package delivery

import (
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"time"
)

type Service interface {
	Create(dto *dto.CreateDeliveryDto) (int64, error)
	Reserve(dto dto.ReserveDeliveryDto) (time.Time, error)
	Delete(deliveryID int64) error
}
