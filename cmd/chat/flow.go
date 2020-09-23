package main

import (
	"context"
	"errors"

	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	pbmanager "github.com/matvoy/chat_server/api/proto/flow_manager"
	"google.golang.org/protobuf/proto"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

func (s *chatService) sendMessageToFlow(conversationID int64, message *pbentity.Message) error {
	confirmationID, err := s.chatCache.ReadConfirmation(conversationID)
	if err != nil {
		return err
	}
	if confirmationID != nil {
		s.log.Debug().
			Int64("conversation_id", conversationID).
			Str("confirmation_id", string(confirmationID)).
			Msg("send confirmed messages")
		messages := []*pbmanager.Message{
			{
				Id:   message.GetId(),
				Type: message.GetType(),
				Value: &pbmanager.Message_TextMessage_{
					TextMessage: &pbmanager.Message_TextMessage{
						Text: message.GetTextMessage().GetText(),
					},
				},
			},
		}
		message := &pbmanager.ConfirmationMessageRequest{
			ConversationId: conversationID,
			ConfirmationId: string(confirmationID),
			Messages:       messages,
		}
		nodeID, err := s.chatCache.ReadConversationNode(conversationID)
		if err != nil {
			return err
		}
		if res, err := s.flowClient.ConfirmationMessage(
			context.Background(),
			message,
			client.WithSelectOption(
				selector.WithFilter(
					FilterNodes(string(nodeID)),
				),
			),
		); err != nil || res.Error != nil {
			if res != nil {
				return errors.New(res.Error.Message)
			}
			return err
		}
		s.chatCache.DeleteConfirmation(conversationID)
		return nil
	}
	s.log.Debug().
		Int64("conversation_id", conversationID).
		Msg("cache messages for confirmation")
	cacheMessage := &pbentity.Message{
		Id:   message.GetId(),
		Type: message.GetType(),
		Value: &pbentity.Message_TextMessage_{
			TextMessage: &pbentity.Message_TextMessage{
				Text: message.GetTextMessage().GetText(),
			},
		},
	}
	messageBytes, err := proto.Marshal(cacheMessage)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if err := s.chatCache.WriteCachedMessage(conversationID, message.GetId(), messageBytes); err != nil {
		s.log.Error().Msg(err.Error())
	}
	return nil
}

func (s *chatService) initFlow(conversationID, profileID, domainID int64, message *pbentity.Message) error {
	s.log.Debug().
		Int64("conversation_id", conversationID).
		Int64("profile_id", profileID).
		Int64("domain_id", domainID).
		Msg("init conversation")
	start := &pbmanager.StartRequest{
		ConversationId: conversationID,
		ProfileId:      profileID,
		DomainId:       domainID,
		Message: &pbmanager.Message{
			Id:   message.GetId(),
			Type: message.GetType(),
			Value: &pbmanager.Message_TextMessage_{
				TextMessage: &pbmanager.Message_TextMessage{
					Text: "start", //req.GetMessage().GetTextMessage().GetText(),
				},
			},
		},
	}
	if res, err := s.flowClient.Start(
		context.Background(),
		start,
		client.WithCallWrapper(
			s.initCallWrapper(conversationID),
		),
	); err != nil ||
		res.Error != nil {
		if res != nil {
			s.log.Error().Msg(res.Error.Message)
		} else {
			s.log.Error().Msg(err.Error())
		}
		return nil
	}
	return nil
}

func (s *chatService) closeFlowConversation(conversationID int64) error {
	nodeID, err := s.chatCache.ReadConversationNode(conversationID)
	if err != nil {
		return err
	}
	if res, err := s.flowClient.Break(
		context.Background(),
		&pbmanager.BreakRequest{
			ConversationId: conversationID,
		},
		client.WithSelectOption(
			selector.WithFilter(
				FilterNodes(string(nodeID)),
			),
		),
	); err != nil {
		return err
	} else if res != nil && res.Error != nil {
		return errors.New(res.Error.Message)
	}
	s.chatCache.DeleteCachedMessages(conversationID)
	s.chatCache.DeleteConfirmation(conversationID)
	s.chatCache.DeleteConversationNode(conversationID)
	return nil
}

func (s *chatService) initCallWrapper(conversationID int64) func(client.CallFunc) client.CallFunc {
	return func(next client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			s.log.Trace().
				Str("id", node.Id).
				Str("address", node.Address).Msg("send request to node")
			err := next(ctx, node, req, rsp, opts)
			if err != nil {
				// s.log.Error().Msg(err.Error())
				return err
			}
			if err := s.chatCache.WriteConversationNode(conversationID, []byte(node.Id)); err != nil {
				// s.log.Error().Msg(err.Error())
				return err
			}
			return nil
		}
	}
}

func FilterNodes(id string) selector.Filter {
	return func(old []*registry.Service) []*registry.Service {
		var services []*registry.Service

		for _, service := range old {
			if service.Name != "workflow" {
				continue
			}

			serv := new(registry.Service)
			var nodes []*registry.Node

			for _, node := range service.Nodes {
				if node.Id == id {
					nodes = append(nodes, node)
					break
				}
			}

			// only add service if there's some nodes
			if len(nodes) > 0 {
				// copy
				*serv = *service
				serv.Nodes = nodes
				services = append(services, serv)
			}
		}

		return services
	}
}
