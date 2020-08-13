package repo

import (
	"context"

	"github.com/matvoy/chat_server/models"
)

// TODO TRANSFORM TO DOMAIN MODELS

type Repository interface {
	GetProfiles(ctx context.Context, limit, offset int) ([]*models.Profile, error)
	GetConversations(ctx context.Context, limit, offset int) ([]*models.Conversation, error)
	GetMessages(ctx context.Context, limit, offset int) ([]*models.Message, error)
	GetClients(ctx context.Context, limit, offset int) ([]*models.Client, error)
	GetUserConversations(ctx context.Context, limit, offset int) ([]*models.UserConversation, error)
	GetAttachments(ctx context.Context, limit, offset int) ([]*models.Attachment, error)
}
