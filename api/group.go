package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	GroupCreatePath            = "/group/create"
	GroupFindGroupInfosPath    = "/group/findGroupInfos"
	GroupParticipantsPath      = "/group/participants"
	GroupUpdateParticipantPath = "/group/updateParticipant"
)

type CreateGroupRequest struct {
	Subject      string   `json:"subject"`
	Participants []string `json:"participants"`
	Description  string   `json:"description,omitempty"`
}

func (r CreateGroupRequest) Validate() error {
	if r.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if len(r.Participants) == 0 {
		return fmt.Errorf("participants is required")
	}
	return nil
}

type GroupResponse struct {
	Success bool           `json:"success"`
	Group   map[string]any `json:"group"`
}

type GroupParticipantsResponse struct {
	Success      bool             `json:"success"`
	Participants []map[string]any `json:"participants"`
}

type ParticipantAction string

const (
	ParticipantActionAdd     ParticipantAction = "add"
	ParticipantActionRemove  ParticipantAction = "remove"
	ParticipantActionPromote ParticipantAction = "promote"
	ParticipantActionDemote  ParticipantAction = "demote"
)

type UpdateParticipantRequest struct {
	GroupJid     string            `json:"groupJid"`
	Action       ParticipantAction `json:"action"`
	Participants []string          `json:"participants"`
}

func (r UpdateParticipantRequest) Validate() error {
	if r.GroupJid == "" {
		return fmt.Errorf("groupJid is required")
	}
	switch r.Action {
	case ParticipantActionAdd, ParticipantActionRemove, ParticipantActionPromote, ParticipantActionDemote:
	default:
		return fmt.Errorf("invalid action: %q", r.Action)
	}
	if len(r.Participants) == 0 {
		return fmt.Errorf("participants is required")
	}
	return nil
}

type GroupService interface {
	Create(ctx context.Context, instanceName string, req CreateGroupRequest) (GroupResponse, error)
	FindGroupInfos(ctx context.Context, instanceName, groupJid string) (GroupResponse, error)
	Participants(ctx context.Context, instanceName, groupJid string) (GroupParticipantsResponse, error)
	UpdateParticipant(ctx context.Context, instanceName string, req UpdateParticipantRequest) (SuccessResponse, error)
}

type groupService struct {
	http   HttpProvider
	apiKey string
}

func NewGroupService(http HttpProvider, apiKey string) GroupService {
	return groupService{http: http, apiKey: apiKey}
}

func (s groupService) Create(ctx context.Context, instanceName string, req CreateGroupRequest) (GroupResponse, error) {
	if instanceName == "" {
		return GroupResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return GroupResponse{}, fmt.Errorf("invalid create group request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return GroupResponse{}, fmt.Errorf("failed to marshal create group request: %w", err)
	}

	path := buildInstancePath(GroupCreatePath, instanceName)
	return executePost[GroupResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s groupService) FindGroupInfos(ctx context.Context, instanceName, groupJid string) (GroupResponse, error) {
	if instanceName == "" {
		return GroupResponse{}, fmt.Errorf("instanceName is required")
	}
	if groupJid == "" {
		return GroupResponse{}, fmt.Errorf("groupJid is required")
	}

	path := buildInstancePath(GroupFindGroupInfosPath, instanceName)
	return executeGet[GroupResponse](ctx, s.http, s.apiKey, path, map[string]string{"groupJid": groupJid})
}

func (s groupService) Participants(ctx context.Context, instanceName, groupJid string) (GroupParticipantsResponse, error) {
	if instanceName == "" {
		return GroupParticipantsResponse{}, fmt.Errorf("instanceName is required")
	}
	if groupJid == "" {
		return GroupParticipantsResponse{}, fmt.Errorf("groupJid is required")
	}

	path := buildInstancePath(GroupParticipantsPath, instanceName)
	return executeGet[GroupParticipantsResponse](ctx, s.http, s.apiKey, path, map[string]string{"groupJid": groupJid})
}

func (s groupService) UpdateParticipant(ctx context.Context, instanceName string, req UpdateParticipantRequest) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return SuccessResponse{}, fmt.Errorf("invalid update participant request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal update participant request: %w", err)
	}

	path := buildInstancePath(GroupUpdateParticipantPath, instanceName)
	return executePost[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}
