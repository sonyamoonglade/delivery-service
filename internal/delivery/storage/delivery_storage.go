package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"go.uber.org/zap"
)

type deliveryStorage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewDeliveryStorage(logger *zap.Logger, db *sqlx.DB) delivery.Storage {
	return &deliveryStorage{db: db, logger: logger}
}

func (d2 deliveryStorage) Delete(id int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d2 deliveryStorage) Create(dto *dto.CreateDeliveryDto) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d2 deliveryStorage) Reserve(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}
