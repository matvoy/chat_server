package pg

import (
	"context"
	"database/sql"

	"github.com/matvoy/chat_server/models"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (repo *PgRepository) GetInviteByID(ctx context.Context, id string) (*models.Invite, error) {
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
	return models.Invites(models.InviteWhere.UserID.EQ(userID)).All(ctx, repo.db)
}

func (repo *PgRepository) CreateInvite(ctx context.Context, m *models.Invite) error {
	m.ID = uuid.New().String()
	if err := m.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) DeleteInvite(ctx context.Context, inviteID string) error {
	_, err := models.Invites(models.InviteWhere.ID.EQ(inviteID)).DeleteAll(ctx, repo.db)
	return err
}
