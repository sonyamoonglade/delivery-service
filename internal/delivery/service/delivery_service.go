package service

import (
	"database/sql"
	"errors"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"go.uber.org/zap"
	"time"
)

type deliveryService struct {
	logger  *zap.SugaredLogger
	storage delivery.Storage
}

func NewDeliveryService(logger *zap.SugaredLogger, storage delivery.Storage) delivery.Service {
	return &deliveryService{logger: logger, storage: storage}
}

func (s *deliveryService) Complete(deliveryID int64) (bool, error) {

	err := s.storage.Complete(deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, tgErrors.DeliveryCouldNotBeCompleted(deliveryID)
		}
		return false, err
	}
	return true, nil
}

func (s *deliveryService) Create(dto *dto.CreateDeliveryDatabaseDto) (int64, error) {

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

func (s *deliveryService) Reserve(dto dto.ReserveDeliveryDto) (time.Time, error) {
	reservedAt, err := s.storage.Reserve(dto)
	if err != nil {
		return time.Time{}, err
	}

	//Signals that delivery has already reserved (see storage reserve impl.)
	if reservedAt.IsZero() == true {
		return time.Time{}, tgErrors.DeliveryHasAlreadyReserved(dto.DeliveryID)
	}

	return reservedAt, nil

}

func (s *deliveryService) Delete(deliveryID int64) error {

	ok, err := s.storage.Delete(deliveryID)

	if err != nil {
		return httpErrors.InternalError()
	}
	// Delivery does not exist
	if !ok {
		return httpErrors.NotFoundError(httpErrors.DeliveryDoesNotExist)
	}

	return nil

}
