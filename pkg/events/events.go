package events

const (
	MessageEventType            = "message"
	CloseConversationEventType  = "close_conversation"
	JoinConversationEventType   = "join_conversation"
	LeaveConversationEventType  = "leave_conversation"
	InviteConversationEventType = "invite_conversation"
	UserInvitationEventType     = "user_invite"
	DeclineInvitationEventType  = "decline_invite"
)

type MessageEvent struct {
	ConversationID int64 `json:"conversation_id"`
	FromChannelID  int64 `json:"from_channel_id"`
	// ToChannelID    int64  `json:"to_channel_id"`
	MessageID int64  `json:"message_id"`
	Type      string `json:"message_type"`
	Value     []byte `json:"message_value"`
}

type CloseConversationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	FromChannelID  int64 `json:"from_channel_id"`
	// ToChannelID    int64  `json:"to_channel_id"`
	Cause string `json:"cause"`
}

type JoinConversationEvent struct {
	ConversationID  int64 `json:"conversation_id"`
	JoinedChannelID int64 `json:"joined_channel_id"`
	// JoinedUserID    int64 `json:"joined_user_id"`
}

type LeaveConversationEvent struct {
	ConversationID  int64 `json:"conversation_id"`
	LeavedChannelID int64 `json:"leaved_channel_id"`
	// LeavedUserID    int64 `json:"leaved_user_id"`
}

type InviteConversationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	UserID         int64 `json:"user_id"`
}

type UserInvitationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	InviteID       int64 `json:"invite_id"`
}
type DeclineInvitationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	UserID         int64 `json:"user_id"`
}
