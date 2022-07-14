package tgdelivery

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func BotWithConfig(v *BotConfig) (*tg.BotAPI, tg.UpdateConfig, error) {

	bot, err := tg.NewBotAPI(v.Token)

	if err != nil {
		return nil, tg.UpdateConfig{}, err
	}
	bot.Debug = v.Debug
	u := tg.NewUpdate(0)

	u.Timeout = 60

	return bot, u, nil
}
