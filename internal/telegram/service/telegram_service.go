package service

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/callback"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/helpers"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
	"strings"
)

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
	msg := tg.NewMessage(bot.GetGroupChatId(), text)

	msg.ReplyMarkup = bot.ReserveDeliveryKeyboard(data)

	_, err := s.bot.Send(msg)
	if err != nil {
		s.logger.Error(err.Error())
		return tgErrors.InternalError
	}
	return nil
}

func (s *telegramService) FormatTemplate(p *dto.CreateDelivery) string {

	template := templates.DeliveryText

	payTranslate := helpers.PayTranslate(p.Order.Pay)
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

	//Comment
	if p.Order.DeliveryDetails.Comment != "" {
		template = strings.Replace(template, "comm", p.Order.DeliveryDetails.Comment, -1)
	} else {
		template = strings.Replace(template, "Комментарий: comm\n", "", -1)
	}

	//Delivered at
	if p.Order.IsDeliveredAsap == true {
		template = strings.Replace(template, "к time", templates.ASAP, -1)
	} else {
		template = strings.Replace(template, "time", p.Order.DeliveryDetails.DeliveredAt.Format("15:04 02.01"), -1)
	}

	return template
}

func (s *telegramService) ExtractDataFromText(text string) *bot.DataFromText {

	//Split text by new line "\n"
	splByNewLine := strings.Split(text, "\n")
	ordLine := splByNewLine[0]
	totalPriceLine := splByNewLine[2]
	userLine := splByNewLine[6]

	data := &bot.DataFromText{
		OrderID:        helpers.ExtractOrderId(ordLine),
		TotalCartPrice: helpers.ExtractTotalPrice(totalPriceLine),
		Username:       helpers.ExtractUsername(userLine),
	}

	return data

}
