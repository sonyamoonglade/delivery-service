package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
)

type deliveryStorage struct {
	db *sqlx.DB
}

var (
	deliveryTable = "delivery"
	reservedTable = "reserved"
)

func NewDeliveryStorage(db *sqlx.DB) delivery.Storage {
	return &deliveryStorage{db: db}
}

func (s *deliveryStorage) Delete(id int64) (bool, error) {

	q := fmt.Sprintf("DELETE FROM %s WHERE delivery_id = $1 RETURNING delivery_id", deliveryTable)
	row := s.db.QueryRowx(q, id)

	var deletedID int64
	if err := row.Scan(&deletedID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *deliveryStorage) Create(dto *dto.CreateDeliveryDto) (int64, error) {

	q := fmt.Sprintf("INSERT INTO %s (order_id, pay) VALUES ($1,$2) ON CONFLICT DO NOTHING RETURNING delivery_id", deliveryTable)
	row := s.db.QueryRowx(q, dto.OrderID, dto.Pay)

	var deliveryID int64

	if err := row.Scan(&deliveryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return deliveryID, nil
}

func (s *deliveryStorage) Reserve(dto dto.ReserveDeliveryDto) (bool, error) {

	var free bool

	tx, err := s.db.BeginTxx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		return false, err
	}

	q := fmt.Sprintf("SELECT is_free FROM %s WHERE delivery_id = $1", deliveryTable)
	row := tx.QueryRowx(q, dto.DeliveryID)
	if err = row.Scan(&free); err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		return false, err
	}
	//Make sure delivery is 100% free to reserve. If not -> rollback
	if free == false {
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		return false, nil
	}

	q1 := fmt.Sprintf("UPDATE %s SET is_free = true THEN FALSE ELSE DO NOTHING WHERE delivery_id = $1", deliveryTable)
	_, err = tx.Exec(q1, dto.DeliveryID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		return false, err
	}
	q2 := fmt.Sprintf("INSERT INTO %s (delivery_id, runner_id) VALUES ($1,$2)", reservedTable)
	_, err = tx.Exec(q2, dto.DeliveryID, dto.RunnerID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		return false, err
	}

	if err = tx.Commit(); err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		return false, err
	}
	return true, nil
}
