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
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
	"strconv"
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

func NewTgHandler(logger *zap.Logger, bot *tg.BotAPI, run runner.Service, del delivery.Service) telegram.Transport {
	return &telegramHandler{logger: logger, bot: bot, runnerService: run, deliveryService: del}
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
	usrID := cb.From.ID
	data := cb.Data

	idStr := strings.Split(data, " ")[1]
	deliveryID, _ := strconv.ParseInt(idStr, 10, 64)

	runnerID, err := h.runnerService.GetByTelegramId(usrID)

	inp := dlvDto.ReserveDeliveryDto{
		RunnerID:   runnerID,
		DeliveryID: deliveryID,
	}
	ok, err := h.deliveryService.Reserve(inp)
	if err != nil {
		//todo: handle err

		return
	}
	if !ok {
		//todo:handle

		return
	}

	//todo:update message with check mark, send delivery to pm by usrID
}

//todo: return err
func (h *telegramHandler) handleMessage(m *tg.Message) {
	fmt.Println(m.Text, m.From.FirstName)
	isGrp := m.Chat.IsGroup()
	chatID := m.Chat.ID
	telegramUsrID := m.From.ID
	if !isGrp {
		switch m.Text {
		case "/work":
			usrID := m.From.ID
			ok, err := h.runnerService.IsKnownByTelegramId(usrID)
			if err != nil {
				h.ResponseWithError(err, chatID)
				return
			}
			if !ok {
				msg := tg.NewMessage(chatID, templates.IsNotKnownByTelegram)
				kb := genKb()
				kb.OneTimeKeyboard = true
				kb.ResizeKeyboard = true
				msg.ReplyMarkup = kb
				h.bot.Send(msg)
				return
			}
			msg := tg.NewMessage(chatID, templates.IsKnownByTelegram)
			kb := genLink()
			msg.ReplyMarkup = kb
			h.bot.Send(msg)
			return
		}
		if m.Contact != nil {
			usrPhNumber := m.Contact.PhoneNumber
			username := m.Contact.FirstName
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
			kb := genLink()
			msg.ReplyMarkup = kb
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
	}
	//todo: parse /start and return command array
	h.bot.Send(tg.NewMessage(m.Chat.ID, "Hello From delivery bot!"))
}

func genKb() tg.ReplyKeyboardMarkup {
	b := tg.KeyboardButton{
		Text:           "Дать телефон",
		RequestContact: true,
	}
	row := []tg.KeyboardButton{b}

	return tg.NewReplyKeyboard(row)
}

func genLink() tg.InlineKeyboardMarkup {
	link := "https://t.me/+Z6oyrZJy2gllYzgy"
	b := tg.InlineKeyboardButton{
		Text: "Перейти в чат",
		URL:  &link,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
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
		h.bot.Send(msg)
		h.logger.Info(err.Error())
		return
	}
	msg := tg.NewMessage(chatID, templates.InternalServiceError)
	msg.ReplyMarkup = genErrKb()
	h.logger.Error(err.Error())
	h.bot.Send(msg)
	return
}

func genErrKb() tg.InlineKeyboardMarkup {
	//todo: move to env/config
	personalLink := "https://t.me/monasweet"
	b := tg.InlineKeyboardButton{
		Text: "Сообщить о проблеме 🔰",
		URL:  &personalLink,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}
