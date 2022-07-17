package delivery

import "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"

type Service interface {
	Create(dto *dto.CreateDeliveryDto) (int64, error)
	Reserve(id int64) (bool, error)
	Delete(id int64) error
}
