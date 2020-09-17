package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro/go-micro/v2/broker"
)

func (s *chatService) routeInvite(conversationID, userID *int64) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	if err := s.sendInviteToWebitelUser(&otherChannels[0].DomainID, conversationID, userID); err != nil {
		return err
	}
	body, _ := json.Marshal(inviteConversationEvent{
		ConversationID: *conversationID,
		UserID:         *userID,
	})
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				s.sendEventToWebitelUser(nil, item, inviteConversationEventType, body)
			}
		default:
		}
	}
	return nil
}

func (s *chatService) sendInviteToWebitelUser(domainID, conversationID, userID *int64) error {
	body, _ := json.Marshal(map[string]int64{
		"conversation_id": *conversationID,
		"user_id":         *userID,
	})
	msg := &broker.Message{
		Header: map[string]string{},
		Body:   body,
	}
	if err := s.broker.Publish(fmt.Sprintf("event.%s.%v.%v", userInvitationEventType, *domainID, *userID), msg); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	return nil
}
