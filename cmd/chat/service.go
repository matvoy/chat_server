package main

import (
	"context"
	"errors"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pb "github.com/matvoy/chat_server/api/proto/chat"
	pbflow "github.com/matvoy/chat_server/api/proto/flow_client"
	cache "github.com/matvoy/chat_server/internal/chat_cache"
	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"

	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
)

type Service interface {
	GetConversationByID(ctx context.Context, req *pb.GetConversationByIDRequest, res *pb.GetConversationByIDResponse) error
	GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error

	SendMessage(ctx context.Context, req *pb.SendMessageRequest, res *pb.SendMessageResponse) error
	StartConversation(ctx context.Context, req *pb.StartConversationRequest, res *pb.StartConversationResponse) error
	CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error
	JoinConversation(ctx context.Context, req *pb.JoinConversationRequest, res *pb.JoinConversationResponse) error
	LeaveConversation(ctx context.Context, req *pb.LeaveConversationRequest, res *pb.LeaveConversationResponse) error
	InviteToConversation(ctx context.Context, req *pb.InviteToConversationRequest, res *pb.InviteToConversationResponse) error
	DeclineInvitation(ctx context.Context, req *pb.DeclineInvitationRequest, res *pb.DeclineInvitationResponse) error

	CheckSession(ctx context.Context, req *pb.CheckSessionRequest, res *pb.CheckSessionResponse) error
}

type chatService struct {
	repo       repo.Repository
	log        *zerolog.Logger
	flowClient pbflow.FlowAdapterService
	botClient  pbbot.BotService
	chatCache  cache.ChatCache
}

func NewChatService(
	repo repo.Repository,
	log *zerolog.Logger,
	flowClient pbflow.FlowAdapterService,
	botClient pbbot.BotService,
	chatCache cache.ChatCache,
) *chatService {
	return &chatService{
		repo,
		log,
		flowClient,
		botClient,
		chatCache,
	}
}

