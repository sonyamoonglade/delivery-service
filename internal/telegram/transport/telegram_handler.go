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
	usrID := cb.From.ID
	msgID := cb.Message.MessageID
	callbackID := cb.ID
	data := cb.Data

	h.logger.Debug(fmt.Sprintf("received a callback data=%s from '%d'", data, usrID))

	switch rcvFromGroup {

	case true:
		var inp callback.ReserveData

		if err := callback.Decode(data, &inp); err != nil {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			h.logger.Error(err.Error())
			return
		}

		runnerID, err := h.runnerService.GetByTelegramId(usrID)
		if err != nil {
			h.ResponseWithError(err, usrID)
			h.AnswerCallback(callbackID)
			return
		}
		rsrvInp := dlvDto.ReserveDeliveryDto{
			RunnerID:   runnerID,
			DeliveryID: inp.DeliveryID,
		}
		reservedAt, err := h.deliveryService.Reserve(rsrvInp)
		if err != nil {
			h.ResponseWithError(err, usrID)
			h.AnswerCallback(callbackID)
			h.logger.Error(err.Error())
			return
		}

		editMsg := tg.NewEditMessageText(grpChatID, msgID, templates.Success)
		editKb := tg.NewEditMessageReplyMarkup(grpChatID, msgID, tg.NewInlineKeyboardMarkup())

		h.Send(editMsg)
		h.Send(editKb)

		h.AnswerCallback(callbackID)
		h.logger.Info(fmt.Sprintf("delivery ID=%d is reserved by runner ID=%d", inp.DeliveryID, runnerID))

		//todo: duplicate delivery text to user's pm
		//m is initialMessage

		//Apply extra reply text to delivery message
		dlvMsg := cb.Message.Text
		dlvMsg = dlvMsg + bot.AfterReserveReplyText(inp.DeliveryID, reservedAt)

		//TODO: table to save chat's locales to print time
		msg := tg.NewMessage(usrID, dlvMsg)
		msg.ReplyMarkup = bot.CompleteDeliveryKeyboard(callback.CompleteData{DeliveryID: inp.DeliveryID})
		h.Send(msg)
		return
	case false:

		var inp callback.CompleteData

		if err := callback.Decode(data, &inp); err != nil {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			h.logger.Error(err.Error())
			return
		}

		ok, err := h.deliveryService.Complete(inp.DeliveryID)
		if err != nil {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			h.logger.Error(err.Error())
		}
		if ok == false {
			h.ResponseWithError(err, grpChatID)
			h.AnswerCallback(callbackID)
			h.logger.Error(err.Error())
		}
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
		if ok, kb := h.PrepareErrorKeyboard(e); ok == true {
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

func (h *telegramHandler) PrepareErrorKeyboard(err error) (bool, tg.ReplyKeyboardMarkup) {
	switch {
	//Sends a message with button to user's pm
	case strings.Contains(strings.ToLower(err.Error()), "вы не курьер!"):
		return true, bot.GreetingKeyboard()
	default:
		return false, tg.NewReplyKeyboard()
	}

}

func (h *telegramHandler) AnswerCallback(callbackID string) {
	response := tg.NewCallback(callbackID, "")
	h.Send(response)
	return
}

func (h *telegramHandler) Send(c tg.Chattable) {
	if sent, err := h.bot.Send(c); err != nil {
		h.ResponseWithError(err, sent.Chat.ID)
		h.logger.Error(err.Error())
		return
	}
}
