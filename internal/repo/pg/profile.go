package pg

import (
	"context"
	"database/sql"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetProfileByID(ctx context.Context, id int64) (*models.Profile, error) {
	result, err := models.Profiles(models.ProfileWhere.ID.EQ(id)).
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

func (repo *PgRepository) GetProfiles(ctx context.Context, profileType string, domainID int64) ([]*models.Profile, error) {
	query := make([]qm.QueryMod, 0, 2)
	if profileType != "" {
		query = append(query, models.ProfileWhere.Type.EQ(profileType))
	}
	if domainID != 0 {
		query = append(query, models.ProfileWhere.DomainID.EQ(domainID))
	}
	return models.Profiles(query...).All(ctx, repo.db)
}

func (repo *PgRepository) CreateProfile(ctx context.Context, p *models.Profile) error {
	if err := p.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) DeleteProfile(ctx context.Context, id int64) error {
	_, err := models.Profiles(models.ProfileWhere.ID.EQ(id)).DeleteAll(ctx, repo.db)
	return err
}
