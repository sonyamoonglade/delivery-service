package telegram

import (
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
)

type Service interface {
	Send(text string, deliveryID int64) error
	FormatTemplate(dto *dto.CreateDelivery) string
	ExtractDataFromText(text string) *bot.DataFromText
}
