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
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
	"strings"
	"time"
)

type telegramHandler struct {
	logger          *zap.Logger
	bot             *tg.BotAPI
	telegramService telegram.Service
	runnerService   runner.Service
	deliveryService delivery.Service
}

func NewTgHandler(logger *zap.Logger, bot *tg.BotAPI, run runner.Service, del delivery.Service, tgService telegram.Service) telegram.Transport {
	return &telegramHandler{logger: logger, bot: bot, runnerService: run, deliveryService: del, telegramService: tgService}
}

func (h *telegramHandler) ListenForUpdates(bot *tg.BotAPI, cfg tg.UpdateConfig) {
	updates := bot.GetUpdatesChan(cfg)
	for u := range updates {
		h.handle(&u)
	}
}

func (h *telegramHandler) handle(u *tg.Update) {
	if u.Message != nil {
		h.handleMessage(u.Message)
	}
	if u.CallbackQuery != nil {
		h.handleCallback(u.CallbackQuery)
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
		ok, err := h.deliveryService.Reserve(rsrvInp)
		if err != nil {
			h.ResponseWithError(err, usrID)
			h.bot.Send(response)
			h.logger.Error(err.Error())
			return
		}
		if !ok {
			//	//todo:handle

			return
		}

		editMsg := tg.NewEditMessageText(grpChatID, msgID, templates.Success)
		editKb := tg.NewEditMessageReplyMarkup(grpChatID, msgID, tg.NewInlineKeyboardMarkup())

		h.bot.Send(editMsg)
		h.bot.Send(editKb)

		resp := tg.NewCallback(cbID, "")

		h.bot.Send(resp)
		h.logger.Info("ok!")
		return
	case false:
		fmt.Println(data)
		h.bot.Send(response)
		return
	}
	////todo:update message with check mark, send delivery to pm by usrID
}

//todo: return err
func (h *telegramHandler) handleMessage(m *tg.Message) {
	fmt.Println(m.Text, m.From.FirstName)
	isGrp := m.Chat.IsGroup()
	chatID := m.Chat.ID
	telegramUsrID := m.From.ID
	if !isGrp {
		switch {
		case strings.ToLower(m.Text) == "/work" || strings.ToLower(m.Text) == "войти":
			usrID := m.From.ID
			ok, err := h.runnerService.IsKnownByTelegramId(usrID)
			if err != nil {
				h.ResponseWithError(err, chatID)
				return
			}
			if !ok {
				msg := tg.NewMessage(chatID, templates.IsNotKnownByTelegram)
				kb := genKb()
				kb.ResizeKeyboard = true
				msg.ReplyMarkup = kb
				h.bot.Send(msg)
				return
			}
			msg := tg.NewMessage(chatID, templates.IsKnownByTelegram)
			h.bot.Send(msg)
			return
		}
		if m.Contact != nil {
			usrPhNumber := m.Contact.PhoneNumber
			username := m.Contact.FirstName
			phWithPlus := "+" + usrPhNumber
			runnerID, err := h.runnerService.IsRunner(phWithPlus)

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
		return
	}
	switch m.Text {
	case "/work":
		msg := tg.NewMessage(chatID, "Тебе сюда 😇")
		kb := genBotLink()
		msg.ReplyMarkup = kb
		sent, _ := h.bot.Send(msg)
		time.Sleep(time.Second * 2)
		m := tg.NewEditMessageText(chatID, sent.MessageID, "abccd")
		h.bot.Send(m)
		return
	default:
		return
	}
	//todo: parse /start and return command array
}

func genKb() tg.ReplyKeyboardMarkup {
	b := tg.KeyboardButton{
		Text:           "Дать телефон",
		RequestContact: true,
	}
	b2 := tg.KeyboardButton{
		Text: "Войти",
	}
	row := []tg.KeyboardButton{b, b2}

	return tg.NewReplyKeyboard(row)
}

func genBotLink() tg.InlineKeyboardMarkup {
	link := "https://t.me/PrivateDeliveryBot"
	b := tg.InlineKeyboardButton{
		Text: "Перейти к боту",
		URL:  &link,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}

func (h *telegramHandler) ResponseWithError(err error, chatID int64) {
	var e tgErrors.TelegramError

	if errors.As(err, &e) {
		msg := tg.NewMessage(chatID, e.Error())
		if ok, kb := prepareErrorKb(e); ok == true {
			msg.ReplyMarkup = kb
		}
		h.bot.Send(msg)
		h.logger.Info(err.Error())
		return
	}
	msg := tg.NewMessage(chatID, templates.InternalServiceError)
	msg.ReplyMarkup = genInternalErrKb()
	h.logger.Error(err.Error())
	h.bot.Send(msg)
	return
}

func genInternalErrKb() tg.InlineKeyboardMarkup {
	//todo: move to env/config
	personalLink := "https://t.me/monasweet"
	b := tg.InlineKeyboardButton{
		Text: "Сообщить о проблеме 🔰",
		URL:  &personalLink,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}

func prepareErrorKb(err error) (bool, tg.ReplyKeyboardMarkup) {
	switch {
	//Sends a message with button to user's pm
	case strings.Contains(strings.ToLower(err.Error()), "вы не курьер!"):
		return true, genKb()
	default:
		return false, tg.NewReplyKeyboard()
	}

}
