package main

import (
	"context"
	"encoding/json"

	"github.com/matvoy/chat_server/pkg/events"
)

func (s *chatService) routeJoinConversation(channelID, conversationID *int64) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, channelID)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	body, _ := json.Marshal(events.JoinConversationEvent{
		ConversationID:  *conversationID,
		JoinedChannelID: *channelID,
	})
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				if err := s.sendEventToWebitelUser(nil, item, events.JoinConversationEventType, body); err != nil {
					s.log.Warn().
						Int64("channel_id", item.ID).
						Bool("internal", item.Internal).
						Int64("user_id", item.UserID).
						Int64("conversation_id", item.ConversationID).
						Str("type", item.Type).
						Str("connection", item.Connection.String).
						Msg("failed to send join conversation event to channel")
				}
			}
		default:
		}
	}
	return nil
}
