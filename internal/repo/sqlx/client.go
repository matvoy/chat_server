package sqlxrepo

import (
	"context"
	"database/sql"
	"time"
)

func (repo *sqlxRepository) GetClientByID(ctx context.Context, id int64) (*Client, error) {
	result := &Client{}
	err := repo.db.GetContext(ctx, result, "SELECT * FROM chat.client WHERE id=$1", id)
	if err != nil {
		repo.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *sqlxRepository) GetClientByExternalID(ctx context.Context, externalID string) (*Client, error) {
	result := &Client{}
	err := repo.db.GetContext(ctx, result, "SELECT * FROM chat.client WHERE external_id=$1", externalID)
	if err != nil {
		repo.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *sqlxRepository) CreateClient(ctx context.Context, c *Client) error {
	c.ID = 0
	c.CreatedAt = sql.NullTime{
		time.Now(),
		true,
	}
	res, err := repo.db.NamedExecContext(ctx, `insert into chat.client (name, number, created_at, external_id, first_name, last_name)
	values (:name, :number, :created_at, :external_id, :first_name, :last_name)`, *c)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = lastID
	return nil
}
