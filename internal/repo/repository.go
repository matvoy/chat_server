package repo

import (
	"context"

	"github.com/matvoy/chat_server/models"
)

// TODO TRANSFORM TO DOMAIN MODELS

type Repository interface {
	GetProfileByID(ctx context.Context, id int64) (*models.Profile, error)
	GetProfiles(ctx context.Context, profileType string) ([]*models.Profile, error)
	GetConversationBySessionID(ctx context.Context, sessionID string) (*models.Conversation, error)
	GetConversationByID(ctx context.Context, id int64) (*models.Conversation, error)
	GetClientByExternalID(ctx context.Context, externalID string) (*models.Client, error)

	CreateProfile(ctx context.Context, p *models.Profile) error
	CreateConversation(ctx context.Context, c *models.Conversation) error
	CreateMessage(ctx context.Context, m *models.Message) error
	CreateClient(ctx context.Context, c *models.Client) error
	CreateUserConversation(ctx context.Context, uc *models.UserConversation) error
	CreateAttachment(ctx context.Context, a *models.Attachment) error

	CloseConversation(ctx context.Context, id int64) error

	GetConversations(ctx context.Context, limit, offset int) ([]*models.Conversation, error)
	GetMessages(ctx context.Context, limit, offset int) ([]*models.Message, error)
	GetClients(ctx context.Context, limit, offset int) ([]*models.Client, error)
	GetUserConversations(ctx context.Context, limit, offset int) ([]*models.UserConversation, error)
	GetAttachments(ctx context.Context, limit, offset int) ([]*models.Attachment, error)
}
