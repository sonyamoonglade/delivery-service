package service

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"go.uber.org/zap"
)

type DeliveryStorage interface {
	Create(d *Delivery) (bool, error)
	Reserve(id int64) (bool, error)
}

type Delivery interface {
	NewDeliveryMessage(v *tgdelivery.Payload) (*tg.Chattable, error)
	Reserve(id int64) (bool, error)
}

type deliveryService struct {
	logger  *zap.Logger
	storage *DeliveryStorage
}

func (d deliveryService) NewDeliveryMessage(v *tgdelivery.Payload) (*tg.Chattable, error) {
	//TODO implement me
	panic("implement me")
}

func (d deliveryService) Reserve(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewDeliveryService(logger *zap.Logger, storage DeliveryStorage) *deliveryService {
	return &deliveryService{logger: logger, storage: &storage}
}
