package main

const (
	messageEventType            = "message"
	closeConversationEventType  = "close_conversation"
	joinConversationEventType   = "join_conversation"
	leaveConversationEventType  = "leave_conversation"
	inviteConversationEventType = "invite_conversation"
	userInvitationEventType     = "user_invite"
	declineInvitationEventType  = "decline_invite"
)

type messageEvent struct {
	ConversationID int64 `json:"conversation_id"`
	FromChannelID  int64 `json:"from_channel_id"`
	// ToChannelID    int64  `json:"to_channel_id"`
	MessageID int64  `json:"message_id"`
	Type      string `json:"message_type"`
	Value     []byte `json:"message_value"`
}

type closeConversationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	FromChannelID  int64 `json:"from_channel_id"`
	// ToChannelID    int64  `json:"to_channel_id"`
	Cause string `json:"cause"`
}

type joinConversationEvent struct {
	ConversationID  int64 `json:"conversation_id"`
	JoinedChannelID int64 `json:"joined_channel_id"`
	// JoinedUserID    int64 `json:"joined_user_id"`
}

type leaveConversationEvent struct {
	ConversationID  int64 `json:"conversation_id"`
	LeavedChannelID int64 `json:"leaved_channel_id"`
	// LeavedUserID    int64 `json:"leaved_user_id"`
}

type inviteConversationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	UserID         int64 `json:"user_id"`
}

type userInvitationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	InviteID       int64 `json:"invite_id"`
}
type declineInvitationEvent struct {
	ConversationID int64 `json:"conversation_id"`
	UserID         int64 `json:"user_id"`
}
