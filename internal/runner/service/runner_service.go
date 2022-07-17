package service

import (
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	"go.uber.org/zap"
)

type runnerService struct {
	logger  *zap.Logger
	storage runner.Storage
}

func (s *runnerService) IsRunner(dto dto.IsRunnerDto) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *runnerService) Register(dto dto.RegisterRunnerDto) error {

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

func (s *runnerService) Ban(id int64) error {
	//TODO implement me
	panic("implement me")
}

func NewRunnerService(logger *zap.Logger, storage runner.Storage) runner.Service {
	return &runnerService{logger: logger, storage: storage}
}
