package events

const (
	MessageEventType            = "message"
	CloseConversationEventType  = "close_conversation"
	JoinConversationEventType   = "join_conversation"
	LeaveConversationEventType  = "leave_conversation"
	InviteConversationEventType = "invite_conversation"
	UserInvitationEventType     = "user_invite"
	ExpireInvitationEventType   = "expire_invite"
	DeclineInvitationEventType  = "decline_invite"
)

type BaseEvent struct {
	ConversationID string `json:"conversation_id"`
	Timestamp      int64  `json:"timestamp"`
}

type MessageEvent struct {
	BaseEvent
	FromChannelID string `json:"from_channel_id"`
	// ToChannelID    int64  `json:"to_channel_id"`
	MessageID int64  `json:"message_id"`
	Type      string `json:"message_type"`
	Value     []byte `json:"message_value"`
}

type CloseConversationEvent struct {
	BaseEvent
	FromChannelID string `json:"from_channel_id"`
	// ToChannelID    int64  `json:"to_channel_id"`
	Cause string `json:"cause"`
}

type JoinConversationEvent struct {
	BaseEvent
	JoinedChannelID string `json:"joined_channel_id"`
	// JoinedUserID    int64 `json:"joined_user_id"`
}

type LeaveConversationEvent struct {
	BaseEvent
	LeavedChannelID string `json:"leaved_channel_id"`
	// LeavedUserID    int64 `json:"leaved_user_id"`
}

type InviteConversationEvent struct {
	BaseEvent
	UserID int64 `json:"user_id"`
}

type UserInvitationEvent struct {
	BaseEvent
	InviteID string `json:"invite_id"`
}

type DeclineInvitationEvent struct {
	BaseEvent
	UserID   int64  `json:"user_id"`
	InviteID string `json:"invite_id"`
}
