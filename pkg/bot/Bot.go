package bot

import (
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sonyamoonglade/delivery-service/pkg/callback"
	"github.com/sonyamoonglade/delivery-service/pkg/templates"
	"strings"
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

type PersonalReserveReplyDto struct {
	DeliveryID int64
	ReservedAt time.Time
}

type GroupReserveReplyDto struct {
	DeliveryID     int64
	OrderID        string
	Username       string
	TotalCartPrice int64
	ReservedAt     time.Time
	RunnerUsername string
}

type PersonalCompleteReplyDto struct {
	DeliveryID     int64
	OrderID        string
	Username       string
	TotalCartPrice int64
}

type DataFromText struct {
	OrderID        string
	TotalCartPrice int64
	Username       string
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

func InternalErrorKeyboard() tg.InlineKeyboardMarkup {
	b := tg.InlineKeyboardButton{
		Text: templates.Report,
		URL:  &adminLink,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}

func CompleteDeliveryKeyboard(data callback.Data) tg.InlineKeyboardMarkup {

	bytes, _ := json.Marshal(data)
	dataStr := string(bytes)

	b := tg.InlineKeyboardButton{
		Text:         templates.Complete,
		CallbackData: &dataStr,
	}
	row := []tg.InlineKeyboardButton{b}
	return tg.NewInlineKeyboardMarkup(row)
}

func ReserveDeliveryKeyboard(data callback.Data) tg.InlineKeyboardMarkup {

	bytes, _ := json.Marshal(data)
	strData := string(bytes)
	fmt.Println(strData)
	reserveButton := tg.InlineKeyboardButton{
		Text:         templates.Reserve,
		CallbackData: &strData,
	}
	row := []tg.InlineKeyboardButton{reserveButton}

	return tg.NewInlineKeyboardMarkup(row)
}

func ParseErrorForKeyboard(err error) (bool, tg.ReplyKeyboardMarkup) {
	switch {
	//Sends a message with button to user's pm
	case strings.Contains(strings.ToLower(err.Error()), "вы не курьер!"):
		return true, GreetingKeyboard()
	default:
		return false, tg.NewReplyKeyboard()
	}
}

func PersonalAfterReserveReply(dto PersonalReserveReplyDto) string {
	return fmt.Sprintf(templates.PersonalAfterReserveText, dto.ReservedAt.Format("15:04 02.01"), dto.DeliveryID)
}

func GroupAfterReserveReply(dto GroupReserveReplyDto) string {
	return fmt.Sprintf(templates.GroupAfterReserveText,
		dto.DeliveryID,
		dto.ReservedAt.Format("15:04 02.01"),
		templates.Success,
		dto.RunnerUsername,
		dto.OrderID,
		dto.Username,
		dto.TotalCartPrice)
}

func AfterCompleteReply(dto PersonalCompleteReplyDto) string {
	return fmt.Sprintf(templates.DeliveryCompletedText, dto.DeliveryID, templates.Success, dto.OrderID, dto.Username, dto.TotalCartPrice)
}
