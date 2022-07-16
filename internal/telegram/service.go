package telegram

import tgdelivery "github.com/sonyamoonglade/delivery-service"

type Telegram interface {
	Send(text string) error
	FromTemplate(p *tgdelivery.Payload) string
}
