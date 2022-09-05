package formatter

import (
	"fmt"
	"strings"
	"time"

	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/helpers"
	botDto "github.com/sonyamoonglade/delivery-service/pkg/telegram/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
)

type ExtractableMessageData struct {
	OrderID  int64
	Amount   int64
	Username string
}

type ExtractFormatter interface {
	FormatTemplate(dto *dto.CreateDelivery, offset int64) string
	ExtractDataFromText(text string) ExtractableMessageData
	PersonalAfterReserveReply(dto botDto.PersonalReserveReplyDto, offset int64) string
	GroupAfterReserveReply(dto botDto.GroupReserveReplyDto, offset int64) string
	AfterCompleteReply(dto botDto.PersonalCompleteReplyDto) string
}

type formatter struct {
	logger *zap.SugaredLogger
}

func NewFormatter(logger *zap.SugaredLogger) ExtractFormatter {
	return &formatter{logger: logger}
}

func (f *formatter) FormatTemplate(dto *dto.CreateDelivery, offset int64) string {
	f.logger.Infof("format a template. offset: %d, orderID: %d, originaltime: %v", offset, dto.Order.OrderID, dto.Order.DeliveryDetails.DeliveredAt)

	template := templates.DeliveryText

	payTranslate := helpers.PayTranslate(dto.Order.Pay)
	idLikeSix := helpers.SixifyOrderId(dto.Order.OrderID)

	usrMarkStr := "Метки пользователя: "
	var sortedByImportance []tgdelivery.Mark
	for _, m := range dto.User.Marks {
		if m.IsImportant {
			sortedByImportance = append([]tgdelivery.Mark{m}, sortedByImportance...)
			continue
		}
		sortedByImportance = append(sortedByImportance, m)
	}

	if len(dto.User.Marks) == 0 {
		usrMarkStr += " Отсутствуют \n\r"
	}

	for i, m := range sortedByImportance {
		if i == 0 {
			usrMarkStr += "\n\r"
		}

		usrMarkStr += fmt.Sprintf(" - %s \n\r", m.Content)
	}

	template = strings.Replace(template, "orderId", idLikeSix, -1)
	template = strings.Replace(template, "sum", fmt.Sprintf("%d", dto.Order.Amount), -1)
	template = strings.Replace(template, "pay", payTranslate, -1)
	template = strings.Replace(template, "username", dto.User.Username, -1)
	template = strings.Replace(template, "phoneNumber", dto.User.PhoneNumber, -1)
	template = strings.Replace(template, "marks", usrMarkStr, -1)
	template = strings.Replace(template, "address", dto.Order.DeliveryDetails.Address, -1)
	template = strings.Replace(template, "ent", fmt.Sprintf("%d", dto.Order.DeliveryDetails.EntranceNumber), -1)
	template = strings.Replace(template, "gr", fmt.Sprintf("%d", dto.Order.DeliveryDetails.Floor), -1)
	template = strings.Replace(template, "fl", fmt.Sprintf("%d", dto.Order.DeliveryDetails.FlatCall), -1)

	//Comment
	if dto.Order.DeliveryDetails.Comment != "" {
		template = strings.Replace(template, "comm", dto.Order.DeliveryDetails.Comment, -1)
	} else {
		template = strings.Replace(template, "Комментарий: comm\n", "", -1)
	}

	//Delivered at
	if dto.Order.IsDeliveredAsap == true {
		template = strings.Replace(template, "к time", templates.ASAP, -1)
	} else {
		template = strings.Replace(template, "time", f.fmtTime(dto.Order.DeliveryDetails.DeliveredAt, offset), -1)
	}
	return template
}

func (f *formatter) ExtractDataFromText(text string) ExtractableMessageData {

	//Split text by new line "\n"
	splByNewLine := strings.Split(text, "\n")
	//See templates.DeliveryText
	ordLine := splByNewLine[0]
	totalPriceLine := splByNewLine[2]
	userLine := splByNewLine[5]

	data := ExtractableMessageData{
		OrderID:  helpers.ExtractOrderId(ordLine),
		Amount:   helpers.ExtractAmount(totalPriceLine),
		Username: helpers.ExtractUsername(userLine),
	}

	return data

}

func (f *formatter) PersonalAfterReserveReply(dto botDto.PersonalReserveReplyDto, offset int64) string {
	return fmt.Sprintf(templates.PersonalAfterReserveText,
		f.fmtTime(dto.ReservedAt, offset),
		dto.DeliveryID)
}
func (f *formatter) GroupAfterReserveReply(dto botDto.GroupReserveReplyDto, offset int64) string {

	return fmt.Sprintf(templates.GroupAfterReserveText,
		dto.DeliveryID,
		templates.Success,
		f.fmtTime(dto.ReservedAt, offset),
		dto.RunnerUsername,
		dto.OrderID,
		dto.Username,
		dto.Amount)
}
func (f *formatter) AfterCompleteReply(dto botDto.PersonalCompleteReplyDto) string {
	return fmt.Sprintf(templates.DeliveryCompletedText,
		dto.DeliveryID,
		templates.Success,
		dto.OrderID,
		dto.Username,
		dto.Amount)
}

func (f *formatter) fmtTime(t time.Time, offset int64) string {
	//extra conversion to int is mandatory( time.Duration returns hell-shit with int64 )
	withOffset := t.Add(time.Hour * time.Duration(int(offset)))
	return withOffset.Format("15:04 02.01")
}
