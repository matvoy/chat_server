package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetClients(ctx context.Context, limit, offset int) ([]*models.Client, error) {
	return models.Clients(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
