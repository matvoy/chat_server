package pg

import (
	"context"
	"database/sql"
	"strings"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetClientByExternalID(ctx context.Context, externalID string) (*models.Client, error) {
	result, err := models.Clients(qm.Where("LOWER(external_id) like ?", strings.ToLower(externalID))).
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

func (repo *PgRepository) CreateClient(ctx context.Context, c *models.Client) error {
	if err := c.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}
