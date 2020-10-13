package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/matvoy/chat_server/internal/repo"
	"github.com/matvoy/chat_server/models"

	"github.com/google/uuid"
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

func (r *PgRepository) GetConversationByID(ctx context.Context, id string) (*repo.Conversation, error) {
	c, err := models.Conversations(
		models.ConversationWhere.ID.EQ(id),
		qm.Load(models.ConversationRels.Channels),
	).One(ctx, r.db)
	if err != nil {
		r.log.Warn().Msg(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
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
				Type:      ch.Type,
				Username:  client.Name.String,
				Firstname: client.FirstName.String,
				Lastname:  client.LastName.String,
				Internal:  ch.Internal,
			})
		} else {
			members = append(members, &repo.Member{
				ChannelID: ch.ID,
				UserID:    ch.UserID,
				Type:      ch.Type,
				Internal:  ch.Internal,
			})
		}

	}
	conv := &repo.Conversation{
		ID:        c.ID,
		Title:     &c.Title.String,
		CreatedAt: &c.CreatedAt.Time,
		ClosedAt:  &c.ClosedAt.Time,
		UpdatedAt: &c.UpdatedAt.Time,
		DomainID:  c.DomainID,
		Members:   members,
	}
	return conv, nil
}

func (r *PgRepository) CreateConversation(ctx context.Context, c *models.Conversation) error {
	c.ID = uuid.New().String()
	if err := c.Insert(ctx, r.db, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) CloseConversation(ctx context.Context, id string) error {
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

func (r *PgRepository) GetConversations(
	ctx context.Context,
	id string,
	size int32,
	page int32,
	fields []string,
	sort []string,
	domainID int64,
	active bool,
	userID int64,
) ([]*repo.Conversation, error) {
	query := make([]qm.QueryMod, 0, 8)
	// query = append(query, qm.Load(models.ConversationRels.Channels))
	if size != 0 {
		query = append(query, qm.Limit(int(size)))
	} else {
		query = append(query, qm.Limit(15))
	}
	if page != 0 {
		query = append(query, qm.Offset(int((page-1)*size)))
	}
	if id != "" {
		query = append(query, models.ConversationWhere.ID.EQ(id))
	}
	if fields != nil && len(fields) > 0 {
		query = append(query, qm.Select(fields...))
	}
	if sort != nil && len(sort) > 0 {
		for _, item := range sort {
			query = append(query, qm.OrderBy(item))
		}
	}
	if domainID != 0 {
		query = append(query, models.ConversationWhere.DomainID.EQ(domainID))
	}
	if active {
		query = append(query, models.ConversationWhere.ClosedAt.IsNull())
	}
	if userID != 0 {
		query = append(query, qm.Load(models.ConversationRels.Channels, models.ChannelWhere.UserID.EQ(userID)))
	} else {
		query = append(query, qm.Load(models.ConversationRels.Channels))
	}
	conversations, err := models.Conversations(query...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	result := make([]*repo.Conversation, 0, len(conversations))
	for _, c := range conversations {
		if len(c.R.Channels) == 0 {
			continue
		}
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
					Type:      ch.Type,
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
