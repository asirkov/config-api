// Code generated by MockGen. DO NOT EDIT.
// Source: httpclient.go

// Package mock_httpclient is a generated GoMock package.
package mock_httpclient

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHttpClient is a mock of HttpClient interface.
type MockHttpClient struct {
	ctrl     *gomock.Controller
	recorder *MockHttpClientMockRecorder
}

// MockHttpClientMockRecorder is the mock recorder for MockHttpClient.
type MockHttpClientMockRecorder struct {
	mock *MockHttpClient
}

// NewMockHttpClient creates a new mock instance.
func NewMockHttpClient(ctrl *gomock.Controller) *MockHttpClient {
	mock := &MockHttpClient{ctrl: ctrl}
	mock.recorder = &MockHttpClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpClient) EXPECT() *MockHttpClientMockRecorder {
	return m.recorder
}

// DoHttpGet mocks base method.
func (m *MockHttpClient) DoHttpGet(url string) ([]byte, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoHttpGet", url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DoHttpGet indicates an expected call of DoHttpGet.
func (mr *MockHttpClientMockRecorder) DoHttpGet(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoHttpGet", reflect.TypeOf((*MockHttpClient)(nil).DoHttpGet), url)
}
