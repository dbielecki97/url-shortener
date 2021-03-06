// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dbielecki97/url-shortener/internal/app (interfaces: Service)

// Package app is a generated GoMock package.
package app

import (
	reflect "reflect"

	api "github.com/dbielecki97/url-shortener/internal/api"
	errs "github.com/dbielecki97/url-shortener/pkg/errs"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Expand mocks base method.
func (m *MockService) Expand(arg0 string) (*api.ShortenInfo, errs.RestErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expand", arg0)
	ret0, _ := ret[0].(*api.ShortenInfo)
	ret1, _ := ret[1].(errs.RestErr)
	return ret0, ret1
}

// Expand indicates an expected call of Expand.
func (mr *MockServiceMockRecorder) Expand(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expand", reflect.TypeOf((*MockService)(nil).Expand), arg0)
}

// Shorten mocks base method.
func (m *MockService) Shorten(arg0 api.ShortenRequest) (*api.ShortenInfo, errs.RestErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shorten", arg0)
	ret0, _ := ret[0].(*api.ShortenInfo)
	ret1, _ := ret[1].(errs.RestErr)
	return ret0, ret1
}

// Shorten indicates an expected call of Shorten.
func (mr *MockServiceMockRecorder) Shorten(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shorten", reflect.TypeOf((*MockService)(nil).Shorten), arg0)
}
