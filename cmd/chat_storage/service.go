package main

import (
	"context"
	"encoding/json"
	"strconv"

	pb "github.com/matvoy/chat_server/api/proto/chat_storage"
	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	pbflow "github.com/matvoy/chat_server/api/proto/flow_client"
	cache "github.com/matvoy/chat_server/internal/chat_cache"
	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"

	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
)

type Service interface {
	ProcessMessage(ctx context.Context, req *pb.ProcessMessageRequest, res *pb.ProcessMessageResponse) error
	GetConversationByID(ctx context.Context, req *pb.GetConversationByIDRequest, res *pb.GetConversationByIDResponse) error
	CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error
	GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error
}

type storageService struct {
	repo       repo.Repository
	log        *zerolog.Logger
	flowClient pbflow.FlowAdapterService
	chatCache  cache.ChatCache
}

func NewStorageService(repo repo.Repository, log *zerolog.Logger, flowClient pbflow.FlowAdapterService, chatCache cache.ChatCache) *storageService {
	return &storageService{
		repo,
		log,
		flowClient,
		chatCache,
	}
}

func (s *storageService) ProcessMessage(ctx context.Context, req *pb.ProcessMessageRequest, res *pb.ProcessMessageResponse) error {
	client, _ := s.repo.GetClientByExternalID(ctx, req.ExternalUserId)
	var conversationID int64
	var err error
	var isNew bool
	if client != nil {
		s.log.Trace().Msg("client found")
		conversationID, isNew, err = s.parseSession(context.Background(), req, client.ID)
		if err != nil {
			s.log.Error().Msg(err.Error())
			return nil
		}
		s.log.Trace().
			Int64("conversation_id", conversationID).
			Int64("client_id", client.ID).
			Msg("info")
	} else {
		s.log.Trace().Msg("creating new client")
		client, err = s.createClient(context.Background(), req)
		if err != nil {
			s.log.Error().Msg(err.Error())
			return nil
		}
		conversationID, _, err = s.parseSession(context.Background(), req, client.ID)
		isNew = true
		if err != nil {
			s.log.Error().Msg(err.Error())
			return nil
		}
		s.log.Trace().
			Int64("conversation_id", conversationID).
			Int64("client_id", client.ID).
			Msg("info")
	}

	message := &models.Message{
		ClientID: null.Int64{
			client.ID,
			true,
		},
		Text: null.String{
			req.Text,
			true,
		},
		ConversationID: conversationID,
	}
	if err := s.repo.CreateMessage(context.Background(), message); err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}

	if isNew {
		s.log.Trace().Msg("init")
		init := &pbflow.InitRequest{
			ConversationId: conversationID,
			ProfileId:      int64(req.GetProfileId()),
			DomainId:       1,
			Message: &pbentity.Message{
				Id:   message.ID,
				Type: "text",
				Value: &pbentity.Message_TextMessage_{
					TextMessage: &pbentity.Message_TextMessage{
						Text: req.Text,
					},
				},
			},
		}
		if res, err := s.flowClient.Init(context.Background(), init); err != nil || res.Error != nil {
			if res != nil {
				s.log.Error().Msg(res.Error.Message)
			} else {
				s.log.Error().Msg(err.Error())
			}
			return nil
		}
	} else {
		s.log.Trace().Msg("send to existing")
		sendMessage := &pbflow.SendMessageToFlowRequest{
			ConversationId: conversationID,
			Message: &pbentity.Message{
				Id:   message.ID,
				Type: "text",
				Value: &pbentity.Message_TextMessage_{
					TextMessage: &pbentity.Message_TextMessage{
						Text: req.Text,
					},
				},
			},
		}
		if res, err := s.flowClient.SendMessageToFlow(context.Background(), sendMessage); err != nil || res.Error != nil {
			if res != nil {
				s.log.Error().Msg(res.Error.Message)
			} else {
				s.log.Error().Msg(err.Error())
			}
			return nil
		}
	}

	res.Created = true
	return nil
}

