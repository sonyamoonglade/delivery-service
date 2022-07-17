package delivery

import "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"

type Storage interface {
	Create(d *dto.CreateDeliveryDto) (int64, error)
	Delete(id int64) (bool, error)
	Reserve(id int64) (bool, error)
}
