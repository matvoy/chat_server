package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetProfiles(ctx context.Context, limit, offset int) ([]*models.Profile, error) {
	return models.Profiles(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
