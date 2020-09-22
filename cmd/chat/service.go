package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pb "github.com/matvoy/chat_server/api/proto/chat"
	pbflow "github.com/matvoy/chat_server/api/proto/flow_client"
	cache "github.com/matvoy/chat_server/internal/chat_cache"
	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"
	"github.com/micro/go-micro/v2/broker"

	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
)

type Service interface {
	GetConversationByID(ctx context.Context, req *pb.GetConversationByIDRequest, res *pb.GetConversationByIDResponse) error
	GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error
	CreateProfile(ctx context.Context, req *pb.CreateProfileRequest, res *pb.CreateProfileResponse) error
	UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest, res *pb.UpdateProfileResponse) error
	DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest, res *pb.DeleteProfileResponse) error

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
	broker     broker.Broker
}

func NewChatService(
	repo repo.Repository,
	log *zerolog.Logger,
	flowClient pbflow.FlowAdapterService,
	botClient pbbot.BotService,
	chatCache cache.ChatCache,
	broker broker.Broker,
) *chatService {
	return &chatService{
		repo,
		log,
		flowClient,
		botClient,
		chatCache,
		broker,
	}
}

func (s *chatService) SendMessage(
	ctx context.Context,
	req *pb.SendMessageRequest,
	res *pb.SendMessageResponse,
) error {
	s.log.Trace().
		Int64("channel_id", req.GetChannelId()).
		Int64("conversation_id", req.GetConversationId()).
		Bool("from_flow", req.GetFromFlow()).
		Msg("send message")
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
		if err := s.routeMessageFromFlow(&req.ConversationId, req.Message); err != nil {
			logger.Error().Msg(err.Error())
			return err
		}
		return nil
	}

	channel, err := s.repo.GetChannelByID(context.Background(), req.ChannelId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if channel == nil {
		s.log.Warn().Msg("channel not found")
		return err //errors.New("channel not found")
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
	if err := s.routeMessage(channel, message); err != nil {
		logger.Warn().Msg(err.Error())
		return err
	}
	return nil
}

func (s *chatService) StartConversation(
	ctx context.Context,
	req *pb.StartConversationRequest,
	res *pb.StartConversationResponse,
) error {
	s.log.Trace().
		Int64("domain_id", req.GetDomainId()).
		Str("user.connection", req.GetUser().GetConnection()).
		Str("user.type", req.GetUser().GetType()).
		Int64("user.id", req.GetUser().GetUserId()).
		Bool("user.internal", req.GetUser().GetInternal()).
		Msg("start conversation")
	channel := &models.Channel{
		Type: req.User.Type,
		// ConversationID: conversation.ID,
		UserID: req.User.UserId,
		Connection: null.String{
			req.User.Connection,
			true,
		},
		Internal: req.User.Internal,
		DomainID: req.DomainId,
	}
	conversation := &models.Conversation{
		DomainID: req.DomainId,
	}
	if err := s.repo.WithTransaction(func(tx *sql.Tx) error {
		if err := s.repo.CreateConversationTx(context.Background(), tx, conversation); err != nil {
			return err
		}
		channel.ConversationID = conversation.ID
		if err := s.repo.CreateChannelTx(context.Background(), tx, channel); err != nil {
			return err
		}
		res.ConversationId = conversation.ID
		res.ChannelId = channel.ID
		return nil
	}); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if !req.User.Internal {
		profileID, err := strconv.ParseInt(req.User.Connection, 10, 64)
		if err != nil {
			return err
		}
		init := &pbflow.InitRequest{
			ConversationId: conversation.ID,
			ProfileId:      profileID,
			DomainId:       req.DomainId,
		}
		if _, err := s.flowClient.Init(context.Background(), init); err != nil {
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
	s.log.Trace().
		Int64("conversation_id", req.GetConversationId()).
		Str("cause", req.GetCause()).
		Int64("closer_channel_id", req.GetCloserChannelId()).
		Msg("close conversation")
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
	if err := s.routeCloseConversation(closerChannel, req.Cause); err != nil {
		logger.Warn().Msg(err.Error())
		return err
	}
	return s.closeConversation(&req.ConversationId)
}

func (s *chatService) JoinConversation(
	ctx context.Context,
	req *pb.JoinConversationRequest,
	res *pb.JoinConversationResponse,
) error {
	s.log.Trace().
		Int64("invite_id", req.GetInviteId()).
		Msg("join conversation")
	invite, err := s.repo.GetInviteByID(context.Background(), req.InviteId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if invite == nil {
		s.log.Warn().Msg("invitation not found")
		return err //errors.New("invitation not found")
	}
	channel := &models.Channel{
		Type:           "webitel",
		Internal:       true,
		ConversationID: invite.ConversationID,
		UserID:         invite.UserID,
	}
	if err := s.repo.WithTransaction(func(tx *sql.Tx) error {
		if err := s.repo.CreateChannelTx(ctx, tx, channel); err != nil {
			return err
		}
		if err := s.repo.DeleteInviteTx(context.Background(), tx, req.GetInviteId()); err != nil {
			return err
		}
		res.ChannelId = channel.ID
		return nil
	}); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if err := s.routeJoinConversation(&channel.ID, &invite.ConversationID); err != nil {
		s.log.Warn().Msg(err.Error())
	}
	return nil
}

func (s *chatService) LeaveConversation(
	ctx context.Context,
	req *pb.LeaveConversationRequest,
	res *pb.LeaveConversationResponse,
) error {
	s.log.Trace().
		Int64("channel_id", req.GetChannelId()).
		Int64("conversation_id", req.GetConversationId()).
		Msg("leave conversation")
	if err := s.repo.CloseChannel(context.Background(), req.ChannelId); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if err := s.routeLeaveConversation(&req.ChannelId, &req.ConversationId); err != nil {
		s.log.Warn().Msg(err.Error())
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
	s.log.Trace().
		Str("user.connection", req.GetUser().GetConnection()).
		Str("user.type", req.GetUser().GetType()).
		Int64("user.id", req.GetUser().GetUserId()).
		Bool("user.internal", req.GetUser().GetInternal()).
		Int64("conversation_id", req.GetConversationId()).
		Msg("invite to conversation")
	invite := &models.Invite{
		ConversationID: req.ConversationId,
		UserID:         req.User.UserId,
	}
	if err := s.repo.CreateInvite(context.Background(), invite); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if err := s.routeInvite(&req.ConversationId, &req.User.UserId); err != nil {
		s.log.Warn().Msg(err.Error())
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
	s.log.Trace().
		Int64("invite_id", req.GetInviteId()).
		Int64("conversation_id", req.GetConversationId()).
		Int64("user_id", req.GetUserId()).
		Msg("decline invitation")
	if err := s.repo.DeleteInvite(context.Background(), req.InviteId); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if err := s.routeDeclineInvite(&req.UserId, &req.ConversationId); err != nil {
		s.log.Warn().Msg(err.Error())
		return err
	}
	return nil
}

func (s *chatService) CreateProfile(
	ctx context.Context,
	req *pb.CreateProfileRequest,
	res *pb.CreateProfileResponse) error {
	s.log.Trace().
		Str("name", req.GetItem().GetName()).
		Str("type", req.GetItem().GetType()).
		Int64("domain_id", req.GetItem().GetDomainId()).
		Int64("schema_id", req.GetItem().GetSchemaId()).
		Str("variables", fmt.Sprintf("%v", req.GetItem().GetVariables())).
		Msg("create profile")
	result, err := transformProfileToRepoModel(req.Item)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	if err := s.repo.CreateProfile(context.Background(), result); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	req.Item.Id = result.ID
	res.Item = req.Item

	addProfileReq := &pbbot.AddProfileRequest{
		Profile: res.Item,
	}
	if _, err := s.botClient.AddProfile(context.Background(), addProfileReq); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (s *chatService) DeleteProfile(
	ctx context.Context,
	req *pb.DeleteProfileRequest,
	res *pb.DeleteProfileResponse) error {
	s.log.Trace().
		Int64("profile_id", req.GetId()).
		Msg("delete profile")
	if err := s.repo.DeleteProfile(context.Background(), req.Id); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	deleteProfileReq := &pbbot.DeleteProfileRequest{
		Id: req.Id,
	}
	if _, err := s.botClient.DeleteProfile(context.Background(), deleteProfileReq); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (s *chatService) UpdateProfile(
	ctx context.Context,
	req *pb.UpdateProfileRequest,
	res *pb.UpdateProfileResponse) error {
	return nil
}

func (s *chatService) closeConversation(conversationID *int64) error {
	if err := s.repo.WithTransaction(func(tx *sql.Tx) error {
		if err := s.repo.CloseConversation(context.Background(), *conversationID); err != nil {
			return err
		}
		if err := s.repo.CloseChannels(context.Background(), *conversationID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	return nil
}
