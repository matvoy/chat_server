package main

import (
	"context"
	"strconv"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pbflow "github.com/matvoy/chat_server/flow_adapter/proto/adapter"
	pb "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

type ChatServer interface {
	Start() error
}

type telegramBot struct {
	token      string
	log        *zerolog.Logger
	client     pbstorage.StorageService
	flowClient pbflow.AdapterService
	profileID  uint64
	bot        *tgbotapi.BotAPI
}

func NewTelegramBot(token string, profileID uint64, log *zerolog.Logger, client pbstorage.StorageService, flowClient pbflow.AdapterService) *telegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Debug().Str("username", bot.Self.UserName).Msg("authorized on account")
	return &telegramBot{
		token,
		log,
		client,
		flowClient,
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

		message := &pbstorage.MessageRequest{
			SessionId:      strconv.FormatInt(update.Message.Chat.ID, 10),
			ExternalUserId: strconv.Itoa(update.Message.From.ID),
			Username:       update.Message.From.UserName,
			FirstName:      update.Message.From.FirstName,
			LastName:       update.Message.From.LastName,
			Text:           update.Message.Text,
			Number:         number,
			ProfileId:      t.profileID,
		}

		messageFlow := &pbflow.MessageToFlow{
			SessionId:      strconv.FormatInt(update.Message.Chat.ID, 10),
			ExternalUserId: strconv.Itoa(update.Message.From.ID),
			Username:       update.Message.From.UserName,
			FirstName:      update.Message.From.FirstName,
			LastName:       update.Message.From.LastName,
			Text:           update.Message.Text,
			Number:         number,
			ProfileId:      t.profileID,
		}

		go func() {
			res, err := t.client.ProcessMessage(context.Background(), message)
			if err != nil || res == nil {
				t.log.Err(err)
			}
		}()

		resFlow, err := t.flowClient.SendMessageToFlow(context.Background(), messageFlow)
		if err != nil || resFlow == nil {
			t.log.Err(err)
			continue
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Created: %v", res.Created))
		// msg.ReplyToMessageID = update.Message.MessageID
		// t.bot.Send(msg)
	}

	return nil
}

func (t *telegramBot) ProcessMessageFromFlow(ctx context.Context, req *pb.MessageFromFlow, res *pb.Response) error {
	id, err := strconv.ParseInt(req.SessionId, 10, 64)
	if err != nil {
		t.log.Err(err)
		return err
	}
	msg := tgbotapi.NewMessage(id, req.Text)
	// msg.ReplyToMessageID = update.Message.MessageID
	_, err = t.bot.Send(msg)
	return err
}
