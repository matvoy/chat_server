package main

import (
	"context"
)

func (s *chatService) routeInvite(conversationID, userID *int64) error {
	if err := s.sendInviteToWebitelUser(conversationID, userID); err != nil {
		return err
	}
	otherChannels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
	if err != nil {
		return err
	}
	if otherChannels == nil {
		return nil
	}
	// TO DO declineInvitationToFlow??
	for _, item := range otherChannels {
		switch item.Type {
		case "webitel":
			{
				s.sendEventToWebitelUser(nil, item, nil)
			}
		default:
		}
	}
	return nil
}

func (s *chatService) sendInviteToWebitelUser(conversationID, userID *int64) error {
	return nil
}
