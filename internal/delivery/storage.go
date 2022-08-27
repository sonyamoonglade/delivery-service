package delivery

import (
	"time"

	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
)

type Storage interface {
	Create(d dto.CreateDeliveryDatabaseDto) (int64, error)
	Delete(id int64) (bool, error)
	Reserve(dto dto.ReserveDeliveryDto) (time.Time, error)
	Complete(deliveryID int64) (bool, error)
	Status(ids []int64) ([]bool, error)
}
