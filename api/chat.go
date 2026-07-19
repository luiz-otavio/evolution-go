package evolution

import (
	"context"
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	ChatWhatsAppNumbersPath      = "/chat/whatsappNumbers"
	ChatMarkMessageAsReadPath    = "/chat/markMessageAsRead"
	ChatArchiveChatPath          = "/chat/archiveChat"
	ChatFindChatsPath            = "/chat/findChats"
	ChatFindContactsPath         = "/chat/findContacts"
	ChatFindMessagesPath         = "/chat/findMessages"
	ChatUpdateProfileNamePath    = "/chat/updateProfileName"
	ChatUpdateProfilePicturePath = "/chat/updateProfilePicture"
	ChatUpdateProfileStatusPath  = "/chat/updateProfileStatus"
)

type Contact struct {
	ID                string `json:"id"`
	PushName          string `json:"pushName"`
	Number            string `json:"number"`
	ProfilePictureURL string `json:"profilePictureUrl"`
}

type WhatsAppNumbersRequest struct {
	Numbers []string `json:"numbers"`
}

func (r WhatsAppNumbersRequest) Validate() error {
	if len(r.Numbers) == 0 {
		return fmt.Errorf("numbers is required")
	}
	return nil
}

type WhatsAppNumberStatus struct {
	Number string `json:"number"`
	Exists bool   `json:"exists"`
}

type WhatsAppNumbersResponse struct {
	Numbers []WhatsAppNumberStatus `json:"numbers"`
}

type MarkMessageAsReadRequest struct {
	ReadMessages []map[string]any `json:"readMessages"`
}

func (r MarkMessageAsReadRequest) Validate() error {
	if len(r.ReadMessages) == 0 {
		return fmt.Errorf("readMessages is required")
	}
	return nil
}

type ArchiveChatRequest struct {
	Number  string `json:"number"`
	Archive bool   `json:"archive"`
}

func (r ArchiveChatRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	return nil
}

type FindMessagesResponse struct {
	Messages FindMessagesPage `json:"messages"`
}

type FindMessagesPage struct {
	Total       int                 `json:"total"`
	Pages       int                 `json:"pages"`
	CurrentPage int                 `json:"currentPage"`
	Records     []FindMessageRecord `json:"records"`
}

type FindMessageRecord struct {
	ID               string          `json:"id,omitempty"`
	Key              MessageKey      `json:"key,omitempty"`
	Message          json.RawMessage `json:"message,omitempty"`
	PushName         string          `json:"pushName,omitempty"`
	MessageTimestamp int64           `json:"messageTimestamp,omitempty"`
}

type ChatSummary struct {
	ID       string `json:"id,omitempty"`
	JID      string `json:"jid,omitempty"`
	Name     string `json:"name,omitempty"`
	Subject  string `json:"subject,omitempty"`
	PushName string `json:"pushName,omitempty"`
}

type UpdateProfileNameRequest struct {
	Name string `json:"name"`
}

func (r UpdateProfileNameRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// UpdateProfilePictureRequest sends the picture as a URL in a JSON body.
type UpdateProfilePictureRequest struct {
	Picture string `json:"picture"`
}

func (r UpdateProfilePictureRequest) Validate() error {
	if r.Picture == "" {
		return fmt.Errorf("picture is required")
	}
	return nil
}

type UpdateProfileStatusRequest struct {
	Status string `json:"status"`
}

func (r UpdateProfileStatusRequest) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	return nil
}

type ChatService interface {
	CheckWhatsAppNumbers(ctx context.Context, instanceName string, req WhatsAppNumbersRequest) (WhatsAppNumbersResponse, error)
	MarkMessageAsRead(ctx context.Context, instanceName string, req MarkMessageAsReadRequest) (SuccessResponse, error)
	ArchiveChat(ctx context.Context, instanceName string, req ArchiveChatRequest) (SuccessResponse, error)
	FindChats(ctx context.Context, instanceName string, query Query) ([]ChatSummary, error)
	FindContacts(ctx context.Context, instanceName string, query Query) ([]Contact, error)
	FindMessages(ctx context.Context, instanceName string, query Query) (FindMessagesResponse, error)
	UpdateProfileName(ctx context.Context, instanceName string, req UpdateProfileNameRequest) (SuccessResponse, error)
	UpdateProfilePicture(ctx context.Context, instanceName string, req UpdateProfilePictureRequest) (SuccessResponse, error)
	UpdateProfileStatus(ctx context.Context, instanceName string, req UpdateProfileStatusRequest) (SuccessResponse, error)
}

