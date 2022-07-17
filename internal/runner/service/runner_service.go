package service

import (
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"go.uber.org/zap"
)

type runnerService struct {
	logger  *zap.Logger
	storage runner.Storage
}

func NewRunnerService(logger *zap.Logger, storage runner.Storage) runner.Service {
	return &runnerService{logger: logger, storage: storage}
}

func (s *runnerService) IsRunner() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *runnerService) Register() error {
	//TODO implement me
	panic("implement me")
}

func (s *runnerService) Ban() error {
	//TODO implement me
	panic("implement me")
}
