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
				if err := s.sendEventToWebitelUser(nil, item, inviteConversationEventType, body); err != nil {
					s.log.Warn().
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
		return err
	}
	return nil
}
