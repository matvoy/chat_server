package main

import (
	"context"
	"encoding/json"
	"strconv"

	pb "github.com/matvoy/chat_server/api/proto/chat"
	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	"github.com/matvoy/chat_server/models"
	"github.com/volatiletech/null/v8"
)

func (s *chatService) CheckSession(ctx context.Context, req *pb.CheckSessionRequest, res *pb.CheckSessionResponse) error {
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
	conversation, err := s.repo.GetConversationByID(context.Background(), req.ConversationId)
	if err != nil {
		s.log.Error().Msg(err.Error())
		return nil
	}
	res.Id = conversation.ID
	res.DomainId = conversation.DomainID
	return nil
}

func (s *chatService) GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error {
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

func (s *chatService) GetProfileByID(ctx context.Context, req *pb.GetProfileByIDRequest, res *pb.GetProfileByIDResponse) error {
	profile, err := s.repo.GetProfileByID(context.Background(), req.ProfileId)
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
	res.Profile = result
	return nil
}

func (s *chatService) createClient(ctx context.Context, req *pb.CheckSessionRequest) (client *models.Client, err error) {
	client = &models.Client{
		ExternalID: null.String{
			req.ExternalId,
			true,
		},
		Name: null.String{
			req.Username,
			true,
		},
	}
	err = s.repo.CreateClient(ctx, client)
	return
}
