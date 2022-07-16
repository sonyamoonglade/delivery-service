package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/sonyamoonglade/delivery-service/internal/handler/dto"
	"go.uber.org/zap"
)

type deliveryStorage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func (d2 deliveryStorage) Create(dto *dto.CreateDeliveryDto) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d2 deliveryStorage) Reserve(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewDeliveryStorage(logger *zap.Logger, db *sqlx.DB) *deliveryStorage {
	return &deliveryStorage{db: db, logger: logger}
}
