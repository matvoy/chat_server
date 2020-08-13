package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (repo *PgRepository) CreateAttachment(ctx context.Context, a *models.Attachment) error {
	if err := a.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}
