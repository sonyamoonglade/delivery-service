package telegram

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Transport interface {
	ListenForUpdates(bot *tg.BotAPI, cfg tg.UpdateConfig)
}
