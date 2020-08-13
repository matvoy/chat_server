package main

import (
	"context"

	pb "github.com/matvoy/chat_server/chat_storage/proto/storage"
	"github.com/matvoy/chat_server/chat_storage/repo"
	"github.com/matvoy/chat_server/models"
	"github.com/volatiletech/null/v8"

	"github.com/rs/zerolog"
)

type Service interface {
	ProcessMessage(ctx context.Context, req *pb.MessageRequest, res *pb.MessageResponse) error
	GetConversationBySessionID(ctx context.Context, req *pb.ConversationRequest, res *pb.ConversationResponse) error
}

type storageService struct {
	repo repo.Repository
	log  *zerolog.Logger
}

func NewStorageService(repo repo.Repository, log *zerolog.Logger) *storageService {
	return &storageService{
		repo,
		log,
	}
}

func (s *storageService) ProcessMessage(ctx context.Context, req *pb.MessageRequest, res *pb.MessageResponse) error {
	client, _ := s.repo.GetClientByExternalID(ctx, req.ExternalUserId)
	if client != nil {
		message := &models.Message{
			ClientID: null.Int64{
				client.ID,
				true,
			},
			Text: null.String{
				req.Text,
				true,
			},
			ConversationID: client.ConversationID.Int64,
		}
		return s.repo.CreateMessage(ctx, message)
	}
	conversation := &models.Conversation{
		ProfileID: int64(req.ProfileId),
		SessionID: null.String{
			req.SessionId,
			true,
		},
	}
	if err := s.repo.CreateConversation(ctx, conversation); err != nil {
		return err
	}
	client = &models.Client{
		ConversationID: null.Int64{
			conversation.ID,
			true,
		},
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
	if err := s.repo.CreateClient(ctx, client); err != nil {
		return err
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
		ConversationID: client.ConversationID.Int64,
	}
	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return err
	}
	res.Created = true
	return nil
}

func (s *storageService) GetConversationBySessionID(ctx context.Context, req *pb.ConversationRequest, res *pb.ConversationResponse) error {
	return nil
}
