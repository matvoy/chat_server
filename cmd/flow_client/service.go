package main

import (
	"context"

	pbstorage "github.com/matvoy/chat_server/api/proto/chat_storage"
	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	pb "github.com/matvoy/chat_server/api/proto/flow_client"
	pbmanager "github.com/matvoy/chat_server/api/proto/flow_manager"
	pbtelegram "github.com/matvoy/chat_server/api/proto/telegram_bot"
	cache "github.com/matvoy/chat_server/internal/chat_cache"

	proto "github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
)

type FlowClient interface {
	SendMessage(ctx context.Context, req *pb.SendMessageRequest, res *pb.SendMessageResponse) error
	WaitMessage(ctx context.Context, req *pb.WaitMessageRequest, res *pb.WaitMessageResponse) error
	CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error
}

type FlowAdapter interface {
	Init(ctx context.Context, req *pb.InitRequest, res *pb.InitResponse) error
	SendMessageToFlow(ctx context.Context, req *pb.SendMessageToFlowRequest, res *pb.SendMessageToFlowResponse) error
}

type Service interface {
	FlowClient
	FlowAdapter
}

type flowService struct {
	log               *zerolog.Logger
	telegramClient    pbtelegram.TelegramBotService
	flowManagerClient pbmanager.FlowChatServerService
	storageClient     pbstorage.StorageService
	chatCache         cache.ChatCache
}

// var cachedMessages []*pb.Message

func NewFlowService(
	log *zerolog.Logger,
	telegramClient pbtelegram.TelegramBotService,
	flowManagerClient pbmanager.FlowChatServerService,
	storageClient pbstorage.StorageService,
	chatCache cache.ChatCache,
) *flowService {
	return &flowService{
		log,
		telegramClient,
		flowManagerClient,
		storageClient,
		chatCache,
	}
}

func (s *flowService) SendMessageToFlow(ctx context.Context, req *pb.SendMessageToFlowRequest, res *pb.SendMessageToFlowResponse) error {
	s.log.Info().Msg("confirmation")
	confirmationID, err := s.chatCache.ReadConfirmation(req.ConversationId)
	if err != nil {
		return nil
	}
	if confirmationID != nil {
		messages := []*pbmanager.Message{
			{
				Id:   req.Message.GetId(),
				Type: req.Message.GetType(),
				Value: &pbmanager.Message_TextMessage_{
					TextMessage: &pbmanager.Message_TextMessage{
						Text: req.GetMessage().GetTextMessage().GetText(),
					},
				},
			},
		}
		message := &pbmanager.ConfirmationMessageRequest{
			ConversationId: req.GetConversationId(),
			ConfirmationId: string(confirmationID),
			Messages:       messages,
		}
		if res, err := s.flowManagerClient.ConfirmationMessage(context.Background(), message); err != nil || res.Error != nil {
			if res != nil {
				s.log.Error().Msg(res.Error.Message)
			} else {
				s.log.Error().Msg(err.Error())
			}
			return nil
		}
		s.chatCache.DeleteConfirmation(req.ConversationId)
		return nil
	}
	s.log.Info().Msg("confirmation messages sent")
	message := &pbentity.Message{
		Id:   req.Message.GetId(),
		Type: req.Message.GetType(),
		Value: &pbentity.Message_TextMessage_{
			TextMessage: &pbentity.Message_TextMessage{
				Text: req.GetMessage().GetTextMessage().GetText(),
			},
		},
	}
	messageBytes, err := proto.Marshal(message)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if err := s.chatCache.WriteCachedMessage(req.ConversationId, req.Message.GetId(), messageBytes); err != nil {
		s.log.Error().Msg(err.Error())
	}
	return nil
}

