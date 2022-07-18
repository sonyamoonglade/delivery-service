package bot

import (
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"time"
)

var botLink string
var adminLink string

type Config struct {
	Token        string
	Timeout      int
	Debug        bool
	TelegramLink string
	AdminLink    string
}

func WithConfig(v *Config) (*tg.BotAPI, tg.UpdateConfig, error) {

	bot, err := tg.NewBotAPI(v.Token)

	if err != nil {
		return nil, tg.UpdateConfig{}, err
	}
	bot.Debug = v.Debug
	u := tg.NewUpdate(0)

	u.Timeout = 60

	//Set bot link for future link gets
	botLink = v.TelegramLink

	//Set admin link to have tech-support
	adminLink = v.AdminLink

	return bot, u, nil
}

func LinkButton() tg.InlineKeyboardMarkup {

	b := tg.InlineKeyboardButton{
		Text: "Перейти к боту",
		URL:  &botLink,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}

func GreetingKeyboard() tg.ReplyKeyboardMarkup {
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

func InternalErrorButton() tg.InlineKeyboardMarkup {
	b := tg.InlineKeyboardButton{
		Text: templates.Report,
		URL:  &adminLink,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}

func CompleteDeliveryButton(data *tgdelivery.CallbackData) tg.InlineKeyboardMarkup {

	bytes, _ := json.Marshal(data)
	dataStr := string(bytes)

	b := tg.InlineKeyboardButton{
		Text:         templates.Complete,
		CallbackData: &dataStr,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}

func AfterReserveReplyText(deliveryID int64, reservedAt time.Time) string {
	return fmt.Sprintf(templates.AfterReserveReply, reservedAt.Format("15:04 02.01"), deliveryID)
}
