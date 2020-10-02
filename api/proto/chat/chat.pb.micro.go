// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: chat.proto

package chat

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for ChatService service

func NewChatServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for ChatService service

type ChatService interface {
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...client.CallOption) (*SendMessageResponse, error)
	StartConversation(ctx context.Context, in *StartConversationRequest, opts ...client.CallOption) (*StartConversationResponse, error)
	CloseConversation(ctx context.Context, in *CloseConversationRequest, opts ...client.CallOption) (*CloseConversationResponse, error)
	JoinConversation(ctx context.Context, in *JoinConversationRequest, opts ...client.CallOption) (*JoinConversationResponse, error)
	LeaveConversation(ctx context.Context, in *LeaveConversationRequest, opts ...client.CallOption) (*LeaveConversationResponse, error)
	InviteToConversation(ctx context.Context, in *InviteToConversationRequest, opts ...client.CallOption) (*InviteToConversationResponse, error)
	DeclineInvitation(ctx context.Context, in *DeclineInvitationRequest, opts ...client.CallOption) (*DeclineInvitationResponse, error)
	CheckSession(ctx context.Context, in *CheckSessionRequest, opts ...client.CallOption) (*CheckSessionResponse, error)
	WaitMessage(ctx context.Context, in *WaitMessageRequest, opts ...client.CallOption) (*WaitMessageResponse, error)
	GetProfilesX(ctx context.Context, in *GetProfilesRequest, opts ...client.CallOption) (*GetProfilesResponse, error)
	UpdateProfileX(ctx context.Context, in *UpdateProfileRequest, opts ...client.CallOption) (*UpdateProfileResponse, error)
	GetConversationByID(ctx context.Context, in *GetConversationByIDRequest, opts ...client.CallOption) (*GetConversationByIDResponse, error)
	GetConversations(ctx context.Context, in *GetConversationsRequest, opts ...client.CallOption) (*GetConversationsResponse, error)
	GetProfiles(ctx context.Context, in *GetProfilesRequest, opts ...client.CallOption) (*GetProfilesResponse, error)
	GetProfileByID(ctx context.Context, in *GetProfileByIDRequest, opts ...client.CallOption) (*GetProfileByIDResponse, error)
	CreateProfile(ctx context.Context, in *CreateProfileRequest, opts ...client.CallOption) (*CreateProfileResponse, error)
	DeleteProfile(ctx context.Context, in *DeleteProfileRequest, opts ...client.CallOption) (*DeleteProfileResponse, error)
	UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...client.CallOption) (*UpdateProfileResponse, error)
	GetHistoryMessages(ctx context.Context, in *GetHistoryMessagesRequest, opts ...client.CallOption) (*GetHistoryMessagesResponse, error)
}

type chatService struct {
	c    client.Client
	name string
}

func NewChatService(name string, c client.Client) ChatService {
	return &chatService{
		c:    c,
		name: name,
	}
}

