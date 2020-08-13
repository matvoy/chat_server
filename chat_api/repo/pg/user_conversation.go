package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetUserConversations(ctx context.Context, limit, offset int) ([]*models.UserConversation, error) {
	return models.UserConversations(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
