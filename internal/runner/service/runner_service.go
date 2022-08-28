package service

import (
	"context"

	"github.com/sonyamoonglade/delivery-service/internal/entity"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"go.uber.org/zap"
)

type runnerService struct {
	logger  *zap.SugaredLogger
	storage runner.Storage
}

func NewRunnerService(logger *zap.SugaredLogger, storage runner.Storage) runner.Service {
	return &runnerService{logger: logger, storage: storage}
}

func (s *runnerService) All(ctx context.Context) ([]*entity.Runner, error) {
	return s.storage.All(ctx)
}
func (s *runnerService) GetByTelegramId(tgUsrID int64) (*entity.Runner, error) {

	rn, err := s.storage.GetByTelegramId(tgUsrID)
	if err != nil {
		return nil, err
	}

	if rn == nil {
		return nil, tgErrors.RunnerDoesNotExistClean()
	}

	return rn, nil

}

func (s *runnerService) IsRunner(usrPhoneNumber string) (int64, error) {

	runnerID, err := s.storage.IsRunner(usrPhoneNumber)
	if err != nil {
		return 0, err
	}

	if runnerID == 0 {
		return 0, tgErrors.RunnerDoesNotExist(usrPhoneNumber)
	}

	return runnerID, nil
}

func (s *runnerService) IsKnownByTelegramId(tgUsrID int64) (bool, error) {
	ok, err := s.storage.IsKnownByTelegramId(tgUsrID)

	if err != nil {
		return false, err
	}

	//Can be false, nil OR true, nil
	return ok, nil

}

func (s *runnerService) Register(dto dto.RegisterRunnerDto) error {

	runnerID, err := s.storage.Register(dto)
	if err != nil {
		s.logger.Error(err.Error())
		return httpErrors.InternalError()
	}

	if runnerID == 0 {
		return httpErrors.ConflictError(httpErrors.RunnerAlreadyExists)
	}
	s.logger.Debug("registered runner successfully")
	return nil

}

func (s *runnerService) BeginWork(dto dto.RunnerBeginWorkDto) error {

	err := s.storage.BeginWork(dto)
	if err != nil {
		return err
	}

	return nil
}

func (s *runnerService) Ban(phoneNumber string) error {
	return s.storage.Ban(phoneNumber)
}
