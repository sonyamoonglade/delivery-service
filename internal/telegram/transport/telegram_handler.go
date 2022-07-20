package transport

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	dlvDto "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/callback"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

type telegramHandler struct {
	logger          *zap.SugaredLogger
	bot             *tg.BotAPI
	telegramService telegram.Service
	runnerService   runner.Service
	deliveryService delivery.Service
}

func NewTgHandler(logger *zap.SugaredLogger, bot *tg.BotAPI, run runner.Service, del delivery.Service, tgService telegram.Service) telegram.Transport {
	return &telegramHandler{logger: logger, bot: bot, runnerService: run, deliveryService: del, telegramService: tgService}
}

func (h *telegramHandler) ListenForUpdates(bot *tg.BotAPI, cfg tg.UpdateConfig) {
	updates := bot.GetUpdatesChan(cfg)
	for u := range updates {
		h.handle(&u)
	}
}

func (h *telegramHandler) handle(u *tg.Update) {
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

func (h *telegramHandler) HandleCallback(cb *tg.CallbackQuery) {

	rcvFromGroup := cb.Message.Chat.IsGroup()
	grpChatID := h.telegramService.GetGroupChatId()
	//usrID can be used as chatID for private messages
	usrID := cb.From.ID
	msgID := cb.Message.MessageID
	callbackID := cb.ID
	callbackData := cb.Data
	inititalMsg := cb.Message.Text

	h.logger.Debugf("received a callback data=%s from '%d'", callbackData, usrID)

	txtData := h.telegramService.ExtractDataFromText(inititalMsg)
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

		editMsg := tg.NewEditMessageText(grpChatID, msgID, templates.Success)
		h.Send(editMsg)

		h.logger.Infof("delivery ID=%d is reserved by runner ID=%d", inp.DeliveryID, rn.RunnerID)

		//Get initial message text

		//Apply extra parameters for runner's comfort
		personalData := bot.PersonalReserveReplyDto{
			DeliveryID: inp.DeliveryID,
			ReservedAt: reservedAt,
		}
		personalMsg := inititalMsg + bot.PersonalAfterReserveReply(personalData)

		//TODO: table to save chat's locales to print time
		//Send to runner's private messages so he follows the delivery
		msg := tg.NewMessage(usrID, personalMsg)
		msg.ReplyMarkup = bot.CompleteDeliveryKeyboard(inp)
		h.Send(msg)
		h.logger.Debug("Sent private after-reserve reply")
		//Extracted data from initial delivery message

		//Prepare new group message (edit initial)
		groupData := bot.GroupReserveReplyDto{
			DeliveryID:     inp.DeliveryID,
			ReservedAt:     reservedAt,
			OrderID:        txtData.OrderID,
			Username:       txtData.Username,
			TotalCartPrice: txtData.TotalCartPrice,
			RunnerUsername: rn.Username,
		}
		//Sent short and informative message to delivery group
		editMsg = tg.NewEditMessageText(grpChatID, msgID, bot.GroupAfterReserveReply(groupData))
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

		/*
			1. Replace old text with new template
			2. Delete markup completely
			3. Send it
		*/

		data := bot.PersonalCompleteReplyDto{
			DeliveryID:     inp.DeliveryID,
			OrderID:        txtData.OrderID,
			Username:       txtData.Username,
			TotalCartPrice: txtData.TotalCartPrice,
		}

		//Edit after-reserve delivery text for short and informative
		aftCompleteMsg := tg.NewEditMessageText(usrID, msgID, bot.AfterCompleteReply(data))
		h.Send(aftCompleteMsg)

		h.AnswerCallback(callbackID)
		return
	}
	//todo:update message with check mark, send delivery to pm by usrID
}

func (h *telegramHandler) HandleMessage(m *tg.Message) {
	//todo: parse /start and return command array
	h.logger.Debug(fmt.Sprintf("received a message '%s' from '%s'", m.Text, m.From.UserName))

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
				msg.ReplyMarkup = bot.GreetingKeyboard()
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
		msg.ReplyMarkup = bot.LinkButton()
		sentMsg, _ := h.bot.Send(msg)

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

func (h *telegramHandler) HandleContact(m *tg.Message) {

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

func (h *telegramHandler) ResponseWithError(err error, chatID int64) {
	var e tgErrors.TelegramError

	if errors.As(err, &e) {

		msg := tg.NewMessage(chatID, e.Error())
		if ok, kb := bot.ParseErrorForKeyboard(e); ok == true {
			msg.ReplyMarkup = kb
		}

		h.Send(msg)
		h.logger.Info(err.Error())

		return
	}
	msg := tg.NewMessage(chatID, templates.InternalServiceError)
	msg.ReplyMarkup = bot.InternalErrorKeyboard()
	h.Send(msg)
	h.logger.Error(err.Error())
	return
}

func (h *telegramHandler) AnswerCallback(callbackID string) {
	response := tg.NewCallback(callbackID, "")
	h.Send(response)
	return
}

func (h *telegramHandler) Send(c tg.Chattable) {

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

	if _, err := h.bot.Request(c); err != nil {
		if chatID != 0 {
			h.ResponseWithError(err, chatID)
		}
		h.logger.Error(err.Error())
		return
	}
}
