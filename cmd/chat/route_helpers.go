package main

import (
	"context"
	"strconv"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	"github.com/matvoy/chat_server/models"
)

func (s *chatService) sendEventToWebitelUser(from *models.Channel, to *models.Channel, message *pbentity.Message) error {
	return nil
}

func (s *chatService) sendMessageToTelegramUser(from *models.Channel, to *models.Channel, message *pbentity.Message) error {
	profileID, err := strconv.ParseInt(to.Connection.String, 10, 64)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	client, err := s.repo.GetClientByID(context.Background(), to.UserID)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	botMessage := &pbbot.SendMessageRequest{
		ProfileId:      profileID,
		ExternalUserId: client.ExternalID.String,
		Message:        message,
	}
	if _, err := s.botClient.SendMessage(context.Background(), botMessage); err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	return nil
}
