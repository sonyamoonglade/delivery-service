package delivery

import (
	"net/http"
	"time"

	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
)

type Service interface {
	Create(dto dto.CreateDeliveryDatabaseDto) (int64, error)
	Reserve(dto dto.ReserveDeliveryDto) (time.Time, error)
	Complete(deliveryID int64) (bool, error)
	Delete(deliveryID int64) error
	Status(dto dto.StatusOfDeliveryDto) ([]tgdelivery.DeliveryStatus, error)
	WriteCheck(dto dto.CheckDtoForCli) error
	CopyFromCheck(w http.ResponseWriter) error
}
