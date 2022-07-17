package telegram

import tgdelivery "github.com/sonyamoonglade/delivery-service"

type Service interface {
	Send(text string) error
	FromTemplate(p *tgdelivery.Payload) string
}
