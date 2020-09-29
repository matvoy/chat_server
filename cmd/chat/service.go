package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pb "github.com/matvoy/chat_server/api/proto/chat"
	cache "github.com/matvoy/chat_server/internal/chat_cache"
	"github.com/matvoy/chat_server/internal/flow"
	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"

	"github.com/micro/go-micro/v2/broker"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"google.golang.org/protobuf/proto"
)

type Service interface {
	GetConversationByID(ctx context.Context, req *pb.GetConversationByIDRequest, res *pb.GetConversationByIDResponse) error
	GetConversations(ctx context.Context, req *pb.GetConversationsRequest, res *pb.GetConversationsResponse) error
	GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error
	CreateProfile(ctx context.Context, req *pb.CreateProfileRequest, res *pb.CreateProfileResponse) error
	UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest, res *pb.UpdateProfileResponse) error
	DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest, res *pb.DeleteProfileResponse) error
	GetHistoryMessages(ctx context.Context, req *pb.GetHistoryMessagesRequest, res *pb.GetHistoryMessagesResponse) error

	SendMessage(ctx context.Context, req *pb.SendMessageRequest, res *pb.SendMessageResponse) error
	StartConversation(ctx context.Context, req *pb.StartConversationRequest, res *pb.StartConversationResponse) error
	CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error
	JoinConversation(ctx context.Context, req *pb.JoinConversationRequest, res *pb.JoinConversationResponse) error
	LeaveConversation(ctx context.Context, req *pb.LeaveConversationRequest, res *pb.LeaveConversationResponse) error
	InviteToConversation(ctx context.Context, req *pb.InviteToConversationRequest, res *pb.InviteToConversationResponse) error
	DeclineInvitation(ctx context.Context, req *pb.DeclineInvitationRequest, res *pb.DeclineInvitationResponse) error
	WaitMessage(ctx context.Context, req *pb.WaitMessageRequest, res *pb.WaitMessageResponse) error
	CheckSession(ctx context.Context, req *pb.CheckSessionRequest, res *pb.CheckSessionResponse) error
}

type chatService struct {
	repo       repo.Repository
	log        *zerolog.Logger
	flowClient flow.Client
	botClient  pbbot.BotService
	chatCache  cache.ChatCache
	broker     broker.Broker
}

func NewChatService(
	repo repo.Repository,
	log *zerolog.Logger,
	flowClient flow.Client,
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
			if err := s.flowClient.CloseConversation(req.GetConversationId()); err != nil {
				s.log.Error().Msg(err.Error())
			}
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
		err = s.flowClient.Init(conversation.ID, profileID, req.DomainId, nil)
		if err != nil {
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
		s.chatCache.DeleteCachedMessages(req.GetConversationId())
		s.chatCache.DeleteConfirmation(req.GetConversationId())
		s.chatCache.DeleteConversationNode(req.GetConversationId())
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

func (s *chatService) WaitMessage(ctx context.Context, req *pb.WaitMessageRequest, res *pb.WaitMessageResponse) error {
	s.log.Debug().
		Int64("conversation_id", req.GetConversationId()).
		Str("confirmation_id", req.GetConfirmationId()).
		Msg("accept confirmation")
	cachedMessages, err := s.chatCache.ReadCachedMessages(req.GetConversationId())
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if cachedMessages != nil {
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

func (s *chatService) CheckSession(ctx context.Context, req *pb.CheckSessionRequest, res *pb.CheckSessionResponse) error {
	s.log.Trace().
		Str("external_id", req.GetExternalId()).
		Int64("profile_id", req.GetProfileId()).
		Msg("check session")
	client, err := s.repo.GetClientByExternalID(context.Background(), req.ExternalId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if client == nil {
		client, err = s.createClient(context.Background(), req)
		if err != nil {
			s.log.Error().Msg(err.Error())
			return nil
		}
		res.ClientId = client.ID
		res.Exists = false
		return nil
	}
	profileStr := strconv.Itoa(int(req.ProfileId))
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	channel, err := s.repo.GetChannels(context.Background(), &client.ID, nil, &profileStr, func() *bool { b := false; return &b }(), nil)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	if len(channel) > 0 {
		res.Exists = true
		res.ChannelId = channel[0].ID
		res.ClientId = client.ID
	} else {
		res.Exists = false
		res.ClientId = client.ID
	}
	return nil
}

func (s *chatService) GetConversationByID(ctx context.Context, req *pb.GetConversationByIDRequest, res *pb.GetConversationByIDResponse) error {
	s.log.Trace().
		Int64("conversation_id", req.GetId()).
		Msg("get conversation by id")
	conversation, err := s.repo.GetConversationByID(context.Background(), req.GetId())
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Item = transformConversationFromRepoModel(conversation)
	return nil
}

func (s *chatService) GetConversations(ctx context.Context, req *pb.GetConversationsRequest, res *pb.GetConversationsResponse) error {
	s.log.Trace().
		Int64("conversation_id", req.GetId()).
		Msg("get conversations")
	conversations, err := s.repo.GetConversations(
		context.Background(),
		req.GetId(),
		req.GetSize(),
		req.GetPage(),
		req.GetFields(),
		req.GetSort(),
		req.GetDomainId(),
	)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Items = transformConversationsFromRepoModel(conversations)
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
	profile, err := s.repo.GetProfileByID(context.Background(), req.GetId())
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	} else if profile == nil {
		return errors.New("profile not found")
	}
	if err := s.repo.DeleteProfile(context.Background(), req.GetId()); err != nil {
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
	res.Item, err = transformProfileFromRepoModel(profile)
	return err
}

func (s *chatService) UpdateProfile(
	ctx context.Context,
	req *pb.UpdateProfileRequest,
	res *pb.UpdateProfileResponse) error {
	return nil
}

func (s *chatService) GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error {
	s.log.Trace().
		Str("type", req.GetType()).
		Int64("domain_id", req.GetDomainId()).
		Msg("get profiles")
	profiles, err := s.repo.GetProfiles(
		context.Background(),
		req.GetId(),
		req.GetSize(),
		req.GetPage(),
		req.GetFields(),
		req.GetSort(),
		req.Type,
		req.DomainId,
	)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	result, err := transformProfilesFromRepoModel(profiles)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Items = result
	return nil
}

func (s *chatService) GetProfileByID(ctx context.Context, req *pb.GetProfileByIDRequest, res *pb.GetProfileByIDResponse) error {
	s.log.Trace().
		Int64("profile_id", req.GetId()).
		Msg("get profile by id")
	profile, err := s.repo.GetProfileByID(context.Background(), req.Id)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	variableBytes, err := profile.Variables.MarshalJSON()
	variables := make(map[string]string)
	err = json.Unmarshal(variableBytes, &variables)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	result, err := transformProfileFromRepoModel(profile)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return err
	}
	res.Item = result
	return nil
}

func (s *chatService) GetHistoryMessages(ctx context.Context, req *pb.GetHistoryMessagesRequest, res *pb.GetHistoryMessagesResponse) error {
	s.log.Trace().
		Int64("conversation_id", req.GetConversationId()).
		Msg("get history")
	messages, err := s.repo.GetMessages(
		context.Background(),
		req.GetId(),
		req.GetSize(),
		req.GetPage(),
		req.GetFields(),
		req.GetSort(),
		req.GetConversationId(),
	)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Items = transformMessagesFromRepoModel(messages)
	return nil
}
