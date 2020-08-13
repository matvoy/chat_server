package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetMessages(ctx context.Context, limit, offset int) ([]*models.Message, error) {
	return models.Messages(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
