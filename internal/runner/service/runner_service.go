package service

import (
	"database/sql"
	"errors"
	"github.com/sonyamoonglade/delivery-service/internal/entity"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/validation"
	"go.uber.org/zap"
)

type runnerService struct {
	logger  *zap.SugaredLogger
	storage runner.Storage
}

func NewRunnerService(logger *zap.SugaredLogger, storage runner.Storage) runner.Service {
	return &runnerService{logger: logger, storage: storage}
}

func (s *runnerService) GetByTelegramId(tgUsrID int64) (*entity.Runner, error) {

	runnerID, err := s.storage.GetByTelegramId(tgUsrID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tgErrors.RunnerDoesNotExistClean()
		}
		return nil, err
	}
	return runnerID, nil

}

func (s *runnerService) IsRunner(usrPhoneNumber string) (int64, error) {

	runnerID, err := s.storage.IsRunner(usrPhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, tgErrors.RunnerDoesNotExist(usrPhoneNumber)
		}
		return 0, err
	}
	//todo: throw custom tg_error here!
	return runnerID, nil
}

func (s *runnerService) IsKnownByTelegramId(tgUsrID int64) (bool, error) {
	ok, err := s.storage.IsKnownByTelegramId(tgUsrID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return ok, nil

}

func (s *runnerService) Register(dto dto.RegisterRunnerDto) error {

	var valRes bool
	if valRes = validation.ValidateUsername(dto.Username); !valRes {
		return httpErrors.BadRequestError(httpErrors.InvalidUsername)
	}
	if valRes = validation.ValidatePhoneNumber(dto.PhoneNumber); !valRes {
		return httpErrors.BadRequestError(httpErrors.InvalidPhoneNumber)
	}

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

func (s *runnerService) Ban(runnerID int64) error {
	//TODO implement me
	panic("implement me")
}
