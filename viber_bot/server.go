package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/viber_bot/proto/bot_message"

	"github.com/gorilla/mux"
	"github.com/mileusna/viber"
	"github.com/rs/zerolog"
)

type ChatServer interface {
	MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error
	StartWebhookServer() error
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

type viberBotServer struct {
	log      *zerolog.Logger
	client   pbstorage.StorageService
	router   *mux.Router
	bots     map[int64]*viber.Viber
	profiles map[string]int64
}

func NewViberBotServer(
	log *zerolog.Logger,
	client pbstorage.StorageService,
) *viberBotServer {
	v := &viberBotServer{
		log:    log,
		client: client,
	}

	// v.router.HandleFunc("/viber/{profile_id}", v.WebhookHandler).
	// 	Methods("POST")

	res, err := v.client.GetProfiles(context.Background(), &pbstorage.GetProfilesRequest{Type: "viber"})
	if err != nil || res == nil {
		v.log.Error().Msg(err.Error())
		return nil
	}

	bots := make(map[int64]*viber.Viber)
	profiles := make(map[string]int64)
	for _, profile := range res.Profiles {
		token, ok := profile.Variables["token"]
		if !ok {
			log.Error().Msg("token not found")
			return nil
		}
		bot := &viber.Viber{
			AppKey: token,
			Sender: viber.Sender{
				Name:   "",
				Avatar: "",
			},
			Message: v.MsgReceivedFunc,
		}
		http.Handle(fmt.Sprintf("/viber/%v", profile.Id), bot)
		bot.SetWebhook(fmt.Sprintf("%s/viber/%v", cfg.ViberWebhook, profile.Id), []string{"message", "subscribed", "unsubscribed", "conversation_started"})
		bots[profile.Id] = bot
		profiles[token] = profile.Id
	}
	v.bots = bots
	v.profiles = profiles
	return v
}

func (b *viberBotServer) StartWebhookServer() error {
	b.log.Info().
		Int("port", cfg.AppPort).
		Msg("webhook started listening on port")
	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.AppPort), nil) // srv.ListenAndServeTLS(cfg.CertPath, cfg.KeyPath)
}

func (b *viberBotServer) MsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	switch m.(type) {

	case *viber.TextMessage:
		txt := m.(*viber.TextMessage).Text
		b.log.Debug().
			Str("id", u.ID).
			Str("username", u.Name).
			Str("text", txt).
			Msg("receive message")

		message := &pbstorage.ProcessMessageRequest{
			SessionId:      u.ID,
			ExternalUserId: u.ID,
			Username:       u.Name,
			Text:           txt,
			ProfileId:      b.profiles[v.AppKey],
		}

		res, err := b.client.ProcessMessage(context.Background(), message)
		if err != nil || res == nil {
			b.log.Error().Msg(err.Error())
		}
		b.log.Debug().Msg("records created in the storage")

		// TO DO OTHER TYPES
	}
}

func (b *viberBotServer) MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error {
	_, err := b.bots[req.ProfileId].SendTextMessage(req.SessionId, req.GetMessage().GetTextMessage().GetText())
	if err != nil {
		b.log.Error().Msg(err.Error())
	}
	return nil
}
