package event_router

import (
	"context"
	"encoding/json"
	"fmt"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pb "github.com/matvoy/chat_server/api/proto/chat"
	"github.com/matvoy/chat_server/internal/flow"
	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"
	"github.com/matvoy/chat_server/pkg/events"

	"github.com/micro/go-micro/v2/broker"
	"github.com/rs/zerolog"
)

type eventRouter struct {
	botClient  pbbot.BotService
	flowClient flow.Client
	broker     broker.Broker
	repo       repo.Repository
	log        *zerolog.Logger
}

type Router interface {
	RouteCloseConversation(channel *models.Channel, cause string) error
	RouteCloseConversationFromFlow(conversationID *int64, cause string) error
	RouteDeclineInvite(userID, conversationID *int64) error
	RouteInvite(conversationID, userID *int64) error
	RouteJoinConversation(channelID, conversationID *int64) error
	RouteLeaveConversation(channelID, conversationID *int64) error
	RouteMessage(channel *models.Channel, message *models.Message) error
	RouteMessageFromFlow(conversationID *int64, message *pb.Message) error
}

func NewRouter(
	botClient pbbot.BotService,
	flowClient flow.Client,
	broker broker.Broker,
	repo repo.Repository,
	log *zerolog.Logger,
) Router {
	return &eventRouter{
		botClient,
		flowClient,
		broker,
		repo,
		log,
	}
}

func (e *eventRouter) RouteCloseConversation(channel *models.Channel, cause string) error {
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, &channel.ConversationID, nil, nil, &channel.ID)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		if !channel.Internal {
			return e.flowClient.CloseConversation(channel.ConversationID)
		}
		return nil
	}
	body, _ := json.Marshal(events.CloseConversationEvent{
		ConversationID: channel.ConversationID,
		FromChannelID:  channel.ID,
		Cause:          cause,
	})
	for _, item := range otherChannels {
		var err error
		switch item.Type {
		case "webitel":
			{
				err = e.sendEventToWebitelUser(channel, item, events.CloseConversationEventType, body)
			}
		case "telegram", "infobip-whatsapp":
			{
				reqMessage := &pb.Message{
					Type: "text",
					Value: &pb.Message_TextMessage_{
						TextMessage: &pb.Message_TextMessage{
							Text: cause,
						},
					},
				}
				err = e.sendMessageToBotUser(channel, item, reqMessage)
			}
		default:
		}
		if err != nil {
			e.log.Warn().
				Int64("channel_id", item.ID).
				Bool("internal", item.Internal).
				Int64("user_id", item.UserID).
				Int64("conversation_id", item.ConversationID).
				Str("type", item.Type).
				Str("connection", item.Connection.String).
				Msg("failed to send close conversation event to channel")
		}
	}
	return nil
}

func (e *eventRouter) RouteCloseConversationFromFlow(conversationID *int64, cause string) error {
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, item := range otherChannels {
		switch item.Type {
		case "telegram", "infobip-whatsapp":
			{
				text := "Conversation closed"
				if cause != "" {
					text = cause
				}
				reqMessage := &pb.Message{
					Type: "text",
					Value: &pb.Message_TextMessage_{
						TextMessage: &pb.Message_TextMessage{
							Text: text,
						},
					},
				}
				if err := e.sendMessageToBotUser(nil, item, reqMessage); err != nil {
					e.log.Warn().
						Int64("channel_id", item.ID).
						Bool("internal", item.Internal).
						Int64("user_id", item.UserID).
						Int64("conversation_id", item.ConversationID).
						Str("type", item.Type).
						Str("connection", item.Connection.String).
						Msg("failed to send close conversation event to channel")
				}
			}
		default:
		}
	}
	return nil
}

func (e *eventRouter) RouteDeclineInvite(userID, conversationID *int64) error {
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	body, _ := json.Marshal(events.DeclineInvitationEvent{
		ConversationID: *conversationID,
		UserID:         *userID,
	})
	// TO DO declineInvitationToFlow??
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				if err := e.sendEventToWebitelUser(nil, item, events.DeclineInvitationEventType, body); err != nil {
					e.log.Warn().
						Int64("channel_id", item.ID).
						Bool("internal", item.Internal).
						Int64("user_id", item.UserID).
						Int64("conversation_id", item.ConversationID).
						Str("type", item.Type).
						Str("connection", item.Connection.String).
						Msg("failed to send invite conversation event to channel")
				}
			}
		default:
		}
	}
	return nil
}

func (e *eventRouter) RouteInvite(conversationID, userID *int64) error {
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	if err := e.sendInviteToWebitelUser(&otherChannels[0].DomainID, conversationID, userID); err != nil {
		return err
	}
	body, _ := json.Marshal(events.InviteConversationEvent{
		ConversationID: *conversationID,
		UserID:         *userID,
	})
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				if err := e.sendEventToWebitelUser(nil, item, events.InviteConversationEventType, body); err != nil {
					e.log.Warn().
						Int64("channel_id", item.ID).
						Bool("internal", item.Internal).
						Int64("user_id", item.UserID).
						Int64("conversation_id", item.ConversationID).
						Str("type", item.Type).
						Str("connection", item.Connection.String).
						Msg("failed to send invite conversation event to channel")
				}
			}
		default:
		}
	}
	return nil
}

