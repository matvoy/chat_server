package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) CreateMessage(ctx context.Context, m *models.Message) error {
	if err := m.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) GetMessages(ctx context.Context, limit, offset int) ([]*models.Message, error) {
	return models.Messages(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
