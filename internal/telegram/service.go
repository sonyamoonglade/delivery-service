package telegram

import (
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
)

type Service interface {
	Send(text string, deliveryID int64) error
	FromTemplate(dto *dto.CreateDelivery) string
	GetGroupChatId() int64
	ExtractDataFromText(text string) *bot.DataFromText
}
