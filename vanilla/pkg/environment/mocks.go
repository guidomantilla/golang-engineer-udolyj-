// Code generated by MockGen. DO NOT EDIT.
// Source: ../pkg/environment/types.go
//
// Generated by this command:
//
//	mockgen -package=environment -source ../pkg/environment/types.go -destination ../pkg/environment/mocks.go
//

// Package environment is a generated GoMock package.
package environment

import (
	reflect "reflect"

	properties "git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/properties"
	gomock "go.uber.org/mock/gomock"
)

// MockEnvironment is a mock of Environment interface.
type MockEnvironment struct {
	ctrl     *gomock.Controller
	recorder *MockEnvironmentMockRecorder
}

// MockEnvironmentMockRecorder is the mock recorder for MockEnvironment.
type MockEnvironmentMockRecorder struct {
	mock *MockEnvironment
}

// NewMockEnvironment creates a new mock instance.
func NewMockEnvironment(ctrl *gomock.Controller) *MockEnvironment {
	mock := &MockEnvironment{ctrl: ctrl}
	mock.recorder = &MockEnvironmentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEnvironment) EXPECT() *MockEnvironmentMockRecorder {
	return m.recorder
}

// AppendPropertySources mocks base method.
func (m *MockEnvironment) AppendPropertySources(propertySources ...properties.PropertySource) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range propertySources {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "AppendPropertySources", varargs...)
}

// AppendPropertySources indicates an expected call of AppendPropertySources.
func (mr *MockEnvironmentMockRecorder) AppendPropertySources(propertySources ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendPropertySources", reflect.TypeOf((*MockEnvironment)(nil).AppendPropertySources), propertySources...)
}

// GetPropertySources mocks base method.
func (m *MockEnvironment) GetPropertySources() []properties.PropertySource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertySources")
	ret0, _ := ret[0].([]properties.PropertySource)
	return ret0
}

// GetPropertySources indicates an expected call of GetPropertySources.
func (mr *MockEnvironmentMockRecorder) GetPropertySources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertySources", reflect.TypeOf((*MockEnvironment)(nil).GetPropertySources))
}

// GetValue mocks base method.
func (m *MockEnvironment) GetValue(property string) EnvVar {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", property)
	ret0, _ := ret[0].(EnvVar)
	return ret0
}

// GetValue indicates an expected call of GetValue.
func (mr *MockEnvironmentMockRecorder) GetValue(property any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockEnvironment)(nil).GetValue), property)
}

// GetValueOrDefault mocks base method.
func (m *MockEnvironment) GetValueOrDefault(property, defaultValue string) EnvVar {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValueOrDefault", property, defaultValue)
	ret0, _ := ret[0].(EnvVar)
	return ret0
}

// GetValueOrDefault indicates an expected call of GetValueOrDefault.
func (mr *MockEnvironmentMockRecorder) GetValueOrDefault(property, defaultValue any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValueOrDefault", reflect.TypeOf((*MockEnvironment)(nil).GetValueOrDefault), property, defaultValue)
}
