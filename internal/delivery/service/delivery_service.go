package service

import (
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	"go.uber.org/zap"
)

type deliveryService struct {
	logger  *zap.Logger
	storage delivery.Storage
}

func NewDeliveryService(logger *zap.Logger, storage delivery.Storage) delivery.Service {
	return &deliveryService{logger: logger, storage: storage}
}

func (s *deliveryService) Create(dto *dto.CreateDeliveryDto) (int64, error) {

	deliveryID, err := s.storage.Create(dto)

	// Delivery already exists
	if deliveryID == 0 {
		return 0, httpErrors.ConflictError(httpErrors.DeliveryAlreadyExists)
	}
	if err != nil {
		return 0, httpErrors.InternalError()
	}

	return deliveryID, nil
}

func (s *deliveryService) Reserve(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *deliveryService) Delete(id int64) error {

	ok, err := s.storage.Delete(id)

	if err != nil {
		return httpErrors.InternalError()
	}
	// Delivery does not exist
	if !ok {
		return httpErrors.NotFoundError(httpErrors.DeliveryDoesNotExist)
	}

	return nil

}
