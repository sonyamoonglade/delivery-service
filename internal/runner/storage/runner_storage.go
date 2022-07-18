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

func NewRunnerStorage(db *sqlx.DB) runner.Storage {
	return &runnerStorage{db: db}
}

const (
	runnerTable         = "runner"
	telegramRunnerTable = "telegram_runner"
)

func (s *runnerStorage) GetByTelegramId(tgUsrID int64) (int64, error) {

	q := fmt.Sprintf("SELECT runner_id FROM %s WHERE telegram_id = $1", telegramRunnerTable)

	var runnerID int64
	row := s.db.QueryRowx(q, tgUsrID)

	if err := row.Scan(&runnerID); err != nil {
		return 0, err
	}
	return runnerID, nil

}

func (s *runnerStorage) IsKnownByTelegramId(tgUsrID int64) (bool, error) {
	q := fmt.Sprintf("SELECT true FROM %s WHERE telegram_id = $1", telegramRunnerTable)

	row := s.db.QueryRowx(q, tgUsrID)
	var ok bool
	if err := row.Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (s *runnerStorage) IsRunner(usrPhoneNumber string) (int64, error) {

	q := fmt.Sprintf("SELECT runner_id FROM %s WHERE phone_number = $1", runnerTable)
	row := s.db.QueryRowx(q, usrPhoneNumber)

	var runnerID int64

	if err := row.Scan(&runnerID); err != nil {
		return 0, err
	}

	return runnerID, nil

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

func (s *runnerStorage) Ban(runnerID int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *runnerStorage) BeginWork(dto dto.RunnerBeginWorkDto) error {

	q := fmt.Sprintf("INSERT INTO %s (runner_id, telegram_id) VALUES ($1,$2) ON CONFLICT DO NOTHING", telegramRunnerTable)
	_, err := s.db.Queryx(q, dto.RunnerID, dto.TelegramUserID)
	if err != nil {
		return err
	}

	return nil

}
