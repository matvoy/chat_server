package sqlxrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func (repo *sqlxRepository) GetChannelByID(ctx context.Context, id string) (*Channel, error) {
	result := &Channel{}
	err := repo.db.GetContext(ctx, result, "SELECT * FROM chat.channel WHERE id=$1", id)
	if err != nil {
		repo.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *sqlxRepository) GetChannels(
	ctx context.Context,
	userID *int64,
	conversationID *string,
	connection *string,
	internal *bool,
	exceptID *string,
) ([]*Channel, error) {
	result := []*Channel{}
	// TO DO FILTERS
	err := repo.db.SelectContext(ctx, &result, "SELECT * FROM chat.channel")
	return result, err
}

func (repo *sqlxRepository) CreateChannel(ctx context.Context, c *Channel) error {
	c.ID = uuid.New().String()
	tmp := sql.NullTime{
		time.Now(),
		true,
	}
	c.CreatedAt = tmp
	c.UpdatedAt = tmp
	_, err := repo.db.NamedExecContext(ctx, `insert into chat.channel (
		id, 
		type, 
		conversation_id, 
		user_id, 
		connection, 
		created_at, 
		internal, 
		closed_at, 
		updated_at, 
		domain_id, 
		flow_bridge,
		name
	)
	values (
		:id, 
		:type, 
		:conversation_id, 
		:user_id, 
		:connection, 
		:created_at, 
		:internal, 
		:closed_at, 
		:updated_at, 
		:domain_id, 
		:flow_bridge,
		:name
		)`, *c)
	return err
}

func (repo *sqlxRepository) CloseChannel(ctx context.Context, id string) (*Channel, error) {
	result := &Channel{}
	err := repo.db.GetContext(ctx, result, "SELECT * FROM chat.channel WHERE id=$1", id)
	if err != nil {
		repo.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	_, err = repo.db.ExecContext(ctx, `update chat.channel set closed_at=$1 where id=$2`, sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}, id)
	return result, err
}

func (repo *sqlxRepository) CloseChannels(ctx context.Context, conversationID string) error {
	_, err := repo.db.ExecContext(ctx, `update chat.channel set closed_at=$1 where conversation_id=$2`, sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}, conversationID)
	return err
}
