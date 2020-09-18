package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/matvoy/chat_server/api/proto/bot"
	pbchat "github.com/matvoy/chat_server/api/proto/chat"
	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	"github.com/rs/zerolog/log"
)

type InfobipWABody struct {
	Results             []*Result `json:"results"`
	MessageCount        int64     `json:"results"`
	PendingMessageCount int64     `json:"results"`
}

type Result struct {
	From            string `json:"from"`
	To              string `json:"to"`
	IntegrationType string `json:"integrationType"`
	ReceivedAt      string `json:"receivedAt"`
	MessageID       string `json:"messageId"`
	Message         `json:"message"`
	Contact         `json:"contact"`
	Price           `json:"price"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Contact struct {
	Name string `json:"name"`
}

type Price struct {
	PricePerMessage float64 `json:"pricePerMessage"`
	Currency        string  `json:"currency"`
}

type infobipWAClient struct {
	apiKey string
}

func NewInfobipWAClient(apiKey string) *infobipWAClient {
	return &infobipWAClient{
		apiKey,
	}
}

func (b *botService) configureInfobipWA(profile *pbentity.Profile) *infobipWAClient {
	apiKey, ok := profile.Variables["api_key"]
	if !ok {
		b.log.Fatal().Msg("token not found")
		return nil
	}
	return NewInfobipWAClient(apiKey)
}

func (b *botService) addProfileInfobipWA(req *pb.AddProfileRequest) error {
	apiKey, ok := req.Profile.Variables["api_key"]
	if !ok {
		b.log.Error().Msg("api_key not found")
		return errors.New("api_key not found")
	}
	bot := NewInfobipWAClient(apiKey)
	b.infobipWABots[req.Profile.Id] = bot
	b.botMap[req.Profile.Id] = "infobip-whatsapp"
	return nil
}

func (b *botService) deleteProfileInfobipWA(req *pb.DeleteProfileRequest) error {
	delete(b.infobipWABots, req.ProfileId)
	delete(b.botMap, req.ProfileId)
	return nil
}

func (b *botService) sendMessageInfobipWA(req *pb.SendMessageRequest) error {
	return nil
}

func (b *botService) InfobipWAWebhookHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/infobip/whatsapp/")
	profileID, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		b.log.Error().Msg(err.Error())
		return
	}
	update := &InfobipWABody{}
	if err := json.NewDecoder(r.Body).Decode(update); err != nil {
		log.Error().Msgf("could not decode request body: %s", err)
		return
	}
	if len(update.Results) == 0 ||
		(Message{}) == update.Results[0].Message {
		log.Warn().Msg("no data")
		return
	}
	b.log.Debug().
		Str("from", update.Results[0].From).
		Str("username", update.Results[0].Contact.Name).
		Str("text", update.Results[0].Message.Text).
		Msg("receive message")

	check := &pbchat.CheckSessionRequest{
		ExternalId: update.Results[0].From,
		ProfileId:  profileID,
		Username:   update.Results[0].Contact.Name,
	}
	resCheck, err := b.client.CheckSession(context.Background(), check)
	if err != nil {
		b.log.Error().Msg(err.Error())
		return
	}
	b.log.Debug().
		Bool("exists", resCheck.Exists).
		Int64("channel_id", resCheck.ChannelId).
		Int64("client_id", resCheck.ClientId).
		Msg("check user")

	if !resCheck.Exists {
		start := &pbchat.StartConversationRequest{
			User: &pbchat.User{
				UserId:     resCheck.ClientId,
				Type:       "telegram",
				Connection: p,
				Internal:   false,
			},
			DomainId: 1,
		}
		_, err := b.client.StartConversation(context.Background(), start)
		if err != nil {
			b.log.Error().Msg(err.Error())
			return
		}
	} else {
		textMessage := &pbentity.Message{
			Type: strings.ToLower(update.Results[0].Message.Type),
			Value: &pbentity.Message_TextMessage_{
				TextMessage: &pbentity.Message_TextMessage{
					Text: update.Results[0].Message.Text,
				},
			},
		}
		message := &pbchat.SendMessageRequest{
			Message:   textMessage,
			ChannelId: resCheck.ChannelId,
			FromFlow:  false,
		}
		_, err := b.client.SendMessage(context.Background(), message)
		if err != nil {
			b.log.Error().Msg(err.Error())
		}
	}
}
