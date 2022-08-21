package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sonyamoonglade/delivery-service/internal/entity"
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

func (s *runnerStorage) All(ctx context.Context) ([]*entity.Runner, error) {

	q := fmt.Sprintf("SELECT phone_number, username FROM %s", runnerTable)
	var runners []*entity.Runner
	rows, err := s.db.QueryxContext(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*entity.Runner{}, nil
		}
		return nil, err
	}

	for rows.Next() {
		var r entity.Runner
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		runners = append(runners, &r)
	}
	return runners, nil
}

func (s *runnerStorage) GetByTelegramId(tgUsrID int64) (*entity.Runner, error) {

	q := fmt.Sprintf("SELECT rn.username, rn.runner_id, rn.phone_number FROM %s rn JOIN %s tgrn ON rn.runner_id = tgrn.runner_id WHERE tgrn.telegram_id = $1", runnerTable, telegramRunnerTable)

	var r entity.Runner
	row := s.db.QueryRowx(q, tgUsrID)

	if err := row.StructScan(&r); err != nil {
		return nil, err
	}
	return &r, nil

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

func (s *runnerStorage) Ban(phoneNumber string) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE phone_number = $1", runnerTable)
	_, err := s.db.Exec(q, phoneNumber)
	if err != nil {
		return err
	}

	return nil

}

func (s *runnerStorage) BeginWork(dto dto.RunnerBeginWorkDto) error {

	q := fmt.Sprintf("INSERT INTO %s (runner_id, telegram_id) VALUES ($1,$2) ON CONFLICT DO NOTHING", telegramRunnerTable)
	_, err := s.db.Queryx(q, dto.RunnerID, dto.TelegramUserID)
	if err != nil {
		return err
	}

	return nil

}
