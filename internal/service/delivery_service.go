package service

import (
	"github.com/sonyamoonglade/delivery-service/internal/handler/dto"
	"go.uber.org/zap"
)

type DeliveryStorage interface {
	Create(d *dto.CreateDeliveryDto) (int64, error)
	Reserve(id int64) (bool, error)
}

type Delivery interface {
	Create(dto *dto.CreateDeliveryDto) (int64, error)
	Reserve(id int64) (bool, error)
}

type deliveryService struct {
	logger  *zap.Logger
	storage DeliveryStorage
}

func (s *deliveryService) Create(dto *dto.CreateDeliveryDto) (int64, error) {

	deliveryID, err := s.storage.Create(dto)
	if err != nil {
		return 0, err
	}

	return deliveryID, nil
}

func (s *deliveryService) Reserve(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewDeliveryService(logger *zap.Logger, storage DeliveryStorage) *deliveryService {
	return &deliveryService{logger: logger, storage: storage}
}
