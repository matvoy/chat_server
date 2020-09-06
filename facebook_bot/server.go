package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/facebook_bot/proto/bot_message"

	"github.com/gorilla/mux"
	"github.com/mileusna/facebook-messenger"
	"github.com/rs/zerolog"
)

type ChatServer interface {
	MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error
	StartWebhookServer() error
}

type facebookBotServer struct {
	log      *zerolog.Logger
	client   pbstorage.StorageService
	router   *mux.Router
	bots     map[int64]*messenger.Messenger
	profiles map[string]int64
}

func NewFacebookBotServer(
	log *zerolog.Logger,
	client pbstorage.StorageService,
) *facebookBotServer {
	v := &facebookBotServer{
		log:    log,
		client: client,
	}

	res, err := v.client.GetProfiles(context.Background(), &pbstorage.GetProfilesRequest{Type: "facebook"})
	if err != nil || res == nil {
		v.log.Error().Msg(err.Error())
		return nil
	}

	bots := make(map[int64]*messenger.Messenger)
	profiles := make(map[string]int64)
	for _, profile := range res.Profiles {
		accessToken, ok := profile.Variables["access_token"]
		if !ok {
			log.Error().Msg("token not found")
			return nil
		}
		verifyToken, ok := profile.Variables["verify_token"]
		if !ok {
			log.Error().Msg("token not found")
			return nil
		}
		pageID, ok := profile.Variables["page_id"]
		if !ok {
			log.Error().Msg("token not found")
			return nil
		}
		bot := &messenger.Messenger{
			AccessToken:     accessToken,
			VerifyToken:     verifyToken,
			PageID:          pageID,
			MessageReceived: v.MsgReceivedFunc,
		}

		http.Handle(fmt.Sprintf("/facebook/%v", profile.Id), bot)

		bots[profile.Id] = bot
		profiles[pageID] = profile.Id
	}
	v.bots = bots
	v.profiles = profiles
	return v
}

func (b *facebookBotServer) StartWebhookServer() error {
	b.log.Info().
		Int("port", cfg.AppPort).
		Msg("webhook started listening on port")
	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.AppPort), nil) // srv.ListenAndServeTLS(cfg.CertPath, cfg.KeyPath)
}

func (b *facebookBotServer) MsgReceivedFunc(msng *messenger.Messenger, userID int64, m messenger.FacebookMessage) {

	b.log.Debug().
		Int64("id", userID).
		Int64("username", userID).
		Str("text", m.Text).
		Msg("receive message")

	strID := strconv.FormatInt(userID, 10)

	message := &pbstorage.ProcessMessageRequest{
		SessionId:      strID,
		ExternalUserId: strID,
		Username:       strID,
		Text:           m.Text,
		ProfileId:      b.profiles[msng.PageID],
	}

	res, err := b.client.ProcessMessage(context.Background(), message)
	if err != nil || res == nil {
		b.log.Error().Msg(err.Error())
	}
	b.log.Debug().Msg("records created in the storage")
}

func (b *facebookBotServer) MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error {
	id, err := strconv.ParseInt(req.SessionId, 10, 64)
	if err != nil {
		b.log.Error().Msg(err.Error())
		return nil
	}
	_, err = b.bots[req.ProfileId].SendTextMessage(id, req.GetMessage().GetTextMessage().GetText())
	if err != nil {
		b.log.Error().Msg(err.Error())
	}
	return nil
}
