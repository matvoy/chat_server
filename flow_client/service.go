package main

import (
	"context"
	"fmt"
	"time"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/flow_client/proto/flow_client"
	pbmanager "github.com/matvoy/chat_server/flow_client/proto/flow_manager"
	pbtelegram "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	proto "github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/v2/store"
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
	redisStore        store.Store
}

// var cachedMessages []*pb.Message

func NewFlowService(
	log *zerolog.Logger,
	telegramClient pbtelegram.TelegramBotService,
	flowManagerClient pbmanager.FlowChatServerService,
	redisStore store.Store,
	storageClient pbstorage.StorageService,
) *flowService {
	return &flowService{
		log,
		telegramClient,
		flowManagerClient,
		storageClient,
		redisStore,
	}
}

func (s *flowService) SendMessageToFlow(ctx context.Context, req *pb.SendMessageToFlowRequest, res *pb.SendMessageToFlowResponse) error {
	s.log.Info().Msg("confirmation")
	confirmationKey := fmt.Sprintf("confirmations:%v", req.ConversationId)
	confirmationID, err := s.redisStore.Read(confirmationKey)
	if err != nil && err.Error() != "not found" {
		return nil
	}
	if confirmationID != nil && len(confirmationID) > 0 && confirmationID[0] != nil {
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
			ConfirmationId: string(confirmationID[0].Value),
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
		s.redisStore.Delete(confirmationKey)
		return nil
	}
	s.log.Info().Msg("confirmation messages sent")
	messagesKey := fmt.Sprintf("cached_messages:%v", req.ConversationId)
	message := &pb.Message{
		Id:   req.Message.GetId(),
		Type: req.Message.GetType(),
		Value: &pb.Message_TextMessage_{
			TextMessage: &pb.Message_TextMessage{
				Text: req.GetMessage().GetTextMessage().GetText(),
			},
		},
	}
	messageBytes, err := proto.Marshal(message)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if err := s.redisStore.Write(&store.Record{
		Key:    fmt.Sprintf("%s:%v", messagesKey, req.Message.GetId()),
		Value:  messageBytes,
		Expiry: time.Hour * time.Duration(24),
	}); err != nil {
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
				Message: &pbtelegram.Message{
					Type: req.Messages.GetType(),
					Value: &pbtelegram.Message_TextMessage_{
						TextMessage: &pbtelegram.Message_TextMessage{
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

	return nil
}

func (s *flowService) WaitMessage(ctx context.Context, req *pb.WaitMessageRequest, res *pb.WaitMessageResponse) error {
	messagesKey := fmt.Sprintf("cached_messages:%v", req.GetConversationId())
	confirmationKey := fmt.Sprintf("confirmations:%v", req.GetConversationId())
	cachedMessages, err := s.redisStore.Read(messagesKey)
	if err != nil && err.Error() != "not found" {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if len(cachedMessages) > 0 {
		messages := make([]*pb.Message, 0, len(cachedMessages))
		var tmp *pb.Message
		var err error
		s.log.Info().Msg("send cached messages")
		for _, m := range cachedMessages {
			err = proto.Unmarshal(m.Value, tmp)
			if err != nil {
				s.log.Error().Msg(err.Error())
				return nil
			}
			messages = append(messages, tmp)
			s.redisStore.Delete(m.Key)
		}
		res.Messages = messages
		s.redisStore.Delete(confirmationKey)
		res.TimeoutSec = int64(timeout)
		return nil
	}
	if err := s.redisStore.Write(&store.Record{
		Key:    confirmationKey,
		Value:  []byte(req.ConfirmationId),
		Expiry: time.Second * time.Duration(timeout),
	}); err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.TimeoutSec = int64(timeout)
	return nil
}

func (s *flowService) CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error {
	messagesKey := fmt.Sprintf("cached_messages:%v", req.GetConversationId())
	confirmationKey := fmt.Sprintf("confirmations:%v", req.GetConversationId())
	cachedMessages, _ := s.redisStore.Read(messagesKey)
	for _, m := range cachedMessages {
		s.redisStore.Delete(m.Key)
	}
	s.redisStore.Delete(confirmationKey)
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
