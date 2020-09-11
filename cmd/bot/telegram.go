package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	pb "github.com/matvoy/chat_server/api/proto/bot"
	pbchat "github.com/matvoy/chat_server/api/proto/chat"
	"github.com/rs/zerolog/log"
)

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

func (b *botService) configureTelegram() {
	b.router.HandleFunc("/telegram/{profile_id}", b.TelegramWebhookHandler).
		Methods("POST")

	res, err := b.client.GetProfiles(context.Background(), &pbchat.GetProfilesRequest{Type: "telegram"})
	if err != nil || res == nil {
		b.log.Fatal().Msg(err.Error())
		return
	}

	bots := make(map[int64]*tgbotapi.BotAPI)
	for _, profile := range res.Profiles {
		b.botMap[profile.Id] = "telegram"
		token, ok := profile.Variables["token"]
		if !ok {
			b.log.Fatal().Msg("token not found")
			return
		}
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			b.log.Fatal().Msg(err.Error())
			return
		}
		// webhookInfo := tgbotapi.NewWebhookWithCert(fmt.Sprintf("%s/telegram/%v", cfg.TgWebhook, profile.Id), cfg.CertPath)
		webhookInfo := tgbotapi.NewWebhook(fmt.Sprintf("%s/telegram/%v", cfg.TgWebhook, profile.Id))
		_, err = bot.SetWebhook(webhookInfo)
		if err != nil {
			b.log.Fatal().Msg(err.Error())
			return
		}
		bots[profile.Id] = bot
	}

	b.telegramBots = bots
}

func (b *botService) addProfileTelegram(req *pb.AddProfileRequest) error {
	token, ok := req.Profile.Variables["token"]
	if !ok {
		b.log.Error().Msg("token not found")
		return errors.New("token not found")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		b.log.Error().Msg(err.Error())
		return err
	}
	// webhookInfo := tgbotapi.NewWebhookWithCert(fmt.Sprintf("%s/telegram/%v", cfg.TgWebhook, profile.Id), cfg.CertPath)
	webhookInfo := tgbotapi.NewWebhook(fmt.Sprintf("%s/telegram/%v", cfg.TgWebhook, req.Profile.Id))
	_, err = bot.SetWebhook(webhookInfo)
	if err != nil {
		b.log.Error().Msg(err.Error())
		return err
	}
	b.telegramBots[req.Profile.Id] = bot
	return nil
}

func (b *botService) deleteProfileTelegram(req *pb.DeleteProfileRequest) error {
	if _, err := b.telegramBots[req.ProfileId].RemoveWebhook(); err != nil {
		b.log.Error().Msg(err.Error())
		return err
	}
	delete(b.telegramBots, req.ProfileId)
	return nil
}

func (b *botService) sendMessageTelegram(req *pb.SendMessageRequest) error {
	id, err := strconv.ParseInt(req.SessionId, 10, 64)
	if err != nil {
		b.log.Error().Msg(err.Error())
		return err
	}
	msg := tgbotapi.NewMessage(id, req.GetMessage().GetTextMessage().GetText())
	// msg.ReplyToMessageID = update.Message.MessageID
	_, err = b.telegramBots[req.ProfileId].Send(msg)
	if err != nil {
		b.log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (b *botService) TelegramWebhookHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/telegram/")
	profileID, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		b.log.Error().Msg(err.Error())
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

	b.log.Debug().
		Int64("id", update.Message.From.ID).
		Str("username", update.Message.From.Username).
		Str("first_name", update.Message.From.FirstName).
		Str("last_name", update.Message.From.LastName).
		Str("text", update.Message.Text).
		Msg("receive message")

	strChatID := strconv.FormatInt(update.Message.Chat.ID, 10)

	message := &pbchat.ProcessMessageRequest{
		SessionId:      strChatID,
		ExternalUserId: strconv.FormatInt(update.Message.From.ID, 10),
		Username:       update.Message.From.Username,
		FirstName:      update.Message.From.FirstName,
		LastName:       update.Message.From.LastName,
		Text:           update.Message.Text,
		ProfileId:      profileID,
	}

	res, err := b.client.ProcessMessage(context.Background(), message)
	if err != nil || res == nil {
		b.log.Error().Msg(err.Error())
	}
	b.log.Debug().Msg("records created in the storage")
}
