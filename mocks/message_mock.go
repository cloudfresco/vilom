// Code generated by MockGen. DO NOT EDIT.
// Source: msg/msgservices/message_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	sql "database/sql"
	msgservices "github.com/cloudfresco/vilom/msg/msgservices"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMessageServiceIntf is a mock of MessageServiceIntf interface
type MockMessageServiceIntf struct {
	ctrl     *gomock.Controller
	recorder *MockMessageServiceIntfMockRecorder
}

// MockMessageServiceIntfMockRecorder is the mock recorder for MockMessageServiceIntf
type MockMessageServiceIntfMockRecorder struct {
	mock *MockMessageServiceIntf
}

// NewMockMessageServiceIntf creates a new mock instance
func NewMockMessageServiceIntf(ctrl *gomock.Controller) *MockMessageServiceIntf {
	mock := &MockMessageServiceIntf{ctrl: ctrl}
	mock.recorder = &MockMessageServiceIntfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessageServiceIntf) EXPECT() *MockMessageServiceIntfMockRecorder {
	return m.recorder
}

// CreateMessage mocks base method
func (m *MockMessageServiceIntf) CreateMessage(ctx context.Context, form *msgservices.Message, UserID string, rplymsg bool, userEmail, requestID string) (*msgservices.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", ctx, form, UserID, rplymsg, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMessage indicates an expected call of CreateMessage
func (mr *MockMessageServiceIntfMockRecorder) CreateMessage(ctx, form, UserID, rplymsg, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockMessageServiceIntf)(nil).CreateMessage), ctx, form, UserID, rplymsg, userEmail, requestID)
}

// CreateUserReply mocks base method
func (m *MockMessageServiceIntf) CreateUserReply(ctx context.Context, tx *sql.Tx, channelID, messageID, userID, ugroupID uint, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserReply", ctx, tx, channelID, messageID, userID, ugroupID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserReply indicates an expected call of CreateUserReply
func (mr *MockMessageServiceIntfMockRecorder) CreateUserReply(ctx, tx, channelID, messageID, userID, ugroupID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserReply", reflect.TypeOf((*MockMessageServiceIntf)(nil).CreateUserReply), ctx, tx, channelID, messageID, userID, ugroupID, userEmail, requestID)
}

// CreateUserLike mocks base method
func (m *MockMessageServiceIntf) CreateUserLike(ctx context.Context, form *msgservices.UserLike, UserID, userEmail, requestID string) (*msgservices.UserLike, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserLike", ctx, form, UserID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.UserLike)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserLike indicates an expected call of CreateUserLike
func (mr *MockMessageServiceIntfMockRecorder) CreateUserLike(ctx, form, UserID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserLike", reflect.TypeOf((*MockMessageServiceIntf)(nil).CreateUserLike), ctx, form, UserID, userEmail, requestID)
}

// CreateUserVote mocks base method
func (m *MockMessageServiceIntf) CreateUserVote(ctx context.Context, form *msgservices.UserVote, UserID, userEmail, requestID string) (*msgservices.UserVote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserVote", ctx, form, UserID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.UserVote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserVote indicates an expected call of CreateUserVote
func (mr *MockMessageServiceIntfMockRecorder) CreateUserVote(ctx, form, UserID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserVote", reflect.TypeOf((*MockMessageServiceIntf)(nil).CreateUserVote), ctx, form, UserID, userEmail, requestID)
}

// GetMessage mocks base method
func (m *MockMessageServiceIntf) GetMessage(ctx context.Context, ID, userEmail, requestID string) (*msgservices.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(*msgservices.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage
func (mr *MockMessageServiceIntfMockRecorder) GetMessage(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockMessageServiceIntf)(nil).GetMessage), ctx, ID, userEmail, requestID)
}

// GetMessagesWithTextAttach mocks base method
func (m *MockMessageServiceIntf) GetMessagesWithTextAttach(ctx context.Context, messages []*msgservices.Message, userEmail, requestID string) ([]*msgservices.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessagesWithTextAttach", ctx, messages, userEmail, requestID)
	ret0, _ := ret[0].([]*msgservices.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessagesWithTextAttach indicates an expected call of GetMessagesWithTextAttach
func (mr *MockMessageServiceIntfMockRecorder) GetMessagesWithTextAttach(ctx, messages, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessagesWithTextAttach", reflect.TypeOf((*MockMessageServiceIntf)(nil).GetMessagesWithTextAttach), ctx, messages, userEmail, requestID)
}

// GetMessagesTexts mocks base method
func (m *MockMessageServiceIntf) GetMessagesTexts(ctx context.Context, messageID uint, userEmail, requestID string) ([]*msgservices.MessageText, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessagesTexts", ctx, messageID, userEmail, requestID)
	ret0, _ := ret[0].([]*msgservices.MessageText)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessagesTexts indicates an expected call of GetMessagesTexts
func (mr *MockMessageServiceIntfMockRecorder) GetMessagesTexts(ctx, messageID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessagesTexts", reflect.TypeOf((*MockMessageServiceIntf)(nil).GetMessagesTexts), ctx, messageID, userEmail, requestID)
}

// GetMessageAttachments mocks base method
func (m *MockMessageServiceIntf) GetMessageAttachments(ctx context.Context, messageID uint, userEmail, requestID string) ([]*msgservices.MessageAttachment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageAttachments", ctx, messageID, userEmail, requestID)
	ret0, _ := ret[0].([]*msgservices.MessageAttachment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageAttachments indicates an expected call of GetMessageAttachments
func (mr *MockMessageServiceIntfMockRecorder) GetMessageAttachments(ctx, messageID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageAttachments", reflect.TypeOf((*MockMessageServiceIntf)(nil).GetMessageAttachments), ctx, messageID, userEmail, requestID)
}

// UpdateMessage mocks base method
func (m *MockMessageServiceIntf) UpdateMessage(ctx context.Context, ID string, form *msgservices.Message, UserID, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessage", ctx, ID, form, UserID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMessage indicates an expected call of UpdateMessage
func (mr *MockMessageServiceIntfMockRecorder) UpdateMessage(ctx, ID, form, UserID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessage", reflect.TypeOf((*MockMessageServiceIntf)(nil).UpdateMessage), ctx, ID, form, UserID, userEmail, requestID)
}

// DeleteMessage mocks base method
func (m *MockMessageServiceIntf) DeleteMessage(ctx context.Context, ID, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessage indicates an expected call of DeleteMessage
func (mr *MockMessageServiceIntfMockRecorder) DeleteMessage(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockMessageServiceIntf)(nil).DeleteMessage), ctx, ID, userEmail, requestID)
}
