package chat_cache

import (
	"fmt"
	"time"

	"github.com/micro/go-micro/v2/store"
)

const (
	sessionStr            = "session_id:%s"    // %s - session id, value - conversation id
	confirmationStr       = "confirmations:%v" // %s - conversation id, value - confirmation id
	writeCachedMessageStr = "cached_messages:%v:%v"
	readCachedMessageStr  = "cached_messages:%v"
)

type ChatCache interface {
	ReadSession(sessionID string) ([]byte, error)
	WriteSession(sessionID string, conversationIDBytes []byte) error
	DeleteSession(sessionID string) error

	ReadConfirmation(conversationID int64) ([]byte, error)
	WriteConfirmation(conversationID int64, confirmationIDBytes []byte) error
	DeleteConfirmation(conversationID int64) error

	ReadCachedMessages(conversationID int64) ([]*store.Record, error)
	WriteCachedMessage(conversationID int64, messageID int64, messageBytes []byte) error
	DeleteCachedMessages(conversationID int64) error
	DeleteCachedMessage(key string) error
}

type chatCache struct {
	redisStore store.Store
}

func NewChatCache(redisStore store.Store) ChatCache {
	return &chatCache{
		redisStore,
	}
}

func (c *chatCache) ReadSession(sessionID string) ([]byte, error) {
	sessionKey := fmt.Sprintf(sessionStr, sessionID)
	session, err := c.redisStore.Read(sessionKey)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	if len(session) > 0 {
		return session[0].Value, nil
	} else {
		return nil, nil
	}
}

func (c *chatCache) WriteSession(sessionID string, conversationIDBytes []byte) error {
	sessionKey := fmt.Sprintf(sessionStr, sessionID)
	return c.redisStore.Write(&store.Record{
		Key:    sessionKey,
		Value:  conversationIDBytes,
		Expiry: time.Hour * time.Duration(24),
	})
}

func (c *chatCache) DeleteSession(sessionID string) error {
	sessionKey := fmt.Sprintf(sessionStr, sessionID)
	return c.redisStore.Delete(sessionKey)
}

func (c *chatCache) ReadConfirmation(conversationID int64) ([]byte, error) {
	confirmationKey := fmt.Sprintf(confirmationStr, conversationID)
	confirmationID, err := c.redisStore.Read(confirmationKey)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	if len(confirmationID) > 0 {
		return confirmationID[0].Value, nil
	} else {
		return nil, nil
	}
}

func (c *chatCache) WriteConfirmation(conversationID int64, confirmationIDBytes []byte) error {
	confirmationKey := fmt.Sprintf(confirmationStr, conversationID)
	return c.redisStore.Write(&store.Record{
		Key:    confirmationKey,
		Value:  confirmationIDBytes,
		Expiry: time.Hour * time.Duration(24),
	})
}

func (c *chatCache) DeleteConfirmation(conversationID int64) error {
	confirmationKey := fmt.Sprintf(confirmationStr, conversationID)
	return c.redisStore.Delete(confirmationKey)
}

func (c *chatCache) ReadCachedMessages(conversationID int64) ([]*store.Record, error) {
	messagesKey := fmt.Sprintf(readCachedMessageStr, conversationID)
	cachedMessages, err := c.redisStore.Read(messagesKey)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	if len(cachedMessages) > 0 {
		return cachedMessages, nil
	} else {
		return nil, nil
	}
}

func (c *chatCache) WriteCachedMessage(conversationID int64, messageID int64, messageBytes []byte) error {
	messagesKey := fmt.Sprintf(writeCachedMessageStr, conversationID, messageID)
	return c.redisStore.Write(&store.Record{
		Key:    messagesKey,
		Value:  messageBytes,
		Expiry: time.Hour * time.Duration(24),
	})
}

func (c *chatCache) DeleteCachedMessages(conversationID int64) error {
	messagesKey := fmt.Sprintf(readCachedMessageStr, conversationID)
	cachedMessages, _ := c.redisStore.Read(messagesKey)
	for _, m := range cachedMessages {
		if err := c.redisStore.Delete(m.Key); err != nil {
			return err
		}
	}
	return nil
}

func (c *chatCache) DeleteCachedMessage(key string) error {
	return c.redisStore.Delete(key)
}
