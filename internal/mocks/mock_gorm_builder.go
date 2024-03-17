// Code generated by MockGen. DO NOT EDIT.
// Source: gorm.io/gorm/clause (interfaces: Builder)
//
// Generated by this command:
//
//	mockgen -destination=internal/mocks/mock_gorm_builder.go -package=mocks gorm.io/gorm/clause Builder
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	clause "gorm.io/gorm/clause"
)

// MockBuilder is a mock of Builder interface.
type MockBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderMockRecorder
}

// MockBuilderMockRecorder is the mock recorder for MockBuilder.
type MockBuilderMockRecorder struct {
	mock *MockBuilder
}

// NewMockBuilder creates a new mock instance.
func NewMockBuilder(ctrl *gomock.Controller) *MockBuilder {
	mock := &MockBuilder{ctrl: ctrl}
	mock.recorder = &MockBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuilder) EXPECT() *MockBuilderMockRecorder {
	return m.recorder
}

// AddError mocks base method.
func (m *MockBuilder) AddError(arg0 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddError", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddError indicates an expected call of AddError.
func (mr *MockBuilderMockRecorder) AddError(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddError", reflect.TypeOf((*MockBuilder)(nil).AddError), arg0)
}

// AddVar mocks base method.
func (m *MockBuilder) AddVar(arg0 clause.Writer, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "AddVar", varargs...)
}

// AddVar indicates an expected call of AddVar.
func (mr *MockBuilderMockRecorder) AddVar(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVar", reflect.TypeOf((*MockBuilder)(nil).AddVar), varargs...)
}

// WriteByte mocks base method.
func (m *MockBuilder) WriteByte(arg0 byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteByte", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteByte indicates an expected call of WriteByte.
func (mr *MockBuilderMockRecorder) WriteByte(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteByte", reflect.TypeOf((*MockBuilder)(nil).WriteByte), arg0)
}

// WriteQuoted mocks base method.
func (m *MockBuilder) WriteQuoted(arg0 any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteQuoted", arg0)
}

// WriteQuoted indicates an expected call of WriteQuoted.
func (mr *MockBuilderMockRecorder) WriteQuoted(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteQuoted", reflect.TypeOf((*MockBuilder)(nil).WriteQuoted), arg0)
}

// WriteString mocks base method.
func (m *MockBuilder) WriteString(arg0 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteString", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteString indicates an expected call of WriteString.
func (mr *MockBuilderMockRecorder) WriteString(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteString", reflect.TypeOf((*MockBuilder)(nil).WriteString), arg0)
}
