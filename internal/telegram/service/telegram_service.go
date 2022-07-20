package service

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/callback"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/helpers"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
	"strings"
)

const chatId int64 = -693829286

type telegramService struct {
	bot    *tg.BotAPI
	logger *zap.SugaredLogger
}

func NewTelegramService(logger *zap.SugaredLogger, bot *tg.BotAPI) telegram.Service {
	return &telegramService{bot: bot, logger: logger}
}

func (s *telegramService) Send(text string, deliveryID int64) error {
	data := callback.Data{
		DeliveryID: deliveryID,
	}
	msg := tg.NewMessage(chatId, text)

	msg.ReplyMarkup = bot.ReserveDeliveryKeyboard(data)

	_, err := s.bot.Send(msg)
	if err != nil {
		s.logger.Error(err.Error())
		return tgErrors.InternalError
	}

	return nil
}

func (s *telegramService) FromTemplate(p *tgdelivery.Payload) string {
	template := templates.DeliveryText

	var payTranslate string

	switch p.Order.Pay {
	case tgdelivery.Cash:
		payTranslate = "Наличные"
	case tgdelivery.Paid:
		payTranslate = "Оплачен"
	case tgdelivery.WithCard:
		payTranslate = "Банковская карта"
	}

	idLikeSix := helpers.SixifyOrderId(p.Order.OrderID)

	usrMarkStr := "Метки пользователя: "
	var sortedByImportance []tgdelivery.Mark
	for _, m := range p.User.Marks {
		if m.IsImportant {
			sortedByImportance = append([]tgdelivery.Mark{m}, sortedByImportance...)
			continue
		}
		sortedByImportance = append(sortedByImportance, m)
	}

	if len(p.User.Marks) == 0 {
		usrMarkStr += " Отсутствуют \n\r"
	}

	for i, m := range sortedByImportance {
		if i == 0 {
			usrMarkStr += "\n\r"
		}

		usrMarkStr += fmt.Sprintf(" - %s \n\r", m.Content)
	}

	template = strings.Replace(template, "orderId", idLikeSix, -1)
	template = strings.Replace(template, "sum", fmt.Sprintf("%d", p.Order.TotalCartPrice), -1)
	template = strings.Replace(template, "pay", payTranslate, -1)
	template = strings.Replace(template, "username", p.User.Username, -1)
	template = strings.Replace(template, "phoneNumber", p.User.PhoneNumber, -1)
	template = strings.Replace(template, "marks", usrMarkStr, -1)
	template = strings.Replace(template, "address", p.Order.DeliveryDetails.Address, -1)
	template = strings.Replace(template, "ent", fmt.Sprintf("%d", p.Order.DeliveryDetails.EntranceNumber), -1)
	template = strings.Replace(template, "gr", fmt.Sprintf("%d", p.Order.DeliveryDetails.Floor), -1)
	template = strings.Replace(template, "fl", fmt.Sprintf("%d", p.Order.DeliveryDetails.FlatCall), -1)
	template = strings.Replace(template, "time", p.Order.DeliveryDetails.DeliveredAt.Format("15:04 02.01"), -1)

	return template
}

func (s *telegramService) ExtractDataFromText(text string) *bot.DataFromText {

	var data bot.DataFromText

	fmt.Println(text, data)
	return &data

}

func (s *telegramService) GetGroupChatId() int64 {
	return chatId
}
