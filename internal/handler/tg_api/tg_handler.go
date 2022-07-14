package tghandler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgHandler struct {
	logger *zap.Logger
	bot    *tg.BotAPI
}

func NewTgHandler(logger *zap.Logger, bot *tg.BotAPI) *TgHandler {
	return &TgHandler{logger: logger, bot: bot}
}

func (h *TgHandler) ListenForUpdates(bot *tg.BotAPI, cfg tg.UpdateConfig) {
	updates := bot.GetUpdatesChan(cfg)
	for u := range updates {
		h.Map(&u)
	}
}

func (h *TgHandler) Map(u *tg.Update) {
	if u.Message != nil {
		h.HandleMessage(u.Message)
	}
	if u.CallbackQuery != nil {
		h.HandleCallback(u.CallbackQuery)
	}
}

func (h *TgHandler) HandleCallback(cb *tg.CallbackQuery) {

}

func (h *TgHandler) HandleMessage(m *tg.Message) {
	h.bot.Send(tg.NewMessage(m.Chat.ID, "Hello:D"))
}
