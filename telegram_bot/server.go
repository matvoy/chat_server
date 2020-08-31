package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ChatServer interface {
	WebhookHandler(w http.ResponseWriter, r *http.Request)
	MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error
}

type webhookReqBody struct {
	Message struct {
		MessageID int64  `json:"message_id"`
		Text      string `json:"text"`
		From      struct {
			Username  string `json:"username"`
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"from"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type telegramBot struct {
	log    *zerolog.Logger
	client pbstorage.StorageService
	bots   map[int64]*tgbotapi.BotAPI
}

func NewTelegramBot(
	log *zerolog.Logger,
	client pbstorage.StorageService,
) *telegramBot {
	t := &telegramBot{
		log:    log,
		client: client,
	}

	res, err := t.client.GetProfiles(context.Background(), &pbstorage.GetProfilesRequest{Type: "telegram"})
	if err != nil || res == nil {
		t.log.Error().Msg(err.Error())
		return nil
	}

	bots := make(map[int64]*tgbotapi.BotAPI)
	for _, profile := range res.Profiles {
		token, ok := profile.Variables["token"]
		if !ok {
			log.Error().Msg("token not found")
			return nil
		}
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Error().Msg(err.Error())
			return nil
		}
		webhookInfo := tgbotapi.NewWebhookWithCert(fmt.Sprintf("https://sdafsdf/telegram/%v", profile.Id), "")
		_, err = bot.SetWebhook(webhookInfo)
		bots[profile.Id] = bot
	}
	t.bots = bots
	return t
}

func (t *telegramBot) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/telegram/")
	profileID, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		t.log.Error().Msg(err.Error())
		return
	}
	update := &webhookReqBody{}
	if err := json.NewDecoder(r.Body).Decode(update); err != nil {
		log.Error().Msgf("could not decode request body: %s", err)
		return
	}

	if update.Message.Text == "" { // ignore any non-Message Updates
		return
	}

	t.log.Debug().
		Int64("id", update.Message.From.ID).
		Str("username", update.Message.From.Username).
		Str("first_name", update.Message.From.FirstName).
		Str("last_name", update.Message.From.LastName).
		Str("text", update.Message.Text).
		Msg("receive message")

	strChatID := strconv.FormatInt(update.Message.Chat.ID, 10)

	message := &pbstorage.ProcessMessageRequest{
		SessionId:      strChatID,
		ExternalUserId: strconv.FormatInt(update.Message.From.ID, 10),
		Username:       update.Message.From.Username,
		FirstName:      update.Message.From.FirstName,
		LastName:       update.Message.From.LastName,
		Text:           update.Message.Text,
		ProfileId:      profileID,
	}

	res, err := t.client.ProcessMessage(context.Background(), message)
	if err != nil || res == nil {
		t.log.Error().Msg(err.Error())
	}
	t.log.Debug().Msg("records created in the storage")

}

func (t *telegramBot) MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error {
	id, err := strconv.ParseInt(req.SessionId, 10, 64)
	if err != nil {
		t.log.Error().Msg(err.Error())
		return nil
	}
	msg := tgbotapi.NewMessage(id, req.GetMessage().GetTextMessage().GetText())
	// msg.ReplyToMessageID = update.Message.MessageID
	_, err = t.bots[req.ProfileId].Send(msg)
	if err != nil {
		t.log.Error().Msg(err.Error())
	}
	return nil
}
