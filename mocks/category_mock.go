// Code generated by MockGen. DO NOT EDIT.
// Source: msg/msgservices/workspace_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	sql "database/sql"
	msgservices "github.com/cloudfresco/vilom/msg/msgservices"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockWorkspaceServiceIntf is a mock of WorkspaceServiceIntf interface
type MockWorkspaceServiceIntf struct {
	ctrl     *gomock.Controller
	recorder *MockWorkspaceServiceIntfMockRecorder
}

// MockWorkspaceServiceIntfMockRecorder is the mock recorder for MockWorkspaceServiceIntf
type MockWorkspaceServiceIntfMockRecorder struct {
	mock *MockWorkspaceServiceIntf
}

// NewMockWorkspaceServiceIntf creates a new mock instance
func NewMockWorkspaceServiceIntf(ctrl *gomock.Controller) *MockWorkspaceServiceIntf {
	mock := &MockWorkspaceServiceIntf{ctrl: ctrl}
	mock.recorder = &MockWorkspaceServiceIntfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkspaceServiceIntf) EXPECT() *MockWorkspaceServiceIntfMockRecorder {
	return m.recorder
}

// CreateWorkspace mocks base method
func (m *MockWorkspaceServiceIntf) CreateWorkspace(ctx context.Context, form *msgservices.Workspace, UserID, userEmail, requestID string) (*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWorkspace", ctx, form, UserID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWorkspace indicates an expected call of CreateWorkspace
func (mr *MockWorkspaceServiceIntfMockRecorder) CreateWorkspace(ctx, form, UserID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWorkspace", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).CreateWorkspace), ctx, form, UserID, userEmail, requestID)
}

// CreateChild mocks base method
func (m *MockWorkspaceServiceIntf) CreateChild(ctx context.Context, form *msgservices.Workspace, UserID, userEmail, requestID string) (*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChild", ctx, form, UserID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChild indicates an expected call of CreateChild
func (mr *MockWorkspaceServiceIntfMockRecorder) CreateChild(ctx, form, UserID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChild", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).CreateChild), ctx, form, UserID, userEmail, requestID)
}

// GetWorkspaces mocks base method
func (m *MockWorkspaceServiceIntf) GetWorkspaces(ctx context.Context, limit, nextCursor, userEmail, requestID string) (*msgservices.WorkspaceCursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkspaces", ctx, limit, nextCursor, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.WorkspaceCursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkspaces indicates an expected call of GetWorkspaces
func (mr *MockWorkspaceServiceIntfMockRecorder) GetWorkspaces(ctx, limit, nextCursor, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkspaces", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).GetWorkspaces), ctx, limit, nextCursor, userEmail, requestID)
}

// GetWorkspaceWithChannels mocks base method
func (m *MockWorkspaceServiceIntf) GetWorkspaceWithChannels(ctx context.Context, ID, userEmail, requestID string) (*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkspaceWithChannels", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkspaceWithChannels indicates an expected call of GetWorkspaceWithChannels
func (mr *MockWorkspaceServiceIntfMockRecorder) GetWorkspaceWithChannels(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkspaceWithChannels", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).GetWorkspaceWithChannels), ctx, ID, userEmail, requestID)
}

// GetWorkspace mocks base method
func (m *MockWorkspaceServiceIntf) GetWorkspace(ctx context.Context, ID, userEmail, requestID string) (*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkspace", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkspace indicates an expected call of GetWorkspace
func (mr *MockWorkspaceServiceIntfMockRecorder) GetWorkspace(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkspace", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).GetWorkspace), ctx, ID, userEmail, requestID)
}

// GetWorkspaceByID mocks base method
func (m *MockWorkspaceServiceIntf) GetWorkspaceByID(ctx context.Context, ID uint, userEmail, requestID string) (*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkspaceByID", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkspaceByID indicates an expected call of GetWorkspaceByID
func (mr *MockWorkspaceServiceIntfMockRecorder) GetWorkspaceByID(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkspaceByID", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).GetWorkspaceByID), ctx, ID, userEmail, requestID)
}

// GetTopLevelWorkspaces mocks base method
func (m *MockWorkspaceServiceIntf) GetTopLevelWorkspaces(ctx context.Context, userEmail, requestID string) ([]*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopLevelWorkspaces", ctx, userEmail, requestID)
	ret0, _ := ret[0].([]*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopLevelWorkspaces indicates an expected call of GetTopLevelWorkspaces
func (mr *MockWorkspaceServiceIntfMockRecorder) GetTopLevelWorkspaces(ctx, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopLevelWorkspaces", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).GetTopLevelWorkspaces), ctx, userEmail, requestID)
}

// GetChildWorkspaces mocks base method
func (m *MockWorkspaceServiceIntf) GetChildWorkspaces(ctx context.Context, ID, userEmail, requestID string) ([]*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChildWorkspaces", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].([]*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChildWorkspaces indicates an expected call of GetChildWorkspaces
func (mr *MockWorkspaceServiceIntfMockRecorder) GetChildWorkspaces(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChildWorkspaces", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).GetChildWorkspaces), ctx, ID, userEmail, requestID)
}

// GetParentWorkspace mocks base method
func (m *MockWorkspaceServiceIntf) GetParentWorkspace(ctx context.Context, ID, userEmail, requestID string) (*msgservices.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParentWorkspace", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetParentWorkspace indicates an expected call of GetParentWorkspace
func (mr *MockWorkspaceServiceIntfMockRecorder) GetParentWorkspace(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParentWorkspace", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).GetParentWorkspace), ctx, ID, userEmail, requestID)
}

// UpdateWorkspace mocks base method
func (m *MockWorkspaceServiceIntf) UpdateWorkspace(ctx context.Context, ID string, form *msgservices.Workspace, UserID, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWorkspace", ctx, ID, form, UserID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWorkspace indicates an expected call of UpdateWorkspace
func (mr *MockWorkspaceServiceIntfMockRecorder) UpdateWorkspace(ctx, ID, form, UserID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWorkspace", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).UpdateWorkspace), ctx, ID, form, UserID, userEmail, requestID)
}

// UpdateNumChannels mocks base method
func (m *MockWorkspaceServiceIntf) UpdateNumChannels(ctx context.Context, tx *sql.Tx, numChannels, ID uint, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNumChannels", ctx, tx, numChannels, ID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNumChannels indicates an expected call of UpdateNumChannels
func (mr *MockWorkspaceServiceIntfMockRecorder) UpdateNumChannels(ctx, tx, numChannels, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNumChannels", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).UpdateNumChannels), ctx, tx, numChannels, ID, userEmail, requestID)
}

// DeleteWorkspace mocks base method
func (m *MockWorkspaceServiceIntf) DeleteWorkspace(ctx context.Context, ID, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWorkspace", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWorkspace indicates an expected call of DeleteWorkspace
func (mr *MockWorkspaceServiceIntfMockRecorder) DeleteWorkspace(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWorkspace", reflect.TypeOf((*MockWorkspaceServiceIntf)(nil).DeleteWorkspace), ctx, ID, userEmail, requestID)
}
