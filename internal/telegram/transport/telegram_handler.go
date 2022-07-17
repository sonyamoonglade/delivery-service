package transport

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"go.uber.org/zap"
)

type telegramHandler struct {
	logger          *zap.Logger
	bot             *tg.BotAPI
	telegramService telegram.Service
}

func NewTgHandler(logger *zap.Logger, bot *tg.BotAPI) telegram.Transport {
	return &telegramHandler{logger: logger, bot: bot}
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

}

//todo: return err
func (h *telegramHandler) handleMessage(m *tg.Message) {
	h.bot.Send(tg.NewMessage(m.Chat.ID, "Hello:D"))
}
