package service

import (
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"go.uber.org/zap"
)

type deliveryService struct {
	logger  *zap.Logger
	storage delivery.Storage
}

func NewDeliveryService(logger *zap.Logger, storage delivery.Storage) delivery.Delivery {
	return &deliveryService{logger: logger, storage: storage}
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

func (s *deliveryService) Delete(id int64) error {

	deliveryID, err := s.storage.Delete(id)
	if err != nil {
		return err
	}
	if deliveryID == 0 {
		// throw err here
	}

	return nil

}