func (c *chatService) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...client.CallOption) (*SendMessageResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.SendMessage", in)
	out := new(SendMessageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) StartConversation(ctx context.Context, in *StartConversationRequest, opts ...client.CallOption) (*StartConversationResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.StartConversation", in)
	out := new(StartConversationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) CloseConversation(ctx context.Context, in *CloseConversationRequest, opts ...client.CallOption) (*CloseConversationResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.CloseConversation", in)
	out := new(CloseConversationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) JoinConversation(ctx context.Context, in *JoinConversationRequest, opts ...client.CallOption) (*JoinConversationResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.JoinConversation", in)
	out := new(JoinConversationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) LeaveConversation(ctx context.Context, in *LeaveConversationRequest, opts ...client.CallOption) (*LeaveConversationResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.LeaveConversation", in)
	out := new(LeaveConversationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) InviteToConversation(ctx context.Context, in *InviteToConversationRequest, opts ...client.CallOption) (*InviteToConversationResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.InviteToConversation", in)
	out := new(InviteToConversationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) DeclineInvitation(ctx context.Context, in *DeclineInvitationRequest, opts ...client.CallOption) (*DeclineInvitationResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.DeclineInvitation", in)
	out := new(DeclineInvitationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) CheckSession(ctx context.Context, in *CheckSessionRequest, opts ...client.CallOption) (*CheckSessionResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.CheckSession", in)
	out := new(CheckSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) WaitMessage(ctx context.Context, in *WaitMessageRequest, opts ...client.CallOption) (*WaitMessageResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.WaitMessage", in)
	out := new(WaitMessageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) GetProfilesX(ctx context.Context, in *GetProfilesRequest, opts ...client.CallOption) (*GetProfilesResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.GetProfilesX", in)
	out := new(GetProfilesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) UpdateProfileX(ctx context.Context, in *UpdateProfileRequest, opts ...client.CallOption) (*UpdateProfileResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.UpdateProfileX", in)
	out := new(UpdateProfileResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) GetConversationByID(ctx context.Context, in *GetConversationByIDRequest, opts ...client.CallOption) (*GetConversationByIDResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.GetConversationByID", in)
	out := new(GetConversationByIDResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) GetConversations(ctx context.Context, in *GetConversationsRequest, opts ...client.CallOption) (*GetConversationsResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.GetConversations", in)
	out := new(GetConversationsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) GetProfiles(ctx context.Context, in *GetProfilesRequest, opts ...client.CallOption) (*GetProfilesResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.GetProfiles", in)
	out := new(GetProfilesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) GetProfileByID(ctx context.Context, in *GetProfileByIDRequest, opts ...client.CallOption) (*GetProfileByIDResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.GetProfileByID", in)
	out := new(GetProfileByIDResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) CreateProfile(ctx context.Context, in *CreateProfileRequest, opts ...client.CallOption) (*CreateProfileResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.CreateProfile", in)
	out := new(CreateProfileResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) DeleteProfile(ctx context.Context, in *DeleteProfileRequest, opts ...client.CallOption) (*DeleteProfileResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.DeleteProfile", in)
	out := new(DeleteProfileResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...client.CallOption) (*UpdateProfileResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.UpdateProfile", in)
	out := new(UpdateProfileResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) GetHistoryMessages(ctx context.Context, in *GetHistoryMessagesRequest, opts ...client.CallOption) (*GetHistoryMessagesResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.GetHistoryMessages", in)
	out := new(GetHistoryMessagesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ChatService service

type ChatServiceHandler interface {
	SendMessage(context.Context, *SendMessageRequest, *SendMessageResponse) error
	StartConversation(context.Context, *StartConversationRequest, *StartConversationResponse) error
	CloseConversation(context.Context, *CloseConversationRequest, *CloseConversationResponse) error
	JoinConversation(context.Context, *JoinConversationRequest, *JoinConversationResponse) error
	LeaveConversation(context.Context, *LeaveConversationRequest, *LeaveConversationResponse) error
	InviteToConversation(context.Context, *InviteToConversationRequest, *InviteToConversationResponse) error
	DeclineInvitation(context.Context, *DeclineInvitationRequest, *DeclineInvitationResponse) error
	CheckSession(context.Context, *CheckSessionRequest, *CheckSessionResponse) error
	WaitMessage(context.Context, *WaitMessageRequest, *WaitMessageResponse) error
	GetProfilesX(context.Context, *GetProfilesRequest, *GetProfilesResponse) error
	UpdateProfileX(context.Context, *UpdateProfileRequest, *UpdateProfileResponse) error
	GetConversationByID(context.Context, *GetConversationByIDRequest, *GetConversationByIDResponse) error
	GetConversations(context.Context, *GetConversationsRequest, *GetConversationsResponse) error
	GetProfiles(context.Context, *GetProfilesRequest, *GetProfilesResponse) error
	GetProfileByID(context.Context, *GetProfileByIDRequest, *GetProfileByIDResponse) error
	CreateProfile(context.Context, *CreateProfileRequest, *CreateProfileResponse) error
	DeleteProfile(context.Context, *DeleteProfileRequest, *DeleteProfileResponse) error
	UpdateProfile(context.Context, *UpdateProfileRequest, *UpdateProfileResponse) error
	GetHistoryMessages(context.Context, *GetHistoryMessagesRequest, *GetHistoryMessagesResponse) error
}

func RegisterChatServiceHandler(s server.Server, hdlr ChatServiceHandler, opts ...server.HandlerOption) error {
	type chatService interface {
		SendMessage(ctx context.Context, in *SendMessageRequest, out *SendMessageResponse) error
		StartConversation(ctx context.Context, in *StartConversationRequest, out *StartConversationResponse) error
		CloseConversation(ctx context.Context, in *CloseConversationRequest, out *CloseConversationResponse) error
		JoinConversation(ctx context.Context, in *JoinConversationRequest, out *JoinConversationResponse) error
		LeaveConversation(ctx context.Context, in *LeaveConversationRequest, out *LeaveConversationResponse) error
		InviteToConversation(ctx context.Context, in *InviteToConversationRequest, out *InviteToConversationResponse) error
		DeclineInvitation(ctx context.Context, in *DeclineInvitationRequest, out *DeclineInvitationResponse) error
		CheckSession(ctx context.Context, in *CheckSessionRequest, out *CheckSessionResponse) error
		WaitMessage(ctx context.Context, in *WaitMessageRequest, out *WaitMessageResponse) error
		GetProfilesX(ctx context.Context, in *GetProfilesRequest, out *GetProfilesResponse) error
		UpdateProfileX(ctx context.Context, in *UpdateProfileRequest, out *UpdateProfileResponse) error
		GetConversationByID(ctx context.Context, in *GetConversationByIDRequest, out *GetConversationByIDResponse) error
		GetConversations(ctx context.Context, in *GetConversationsRequest, out *GetConversationsResponse) error
		GetProfiles(ctx context.Context, in *GetProfilesRequest, out *GetProfilesResponse) error
		GetProfileByID(ctx context.Context, in *GetProfileByIDRequest, out *GetProfileByIDResponse) error
		CreateProfile(ctx context.Context, in *CreateProfileRequest, out *CreateProfileResponse) error
		DeleteProfile(ctx context.Context, in *DeleteProfileRequest, out *DeleteProfileResponse) error
		UpdateProfile(ctx context.Context, in *UpdateProfileRequest, out *UpdateProfileResponse) error
		GetHistoryMessages(ctx context.Context, in *GetHistoryMessagesRequest, out *GetHistoryMessagesResponse) error
	}
	type ChatService struct {
		chatService
	}
	h := &chatServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ChatService{h}, opts...))
}

type chatServiceHandler struct {
	ChatServiceHandler
}

func (h *chatServiceHandler) SendMessage(ctx context.Context, in *SendMessageRequest, out *SendMessageResponse) error {
	return h.ChatServiceHandler.SendMessage(ctx, in, out)
}

func (h *chatServiceHandler) StartConversation(ctx context.Context, in *StartConversationRequest, out *StartConversationResponse) error {
	return h.ChatServiceHandler.StartConversation(ctx, in, out)
}

func (h *chatServiceHandler) CloseConversation(ctx context.Context, in *CloseConversationRequest, out *CloseConversationResponse) error {
	return h.ChatServiceHandler.CloseConversation(ctx, in, out)
}

func (h *chatServiceHandler) JoinConversation(ctx context.Context, in *JoinConversationRequest, out *JoinConversationResponse) error {
	return h.ChatServiceHandler.JoinConversation(ctx, in, out)
}

func (h *chatServiceHandler) LeaveConversation(ctx context.Context, in *LeaveConversationRequest, out *LeaveConversationResponse) error {
	return h.ChatServiceHandler.LeaveConversation(ctx, in, out)
}

func (h *chatServiceHandler) InviteToConversation(ctx context.Context, in *InviteToConversationRequest, out *InviteToConversationResponse) error {
	return h.ChatServiceHandler.InviteToConversation(ctx, in, out)
}

func (h *chatServiceHandler) DeclineInvitation(ctx context.Context, in *DeclineInvitationRequest, out *DeclineInvitationResponse) error {
	return h.ChatServiceHandler.DeclineInvitation(ctx, in, out)
}

func (h *chatServiceHandler) CheckSession(ctx context.Context, in *CheckSessionRequest, out *CheckSessionResponse) error {
	return h.ChatServiceHandler.CheckSession(ctx, in, out)
}

func (h *chatServiceHandler) WaitMessage(ctx context.Context, in *WaitMessageRequest, out *WaitMessageResponse) error {
	return h.ChatServiceHandler.WaitMessage(ctx, in, out)
}

func (h *chatServiceHandler) GetProfilesX(ctx context.Context, in *GetProfilesRequest, out *GetProfilesResponse) error {
	return h.ChatServiceHandler.GetProfilesX(ctx, in, out)
}

func (h *chatServiceHandler) UpdateProfileX(ctx context.Context, in *UpdateProfileRequest, out *UpdateProfileResponse) error {
	return h.ChatServiceHandler.UpdateProfileX(ctx, in, out)
}

func (h *chatServiceHandler) GetConversationByID(ctx context.Context, in *GetConversationByIDRequest, out *GetConversationByIDResponse) error {
	return h.ChatServiceHandler.GetConversationByID(ctx, in, out)
}

func (h *chatServiceHandler) GetConversations(ctx context.Context, in *GetConversationsRequest, out *GetConversationsResponse) error {
	return h.ChatServiceHandler.GetConversations(ctx, in, out)
}

func (h *chatServiceHandler) GetProfiles(ctx context.Context, in *GetProfilesRequest, out *GetProfilesResponse) error {
	return h.ChatServiceHandler.GetProfiles(ctx, in, out)
}

func (h *chatServiceHandler) GetProfileByID(ctx context.Context, in *GetProfileByIDRequest, out *GetProfileByIDResponse) error {
	return h.ChatServiceHandler.GetProfileByID(ctx, in, out)
}

func (h *chatServiceHandler) CreateProfile(ctx context.Context, in *CreateProfileRequest, out *CreateProfileResponse) error {
	return h.ChatServiceHandler.CreateProfile(ctx, in, out)
}

func (h *chatServiceHandler) DeleteProfile(ctx context.Context, in *DeleteProfileRequest, out *DeleteProfileResponse) error {
	return h.ChatServiceHandler.DeleteProfile(ctx, in, out)
}

func (h *chatServiceHandler) UpdateProfile(ctx context.Context, in *UpdateProfileRequest, out *UpdateProfileResponse) error {
	return h.ChatServiceHandler.UpdateProfile(ctx, in, out)
}

func (h *chatServiceHandler) GetHistoryMessages(ctx context.Context, in *GetHistoryMessagesRequest, out *GetHistoryMessagesResponse) error {
	return h.ChatServiceHandler.GetHistoryMessages(ctx, in, out)
}
