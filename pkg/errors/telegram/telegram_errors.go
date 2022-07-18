package tgErrors

import (
	"errors"
	"fmt"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
)

var (
	InternalError = errors.New("internal telegram error")
)

type TelegramError struct {
	err string
}

func (e TelegramError) Error() string {
	return e.err
}

func NewTelegramError(e string) TelegramError {
	return TelegramError{
		err: e,
	}
}

func RunnerDoesNotExist(phoneNumber string) TelegramError {
	return NewTelegramError(fmt.Sprintf(templates.RunnerDoesNotExist, phoneNumber))
}

func DeliveryHasAlreadyReserved(deliveryID int64) TelegramError {
	return NewTelegramError(fmt.Sprintf(templates.DeliveryHasAlreadyReserved, deliveryID))
}
