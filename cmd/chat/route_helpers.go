package main

import (
	"context"
	"fmt"
	"strconv"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	"github.com/matvoy/chat_server/models"

	"github.com/micro/go-micro/v2/broker"
)

func (s *chatService) sendEventToWebitelUser(from *models.Channel, to *models.Channel, eventType string, body []byte) error {
	msg := &broker.Message{
		Header: map[string]string{},
		Body:   body,
	}
	if err := s.broker.Publish(fmt.Sprintf("event.%s.%v.%v", eventType, to.DomainID, to.UserID), msg); err != nil {
		return err
	}
	return nil
}

func (s *chatService) sendMessageToBotUser(from *models.Channel, to *models.Channel, message *pbentity.Message) error {
	profileID, err := strconv.ParseInt(to.Connection.String, 10, 64)
	if err != nil {
		return err
	}
	client, err := s.repo.GetClientByID(context.Background(), to.UserID)
	if err != nil {
		return err
	}
	botMessage := &pbbot.SendMessageRequest{
		ProfileId:      profileID,
		ExternalUserId: client.ExternalID.String,
		Message:        message,
	}
	if _, err := s.botClient.SendMessage(context.Background(), botMessage); err != nil {
		return err
	}
	return nil
}
