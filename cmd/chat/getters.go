package main

import (
	"context"
	"encoding/json"

	pb "github.com/matvoy/chat_server/api/proto/chat"
	pbentity "github.com/matvoy/chat_server/api/proto/entity"
	"github.com/matvoy/chat_server/models"
)

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

// func (s *storageService) GetProfileByID(ctx context.Context, req *pb.GetProfileByIDRequest, res *pb.GetProfileByIDResponse) error {
// 	profile, err := s.repo.GetProfileByID(context.Background(), req.ProfileId)
// 	if err != nil {
// 		s.log.Error().Msg(err.Error())
// 		return nil
// 	}
// 	variableBytes, err := profile.Variables.MarshalJSON()
// 	variables := make(map[string]string)
// 	err = json.Unmarshal(variableBytes, &variables)
// 	if err != nil {
// 		s.log.Error().Msg(err.Error())
// 		return nil
// 	}
// 	res.Profile = &pbentity.Profile{
// 		Id:        profile.ID,
// 		Name:      profile.Name,
// 		Type:      profile.Type,
// 		DomainId:  profile.DomainID,
// 		Variables: variables,
// 	}
// 	return nil
// }
