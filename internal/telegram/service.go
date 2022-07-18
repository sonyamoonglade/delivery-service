package telegram

import tgdelivery "github.com/sonyamoonglade/delivery-service"

type Service interface {
	Send(text string, deliveryID int64) error
	FromTemplate(p *tgdelivery.Payload) string
	GetGroupChatId() int64
}
