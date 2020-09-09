package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) CreateUserConversation(ctx context.Context, uc *models.UserConversation) error {
	if err := uc.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) GetUserConversations(ctx context.Context, limit, offset int) ([]*models.UserConversation, error) {
	return models.UserConversations(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
