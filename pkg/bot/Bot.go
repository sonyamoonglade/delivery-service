package bot

import (
	"encoding/json"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sonyamoonglade/delivery-service/pkg/callback"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"go.uber.org/zap"
	"strings"
)

type Config struct {
	Timeout     int
	Debug       bool
	URL         string
	GroupChatID int64
	BotToken    string
	AdminLink   string
}

type Bot interface {
	GetGroupChatId() int64
	GetTelegramClient() *tg.BotAPI
	LinkButton() tg.InlineKeyboardMarkup
	GetUpdatesConfig() tg.UpdateConfig
	GreetingKeyboard() tg.ReplyKeyboardMarkup
	InternalErrorKeyboard() tg.InlineKeyboardMarkup
	CompleteDeliveryKeyboard(data callback.Data) tg.InlineKeyboardMarkup
	ReserveDeliveryKeyboard(data callback.Data) tg.InlineKeyboardMarkup
	ParseErrorForKeyboard(err error) (bool, tg.ReplyKeyboardMarkup)
	PostDeliveryMessage(msg string, deliveryID int64) error
}

type bot struct {
	logger         *zap.SugaredLogger
	telegramClient *tg.BotAPI
	updateConfig   tg.UpdateConfig
	botLink        string
	adminLink      string
	groupChatID    int64
}

func NewBot(v *Config, logger *zap.SugaredLogger) (Bot, error) {

	client, err := tg.NewBotAPI(v.BotToken)
	if err != nil {
		return nil, err
	}
	client.Debug = v.Debug

	updateCfg := tg.NewUpdate(0)

	updateCfg.Offset = 0
	updateCfg.Timeout = v.Timeout
	return &bot{
		logger:         logger,
		telegramClient: client,
		updateConfig:   updateCfg,
		botLink:        v.URL,
		adminLink:      v.AdminLink,
		groupChatID:    v.GroupChatID,
	}, nil
}

func (t *bot) PostDeliveryMessage(text string, deliveryID int64) error {

	data := callback.Data{DeliveryID: deliveryID}
	keyboard := t.ReserveDeliveryKeyboard(data)

	msg := tg.NewMessage(t.groupChatID, text)
	msg.ReplyMarkup = keyboard

	_, err := t.telegramClient.Send(msg)
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t *bot) GetUpdatesConfig() tg.UpdateConfig {
	return t.updateConfig
}
func (t *bot) GetTelegramClient() *tg.BotAPI {
	return t.telegramClient
}
func (t *bot) LinkButton() tg.InlineKeyboardMarkup {
	b := tg.InlineKeyboardButton{
		Text: "Перейти к боту",
		URL:  &t.botLink,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}
func (t *bot) GreetingKeyboard() tg.ReplyKeyboardMarkup {
	b := tg.KeyboardButton{
		Text:           "Дать телефон",
		RequestContact: true,
	}
	b2 := tg.KeyboardButton{
		Text: "Войти",
	}
	row := []tg.KeyboardButton{b, b2}

	kb := tg.NewReplyKeyboard(row)
	kb.OneTimeKeyboard = true
	kb.ResizeKeyboard = true
	return kb
}
func (t *bot) InternalErrorKeyboard() tg.InlineKeyboardMarkup {
	b := tg.InlineKeyboardButton{
		Text: templates.Report,
		URL:  &t.adminLink,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}
func (t *bot) CompleteDeliveryKeyboard(data callback.Data) tg.InlineKeyboardMarkup {

	bytes, _ := json.Marshal(data)
	dataStr := string(bytes)

	b := tg.InlineKeyboardButton{
		Text:         templates.Complete,
		CallbackData: &dataStr,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}
func (t *bot) ReserveDeliveryKeyboard(data callback.Data) tg.InlineKeyboardMarkup {

	bytes, _ := json.Marshal(data)
	strData := string(bytes)
	reserveButton := tg.InlineKeyboardButton{
		Text:         templates.Reserve,
		CallbackData: &strData,
	}
	row := []tg.InlineKeyboardButton{reserveButton}

	return tg.NewInlineKeyboardMarkup(row)
}
func (t *bot) ParseErrorForKeyboard(err error) (bool, tg.ReplyKeyboardMarkup) {
	switch {
	//Sends a message with button to user's pm
	case strings.Contains(strings.ToLower(err.Error()), "вы не курьер!"):
		return true, t.GreetingKeyboard()
	default:
		return false, tg.NewReplyKeyboard()
	}
}
func (t *bot) GetGroupChatId() int64 {
	return t.groupChatID
}
