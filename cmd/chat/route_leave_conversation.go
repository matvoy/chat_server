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
				if err := s.sendEventToWebitelUser(nil, item, joinConversationEventType, body); err != nil {
					s.log.Warn().
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
