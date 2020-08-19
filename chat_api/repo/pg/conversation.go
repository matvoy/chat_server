package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetConversations(ctx context.Context, limit, offset int) ([]*models.Conversation, error) {
	return models.Conversations(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}