package main

import (
	"context"
	"encoding/json"
)

func (s *chatService) routeLeaveConversation(channelID, conversationID *int64) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, channelID)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	body, _ := json.Marshal(leaveConversationEvent{
		ConversationID:  *conversationID,
		LeavedChannelID: *channelID,
	})
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				s.sendEventToWebitelUser(nil, item, joinConversationEventType, body)
			}
		default:
		}
	}
	return nil
}
