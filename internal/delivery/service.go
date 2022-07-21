package delivery

import (
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"time"
)

type Service interface {
	Create(dto *dto.CreateDeliveryDatabaseDto) (int64, error)
	Reserve(dto dto.ReserveDeliveryDto) (time.Time, error)
	Complete(deliveryID int64) (bool, error)
	Delete(deliveryID int64) error
}
