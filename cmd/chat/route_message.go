package main

import (
	"context"
	"strconv"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
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
			s.sendToFlow(channel, reqMessage)
		}
	}
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				s.sendToWebitelUser(channel, item, reqMessage)
			}
		case "telegram":
			{
				s.sendToTelegramUser(channel, item, reqMessage)
			}
		}
	}
	return nil
}

func (s *chatService) sendToFlow(channel *models.Channel, message *pbentity.Message) error {
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

func (s *chatService) sendToWebitelUser(from *models.Channel, to *models.Channel, message *pbentity.Message) error {
	return nil
}

func (s *chatService) sendToTelegramUser(from *models.Channel, to *models.Channel, message *pbentity.Message) error {
	profileID, err := strconv.ParseInt(to.Connection.String, 10, 64)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	client, err := s.repo.GetClientByID(context.Background(), to.UserID)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	botMessage := &pbbot.SendMessageRequest{
		ProfileId:      profileID,
		ExternalUserId: client.ExternalID.String,
		Message:        message,
	}
	if _, err := s.botClient.SendMessage(context.Background(), botMessage); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	return nil
}
