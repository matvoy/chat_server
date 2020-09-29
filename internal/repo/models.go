package repo

import (
	"time"
)

type Conversation struct {
	ID        int64      `json:"id"`
	Title     *string    `json:"title,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DomainID  int64      `json:"domain_id"`
	Members   []*Member  `json:"members"`
}

type Member struct {
	ChannelID int64  `json:"channel_id"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Type      string `json:"type"`
	Internal  bool   `json:"internal"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}