func (s *storageService) GetConversationByID(ctx context.Context, req *pb.GetConversationByIDRequest, res *pb.GetConversationByIDResponse) error {
	conversation, err := s.repo.GetConversationByID(context.Background(), req.ConversationId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Id = conversation.ID
	res.ProfileId = conversation.ProfileID
	res.SessionId = conversation.SessionID.String
	profile := conversation.R.Profile
	res.Profile = &pbentity.Profile{
		Id:       profile.ID,
		Name:     profile.Name,
		Type:     profile.Type,
		DomainId: profile.DomainID,
	}
	return nil
}

func (s *storageService) CloseConversation(ctx context.Context, req *pb.CloseConversationRequest, res *pb.CloseConversationResponse) error {
	c, err := s.repo.GetConversationByID(context.Background(), req.ConversationId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	s.chatCache.DeleteSession(c.SessionID.String)
	if err := s.repo.CloseConversation(context.Background(), req.ConversationId); err != nil {
		s.log.Error().Msg(err.Error())
	}
	return nil
}

func (s *storageService) GetProfileByID(ctx context.Context, req *pb.GetProfileByIDRequest, res *pb.GetProfileByIDResponse) error {
	profile, err := s.repo.GetProfileByID(context.Background(), req.ProfileId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	variableBytes, err := profile.Variables.MarshalJSON()
	variables := make(map[string]string)
	err = json.Unmarshal(variableBytes, &variables)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Profile = &pbentity.Profile{
		Id:        profile.ID,
		Name:      profile.Name,
		Type:      profile.Type,
		DomainId:  profile.DomainID,
		Variables: variables,
	}
	return nil
}

func (s *storageService) GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error {
	profiles, err := s.repo.GetProfiles(context.Background(), req.Type)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	result, err := transformProfilesFromRepoModel(profiles)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Profiles = result
	return nil
}

func (s *storageService) SaveMessageFromFlow(ctx context.Context, req *pb.SaveMessageFromFlowRequest, res *pb.SaveMessageFromFlowResponse) error {
	message := &models.Message{
		Text: null.String{
			req.GetMessage().GetTextMessage().GetText(),
			true,
		},
		ConversationID: req.GetConversationId(),
	}
	if err := s.repo.CreateMessage(context.Background(), message); err != nil {
		s.log.Error().Msg(err.Error())
		res.Error = &pbentity.Error{
			Message: err.Error(),
		}
	}
	return nil
}

func (s *storageService) parseSession(ctx context.Context, req *pb.ProcessMessageRequest, clientID int64) (conversationID int64, isNew bool, err error) {
	cachedConversationID, err := s.chatCache.ReadSession(req.SessionId)
	if err != nil {
		return
	}
	var conversation *models.Conversation
	if cachedConversationID != nil {
		conversationID, _ = strconv.ParseInt(string(cachedConversationID), 10, 64)
		if err = s.chatCache.WriteSession(req.SessionId, cachedConversationID); err != nil {
			return
		}
	} else {
		conversation = &models.Conversation{
			ProfileID: int64(req.ProfileId),
			SessionID: null.String{
				req.SessionId,
				true,
			},
			ClientID: null.Int64{
				clientID,
				true,
			},
		}
		if err = s.repo.CreateConversation(ctx, conversation); err != nil {
			return
		}
		isNew = true
		conversationID = conversation.ID

		if err = s.chatCache.WriteSession(req.SessionId, []byte(strconv.Itoa(int(conversationID)))); err != nil {
			return
		}
	}
	return
}

func (s *storageService) createClient(ctx context.Context, req *pb.ProcessMessageRequest) (client *models.Client, err error) {
	client = &models.Client{
		ExternalID: null.String{
			req.ExternalUserId,
			true,
		},
		Name: null.String{
			req.Username,
			true,
		},
		FirstName: null.String{
			req.FirstName,
			true,
		},
		LastName: null.String{
			req.LastName,
			true,
		},
		Number: null.String{
			req.Number,
			true,
		},
	}
	err = s.repo.CreateClient(ctx, client)
	return
}

func transformProfileFromRepoModel(profile *models.Profile) (*pbentity.Profile, error) {
	variableBytes, err := profile.Variables.MarshalJSON()
	variables := make(map[string]string)
	err = json.Unmarshal(variableBytes, &variables)
	if err != nil {
		return nil, err
	}
	result := &pbentity.Profile{
		Id:        profile.ID,
		Name:      profile.Name,
		Type:      profile.Type,
		DomainId:  profile.DomainID,
		Variables: variables,
	}
	return result, nil
}

func transformProfilesFromRepoModel(profiles []*models.Profile) ([]*pbentity.Profile, error) {
	result := make([]*pbentity.Profile, 0, len(profiles))
	var tmp *pbentity.Profile
	var err error
	for _, item := range profiles {
		tmp, err = transformProfileFromRepoModel(item)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp)
	}
	return result, nil
}
