package main

import (
	"context"
	"encoding/json"
)

func (s *chatService) routeDeclineInvite(userID, conversationID *int64) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	body, _ := json.Marshal(declineInvitationEvent{
		ConversationID: *conversationID,
		UserID:         *userID,
	})
	// TO DO declineInvitationToFlow??
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				s.sendEventToWebitelUser(nil, item, declineInvitationEventType, body)
			}
		default:
		}
	}
	return nil
}