func (s *chatService) SendMessage(
	ctx context.Context,
	req *pb.SendMessageRequest,
	res *pb.SendMessageResponse,
) error {
	channel, err := s.repo.GetChannelByID(context.Background(), req.ChannelId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if channel == nil && !req.FromFlow {
		s.log.Warn().Msg("invitation not found")
		return errors.New("invitation not found")
	}
	var channelID null.Int64
	if !req.FromFlow {
		channelID = null.Int64{
			channel.ID,
			true,
		}
	}
	message := &models.Message{
		Type:           "text",
		ChannelID:      channelID,
		ConversationID: channel.ConversationID,
		Text: null.String{
			req.Message.GetTextMessage().GetText(),
			true,
		},
	}
	if err := s.repo.CreateMessage(context.Background(), message); err != nil {
		logger.Error().Msg(err.Error())
		return err
	}
	s.routeMessage(channel, message)
	return nil
}

func (s *chatService) StartConversation(
	ctx context.Context,
	req *pb.StartConversationRequest,
	res *pb.StartConversationResponse,
) error {
	conversation := &models.Conversation{
		DomainID: req.DomainId,
	}
	if err := s.repo.CreateConversation(context.Background(), conversation); err != nil {
		logger.Error().Msg(err.Error())
		return err
	}
	channel := &models.Channel{
		Type:           req.User.Type,
		ConversationID: conversation.ID,
		UserID:         req.User.UserId,
		Connection: null.String{
			req.User.Connection,
			true,
		},
		Internal: req.User.Internal,
	}
	if err := s.repo.CreateChannel(context.Background(), channel); err != nil {
		logger.Error().Msg(err.Error())
		return err
	}
	res.ConversationId = conversation.ID
	res.ChannelId = channel.ID
	return nil
}

func (s *chatService) CloseConversation(
	ctx context.Context,
	req *pb.CloseConversationRequest,
	res *pb.CloseConversationResponse,
) error {
	if err := s.repo.CloseConversation(context.Background(), req.ConversationId); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	channels, err := s.repo.GetChannels(context.Background(), nil, &req.ConversationId, nil, nil, nil)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	for _, channel := range channels {
		if err := s.repo.CloseChannel(context.Background(), channel.ID); err != nil {
			s.log.Error().Msg(err.Error())
			return err
		}
	}
	return nil
}

func (s *chatService) JoinConversation(
	ctx context.Context,
	req *pb.JoinConversationRequest,
	res *pb.JoinConversationResponse,
) error {
	invite, err := s.repo.GetInviteByID(context.Background(), req.InviteId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if invite == nil {
		s.log.Warn().Msg("invitation not found")
		return errors.New("invitation not found")
	}
	channel := &models.Channel{
		Type:           "webitel",
		Internal:       true,
		ConversationID: invite.ConversationID,
		UserID:         invite.UserID,
	}
	if err := s.repo.CreateChannel(ctx, channel); err != nil {
		logger.Error().Msg(err.Error())
		return err
	}
	res.ChannelId = channel.ID
	return nil
}

func (s *chatService) LeaveConversation(
	ctx context.Context,
	req *pb.LeaveConversationRequest,
	res *pb.LeaveConversationResponse,
) error {
	if err := s.repo.CloseChannel(context.Background(), req.ChannelId); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (s *chatService) InviteToConversation(
	ctx context.Context,
	req *pb.InviteToConversationRequest,
	res *pb.InviteToConversationResponse,
) error {
	invite := &models.Invite{
		ConversationID: req.ConversationId,
		UserID:         req.User.UserId,
	}
	if err := s.repo.CreateInvite(context.Background(), invite); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	res.InviteId = invite.ID
	return nil
}

func (s *chatService) DeclineInvitation(
	ctx context.Context,
	req *pb.DeclineInvitationRequest,
	res *pb.DeclineInvitationResponse,
) error {
	if err := s.repo.DeleteInvite(context.Background(), req.InviteId); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	return nil
}

// func (s *chatService) CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error {
// 	c, err := s.repo.GetConversationByID(context.Background(), req.ConversationId)
// 	if err != nil {
// 		s.log.Error().Msg(err.Error())
// 		return nil
// 	}
// 	s.chatCache.DeleteSession(c.SessionID.String)
// 	if err := s.repo.CloseConversation(context.Background(), req.ConversationId); err != nil {
// 		s.log.Error().Msg(err.Error())
// 	}
// 	return nil
// }

// func (s *chatService) parseSession(ctx context.Context, req *pb.ProcessMessageRequest, clientID int64) (conversationID int64, isNew bool, err error) {
// 	cachedConversationID, err := s.chatCache.ReadSession(req.SessionId)
// 	if err != nil {
// 		return
// 	}
// 	var conversation *models.Conversation
// 	if cachedConversationID != nil {
// 		conversationID, _ = strconv.ParseInt(string(cachedConversationID), 10, 64)
// 		if err = s.chatCache.WriteSession(req.SessionId, cachedConversationID); err != nil {
// 			return
// 		}
// 	} else {
// 		conversation = &models.Conversation{
// 			ProfileID: int64(req.ProfileId),
// 			SessionID: null.String{
// 				req.SessionId,
// 				true,
// 			},
// 			ClientID: null.Int64{
// 				clientID,
// 				true,
// 			},
// 		}
// 		if err = s.repo.CreateConversation(ctx, conversation); err != nil {
// 			return
// 		}
// 		isNew = true
// 		conversationID = conversation.ID

// 		if err = s.chatCache.WriteSession(req.SessionId, []byte(strconv.Itoa(int(conversationID)))); err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// func (s *chatService) createClient(ctx context.Context, req *pb.ProcessMessageRequest) (client *models.Client, err error) {
// 	client = &models.Client{
// 		ExternalID: null.String{
// 			req.ExternalUserId,
// 			true,
// 		},
// 		Name: null.String{
// 			req.Username,
// 			true,
// 		},
// 		FirstName: null.String{
// 			req.FirstName,
// 			true,
// 		},
// 		LastName: null.String{
// 			req.LastName,
// 			true,
// 		},
// 		Number: null.String{
// 			req.Number,
// 			true,
// 		},
// 	}
// 	err = s.repo.CreateClient(ctx, client)
// 	return
// }

// func (s *storageService) ProcessMessage(ctx context.Context, req *pb.ProcessMessageRequest, res *pb.ProcessMessageResponse) error {
// 	client, _ := s.repo.GetClientByExternalID(ctx, req.ExternalUserId)
// 	var conversationID int64
// 	var err error
// 	var isNew bool
// 	if client != nil {
// 		s.log.Trace().Msg("client found")
// 		conversationID, isNew, err = s.parseSession(context.Background(), req, client.ID)
// 		if err != nil {
// 			s.log.Error().Msg(err.Error())
// 			return nil
// 		}
// 		s.log.Trace().
// 			Int64("conversation_id", conversationID).
// 			Int64("client_id", client.ID).
// 			Msg("info")
// 	} else {
// 		s.log.Trace().Msg("creating new client")
// 		client, err = s.createClient(context.Background(), req)
// 		if err != nil {
// 			s.log.Error().Msg(err.Error())
// 			return nil
// 		}
// 		conversationID, _, err = s.parseSession(context.Background(), req, client.ID)
// 		isNew = true
// 		if err != nil {
// 			s.log.Error().Msg(err.Error())
// 			return nil
// 		}
// 		s.log.Trace().
// 			Int64("conversation_id", conversationID).
// 			Int64("client_id", client.ID).
// 			Msg("info")
// 	}

// 	message := &models.Message{
// 		ClientID: null.Int64{
// 			client.ID,
// 			true,
// 		},
// 		Text: null.String{
// 			req.Text,
// 			true,
// 		},
// 		ConversationID: conversationID,
// 	}
// 	if err := s.repo.CreateMessage(context.Background(), message); err != nil {
// 		s.log.Error().Msg(err.Error())
// 		return nil
// 	}

// 	if isNew {
// 		s.log.Trace().Msg("init")
// 		init := &pbflow.InitRequest{
// 			ConversationId: conversationID,
// 			ProfileId:      int64(req.GetProfileId()),
// 			DomainId:       1,
// 			Message: &pbentity.Message{
// 				Id:   message.ID,
// 				Type: "text",
// 				Value: &pbentity.Message_TextMessage_{
// 					TextMessage: &pbentity.Message_TextMessage{
// 						Text: req.Text,
// 					},
// 				},
// 			},
// 		}
// 		if res, err := s.flowClient.Init(context.Background(), init); err != nil || res.Error != nil {
// 			if res != nil {
// 				s.log.Error().Msg(res.Error.Message)
// 			} else {
// 				s.log.Error().Msg(err.Error())
// 			}
// 			return nil
// 		}
// 	} else {
// 		s.log.Trace().Msg("send to existing")
// 		sendMessage := &pbflow.SendMessageToFlowRequest{
// 			ConversationId: conversationID,
// 			Message: &pbentity.Message{
// 				Id:   message.ID,
// 				Type: "text",
// 				Value: &pbentity.Message_TextMessage_{
// 					TextMessage: &pbentity.Message_TextMessage{
// 						Text: req.Text,
// 					},
// 				},
// 			},
// 		}
// 		if res, err := s.flowClient.SendMessageToFlow(context.Background(), sendMessage); err != nil || res.Error != nil {
// 			if res != nil {
// 				s.log.Error().Msg(res.Error.Message)
// 			} else {
// 				s.log.Error().Msg(err.Error())
// 			}
// 			return nil
// 		}
// 	}

// 	res.Created = true
// 	return nil
// }

// func (s *storageService) SaveMessageFromFlow(ctx context.Context, req *pb.SaveMessageFromFlowRequest, res *pb.SaveMessageFromFlowResponse) error {
// 	message := &models.Message{
// 		Text: null.String{
// 			req.GetMessage().GetTextMessage().GetText(),
// 			true,
// 		},
// 		ConversationID: req.GetConversationId(),
// 	}
// 	if err := s.repo.CreateMessage(context.Background(), message); err != nil {
// 		s.log.Error().Msg(err.Error())
// 		res.Error = &pbentity.Error{
// 			Message: err.Error(),
// 		}
// 	}
// 	return nil
// }
