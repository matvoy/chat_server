package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	pb "github.com/matvoy/chat_server/api/proto/bot"
	pbchat "github.com/matvoy/chat_server/api/proto/chat"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

type ChatServer interface {
	TelegramWebhookHandler(w http.ResponseWriter, r *http.Request)
	SendMessage(ctx context.Context, req *pb.SendMessageRequest, res *pb.SendMessageResponse) error
	AddProfile(ctx context.Context, req *pb.AddProfileRequest, res *pb.AddProfileResponse) error
	DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest, res *pb.DeleteProfileResponse) error
	StartWebhookServer() error
	StopWebhookServer() error
}

type botService struct {
	log          *zerolog.Logger
	client       pbchat.ChatService
	router       *mux.Router
	telegramBots map[int64]*tgbotapi.BotAPI
	botMap       map[int64]string
}

func NewBotService(
	log *zerolog.Logger,
	client pbchat.ChatService,
	router *mux.Router,
) *botService {
	b := &botService{
		log:    log,
		client: client,
		router: router,
	}
	b.botMap = make(map[int64]string)
	b.configureTelegram()
	return b
}

func (b *botService) StartWebhookServer() error {
	b.log.Info().
		Int("port", cfg.AppPort).
		Msg("webhook started listening on port")
	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.AppPort), b.router) // srv.ListenAndServeTLS(cfg.CertPath, cfg.KeyPath)
}

func (b *botService) StopWebhookServer() error {
	b.log.Info().
		Msg("removing webhooks")
	for k := range b.telegramBots {
		if _, err := b.telegramBots[k].RemoveWebhook(); err != nil {
			b.log.Error().Msg(err.Error())
		}
		delete(b.telegramBots, k)
	}
	return nil
}

func (b *botService) SendMessage(ctx context.Context, req *pb.SendMessageRequest, res *pb.SendMessageResponse) error {
	switch b.botMap[req.ProfileId] {
	case "telegram":
		{
			if err := b.sendMessageTelegram(req); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *botService) AddProfile(ctx context.Context, req *pb.AddProfileRequest, res *pb.AddProfileResponse) error {
	switch req.Profile.Type {
	case "telegram":
		{
			if err := b.addProfileTelegram(req); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *botService) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest, res *pb.DeleteProfileResponse) error {
	switch b.botMap[req.ProfileId] {
	case "telegram":
		{
			if err := b.deleteProfileTelegram(req); err != nil {
				return err
			}
		}
	}
	return nil
}
