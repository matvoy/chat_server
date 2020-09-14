package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetChannelByID(ctx context.Context, id int64) (*models.Channel, error) {
	result, err := models.Channels(
		models.ChannelWhere.ID.EQ(id),
		qm.Load(models.ChannelRels.Conversation),
	).
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

func (repo *PgRepository) GetChannels(
	ctx context.Context,
	userID *int64,
	conversationID *int64,
	connection *string,
	internal *bool,
	exceptID *int64,
) ([]*models.Channel, error) {
	query := make([]qm.QueryMod, 0, 6)
	query = append(query, models.ChannelWhere.ClosedAt.IsNull())
	if userID != nil {
		query = append(query, models.ChannelWhere.UserID.EQ(*userID))
	}
	if conversationID != nil {
		query = append(query, models.ChannelWhere.UserID.EQ(*conversationID))
	}
	if connection != nil {
		query = append(query, models.ChannelWhere.Connection.EQ(
			null.String{
				*connection,
				true,
			},
		))
	}
	if internal != nil {
		query = append(query, models.ChannelWhere.Internal.EQ(*internal))
	}
	if exceptID != nil {
		query = append(query, models.ChannelWhere.ID.NEQ(*exceptID))
	}
	return models.Channels(query...).All(ctx, repo.db)
}

func (repo *PgRepository) CreateChannel(ctx context.Context, c *models.Channel) error {
	if err := c.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) CloseChannel(ctx context.Context, id int64) error {
	result, err := models.Conversations(models.ConversationWhere.ID.EQ(id)).
		One(ctx, repo.db)
	if err != nil {
		repo.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	result.ClosedAt = null.Time{
		Valid: true,
		Time:  time.Now(),
	}
	_, err = result.Update(ctx, repo.db, boil.Infer())
	return err
}
