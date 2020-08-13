package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetAttachments(ctx context.Context, limit, offset int) ([]*models.Attachment, error) {
	return models.Attachments(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
