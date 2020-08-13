package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (repo *PgRepository) CreateMessage(ctx context.Context, m *models.Message) error {
	if err := m.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}
