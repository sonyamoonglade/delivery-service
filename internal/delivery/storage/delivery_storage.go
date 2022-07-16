package storage

import (
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

func (s *deliveryStorage) Delete(id int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *deliveryStorage) Create(dto *dto.CreateDeliveryDto) (int64, error) {

	sql := fmt.Sprintf("INSERT INTO %s (order_id, pay) VALUES ($1,$2) ON CONFLICT DO NOTHING RETURNING delivery_id", deliveryTable)
	row := s.db.QueryRowx(sql, dto.OrderID, dto.Pay)

	out := make(map[string]interface{})
	cols, _ := row.Columns()
	key := cols[0]
	if err := row.MapScan(out); err != nil {
		if len(out) == 0 {
			return 0, nil
		}
		return 0, err
	}

	return out[key].(int64), nil
}

func (s *deliveryStorage) Reserve(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}
