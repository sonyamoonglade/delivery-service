package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"go.uber.org/zap"
)

type deliveryStorage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

var (
	deliveryTable = "delivery"
	reservedTable = "reserved"
)

func NewDeliveryStorage(logger *zap.Logger, db *sqlx.DB) delivery.Storage {
	return &deliveryStorage{db: db, logger: logger}
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

func (s *deliveryStorage) Reserve(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}
