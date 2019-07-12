// Code generated by MockGen. DO NOT EDIT.
// Source: search/searchservices/search_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	bleve "github.com/blevesearch/bleve"
	searchservices "github.com/cloudfresco/vilom/search/searchservices"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSearchServiceIntf is a mock of SearchServiceIntf interface
type MockSearchServiceIntf struct {
	ctrl     *gomock.Controller
	recorder *MockSearchServiceIntfMockRecorder
}

// MockSearchServiceIntfMockRecorder is the mock recorder for MockSearchServiceIntf
type MockSearchServiceIntfMockRecorder struct {
	mock *MockSearchServiceIntf
}

// NewMockSearchServiceIntf creates a new mock instance
func NewMockSearchServiceIntf(ctrl *gomock.Controller) *MockSearchServiceIntf {
	mock := &MockSearchServiceIntf{ctrl: ctrl}
	mock.recorder = &MockSearchServiceIntfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSearchServiceIntf) EXPECT() *MockSearchServiceIntfMockRecorder {
	return m.recorder
}

// Search mocks base method
func (m *MockSearchServiceIntf) Search(form *searchservices.BleveForm, userEmail, requestID string) (*bleve.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", form, userEmail, requestID)
	ret0, _ := ret[0].(*bleve.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockSearchServiceIntfMockRecorder) Search(form, userEmail, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockSearchServiceIntf)(nil).Search), form, userEmail, requestID)
}
