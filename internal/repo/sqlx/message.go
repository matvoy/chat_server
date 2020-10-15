package sqlxrepo

import (
	"context"
	"database/sql"
	"time"
)

func (repo *sqlxRepository) CreateMessage(ctx context.Context, m *Message) error {
	m.ID = 0
	tmp := sql.NullTime{
		time.Now(),
		true,
	}
	m.CreatedAt = tmp
	m.UpdatedAt = tmp
	res, err := repo.db.NamedExecContext(ctx, `insert into chat.message (channel_id, conversation_id, text, variables, created_at, updated_at, type)
	values (:channel_id, :conversation_id, :text, :variables, :created_at, :updated_at, :type)`, *m)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = lastID
	return nil
}

func (repo *sqlxRepository) GetMessages(ctx context.Context, id int64, size, page int32, fields, sort []string, conversationID string) ([]*Message, error) {
	result := []*Message{}
	// TO DO FILTERS
	err := repo.db.SelectContext(ctx, &result, "SELECT * FROM chat.message")
	return result, err
	// query := make([]qm.QueryMod, 0, 6)
	// if size != 0 {
	// 	query = append(query, qm.Limit(int(size)))
	// } else {
	// 	query = append(query, qm.Limit(15))
	// }
	// if page != 0 {
	// 	query = append(query, qm.Offset(int((page-1)*size)))
	// }
	// if id != 0 {
	// 	query = append(query, models.MessageWhere.ID.EQ(id))
	// }
	// if fields != nil && len(fields) > 0 {
	// 	query = append(query, qm.Select(fields...))
	// }
	// if sort != nil && len(sort) > 0 {
	// 	for _, item := range sort {
	// 		query = append(query, qm.OrderBy(item))
	// 	}
	// } else {
	// 	query = append(query, qm.OrderBy("created_at"))
	// }
	// if conversationID != "" {
	// 	query = append(query, models.MessageWhere.ConversationID.EQ(conversationID))
	// }
	// return models.Messages(query...).All(ctx, repo.db)
}
