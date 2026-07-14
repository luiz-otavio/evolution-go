package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	MessageSendTextPath     = "/message/sendText"
	MessageSendMediaPath    = "/message/sendMedia"
	MessageSendButtonsPath  = "/message/sendButtons"
	MessageSendListPath     = "/message/sendList"
	MessageSendContactPath  = "/message/sendContact"
	MessageSendLocationPath = "/message/sendLocation"
	MessageSendPollPath     = "/message/sendPoll"
	MessageSendReactionPath = "/message/sendReaction"
	MessageSendTemplatePath = "/message/sendTemplate"
)

// TextMessage is the text payload shared by text and template messages.
type TextMessage struct {
	Text string `json:"text"`
}

type SendTextRequest struct {
	Number      string         `json:"number"`
	TextMessage TextMessage    `json:"textMessage"`
	Delay       int            `json:"delay,omitempty"`
	Quoted      map[string]any `json:"quoted,omitempty"`
	LinkPreview *bool          `json:"linkPreview,omitempty"`
	Mentioned   []string       `json:"mentioned,omitempty"`
}

func (r SendTextRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	if r.TextMessage.Text == "" {
		return fmt.Errorf("textMessage.text is required")
	}
	return nil
}

type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeDocument MediaType = "document"
)

// SendMediaRequest sends media as a JSON body where Media is a URL or base64 string.
type SendMediaRequest struct {
	Number    string    `json:"number"`
	MediaType MediaType `json:"mediatype"`
	Media     string    `json:"media"`
	Caption   string    `json:"caption,omitempty"`
	FileName  string    `json:"fileName,omitempty"`
	MimeType  string    `json:"mimetype,omitempty"`
	Delay     int       `json:"delay,omitempty"`
}

func (r SendMediaRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	switch r.MediaType {
	case MediaTypeImage, MediaTypeVideo, MediaTypeAudio, MediaTypeDocument:
	default:
		return fmt.Errorf("invalid mediatype: %q", r.MediaType)
	}
	if r.Media == "" {
		return fmt.Errorf("media is required")
	}
	return nil
}

type Button struct {
	ButtonID   string         `json:"buttonId"`
	ButtonText map[string]any `json:"buttonText"`
}

type SendButtonsRequest struct {
	Number     string   `json:"number"`
	Buttons    []Button `json:"buttons"`
	Text       string   `json:"text"`
	FooterText string   `json:"footerText"`
}

func (r SendButtonsRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	if len(r.Buttons) == 0 {
		return fmt.Errorf("buttons is required")
	}
	if r.Text == "" {
		return fmt.Errorf("text is required")
	}
	return nil
}

type SendListRequest struct {
	Number      string           `json:"number"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	ButtonText  string           `json:"buttonText"`
	Sections    []map[string]any `json:"sections"`
}

func (r SendListRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(r.Sections) == 0 {
		return fmt.Errorf("sections is required")
	}
	return nil
}

type SendContactRequest struct {
	Number   string           `json:"number"`
	Contacts []map[string]any `json:"contacts"`
}

func (r SendContactRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	if len(r.Contacts) == 0 {
		return fmt.Errorf("contacts is required")
	}
	return nil
}

type SendLocationRequest struct {
	Number    string  `json:"number"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

func (r SendLocationRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	return nil
}

type SendPollRequest struct {
	Number          string   `json:"number"`
	Name            string   `json:"name"`
	SelectableCount int      `json:"selectableCount"`
	Values          []string `json:"values"`
}

func (r SendPollRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(r.Values) == 0 {
		return fmt.Errorf("values is required")
	}
	return nil
}

type SendReactionRequest struct {
	ReactionKey     map[string]any `json:"reactionKey"`
	ReactionMessage string         `json:"reactionMessage"`
}

func (r SendReactionRequest) Validate() error {
	if len(r.ReactionKey) == 0 {
		return fmt.Errorf("reactionKey is required")
	}
	return nil
}

type SendTemplateRequest struct {
	Number      string      `json:"number"`
	TextMessage TextMessage `json:"textMessage"`
	Delay       int         `json:"delay,omitempty"`
}

func (r SendTemplateRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	if r.TextMessage.Text == "" {
		return fmt.Errorf("textMessage.text is required")
	}
	return nil
}

type MessageService interface {
	SendText(ctx context.Context, instanceName string, req SendTextRequest) (MessageResponse, error)
	SendMedia(ctx context.Context, instanceName string, req SendMediaRequest) (MessageResponse, error)
	SendButtons(ctx context.Context, instanceName string, req SendButtonsRequest) (MessageResponse, error)
	SendList(ctx context.Context, instanceName string, req SendListRequest) (MessageResponse, error)
	SendContact(ctx context.Context, instanceName string, req SendContactRequest) (MessageResponse, error)
	SendLocation(ctx context.Context, instanceName string, req SendLocationRequest) (MessageResponse, error)
	SendPoll(ctx context.Context, instanceName string, req SendPollRequest) (MessageResponse, error)
	SendReaction(ctx context.Context, instanceName string, req SendReactionRequest) (MessageResponse, error)
	SendTemplate(ctx context.Context, instanceName string, req SendTemplateRequest) (MessageResponse, error)
}

type messageService struct {
	http   HttpProvider
	apiKey string
}

func NewMessageService(http HttpProvider, apiKey string) MessageService {
	return messageService{http: http, apiKey: apiKey}
}

type validatable interface {
	Validate() error
}

// sendMessage validates the request, marshals it, and posts to the given base path.
func sendMessage[T validatable](ctx context.Context, s messageService, base, instanceName string, req T) (MessageResponse, error) {
	if instanceName == "" {
		return MessageResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return MessageResponse{}, fmt.Errorf("invalid message request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return MessageResponse{}, fmt.Errorf("failed to marshal message request: %w", err)
	}

	path := buildInstancePath(base, instanceName)
	return executePost[MessageResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s messageService) SendText(ctx context.Context, instanceName string, req SendTextRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendTextPath, instanceName, req)
}

func (s messageService) SendMedia(ctx context.Context, instanceName string, req SendMediaRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendMediaPath, instanceName, req)
}

func (s messageService) SendButtons(ctx context.Context, instanceName string, req SendButtonsRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendButtonsPath, instanceName, req)
}

func (s messageService) SendList(ctx context.Context, instanceName string, req SendListRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendListPath, instanceName, req)
}

func (s messageService) SendContact(ctx context.Context, instanceName string, req SendContactRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendContactPath, instanceName, req)
}

func (s messageService) SendLocation(ctx context.Context, instanceName string, req SendLocationRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendLocationPath, instanceName, req)
}

func (s messageService) SendPoll(ctx context.Context, instanceName string, req SendPollRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendPollPath, instanceName, req)
}

func (s messageService) SendReaction(ctx context.Context, instanceName string, req SendReactionRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendReactionPath, instanceName, req)
}

func (s messageService) SendTemplate(ctx context.Context, instanceName string, req SendTemplateRequest) (MessageResponse, error) {
	return sendMessage(ctx, s, MessageSendTemplatePath, instanceName, req)
}
