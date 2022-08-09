package telegram

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/sonyamoonglade/delivery-service/config"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	dlvDto "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/callback"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/formatter"
	"github.com/sonyamoonglade/delivery-service/pkg/helpers"
	tgDto "github.com/sonyamoonglade/delivery-service/pkg/telegram/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

type Transport interface {
	ListenForUpdates()
}

type telegramTransport struct {
	logger          *zap.SugaredLogger
	bot             bot.Bot
	runnerService   runner.Service
	deliveryService delivery.Service
	extractFmt      formatter.ExtractFormatter
}

func NewTelegramTransport(logger *zap.SugaredLogger, bot bot.Bot, run runner.Service, del delivery.Service, extractFormatter formatter.ExtractFormatter) Transport {
	return &telegramTransport{logger: logger, bot: bot, runnerService: run, deliveryService: del, extractFmt: extractFormatter}
}

func (h *telegramTransport) ListenForUpdates() {

	client := h.bot.GetTelegramClient()
	cfg := h.bot.GetUpdatesConfig()

	updates := client.GetUpdatesChan(cfg)

	for u := range updates {
		h.handle(&u)
	}
}

func (h *telegramTransport) handle(u *tg.Update) {

	switch {
	case u.Message != nil && u.Message.Contact != nil:
		h.HandleContact(u.Message)
		return
	case u.Message != nil:
		h.HandleMessage(u.Message)
		return
	case u.CallbackQuery != nil:
		h.HandleCallback(u.CallbackQuery)
		return
	default:
		return
	}

}

func (h *telegramTransport) HandleCallback(cb *tg.CallbackQuery) {

	rcvFromGroup := cb.Message.Chat.IsGroup()
	grpChatID := h.bot.GetGroupChatId()
	//usrID can be used as chatID for private messages
	usrID := cb.From.ID
	msgID := cb.Message.MessageID
	callbackID := cb.ID
	callbackData := cb.Data
	inititalMsg := cb.Message.Text

	h.logger.Debugf("received a callback data=%s from '%d'", callbackData, usrID)

	txtData := h.extractFmt.ExtractDataFromText(inititalMsg)
	h.logger.Debugf("Extracted text-data %v", txtData)

	switch rcvFromGroup {

	case true:
		h.logger.Debug("Recieved group callback")
		var inp callback.Data

		if err := callback.Decode(callbackData, &inp); err != nil {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			return
		}
		h.logger.Debugf("Recieved data %v", inp)

		rn, err := h.runnerService.GetByTelegramId(usrID)
		if err != nil {
			h.ResponseWithError(err, usrID)
			h.AnswerCallback(callbackID)
			return
		}
		h.logger.Debugf("Got runner by telegram id %v", rn)

		rsrvDto := dlvDto.ReserveDeliveryDto{
			RunnerID:   rn.RunnerID,
			DeliveryID: inp.DeliveryID,
		}

		reservedAt, err := h.deliveryService.Reserve(rsrvDto)
		if err != nil {
			h.ResponseWithError(err, usrID)
			h.AnswerCallback(callbackID)
			return
		}
		h.logger.Debug("Reserved a delivery %d", inp.DeliveryID)

		h.logger.Infof("delivery ID=%d is reserved by runner ID=%d", inp.DeliveryID, rn.RunnerID)

		//Get initial message text

		//Apply extra parameters for runner's comfort
		personalData := tgDto.PersonalReserveReplyDto{
			DeliveryID: inp.DeliveryID,
			ReservedAt: reservedAt,
		}
		personalMsg := inititalMsg + h.extractFmt.PersonalAfterReserveReply(personalData, config.TempOffset)

		//TODO: table to save chat's locales to print time
		//Send to runner's private messages so he follows the delivery
		msg := tg.NewMessage(usrID, personalMsg)
		msg.ReplyMarkup = h.bot.CompleteDeliveryKeyboard(inp)
		h.Send(msg)
		h.logger.Debug("Sent private after-reserve reply")
		//Extracted data from initial delivery message

		//Prepare new group message (edit initial)
		groupData := tgDto.GroupReserveReplyDto{
			DeliveryID:     inp.DeliveryID,
			ReservedAt:     reservedAt,
			OrderID:        helpers.SixifyOrderId(txtData.OrderID),
			Username:       txtData.Username,
			TotalCartPrice: txtData.TotalCartPrice,
			RunnerUsername: rn.Username,
		}
		//Sent short and informative message to delivery group
		editMsg := tg.NewEditMessageText(grpChatID, msgID, h.extractFmt.GroupAfterReserveReply(groupData, config.TempOffset))
		h.Send(editMsg)
		h.logger.Debug("Sent group after-reserve reply")

		h.AnswerCallback(callbackID)
		h.logger.Debugf("Answered callback %s successfully", callbackID)
		return
	case false:

		var inp callback.Data

		if err := callback.Decode(callbackData, &inp); err != nil {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			return
		}

		//Complete the delivery
		ok, err := h.deliveryService.Complete(inp.DeliveryID)
		if err != nil {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			return
		}
		//Delivery could not be completed
		if ok == false {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			return
		}
		//Delivery is completed
		h.logger.Debugf("Completed delivery %d", inp.DeliveryID)

		data := tgDto.PersonalCompleteReplyDto{
			DeliveryID:     inp.DeliveryID,
			OrderID:        helpers.SixifyOrderId(txtData.OrderID),
			Username:       txtData.Username,
			TotalCartPrice: txtData.TotalCartPrice,
		}

		//Edit after-reserve delivery text for short and informative after-complete message
		aftCompleteMsg := tg.NewEditMessageText(usrID, msgID, h.extractFmt.AfterCompleteReply(data))
		h.Send(aftCompleteMsg)
		h.logger.Debugf("Sent after-complete message to %d", usrID)

		h.AnswerCallback(callbackID)
		h.logger.Debugf("Answered callback %s successfully", callbackID)
		return
	}
	//todo:update message with check mark, send delivery to pm by usrID
}

