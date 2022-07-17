package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
)

type runnerStorage struct {
	db *sqlx.DB
}

func NewRunnerStorage(db *sqlx.DB) runner.Storage {
	return &runnerStorage{db: db}
}

func (s *runnerStorage) IsRunner() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *runnerStorage) Register() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *runnerStorage) Ban() (int64, error) {
	//TODO implement me
	panic("implement me")
}
