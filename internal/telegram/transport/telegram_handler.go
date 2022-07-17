package transport

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"go.uber.org/zap"
)

type telegramHandler struct {
	logger          *zap.Logger
	bot             *tg.BotAPI
	telegramService telegram.Service
	runnerService   runner.Service
}

func NewTgHandler(logger *zap.Logger, bot *tg.BotAPI, run runner.Service) telegram.Transport {
	return &telegramHandler{logger: logger, bot: bot, runnerService: run}
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
	fmt.Println(cb.From)

}

//todo: return err
func (h *telegramHandler) handleMessage(m *tg.Message) {
	isGrp := m.Chat.IsGroup()

	if !isGrp {
		chatID := m.Chat.ID
		telegramUsrID := m.From.ID
		switch m.Text {
		case "/work":
			//usrID := m.From.ID
			//ok, err := h.runnerService.IsRunner(usrID)
			ok := false

			if !ok {
				msg := tg.NewMessage(chatID, "Извини, пока-что я не могу работать с тобой,\nно если ты дашь мне свой номер телефона - мы в расчете!")
				kb := genKb()
				kb.OneTimeKeyboard = true
				msg.ReplyMarkup = kb
				_, err := h.bot.Send(msg)
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}
		if m.Contact != nil {
			usrPhNumber := m.Contact.PhoneNumber
			username := m.Contact.FirstName
			runnerID, err := h.runnerService.IsRunner(usrPhNumber)

			if err != nil {
				//handle error
				msg := tg.NewMessage(chatID, err.Error())
				h.bot.Send(msg)
				h.logger.Error(err.Error())
				return
			}
			if runnerID == 0 {
				msg := tg.NewMessage(chatID, "Спасибо за твои данные,\nПока что, ты не зарегистрирован как доставщик,\nя вынужден отклонить твой запрос начать работу")
				h.bot.Send(msg)
				return
			}
			err = h.runnerService.BeginWork(dto.RunnerBeginWorkDto{
				TelegramUserID: telegramUsrID,
				RunnerID:       runnerID,
			})
			if err != nil {
				//handle error
				h.logger.Error(err.Error())
				return
			}
			msg := tg.NewMessage(chatID, fmt.Sprintf("Спасибо за твои данные.\nЯ нашел тебя, %s. Теперь мы можем работать вместе!\nТы можешь брать заказы в чате", username))
			h.bot.Send(msg)
			return
		}
	}
	h.bot.Send(tg.NewMessage(m.Chat.ID, "Hello From delivery bot!"))
}

func genKb() tg.ReplyKeyboardMarkup {
	//b := tg.InlineKeyboardButton{
	//	Text:                         "",
	//	URL:                          nil,
	//	LoginURL:                     nil,
	//	CallbackData:                 nil,
	//	SwitchInlineQuery:            nil,
	//	SwitchInlineQueryCurrentChat: nil,
	//	CallbackGame:                 nil,
	//	Pay:                          false,
	//}
	b := tg.KeyboardButton{
		Text:           "Дать телефон",
		RequestContact: true,
	}
	row := []tg.KeyboardButton{b}

	return tg.NewReplyKeyboard(row)
}
