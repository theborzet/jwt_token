// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\artem\Documents\Go_scripts\test_v_galery\pkg\auth\manager.go

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	auth "github.com/theborzet/jwt_token/pkg/auth"
)

// MockTokenManager is a mock of TokenManager interface.
type MockTokenManager struct {
	ctrl     *gomock.Controller
	recorder *MockTokenManagerMockRecorder
}

// MockTokenManagerMockRecorder is the mock recorder for MockTokenManager.
type MockTokenManagerMockRecorder struct {
	mock *MockTokenManager
}

// NewMockTokenManager creates a new mock instance.
func NewMockTokenManager(ctrl *gomock.Controller) *MockTokenManager {
	mock := &MockTokenManager{ctrl: ctrl}
	mock.recorder = &MockTokenManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenManager) EXPECT() *MockTokenManagerMockRecorder {
	return m.recorder
}

// NewJWT mocks base method.
func (m *MockTokenManager) NewJWT(userId, ipAddress string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewJWT", userId, ipAddress)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewJWT indicates an expected call of NewJWT.
func (mr *MockTokenManagerMockRecorder) NewJWT(userId, ipAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewJWT", reflect.TypeOf((*MockTokenManager)(nil).NewJWT), userId, ipAddress)
}

// NewRefreshToken mocks base method.
func (m *MockTokenManager) NewRefreshToken() (*auth.RefreshToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRefreshToken")
	ret0, _ := ret[0].(*auth.RefreshToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRefreshToken indicates an expected call of NewRefreshToken.
func (mr *MockTokenManagerMockRecorder) NewRefreshToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRefreshToken", reflect.TypeOf((*MockTokenManager)(nil).NewRefreshToken))
}

// ParseJWT mocks base method.
func (m *MockTokenManager) ParseJWT(accessToken string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseJWT", accessToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ParseJWT indicates an expected call of ParseJWT.
func (mr *MockTokenManagerMockRecorder) ParseJWT(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseJWT", reflect.TypeOf((*MockTokenManager)(nil).ParseJWT), accessToken)
}
