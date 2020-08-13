package pg

import (
	"context"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (repo *PgRepository) CreateUserConversation(ctx context.Context, uc *models.UserConversation) error {
	if err := uc.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}
