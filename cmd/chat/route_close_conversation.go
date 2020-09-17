package main

import (
	"context"
	"encoding/json"

	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	"github.com/matvoy/chat_server/models"
)

func (s *chatService) routeCloseConversation(channel *models.Channel, cause string) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, &channel.ConversationID, nil, nil, &channel.ID)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		if !channel.Internal {
			return s.closeConversationToFlow(channel, cause)
		}
		return nil
	}
	body, _ := json.Marshal(closeConversationEvent{
		ConversationID: channel.ConversationID,
		FromChannelID:  channel.ID,
		Cause:          cause,
	})
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				s.sendEventToWebitelUser(channel, item, messageEventType, body)
			}
		case "telegram":
			{
				reqMessage := &pbentity.Message{
					Type: "text",
					Value: &pbentity.Message_TextMessage_{
						TextMessage: &pbentity.Message_TextMessage{
							Text: cause,
						},
					},
				}
				s.sendMessageToBotUser(channel, item, reqMessage)
			}
		default:
		}
	}
	return nil
}

func (s *chatService) routeCloseConversationFromFlow(conversationID *int64, cause string) error {
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, item := range otherChannels {
		switch item.Type {
		case "telegram":
			{
				reqMessage := &pbentity.Message{
					Type: "text",
					Value: &pbentity.Message_TextMessage_{
						TextMessage: &pbentity.Message_TextMessage{
							Text: cause,
						},
					},
				}
				s.sendMessageToBotUser(nil, item, reqMessage)
			}
		default:
		}
	}
	return nil
}

func (s *chatService) closeConversationToFlow(channel *models.Channel, cause string) error {
	return nil
}

// func (s *chatService) closeConversationToWebitelUser(from *models.Channel, to *models.Channel, cause string) error {
// 	return nil
// }