type chatService struct {
	http   HttpProvider
	apiKey string
}

func NewChatService(http HttpProvider, apiKey string) ChatService {
	return chatService{http: http, apiKey: apiKey}
}

// postChat validates, marshals and posts a chat request to an instance path.
func postChat[Req validatable, Res any](ctx context.Context, s chatService, base, instanceName string, req Req) (Res, error) {
	var out Res
	if instanceName == "" {
		return out, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return out, fmt.Errorf("invalid chat request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return out, fmt.Errorf("failed to marshal chat request: %w", err)
	}

	path := buildInstancePath(base, instanceName)
	return executePost[Res](ctx, s.http, s.apiKey, path, payload)
}

func (s chatService) CheckWhatsAppNumbers(ctx context.Context, instanceName string, req WhatsAppNumbersRequest) (WhatsAppNumbersResponse, error) {
	return postChat[WhatsAppNumbersRequest, WhatsAppNumbersResponse](ctx, s, ChatWhatsAppNumbersPath, instanceName, req)
}

func (s chatService) MarkMessageAsRead(ctx context.Context, instanceName string, req MarkMessageAsReadRequest) (SuccessResponse, error) {
	return postChat[MarkMessageAsReadRequest, SuccessResponse](ctx, s, ChatMarkMessageAsReadPath, instanceName, req)
}

func (s chatService) ArchiveChat(ctx context.Context, instanceName string, req ArchiveChatRequest) (SuccessResponse, error) {
	return postChat[ArchiveChatRequest, SuccessResponse](ctx, s, ChatArchiveChatPath, instanceName, req)
}

func (s chatService) FindChats(ctx context.Context, instanceName string, query Query) ([]ChatSummary, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}

	payload, err := jsoniter.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	path := buildInstancePath(ChatFindChatsPath, instanceName)
	return executePost[[]ChatSummary](ctx, s.http, s.apiKey, path, payload)
}

func (s chatService) FindContacts(ctx context.Context, instanceName string, query Query) ([]Contact, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}

	payload, err := jsoniter.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	path := buildInstancePath(ChatFindContactsPath, instanceName)
	return executePost[[]Contact](ctx, s.http, s.apiKey, path, payload)
}

func (s chatService) FindMessages(ctx context.Context, instanceName string, query Query) (FindMessagesResponse, error) {
	if instanceName == "" {
		return FindMessagesResponse{}, fmt.Errorf("instanceName is required")
	}

	payload, err := jsoniter.Marshal(query)
	if err != nil {
		return FindMessagesResponse{}, fmt.Errorf("failed to marshal query: %w", err)
	}

	path := buildInstancePath(ChatFindMessagesPath, instanceName)
	return executePost[FindMessagesResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s chatService) UpdateProfileName(ctx context.Context, instanceName string, req UpdateProfileNameRequest) (SuccessResponse, error) {
	return postChat[UpdateProfileNameRequest, SuccessResponse](ctx, s, ChatUpdateProfileNamePath, instanceName, req)
}

func (s chatService) UpdateProfilePicture(ctx context.Context, instanceName string, req UpdateProfilePictureRequest) (SuccessResponse, error) {
	return postChat[UpdateProfilePictureRequest, SuccessResponse](ctx, s, ChatUpdateProfilePicturePath, instanceName, req)
}

func (s chatService) UpdateProfileStatus(ctx context.Context, instanceName string, req UpdateProfileStatusRequest) (SuccessResponse, error) {
	return postChat[UpdateProfileStatusRequest, SuccessResponse](ctx, s, ChatUpdateProfileStatusPath, instanceName, req)
}
