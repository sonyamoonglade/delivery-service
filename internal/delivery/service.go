package delivery

import "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"

type Service interface {
	Create(dto *dto.CreateDeliveryDto) (int64, error)
	Reserve(dto dto.ReserveDeliveryDto) (bool, error)
	Delete(deliveryID int64) error
}