func (e *eventRouter) sendInviteToWebitelUser(domainID, conversationID, userID *int64) error {
	body, _ := json.Marshal(map[string]int64{
		"conversation_id": *conversationID,
		"user_id":         *userID,
	})
	msg := &broker.Message{
		Header: map[string]string{},
		Body:   body,
	}
	if err := e.broker.Publish(fmt.Sprintf("event.%s.%v.%v", events.UserInvitationEventType, *domainID, *userID), msg); err != nil {
		return err
	}
	return nil
}

func (e *eventRouter) RouteJoinConversation(channelID, conversationID *int64) error {
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, channelID)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	body, _ := json.Marshal(events.JoinConversationEvent{
		ConversationID:  *conversationID,
		JoinedChannelID: *channelID,
	})
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				if err := e.sendEventToWebitelUser(nil, item, events.JoinConversationEventType, body); err != nil {
					e.log.Warn().
						Int64("channel_id", item.ID).
						Bool("internal", item.Internal).
						Int64("user_id", item.UserID).
						Int64("conversation_id", item.ConversationID).
						Str("type", item.Type).
						Str("connection", item.Connection.String).
						Msg("failed to send join conversation event to channel")
				}
			}
		default:
		}
	}
	return nil
}

func (e *eventRouter) RouteLeaveConversation(channelID, conversationID *int64) error {
	if err := e.flowClient.LeaveConversation(*conversationID); err != nil {
		return err
	}
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, channelID)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	body, _ := json.Marshal(events.LeaveConversationEvent{
		ConversationID:  *conversationID,
		LeavedChannelID: *channelID,
	})
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				if err := e.sendEventToWebitelUser(nil, item, events.JoinConversationEventType, body); err != nil {
					e.log.Warn().
						Int64("channel_id", item.ID).
						Bool("internal", item.Internal).
						Int64("user_id", item.UserID).
						Int64("conversation_id", item.ConversationID).
						Str("type", item.Type).
						Str("connection", item.Connection.String).
						Msg("failed to send leave conversation event to channel")
				}
			}
		default:
		}
	}
	return nil
}

func (e *eventRouter) RouteMessage(channel *models.Channel, message *models.Message) error {
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, &channel.ConversationID, nil, nil, &channel.ID)
	if err != nil {
		return err
	}
	reqMessage := &pb.Message{
		Id:   message.ID,
		Type: message.Type,
		Value: &pb.Message_TextMessage_{
			TextMessage: &pb.Message_TextMessage{
				Text: message.Text.String,
			},
		},
	}
	if otherChannels == nil {
		if !channel.Internal {
			return e.flowClient.SendMessage(channel.ConversationID, reqMessage)
		}
		return nil
	}
	body, _ := json.Marshal(events.MessageEvent{
		ConversationID: channel.ConversationID,
		FromChannelID:  channel.ID,
		// ToChannelID:    item.ID,
		MessageID: message.ID,
		Type:      message.Type,
		Value:     []byte(message.Text.String),
	})
	for _, item := range otherChannels {
		var err error
		switch item.Type {
		case "webitel":
			{
				err = e.sendEventToWebitelUser(channel, item, events.MessageEventType, body)
			}
		case "telegram", "infobip-whatsapp":
			{
				err = e.sendMessageToBotUser(channel, item, reqMessage)
			}
		default:
		}
		if err != nil {
			e.log.Warn().
				Int64("channel_id", item.ID).
				Bool("internal", item.Internal).
				Int64("user_id", item.UserID).
				Int64("conversation_id", item.ConversationID).
				Str("type", item.Type).
				Str("connection", item.Connection.String).
				Msg("failed to send message to channel")
		}
	}
	return nil
}

func (e *eventRouter) RouteMessageFromFlow(conversationID *int64, message *pb.Message) error {
	otherChannels, err := e.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, item := range otherChannels {
		var err error
		switch item.Type {
		// case "webitel":
		// 	{
		// 		e.sendToWebitelUser(channel, item, reqMessage)
		// 	}
		case "telegram", "infobip-whatsapp":
			{
				err = e.sendMessageToBotUser(nil, item, message)
			}
		default:
		}
		if err != nil {
			e.log.Warn().
				Int64("channel_id", item.ID).
				Bool("internal", item.Internal).
				Int64("user_id", item.UserID).
				Int64("conversation_id", item.ConversationID).
				Str("type", item.Type).
				Str("connection", item.Connection.String).
				Msg("failed to send message to channel [from flow]")
		}
	}
	return nil
}