func (s *flowService) Init(ctx context.Context, req *pb.InitRequest, res *pb.InitResponse) error {
	s.log.Info().Msg("init")
	start := &pbmanager.StartRequest{
		ConversationId: req.GetConversationId(),
		ProfileId:      req.GetProfileId(),
		DomainId:       req.GetDomainId(),
		Message: &pbmanager.Message{
			Id:   req.Message.GetId(),
			Type: req.Message.GetType(),
			Value: &pbmanager.Message_TextMessage_{
				TextMessage: &pbmanager.Message_TextMessage{
					Text: req.GetMessage().GetTextMessage().GetText(),
				},
			},
		},
	}
	if res, err := s.flowManagerClient.Start(context.Background(), start); err != nil || res.Error != nil {
		if res != nil {
			s.log.Error().Msg(res.Error.Message)
		} else {
			s.log.Error().Msg(err.Error())
		}
		return nil
	}
	s.log.Info().Msg("init sent")
	return nil
}

func (s *flowService) SendMessage(ctx context.Context, req *pb.SendMessageRequest, res *pb.SendMessageResponse) error {
	conversation, err := s.storageClient.GetConversationByID(context.Background(), &pbstorage.GetConversationByIDRequest{
		ConversationId: req.GetConversationId(),
	})
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if conversation.Profile == nil {
		s.log.Error().Msg("empty profile")
		return nil
	}
	switch conversation.Profile.Type {
	case "telegram":
		{
			message := &pbtelegram.MessageFromFlowRequest{
				ProfileId:      conversation.ProfileId,
				ConversationId: req.GetConversationId(),
				SessionId:      conversation.SessionId,
				Message: &pbentity.Message{
					Type: req.Messages.GetType(),
					Value: &pbentity.Message_TextMessage_{
						TextMessage: &pbentity.Message_TextMessage{
							Text: req.GetMessages().GetTextMessage().GetText(),
						},
					},
				},
			}
			if _, err := s.telegramClient.MessageFromFlow(context.Background(), message); err != nil {
				s.log.Error().Msg(err.Error())
			}
		}
	}
	storageMessage := &pbstorage.SaveMessageFromFlowRequest{
		ConversationId: req.GetConversationId(),
		Message: &pbentity.Message{
			Type: req.Messages.GetType(),
			Value: &pbentity.Message_TextMessage_{
				TextMessage: &pbentity.Message_TextMessage{
					Text: req.GetMessages().GetTextMessage().GetText(),
				},
			},
		},
	}
	if _, err := s.storageClient.SaveMessageFromFlow(context.Background(), storageMessage); err != nil {
		s.log.Error().Msg(err.Error())
	}
	return nil
}

func (s *flowService) WaitMessage(ctx context.Context, req *pb.WaitMessageRequest, res *pb.WaitMessageResponse) error {
	cachedMessages, err := s.chatCache.ReadCachedMessages(req.GetConversationId())
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if cachedMessages != nil {
		messages := make([]*pbentity.Message, 0, len(cachedMessages))
		var tmp *pbentity.Message
		var err error
		s.log.Info().Msg("send cached messages")
		for _, m := range cachedMessages {
			err = proto.Unmarshal(m.Value, tmp)
			if err != nil {
				s.log.Error().Msg(err.Error())
				return nil
			}
			messages = append(messages, tmp)
			s.chatCache.DeleteCachedMessage(m.Key)
		}
		res.Messages = messages
		s.chatCache.DeleteConfirmation(req.GetConversationId())
		res.TimeoutSec = int64(timeout)
		return nil
	}
	if err := s.chatCache.WriteConfirmation(req.GetConversationId(), []byte(req.ConfirmationId)); err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.TimeoutSec = int64(timeout)
	return nil
}

func (s *flowService) CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error {
	s.chatCache.DeleteCachedMessages(req.GetConversationId())
	s.chatCache.DeleteConfirmation(req.GetConversationId())
	if _, err := s.storageClient.CloseConversation(
		context.Background(),
		&pbstorage.CloseConversationRequest{
			ConversationId: req.ConversationId,
		}); err != nil {
		s.log.Error().Msg(err.Error())
	}
	s.log.Info().Msg("close conversation sent")
	return nil
}