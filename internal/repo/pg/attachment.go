package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) CreateAttachment(ctx context.Context, a *models.Attachment) error {
	if err := a.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) GetAttachments(ctx context.Context, limit, offset int) ([]*models.Attachment, error) {
	return models.Attachments(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
