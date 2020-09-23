package main

import (
	"context"
	"encoding/json"

	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	"github.com/matvoy/chat_server/models"
	"github.com/matvoy/chat_server/pkg/events"
)

func (s *chatService) routeMessage(channel *models.Channel, message *models.Message) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, &channel.ConversationID, nil, nil, &channel.ID)
	if err != nil {
		return err
	}
	reqMessage := &pbentity.Message{
		Id:   message.ID,
		Type: message.Type,
		Value: &pbentity.Message_TextMessage_{
			TextMessage: &pbentity.Message_TextMessage{
				Text: message.Text.String,
			},
		},
	}
	if otherChannels == nil {
		if !channel.Internal {
			return s.sendMessageToFlow(channel.ConversationID, reqMessage)
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
				err = s.sendEventToWebitelUser(channel, item, events.MessageEventType, body)
			}
		case "telegram", "infobip-whatsapp":
			{
				err = s.sendMessageToBotUser(channel, item, reqMessage)
			}
		default:
		}
		if err != nil {
			s.log.Warn().
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

func (s *chatService) routeMessageFromFlow(conversationID *int64, message *pbentity.Message) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, item := range otherChannels {
		var err error
		switch item.Type {
		// case "webitel":
		// 	{
		// 		s.sendToWebitelUser(channel, item, reqMessage)
		// 	}
		case "telegram", "infobip-whatsapp":
			{
				err = s.sendMessageToBotUser(nil, item, message)
			}
		default:
		}
		if err != nil {
			s.log.Warn().
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
