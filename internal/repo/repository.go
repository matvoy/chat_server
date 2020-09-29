package repo

import (
	"context"
	"database/sql"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// TODO TRANSFORM TO DOMAIN MODELS

type Repository interface {
	ProfileRepository
	ConversationRepository
	ChannelRepository
	ClientRepository
	InviteRepository
	MessageRepository
	WithTransaction(txFunc func(*sql.Tx) error) (err error)
	CreateConversationTx(ctx context.Context, tx boil.ContextExecutor, c *models.Conversation) error
	CreateMessageTx(ctx context.Context, tx boil.ContextExecutor, m *models.Message) error
	GetChannelByIDTx(ctx context.Context, tx boil.ContextExecutor, id int64) (*models.Channel, error)
	GetChannelsTx(
		ctx context.Context,
		tx boil.ContextExecutor,
		userID *int64,
		conversationID *int64,
		connection *string,
		internal *bool,
		exceptID *int64,
	) ([]*models.Channel, error)
	CloseChannelTx(ctx context.Context, tx boil.ContextExecutor, id int64) error
	CreateChannelTx(ctx context.Context, tx boil.ContextExecutor, c *models.Channel) error
	CloseChannelsTx(ctx context.Context, tx boil.ContextExecutor, conversationID int64) error
	DeleteInviteTx(ctx context.Context, tx boil.ContextExecutor, inviteID int64) error
}

type ProfileRepository interface {
	GetProfileByID(ctx context.Context, id int64) (*models.Profile, error)
	GetProfiles(ctx context.Context, id int64, size, page int32, fields, sort []string, profileType string, domainID int64) ([]*models.Profile, error)
	CreateProfile(ctx context.Context, p *models.Profile) error
	DeleteProfile(ctx context.Context, id int64) error
}

type ConversationRepository interface {
	CloseConversation(ctx context.Context, id int64) error
	GetConversations(ctx context.Context, id int64, size, page int32, fields, sort []string, domainID int64) ([]*Conversation, error)
	CreateConversation(ctx context.Context, c *models.Conversation) error
	GetConversationByID(ctx context.Context, id int64) (*Conversation, error)
}

type ChannelRepository interface {
	CloseChannel(ctx context.Context, id int64) error
	CloseChannels(ctx context.Context, conversationID int64) error
	GetChannels(
		ctx context.Context,
		userID *int64,
		conversationID *int64,
		connection *string,
		internal *bool,
		exceptID *int64,
	) ([]*models.Channel, error)
	CreateChannel(ctx context.Context, c *models.Channel) error
	GetChannelByID(ctx context.Context, id int64) (*models.Channel, error)
}

type ClientRepository interface {
	GetClientByID(ctx context.Context, id int64) (*models.Client, error)
	GetClientByExternalID(ctx context.Context, externalID string) (*models.Client, error)
	CreateClient(ctx context.Context, c *models.Client) error
	GetClients(ctx context.Context, limit, offset int) ([]*models.Client, error)
}

type InviteRepository interface {
	CreateInvite(ctx context.Context, m *models.Invite) error
	DeleteInvite(ctx context.Context, inviteID int64) error
	GetInviteByID(ctx context.Context, id int64) (*models.Invite, error)
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, m *models.Message) error
	GetMessages(ctx context.Context, id int64, size, page int32, fields, sort []string, conversationID int64) ([]*models.Message, error)
}