func (h *telegramTransport) HandleMessage(m *tg.Message) {
	//todo: parse /start and return command array
	h.logger.Debugf("received a message '%s' from '%s' chat '%d'", m.Text, m.From.UserName, m.Chat.ID)

	isGrp := m.Chat.IsGroup()
	chatID := m.Chat.ID
	usrID := m.From.ID
	lwrCaseMsg := strings.ToLower(m.Text)

	//For bot's personal messages
	if isGrp == false {
		switch {
		case lwrCaseMsg == "/work" || lwrCaseMsg == "войти":
			ok, err := h.runnerService.IsKnownByTelegramId(usrID)
			if err != nil {
				h.ResponseWithError(err, chatID)
				return
			}

			if !ok {
				msg := tg.NewMessage(chatID, templates.UsrIsNotKnownByTelegram)
				msg.ReplyMarkup = h.bot.GreetingKeyboard()
				h.Send(msg)
				return
			}

			msg := tg.NewMessage(chatID, templates.UsrIsKnownByTelegram)
			h.Send(msg)
			return
		}

		return
	}
	//For group chat
	switch m.Text {
	case "/work":
		msg := tg.NewMessage(chatID, templates.ComeHere)
		msg.ReplyMarkup = h.bot.LinkButton()
		//use of native send
		sentMsg, _ := h.bot.GetTelegramClient().Send(msg)

		//Introduce some delay for teapots
		time.Sleep(time.Second * 3)

		//Clear message to prevent spam in group
		delMsg := tg.NewDeleteMessage(chatID, sentMsg.MessageID)
		h.Send(delMsg)
		return
	default:
		//Bot doesn't respond for non '/work' messages in group chat
		return
	}
}

func (h *telegramTransport) HandleContact(m *tg.Message) {

	c := m.Contact
	chatID := m.Chat.ID
	telegramUsrID := m.From.ID
	username := c.FirstName
	usrPhNumber := c.PhoneNumber
	if strings.Split(usrPhNumber, "")[0] != "+" {
		usrPhNumber = "+" + usrPhNumber
	}
	h.logger.Debug(fmt.Sprintf("received a contact contact=(%s, %s, %s, %d) ", c.PhoneNumber, c.FirstName, c.LastName, c.UserID))

	runnerID, err := h.runnerService.IsRunner(usrPhNumber)
	if err != nil {
		h.ResponseWithError(err, chatID)
		h.logger.Error(err.Error())
		return
	}

	if runnerID == 0 {
		msg := tg.NewMessage(chatID, templates.UsrIsNotRunner)
		h.logger.Error(err.Error())
		h.Send(msg)
		return
	}

	err = h.runnerService.BeginWork(dto.RunnerBeginWorkDto{
		TelegramUserID: telegramUsrID,
		RunnerID:       runnerID,
	})
	if err != nil {
		h.ResponseWithError(err, chatID)
		h.logger.Error(err.Error())
		return
	}

	msg := tg.NewMessage(chatID, fmt.Sprintf(templates.BeginWorkSuccess, username))
	h.Send(msg)

	return
}

func (h *telegramTransport) ResponseWithError(err error, chatID int64) {
	var e tgErrors.TelegramError

	if errors.As(err, &e) {

		msg := tg.NewMessage(chatID, e.Error())
		if ok, kb := h.bot.ParseErrorForKeyboard(e); ok == true {
			msg.ReplyMarkup = kb
		}

		h.Send(msg)
		h.logger.Info(err.Error())

		return
	}
	msg := tg.NewMessage(chatID, templates.InternalServiceError)
	msg.ReplyMarkup = h.bot.InternalErrorKeyboard()
	h.Send(msg)
	h.logger.Error(err.Error())
	return
}

func (h *telegramTransport) AnswerCallback(callbackID string) {
	response := tg.NewCallback(callbackID, "")
	h.Send(response)
	return
}

func (h *telegramTransport) Send(c tg.Chattable) {

	inpTyp := reflect.TypeOf(c)
	delTyp := reflect.TypeOf(tg.DeleteMessageConfig{})
	msgTyp := reflect.TypeOf(tg.MessageConfig{})
	callbackTyp := reflect.TypeOf(tg.CallbackConfig{})

	var chatID int64

	switch {
	case inpTyp == delTyp:
		chatID = c.(tg.DeleteMessageConfig).ChatID
	case inpTyp == msgTyp:
		chatID = c.(tg.MessageConfig).ChatID
	case inpTyp == callbackTyp:
		//todo: read from cache
		chatID = 0
	default:
		chatID = 0
	}

	if _, err := h.bot.GetTelegramClient().Request(c); err != nil {
		if chatID != 0 {
			h.ResponseWithError(err, chatID)
		}
		h.logger.Error(err.Error())
		return
	}
}
