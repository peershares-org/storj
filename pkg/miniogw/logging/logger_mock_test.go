// Code generated by MockGen. DO NOT EDIT.
// Source: czarcoin.org/czarcoin/pkg/miniogw/logging (interfaces: ErrorLogger)

// Package logging is a generated GoMock package.
package logging

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockErrorLogger is a mock of ErrorLogger interface
type MockErrorLogger struct {
	ctrl     *gomock.Controller
	recorder *MockErrorLoggerMockRecorder
}

// MockErrorLoggerMockRecorder is the mock recorder for MockErrorLogger
type MockErrorLoggerMockRecorder struct {
	mock *MockErrorLogger
}

// NewMockErrorLogger creates a new mock instance
func NewMockErrorLogger(ctrl *gomock.Controller) *MockErrorLogger {
	mock := &MockErrorLogger{ctrl: ctrl}
	mock.recorder = &MockErrorLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockErrorLogger) EXPECT() *MockErrorLoggerMockRecorder {
	return m.recorder
}

// Debugf mocks base method
func (m *MockErrorLogger) Debugf(arg0 string, arg1 ...interface{}) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf
func (mr *MockErrorLoggerMockRecorder) Debugf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*MockErrorLogger)(nil).Debugf), varargs...)
}

// Errorf mocks base method
func (m *MockErrorLogger) Errorf(arg0 string, arg1 ...interface{}) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf
func (mr *MockErrorLoggerMockRecorder) Errorf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockErrorLogger)(nil).Errorf), varargs...)
}
