package main

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	pb "github.com/matvoy/chat_server/chat_storage/proto/storage"
	"github.com/matvoy/chat_server/chat_storage/repo"
	pbflow "github.com/matvoy/chat_server/flow_client/proto/flow_client"
	"github.com/matvoy/chat_server/models"

	"github.com/micro/go-micro/v2/store"
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
	redisStore store.Store
	flowClient pbflow.FlowAdapterService
}

func NewStorageService(repo repo.Repository, log *zerolog.Logger, redisStore store.Store, flowClient pbflow.FlowAdapterService) *storageService {
	return &storageService{
		repo,
		log,
		redisStore,
		flowClient,
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
			Message: &pbflow.Message{
				Id:   message.ID,
				Type: "text",
				Value: &pbflow.Message_TextMessage_{
					TextMessage: &pbflow.Message_TextMessage{
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
			Message: &pbflow.Message{
				Id:   message.ID,
				Type: "text",
				Value: &pbflow.Message_TextMessage_{
					TextMessage: &pbflow.Message_TextMessage{
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
	res.Profile = &pb.Profile{
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
	sessionKey := "session_id:" + c.SessionID.String
	s.redisStore.Delete(sessionKey)
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
	res.Profile = &pb.Profile{
		Id:        profile.ID,
		Name:      profile.Name,
		Type:      profile.Type,
		DomainId:  profile.DomainID,
		Variables: variables,
	}
	return nil
}

func (s *storageService) GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error {
	profiles, err := s.repo.GetProfiles(context.Background())
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

func (s *storageService) parseSession(ctx context.Context, req *pb.ProcessMessageRequest, clientID int64) (conversationID int64, isNew bool, err error) {
	sessionKey := "session_id:" + req.SessionId
	session, err := s.redisStore.Read(sessionKey)
	if err != nil && err.Error() != "not found" {
		return
	}
	var conversation *models.Conversation
	if session != nil && len(session) > 0 && session[0] != nil {
		conversationID, _ = strconv.ParseInt(string(session[0].Value), 10, 64)
		if err = s.redisStore.Write(&store.Record{
			Key:    sessionKey,
			Value:  session[0].Value,
			Expiry: time.Hour * time.Duration(24),
		}); err != nil {
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
		if err = s.redisStore.Write(&store.Record{
			Key:    sessionKey,
			Value:  []byte(strconv.Itoa(int(conversationID))),
			Expiry: time.Hour * time.Duration(24),
		}); err != nil {
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

func transformProfileFromRepoModel(profile *models.Profile) (*pb.Profile, error) {
	variableBytes, err := profile.Variables.MarshalJSON()
	variables := make(map[string]string)
	err = json.Unmarshal(variableBytes, &variables)
	if err != nil {
		return nil, err
	}
	result := &pb.Profile{
		Id:        profile.ID,
		Name:      profile.Name,
		Type:      profile.Type,
		DomainId:  profile.DomainID,
		Variables: variables,
	}
	return result, nil
}

func transformProfilesFromRepoModel(profiles []*models.Profile) ([]*pb.Profile, error) {
	result := make([]*pb.Profile, 0, len(profiles))
	var tmp *pb.Profile
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
