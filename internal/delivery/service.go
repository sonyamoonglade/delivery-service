package delivery

import (
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"net/http"
	"time"
)

type Service interface {
	Create(dto dto.CreateDeliveryDatabaseDto) (int64, error)
	Reserve(dto dto.ReserveDeliveryDto) (time.Time, error)
	Complete(deliveryID int64) (bool, error)
	Delete(deliveryID int64) error
	Status(dto dto.StatusOfDeliveryDto) ([]tgdelivery.DeliveryStatus, error)
	WriteCheck(dto dto.CheckDtoForCli) error
	ReadFromCheck(w http.ResponseWriter) error
}
