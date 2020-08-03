package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

type telegramBot struct {
	token string
	log   *zerolog.Logger
}

func NewTelegramBot(token string, log *zerolog.Logger) *telegramBot {
	return &telegramBot{
		token,
		log,
	}
}

func (t telegramBot) Start() error {
	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		t.log.Fatal().Msg(err.Error())
	}

	// bot.Debug = true

	t.log.Debug().Str("username", bot.Self.UserName).Msg("authorized on account")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		t.log.Debug().
			Int("id", update.Message.From.ID).
			Str("username", update.Message.From.UserName).
			Str("first_name", update.Message.From.FirstName).
			Str("last_name", update.Message.From.LastName).
			Str("text", update.Message.Text).
			Msg("receive message")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

	return nil
}
