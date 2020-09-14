package main

import (
	"context"

	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	pbflow "github.com/matvoy/chat_server/api/proto/flow_client"
	"github.com/matvoy/chat_server/models"
)

func (s *chatService) routeMessage(channel *models.Channel, message *models.Message) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, &channel.ConversationID, nil, nil, &channel.ID)
	if err != nil {
		return err
	}
	reqMessage := &pbentity.Message{
		Type: message.Type,
		Value: &pbentity.Message_TextMessage_{
			TextMessage: &pbentity.Message_TextMessage{
				Text: message.Text.String,
			},
		},
	}
	if otherChannels == nil {
		if !channel.Internal {
			return s.sendMessageToFlow(channel, reqMessage)
		}
		return nil
	}
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				s.sendEventToWebitelUser(channel, item, reqMessage)
			}
		case "telegram":
			{
				s.sendMessageToTelegramUser(channel, item, reqMessage)
			}
		default:
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
		switch item.Type {
		// case "webitel":
		// 	{
		// 		s.sendToWebitelUser(channel, item, reqMessage)
		// 	}
		case "telegram":
			{
				s.sendMessageToTelegramUser(nil, item, message)
			}
		default:
		}
	}
	return nil
}

func (s *chatService) sendMessageToFlow(channel *models.Channel, message *pbentity.Message) error {
	sendMessage := &pbflow.SendMessageToFlowRequest{
		ConversationId: channel.ConversationID,
		Message:        message,
	}
	if res, err := s.flowClient.SendMessageToFlow(context.Background(), sendMessage); err != nil || res.Error != nil {
		if res != nil {
			s.log.Error().Msg(res.Error.Message)
		} else {
			s.log.Error().Msg(err.Error())
		}
		return err
	}
	return nil
}
