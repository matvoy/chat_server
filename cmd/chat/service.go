package main

import (
	"context"
	"errors"
	"strconv"

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
	if req.FromFlow {
		message := &models.Message{
			Type:           "text",
			ConversationID: req.ConversationId,
			Text: null.String{
				req.Message.GetTextMessage().GetText(),
				true,
			},
		}
		if err := s.repo.CreateMessage(context.Background(), message); err != nil {
			logger.Error().Msg(err.Error())
			return err
		}
		s.routeMessageFromFlow(&req.ConversationId, req.Message)
		return nil
	}

	channel, err := s.repo.GetChannelByID(context.Background(), req.ChannelId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if channel == nil {
		s.log.Warn().Msg("channel not found")
		return errors.New("channel not found")
	}

	message := &models.Message{
		Type: "text",
		ChannelID: null.Int64{
			channel.ID,
			true,
		},
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
	if !req.User.Internal {
		profileID, err := strconv.ParseInt(req.User.Connection, 10, 64)
		if err != nil {
			s.log.Error().Msg(err.Error())
			return err
		}
		init := &pbflow.InitRequest{
			ConversationId: conversation.ID,
			ProfileId:      profileID,
			DomainId:       req.DomainId,
		}
		if res, err := s.flowClient.Init(context.Background(), init); err != nil {
			s.log.Error().Msg(res.Error.Message)
			return err
		}
	}
	return nil
}

func (s *chatService) CloseConversation(
	ctx context.Context,
	req *pb.CloseConversationRequest,
	res *pb.CloseConversationResponse,
) error {
	if req.FromFlow {
		if err := s.routeCloseConversationFromFlow(&req.ConversationId, req.Cause); err != nil {
			s.log.Error().Msg(err.Error())
			return err
		}
		return s.closeConversation(&req.ConversationId)
	}
	closerChannel, err := s.repo.GetChannelByID(context.Background(), req.CloserChannelId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	s.routeCloseConversation(closerChannel, req.Cause)
	return s.closeConversation(&req.ConversationId)
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
	if err := s.routeJoinConversation(&channel.ID, &invite.ConversationID); err != nil {
		s.log.Error().Msg(err.Error())
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
	if err := s.routeLeaveConversation(&req.ChannelId, &req.ConversationId); err != nil {
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
	// _, err := s.repo.GetChannelByID(context.Background(), req.InviterChannelId)
	// if err != nil {
	// 	s.log.Error().Msg(err.Error())
	// 	return err
	// }
	invite := &models.Invite{
		ConversationID: req.ConversationId,
		UserID:         req.User.UserId,
	}
	if err := s.repo.CreateInvite(context.Background(), invite); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if err := s.routeInvite(&req.ConversationId, &req.User.UserId); err != nil {
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
	s.routeDeclineInvite(&req.ConversationId)
	return nil
}

func (s *chatService) closeConversation(conversationID *int64) error {
	if err := s.repo.CloseConversation(context.Background(), *conversationID); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	channels, err := s.repo.GetChannels(context.Background(), nil, conversationID, nil, nil, nil)
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
