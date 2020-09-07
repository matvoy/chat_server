package pg

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (repo *PgRepository) GetConversationBySessionID(ctx context.Context, sessionID string) (*models.Conversation, error) {
	result, err := models.Conversations(qm.Where("LOWER(session_id) like ?", strings.ToLower(sessionID)), qm.Where("closed_at is null")).
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

func (repo *PgRepository) GetConversationByID(ctx context.Context, id int64) (*models.Conversation, error) {
	result, err := models.Conversations(
		models.ConversationWhere.ID.EQ(id),
		qm.Load(models.ConversationRels.Profile),
	).One(ctx, repo.db)
	if err != nil {
		repo.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *PgRepository) CreateConversation(ctx context.Context, c *models.Conversation) error {
	if err := c.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) CloseConversation(ctx context.Context, id int64) error {
	// result, err := models.Conversations(qm.Where("LOWER(session_id) like ?", strings.ToLower(sessionID)), qm.Where("closed_at is null")).
	// 	One(ctx, repo.db)
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

func (repo *PgRepository) GetConversations(ctx context.Context, limit, offset int) ([]*models.Conversation, error) {
	return models.Conversations(qm.Limit(limit), qm.Offset(offset)).All(ctx, repo.db)
}
