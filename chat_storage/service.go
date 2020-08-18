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
	var conversation *models.Conversation
	var err error

	if client != nil {
		s.log.Trace().Int64("client_id", client.ID).Msg("client found")
		if req.IsNew {
			s.log.Trace().Msg("creating new conversation")
			if err = s.repo.CloseConversation(ctx, req.SessionId); err != nil {
				s.log.Err(err)
				return err
			}
			conversation = &models.Conversation{
				ProfileID: int64(req.ProfileId),
				SessionID: null.String{
					req.SessionId,
					true,
				},
				ClientID: null.Int64{
					client.ID,
					true,
				},
			}
			if err = s.repo.CreateConversation(ctx, conversation); err != nil {
				s.log.Err(err)
				return err
			}
		} else {
			s.log.Trace().Msg("fetching existing conversation")
			conversation, err = s.repo.GetConversationBySessionID(ctx, req.SessionId)
			if err != nil {
				s.log.Err(err)
				return err
			}
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
			ConversationID: conversation.ID,
		}
		if err = s.repo.CreateMessage(ctx, message); err != nil {
			s.log.Err(err)
			return err
		}
		s.log.Trace().Msg("message processed")

		return nil
	}
	s.log.Trace().Msg("creating new client")
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
	if err := s.repo.CreateClient(ctx, client); err != nil {
		s.log.Err(err)
		return err
	}
	conversation = &models.Conversation{
		ProfileID: int64(req.ProfileId),
		SessionID: null.String{
			req.SessionId,
			true,
		},
		ClientID: null.Int64{
			client.ID,
			true,
		},
	}
	s.log.Trace().Msg("creating new conversation")
	if err := s.repo.CreateConversation(ctx, conversation); err != nil {
		s.log.Err(err)
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
		ConversationID: conversation.ID,
	}
	if err := s.repo.CreateMessage(ctx, message); err != nil {
		s.log.Err(err)
		return err
	}
	s.log.Trace().Msg("message processed")
	res.Created = true
	return nil
}

func (s *storageService) GetConversationBySessionID(ctx context.Context, req *pb.ConversationRequest, res *pb.ConversationResponse) error {
	return nil
}
