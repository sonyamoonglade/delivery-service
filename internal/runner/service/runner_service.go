package service

import (
	"database/sql"
	"errors"
	"fmt"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"go.uber.org/zap"
)

type runnerService struct {
	logger  *zap.Logger
	storage runner.Storage
}

func (s *runnerService) IsRunner(usrPhoneNumber string) (int64, error) {

	runnerID, err := s.storage.IsRunner(usrPhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//todo: custom tg error
			return 0, errors.New(tgErrors.RunnerDoesNotExist)
		}
		return 0, err
	}
	//todo: throw custom tg_error here!
	return runnerID, nil
}

func (s *runnerService) Register(dto dto.RegisterRunnerDto) error {

	var valRes bool
	fmt.Println(tgdelivery.ValidateUsername(dto.Username))
	if valRes = tgdelivery.ValidateUsername(dto.Username); !valRes {

		return httpErrors.BadRequestError(httpErrors.InvalidUsername)
	}
	if valRes = tgdelivery.ValidatePhoneNumber(dto.PhoneNumber); !valRes {
		return httpErrors.BadRequestError(httpErrors.InvalidPhoneNumber)
	}

	runnerID, err := s.storage.Register(dto)
	if err != nil {
		return httpErrors.InternalError()
	}

	if runnerID == 0 {
		return httpErrors.ConflictError(httpErrors.RunnerAlreadyExists)
	}
	s.logger.Debug("registered user")
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

func NewRunnerService(logger *zap.Logger, storage runner.Storage) runner.Service {
	return &runnerService{logger: logger, storage: storage}
}
