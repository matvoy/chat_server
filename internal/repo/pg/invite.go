package pg

import (
	"context"
	"database/sql"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (repo *PgRepository) GetInviteByID(ctx context.Context, id int64) (*models.Invite, error) {
	result, err := models.Invites(models.InviteWhere.ID.EQ(id)).
		One(ctx, repo.db)
	if err != nil {
		repo.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *PgRepository) GetInvites(ctx context.Context, userID int64) ([]*models.Invite, error) {
	return models.Invites(models.InviteWhere.UserID.EQ(
		null.Int64{
			Int64: userID,
			Valid: true,
		})).All(ctx, repo.db)
}

func (repo *PgRepository) CreateInvite(ctx context.Context, p *models.Invite) error {
	if err := p.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}
