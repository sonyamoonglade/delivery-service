package service

import (
	"database/sql"
	"errors"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	"github.com/sonyamoonglade/delivery-service/pkg/cli"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type deliveryService struct {
	logger       *zap.SugaredLogger
	storage      delivery.Storage
	cli          cli.Cli
	checkService check.Service
}

func NewDeliveryService(logger *zap.SugaredLogger, storage delivery.Storage, cli cli.Cli, check check.Service) delivery.Service {
	return &deliveryService{logger: logger, storage: storage, cli: cli, checkService: check}
}

func (s *deliveryService) ReadFromCheck(w http.ResponseWriter) error {
	return s.checkService.Copy(w)
}

func (s *deliveryService) WriteCheck(dto dto.CheckDtoForCli) error {

	//Two Iterations to make sure key has restored and used
	for i := 0; i < 2; i++ {
		s.logger.Debugf("attempt #%d to write check", i+1)
		//Write .docx check file
		err := s.cli.WriteCheck(dto)
		if err != nil {
			if errors.Is(err, check.ApiKeyHasExpired) {
				//Refresh key here
				if err := s.checkService.RestoreKey(); err != nil {
					//Some internal error

					return cli.TimeoutError
				}
				//Restore is successful
				s.logger.Debug("restored key successfully")
				continue
			}
			if errors.Is(err, cli.TimeoutError) {

				return cli.TimeoutError
			}

			//Some internal error
			return err
		}

		//Break if it wrote check at 1st attempt
		break
	}
	return nil

}

func (s *deliveryService) Status(dto dto.StatusOfDeliveryDto) ([]tgdelivery.DeliveryStatus, error) {

	var statuses []tgdelivery.DeliveryStatus

	bools, err := s.storage.Status(dto.OrderIDs)
	if err != nil {
		return nil, httpErrors.InternalError()
	}

	for i, v := range bools {
		//Length of dto.OrderIDs will be always the same as statuses.
		orderId := dto.OrderIDs[i]
		status := tgdelivery.DeliveryStatus{
			OrderID: orderId,
			Status:  v,
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (s *deliveryService) Complete(deliveryID int64) (bool, error) {

	err := s.storage.Complete(deliveryID)
	if err != nil {
		s.logger.Error(err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return false, tgErrors.DeliveryCouldNotBeCompleted(deliveryID)
		}
		return false, err
	}
	return true, nil
}

func (s *deliveryService) Create(dto dto.CreateDeliveryDatabaseDto) (int64, error) {

	deliveryID, err := s.storage.Create(dto)

	// Delivery already exists
	if err != nil {
		s.logger.Error(err.Error())
		return 0, httpErrors.InternalError()
	}

	if deliveryID == 0 {
		return 0, httpErrors.ConflictError(httpErrors.DeliveryAlreadyExists)
	}
	return deliveryID, nil
}

func (s *deliveryService) Reserve(dto dto.ReserveDeliveryDto) (time.Time, error) {
	reservedAt, err := s.storage.Reserve(dto)

	if err != nil {
		s.logger.Error(err.Error())
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
		s.logger.Error(err.Error())
		return httpErrors.InternalError()
	}
	// Delivery does not exist
	if !ok {
		return httpErrors.NotFoundError(httpErrors.DeliveryDoesNotExist)
	}

	return nil
}
