// Code generated by MockGen. DO NOT EDIT.
// Source: user/userservices/user_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	common "github.com/cloudfresco/vilom/common"
	userservices "github.com/cloudfresco/vilom/user/userservices"
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockUserServiceIntf is a mock of UserServiceIntf interface
type MockUserServiceIntf struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceIntfMockRecorder
}

// MockUserServiceIntfMockRecorder is the mock recorder for MockUserServiceIntf
type MockUserServiceIntfMockRecorder struct {
	mock *MockUserServiceIntf
}

// NewMockUserServiceIntf creates a new mock instance
func NewMockUserServiceIntf(ctrl *gomock.Controller) *MockUserServiceIntf {
	mock := &MockUserServiceIntf{ctrl: ctrl}
	mock.recorder = &MockUserServiceIntfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserServiceIntf) EXPECT() *MockUserServiceIntfMockRecorder {
	return m.recorder
}

// Login mocks base method
func (m *MockUserServiceIntf) Login(ctx context.Context, form *userservices.LoginForm, requestID string) (*userservices.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, form, requestID)
	ret0, _ := ret[0].(*userservices.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login
func (mr *MockUserServiceIntfMockRecorder) Login(ctx, form, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserServiceIntf)(nil).Login), ctx, form, requestID)
}

// CreateUser mocks base method
func (m *MockUserServiceIntf) CreateUser(ctx context.Context, form *userservices.User, hostURL, requestID string) (*userservices.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, form, hostURL, requestID)
	ret0, _ := ret[0].(*userservices.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserServiceIntfMockRecorder) CreateUser(ctx, form, hostURL, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserServiceIntf)(nil).CreateUser), ctx, form, hostURL, requestID)
}

// GetUsers mocks base method
func (m *MockUserServiceIntf) GetUsers(ctx context.Context, limit, nextCursor, userEmail, requestID string) (*userservices.UserCursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx, limit, nextCursor, userEmail, requestID)
	ret0, _ := ret[0].(*userservices.UserCursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers
func (mr *MockUserServiceIntfMockRecorder) GetUsers(ctx, limit, nextCursor, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserServiceIntf)(nil).GetUsers), ctx, limit, nextCursor, userEmail, requestID)
}

// GetUserByEmail mocks base method
func (m *MockUserServiceIntf) GetUserByEmail(ctx context.Context, Email, userEmail, requestID string) (*userservices.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, Email, userEmail, requestID)
	ret0, _ := ret[0].(*userservices.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail
func (mr *MockUserServiceIntfMockRecorder) GetUserByEmail(ctx, Email, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserServiceIntf)(nil).GetUserByEmail), ctx, Email, userEmail, requestID)
}

// GetUser mocks base method
func (m *MockUserServiceIntf) GetUser(ctx context.Context, ID, userEmail, requestID string) (*userservices.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(*userservices.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockUserServiceIntfMockRecorder) GetUser(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserServiceIntf)(nil).GetUser), ctx, ID, userEmail, requestID)
}

// UpdateUser mocks base method
func (m *MockUserServiceIntf) UpdateUser(ctx context.Context, ID string, form *userservices.User, UserID, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, ID, form, UserID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser
func (mr *MockUserServiceIntfMockRecorder) UpdateUser(ctx, ID, form, UserID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserServiceIntf)(nil).UpdateUser), ctx, ID, form, UserID, userEmail, requestID)
}

// DeleteUser mocks base method
func (m *MockUserServiceIntf) DeleteUser(ctx context.Context, ID, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, ID, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockUserServiceIntfMockRecorder) DeleteUser(ctx, ID, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserServiceIntf)(nil).DeleteUser), ctx, ID, userEmail, requestID)
}

// ConfirmEmail mocks base method
func (m *MockUserServiceIntf) ConfirmEmail(ctx context.Context, token, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmEmail", ctx, token, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmEmail indicates an expected call of ConfirmEmail
func (mr *MockUserServiceIntfMockRecorder) ConfirmEmail(ctx, token, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmEmail", reflect.TypeOf((*MockUserServiceIntf)(nil).ConfirmEmail), ctx, token, requestID)
}

// ForgotPassword mocks base method
func (m *MockUserServiceIntf) ForgotPassword(ctx context.Context, form *userservices.ForgotPasswordForm, hostURL, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForgotPassword", ctx, form, hostURL, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForgotPassword indicates an expected call of ForgotPassword
func (mr *MockUserServiceIntfMockRecorder) ForgotPassword(ctx, form, hostURL, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgotPassword", reflect.TypeOf((*MockUserServiceIntf)(nil).ForgotPassword), ctx, form, hostURL, requestID)
}

// ConfirmForgotPassword mocks base method
func (m *MockUserServiceIntf) ConfirmForgotPassword(ctx context.Context, form *userservices.PasswordForm, token, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmForgotPassword", ctx, form, token, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmForgotPassword indicates an expected call of ConfirmForgotPassword
func (mr *MockUserServiceIntfMockRecorder) ConfirmForgotPassword(ctx, form, token, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmForgotPassword", reflect.TypeOf((*MockUserServiceIntf)(nil).ConfirmForgotPassword), ctx, form, token, requestID)
}

// ChangePassword mocks base method
func (m *MockUserServiceIntf) ChangePassword(ctx context.Context, form *userservices.PasswordForm, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", ctx, form, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword
func (mr *MockUserServiceIntfMockRecorder) ChangePassword(ctx, form, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockUserServiceIntf)(nil).ChangePassword), ctx, form, userEmail, requestID)
}

// ChangeEmail mocks base method
func (m *MockUserServiceIntf) ChangeEmail(ctx context.Context, form *userservices.ChangeEmailForm, hostURL, userEmail, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeEmail", ctx, form, hostURL, userEmail, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeEmail indicates an expected call of ChangeEmail
func (mr *MockUserServiceIntfMockRecorder) ChangeEmail(ctx, form, hostURL, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeEmail", reflect.TypeOf((*MockUserServiceIntf)(nil).ChangeEmail), ctx, form, hostURL, userEmail, requestID)
}

// ConfirmChangeEmail mocks base method
func (m *MockUserServiceIntf) ConfirmChangeEmail(ctx context.Context, token, requestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmChangeEmail", ctx, token, requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmChangeEmail indicates an expected call of ConfirmChangeEmail
func (mr *MockUserServiceIntfMockRecorder) ConfirmChangeEmail(ctx, token, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmChangeEmail", reflect.TypeOf((*MockUserServiceIntf)(nil).ConfirmChangeEmail), ctx, token, requestID)
}

// GetAuthUserDetails mocks base method
func (m *MockUserServiceIntf) GetAuthUserDetails(r *http.Request) (*common.ContextData, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthUserDetails", r)
	ret0, _ := ret[0].(*common.ContextData)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAuthUserDetails indicates an expected call of GetAuthUserDetails
func (mr *MockUserServiceIntfMockRecorder) GetAuthUserDetails(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthUserDetails", reflect.TypeOf((*MockUserServiceIntf)(nil).GetAuthUserDetails), r)
}
