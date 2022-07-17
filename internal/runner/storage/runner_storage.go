package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
)

type runnerStorage struct {
	db *sqlx.DB
}

const runnerTable = "runner"

func (s *runnerStorage) IsRunner(dto dto.IsRunnerDto) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *runnerStorage) Register(dto dto.RegisterRunnerDto) (int64, error) {

	q := fmt.Sprintf("INSERT INTO %s (phone_number,username) values($1,$2) ON CONFLICT DO NOTHING RETURNING runner_id", runnerTable)
	row := s.db.QueryRowx(q, dto.PhoneNumber, dto.Username)

	var runnerID int64

	if err := row.Scan(&runnerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return runnerID, nil
}

func (s *runnerStorage) Ban(id int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewRunnerStorage(db *sqlx.DB) runner.Storage {
	return &runnerStorage{db: db}
}
