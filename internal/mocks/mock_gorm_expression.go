// Code generated by MockGen. DO NOT EDIT.
// Source: gorm.io/gorm/clause (interfaces: Expression)
//
// Generated by this command:
//
//	mockgen -destination=internal/mocks/mock_gorm_expression.go -package=mocks gorm.io/gorm/clause Expression
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	clause "gorm.io/gorm/clause"
)

// MockExpression is a mock of Expression interface.
type MockExpression struct {
	ctrl     *gomock.Controller
	recorder *MockExpressionMockRecorder
	isgomock struct{}
}

// MockExpressionMockRecorder is the mock recorder for MockExpression.
type MockExpressionMockRecorder struct {
	mock *MockExpression
}

// NewMockExpression creates a new mock instance.
func NewMockExpression(ctrl *gomock.Controller) *MockExpression {
	mock := &MockExpression{ctrl: ctrl}
	mock.recorder = &MockExpressionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExpression) EXPECT() *MockExpressionMockRecorder {
	return m.recorder
}

// Build mocks base method.
func (m *MockExpression) Build(builder clause.Builder) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Build", builder)
}

// Build indicates an expected call of Build.
func (mr *MockExpressionMockRecorder) Build(builder any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockExpression)(nil).Build), builder)
}
