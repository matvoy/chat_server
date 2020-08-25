package main

import (
	"context"
	"strconv"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

type ChatServer interface {
	Start() error
	MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error
}

type telegramBot struct {
	token     string
	log       *zerolog.Logger
	client    pbstorage.StorageService
	profileID uint64
	bot       *tgbotapi.BotAPI
}

func NewTelegramBot(
	token string,
	profileID uint64,
	log *zerolog.Logger,
	client pbstorage.StorageService,
) *telegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Debug().Str("username", bot.Self.UserName).Msg("authorized on account")
	return &telegramBot{
		token,
		log,
		client,
		profileID,
		bot,
	}
}

func (t *telegramBot) Start() error {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		t.log.Fatal().Msg(err.Error())
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		number := ""
		if update.Message.Contact != nil {
			number = update.Message.Contact.PhoneNumber
		}

		t.log.Debug().
			Int("id", update.Message.From.ID).
			Str("username", update.Message.From.UserName).
			Str("first_name", update.Message.From.FirstName).
			Str("last_name", update.Message.From.LastName).
			Str("number", number).
			Str("text", update.Message.Text).
			Msg("receive message")

		strChatID := strconv.FormatInt(update.Message.Chat.ID, 10)

		message := &pbstorage.ProcessMessageRequest{
			SessionId:      strChatID,
			ExternalUserId: strconv.Itoa(update.Message.From.ID),
			Username:       update.Message.From.UserName,
			FirstName:      update.Message.From.FirstName,
			LastName:       update.Message.From.LastName,
			Text:           update.Message.Text,
			Number:         number,
			ProfileId:      t.profileID,
		}

		res, err := t.client.ProcessMessage(context.Background(), message)
		if err != nil || res == nil {
			t.log.Error().Msg(err.Error())
		}
		t.log.Debug().Msg("records created in the storage")
	}

	return nil
}

func (t *telegramBot) MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error {
	id, err := strconv.ParseInt(req.SessionId, 10, 64)
	if err != nil {
		t.log.Error().Msg(err.Error())
		return nil
	}
	msg := tgbotapi.NewMessage(id, req.GetMessage().GetTextMessage().GetText())
	// msg.ReplyToMessageID = update.Message.MessageID
	_, err = t.bot.Send(msg)
	if err != nil {
		t.log.Error().Msg(err.Error())
	}
	return nil
}
