package telegram

import (
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
)

type Service interface {
	Send(text string, deliveryID int64) error
	FromTemplate(p *tgdelivery.Payload) string
	GetGroupChatId() int64
	ExtractDataFromText(text string) *bot.DataFromText
}
