package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	dlvDto "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
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
		h.handleContact(u.Message)
		return
	case u.Message != nil:
		h.handleMessage(u.Message)
		return
	case u.CallbackQuery != nil:
		h.handleCallback(u.CallbackQuery)
		return
	default:
		return
	}

}

func (h *telegramHandler) handleCallback(cb *tg.CallbackQuery) {
	rcvFromGroup := cb.Message.Chat.IsGroup()
	grpChatID := h.telegramService.GetGroupChatId()
	usrID := cb.From.ID
	msgID := cb.Message.MessageID
	cbID := cb.ID
	data := cb.Data
	response := tg.NewCallback(cbID, "")

	h.logger.Debug(fmt.Sprintf("received a callback data=%s from '%d'", data, usrID))

	switch rcvFromGroup {

	case true:
		var inp tgdelivery.CallbackData

		err := json.Unmarshal([]byte(data), &inp)
		if err != nil {
			h.ResponseWithError(err, grpChatID)
			h.bot.Send(response)
			h.logger.Error(err.Error())
			return
		}

		runnerID, err := h.runnerService.GetByTelegramId(usrID)
		if err != nil {
			h.ResponseWithError(err, usrID)
			h.bot.Send(response)
			return
		}
		rsrvInp := dlvDto.ReserveDeliveryDto{
			RunnerID:   runnerID,
			DeliveryID: inp.DeliveryID,
		}
		reservedAt, err := h.deliveryService.Reserve(rsrvInp)
		if err != nil {
			h.ResponseWithError(err, usrID)
			h.bot.Send(response)
			h.logger.Error(err.Error())
			return
		}

		editMsg := tg.NewEditMessageText(grpChatID, msgID, templates.Success)
		editKb := tg.NewEditMessageReplyMarkup(grpChatID, msgID, tg.NewInlineKeyboardMarkup())

		h.bot.Send(editMsg)
		h.bot.Send(editKb)

		resp := tg.NewCallback(cbID, "")

		h.bot.Send(resp)

		h.logger.Info(fmt.Sprintf("delivery ID=%d is reserved by runner ID=%d", inp.DeliveryID, runnerID))

		//todo: duplicate delivery text to user's pm
		//m is initialMessage

		//Apply extra data to delivery text
		dlvMsg := cb.Message.Text
		dlvMsg = dlvMsg + bot.AfterReserveReplyText(inp.DeliveryID, reservedAt)

		//TODO: table to save chat's locales to print time
		msg := tg.NewMessage(usrID, dlvMsg)
		msg.ReplyMarkup = bot.CompleteDeliveryButton(&inp)
		h.bot.Send(msg)
		return
	case false:
		fmt.Println(data)
		h.bot.Send(response)
		return
	}
	//todo:update message with check mark, send delivery to pm by usrID
}

func (h *telegramHandler) handleMessage(m *tg.Message) {
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
				h.bot.Send(msg)
				return
			}

			msg := tg.NewMessage(chatID, templates.UsrIsKnownByTelegram)
			h.bot.Send(msg)
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
		h.bot.Send(delMsg)
		return
	default:
		//Bot doesn't respond for non '/work' messages in group chat
		return
	}
}

func (h *telegramHandler) handleContact(m *tg.Message) {

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
		return
	}

	if runnerID == 0 {
		msg := tg.NewMessage(chatID, templates.UsrIsNotRunner)
		h.bot.Send(msg)
		return
	}

	err = h.runnerService.BeginWork(dto.RunnerBeginWorkDto{
		TelegramUserID: telegramUsrID,
		RunnerID:       runnerID,
	})
	if err != nil {
		h.ResponseWithError(err, chatID)
		return
	}

	msg := tg.NewMessage(chatID, fmt.Sprintf(templates.BeginWorkSuccess, username))
	h.bot.Send(msg)

	return
}

func (h *telegramHandler) ResponseWithError(err error, chatID int64) {
	var e tgErrors.TelegramError

	if errors.As(err, &e) {

		msg := tg.NewMessage(chatID, e.Error())
		if ok, kb := h.PrepareErrorKeyboard(e); ok == true {
			msg.ReplyMarkup = kb
		}

		h.bot.Send(msg)
		h.logger.Info(err.Error())

		return
	}
	msg := tg.NewMessage(chatID, templates.InternalServiceError)
	msg.ReplyMarkup = bot.InternalErrorButton()
	h.logger.Error(err.Error())
	h.bot.Send(msg)
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
