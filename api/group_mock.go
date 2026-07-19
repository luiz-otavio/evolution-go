package evolution

import "context"

var _ GroupService = (*MockGroupService)(nil)

type MockGroupService struct {
	CreateFn            func(ctx context.Context, instanceName string, req CreateGroupRequest) (GroupResponse, error)
	FindGroupInfosFn    func(ctx context.Context, instanceName, groupJid string) (GroupResponse, error)
	FetchAllGroupsFn    func(ctx context.Context, instanceName string, getParticipants bool) (FetchAllGroupsResponse, error)
	ParticipantsFn      func(ctx context.Context, instanceName, groupJid string) (GroupParticipantsResponse, error)
	UpdateParticipantFn func(ctx context.Context, instanceName string, req UpdateParticipantRequest) (SuccessResponse, error)
}

func (m *MockGroupService) Create(ctx context.Context, instanceName string, req CreateGroupRequest) (GroupResponse, error) {
	return m.CreateFn(ctx, instanceName, req)
}

func (m *MockGroupService) FindGroupInfos(ctx context.Context, instanceName, groupJid string) (GroupResponse, error) {
	return m.FindGroupInfosFn(ctx, instanceName, groupJid)
}

func (m *MockGroupService) FetchAllGroups(ctx context.Context, instanceName string, getParticipants bool) (FetchAllGroupsResponse, error) {
	return m.FetchAllGroupsFn(ctx, instanceName, getParticipants)
}

func (m *MockGroupService) Participants(ctx context.Context, instanceName, groupJid string) (GroupParticipantsResponse, error) {
	return m.ParticipantsFn(ctx, instanceName, groupJid)
}

func (m *MockGroupService) UpdateParticipant(ctx context.Context, instanceName string, req UpdateParticipantRequest) (SuccessResponse, error) {
	return m.UpdateParticipantFn(ctx, instanceName, req)
}
