package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (s *deliveryStorage) Status(ids []int64) ([]bool, error) {

	var statuses []bool

	for _, orderId := range ids {
		var ok bool
		q := fmt.Sprintf("SELECT true FROM %s WHERE order_id = $1", deliveryTable)
		row := s.db.QueryRowx(q, orderId)
		if err := row.Scan(&ok); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				//Append false because delivery does not exist
				statuses = append(statuses, false)
				continue
			}
			return nil, err
		}
		//If delivery on specific orderId does exist -> append true
		statuses = append(statuses, true)
	}

	return statuses, nil
}

func (s *deliveryStorage) Complete(deliveryID int64) (bool, error) {

	q := fmt.Sprintf("UPDATE %s SET is_completed = true WHERE delivery_id = $1 AND is_completed = false", deliveryTable)

	r, err := s.db.Exec(q, deliveryID)
	if err != nil {
		return false, err
	}

	if n, err := r.RowsAffected(); err != nil {
		if n == 0 {
			return false, nil
		}
		return false, err
	}

	return true, nil

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

func (s *deliveryStorage) Create(dto dto.CreateDeliveryDatabaseDto) (int64, error) {

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

func (s *deliveryStorage) Reserve(dto dto.ReserveDeliveryDto) (time.Time, error) {

	var free bool
	var reservedAt time.Time

	tx, err := s.db.BeginTxx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return time.Time{}, err
		}

		return time.Time{}, err
	}

	q := fmt.Sprintf("SELECT is_free FROM %s WHERE delivery_id = $1", deliveryTable)
	row := tx.QueryRowx(q, dto.DeliveryID)
	if err = row.Scan(&free); err != nil {
		if err := tx.Rollback(); err != nil {
			return time.Time{}, err
		}
		return time.Time{}, err
	}
	//Make sure delivery is 100% free to reserve. If not -> rollback
	if free == false {
		if err := tx.Rollback(); err != nil {
			return time.Time{}, err
		}
		return time.Time{}, err
	}

	q1 := fmt.Sprintf("UPDATE %s SET is_free = false WHERE delivery_id = $1", deliveryTable)
	_, err = tx.Exec(q1, dto.DeliveryID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return time.Time{}, err
		}
		return time.Time{}, err
	}
	q2 := fmt.Sprintf("INSERT INTO %s (delivery_id, runner_id) VALUES ($1,$2) RETURNING reserved_at", reservedTable)
	row = tx.QueryRowx(q2, dto.DeliveryID, dto.RunnerID)
	if err = row.Scan(&reservedAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return time.Time{}, err
		}
		return time.Time{}, err
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return time.Time{}, err
		}
		return time.Time{}, err
	}
	return reservedAt, nil
}
