package evolution

import "encoding/json"

// SuccessResponse is the common acknowledgement returned by several endpoints.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// MessageResponse is returned by all /message/* send endpoints.
type MessageKey struct {
	ID        string `json:"id,omitempty"`
	RemoteJid string `json:"remoteJid,omitempty"`
	FromMe    bool   `json:"fromMe,omitempty"`
}

type MessageResponse struct {
	Key     MessageKey      `json:"key,omitempty"`
	Message json.RawMessage `json:"message,omitempty"`
	Status  string          `json:"status,omitempty"`
}

// Query is the generic filter payload used by the /chat/find* endpoints.
type Query struct {
	Where   map[string]any `json:"where,omitempty"`
	Take    int            `json:"take,omitempty"`
	Skip    int            `json:"skip,omitempty"`
	OrderBy map[string]any `json:"orderBy,omitempty"`
}

// Presence enumerates the presence states accepted by SetPresence.
type Presence string

const (
	PresenceAvailable   Presence = "available"
	PresenceUnavailable Presence = "unavailable"
	PresenceComposing   Presence = "composing"
	PresenceRecording   Presence = "recording"
)
