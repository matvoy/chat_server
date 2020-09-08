package main

import (
	"context"
	"encoding/json"

	pb "github.com/matvoy/chat_server/cmd/chat_api/proto/chat_api"
	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"

	"github.com/rs/zerolog"
)

type Service interface {
	CreateProfile(ctx context.Context, req *pb.CreateProfileRequest, res *pb.CreateProfileResponse) error
	DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest, res *pb.DeleteProfileResponse) error
	UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest, res *pb.UpdateProfileResponse) error
	GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error
	GetProfileByID(ctx context.Context, req *pb.GetProfileByIDRequest, res *pb.GetProfileByIDResponse) error
}

type chatApiService struct {
	repo repo.Repository
	log  *zerolog.Logger
}

func NewChatApiService(repo repo.Repository, log *zerolog.Logger) Service {
	return &chatApiService{
		repo,
		log,
	}
}

func (s *chatApiService) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest, res *pb.CreateProfileResponse) error {
	return nil
}

func (s *chatApiService) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest, res *pb.DeleteProfileResponse) error {
	return nil
}

func (s *chatApiService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest, res *pb.UpdateProfileResponse) error {
	return nil
}

func (s *chatApiService) GetProfiles(ctx context.Context, req *pb.GetProfilesRequest, res *pb.GetProfilesResponse) error {
	return nil
}

func (s *chatApiService) GetProfileByID(ctx context.Context, req *pb.GetProfileByIDRequest, res *pb.GetProfileByIDResponse) error {
	return nil
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
