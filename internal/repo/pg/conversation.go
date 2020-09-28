package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// func (repo *PgRepository) GetConversationBySessionID(ctx context.Context, sessionID string) (*models.Conversation, error) {
// 	result, err := models.Conversations(qm.Where("LOWER(session_id) like ?", strings.ToLower(sessionID)), qm.Where("closed_at is null")).
// 		One(ctx, repo.db)
// 	if err != nil {
// 		repo.log.Warn().Msg(err.Error())
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return result, nil
// }

func (r *PgRepository) GetConversationByID(ctx context.Context, id int64) (*models.Conversation, error) {
	result, err := models.Conversations(
		models.ConversationWhere.ID.EQ(id),
	).One(ctx, r.db)
	if err != nil {
		r.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *PgRepository) CreateConversation(ctx context.Context, c *models.Conversation) error {
	if err := c.Insert(ctx, r.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) CloseConversation(ctx context.Context, id int64) error {
	// result, err := models.Conversations(qm.Where("LOWER(session_id) like ?", strings.ToLower(sessionID)), qm.Where("closed_at is null")).
	// 	One(ctx, repo.db)
	result, err := models.Conversations(models.ConversationWhere.ID.EQ(id)).
		One(ctx, r.db)
	if err != nil {
		r.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	result.ClosedAt = null.Time{
		Valid: true,
		Time:  time.Now(),
	}
	_, err = result.Update(ctx, r.db, boil.Infer())
	return err
}

func (r *PgRepository) GetConversations(ctx context.Context, limit, offset int) ([]*repo.Conversation, error) {
	conversations, err := models.Conversations(qm.Limit(limit), qm.Offset(offset), qm.Load(models.ConversationRels.Channels)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	result := make([]*repo.Conversation, 0, len(conversations))
	for _, c := range conversations {
		members := make([]*repo.Member, 0, len(c.R.Channels))
		for _, ch := range c.R.Channels {
			if !ch.Internal {
				client, err := models.Clients(models.ClientWhere.ID.EQ(ch.UserID)).One(ctx, r.db)
				if err != nil {
					return nil, err
				}
				members = append(members, &repo.Member{
					ChannelID: ch.ID,
					UserID:    ch.UserID,
					Username:  client.Name.String,
					Firstname: client.FirstName.String,
					Lastname:  client.LastName.String,
				})
			} else {
				members = append(members, &repo.Member{
					ChannelID: ch.ID,
					UserID:    ch.UserID,
				})
			}

		}
		result = append(result, &repo.Conversation{
			ID:        c.ID,
			Title:     &c.Title.String,
			CreatedAt: &c.CreatedAt.Time,
			ClosedAt:  &c.ClosedAt.Time,
			UpdatedAt: &c.UpdatedAt.Time,
			DomainID:  c.DomainID,
			Members:   members,
		})
	}
	return result, nil
}
