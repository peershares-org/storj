// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

// Code generated by MockGen. DO NOT EDIT.
// Source: czarcoin.org/czarcoin/pkg/pb (interfaces: PointerDBClient)

// Package pdbclient is a generated GoMock package.
package pdbclient

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"

	pb "czarcoin.org/czarcoin/pkg/pb"
)

// MockPointerDBClient is a mock of PointerDBClient interface
type MockPointerDBClient struct {
	ctrl     *gomock.Controller
	recorder *MockPointerDBClientMockRecorder
}

// MockPointerDBClientMockRecorder is the mock recorder for MockPointerDBClient
type MockPointerDBClientMockRecorder struct {
	mock *MockPointerDBClient
}

// NewMockPointerDBClient creates a new mock instance
func NewMockPointerDBClient(ctrl *gomock.Controller) *MockPointerDBClient {
	mock := &MockPointerDBClient{ctrl: ctrl}
	mock.recorder = &MockPointerDBClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPointerDBClient) EXPECT() *MockPointerDBClientMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockPointerDBClient) Delete(arg0 context.Context, arg1 *pb.DeleteRequest, arg2 ...grpc.CallOption) (*pb.DeleteResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(*pb.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete
func (mr *MockPointerDBClientMockRecorder) Delete(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPointerDBClient)(nil).Delete), varargs...)
}

// Get mocks base method
func (m *MockPointerDBClient) Get(arg0 context.Context, arg1 *pb.GetRequest, arg2 ...grpc.CallOption) (*pb.GetResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*pb.GetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockPointerDBClientMockRecorder) Get(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPointerDBClient)(nil).Get), varargs...)
}

// List mocks base method
func (m *MockPointerDBClient) List(arg0 context.Context, arg1 *pb.ListRequest, arg2 ...grpc.CallOption) (*pb.ListResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(*pb.ListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockPointerDBClientMockRecorder) List(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPointerDBClient)(nil).List), varargs...)
}

// PayerBandwidthAllocation mocks base method
func (m *MockPointerDBClient) PayerBandwidthAllocation(arg0 context.Context, arg1 *pb.PayerBandwidthAllocationRequest, arg2 ...grpc.CallOption) (*pb.PayerBandwidthAllocationResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PayerBandwidthAllocation", varargs...)
	ret0, _ := ret[0].(*pb.PayerBandwidthAllocationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PayerBandwidthAllocation indicates an expected call of PayerBandwidthAllocation
func (mr *MockPointerDBClientMockRecorder) PayerBandwidthAllocation(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PayerBandwidthAllocation", reflect.TypeOf((*MockPointerDBClient)(nil).PayerBandwidthAllocation), varargs...)
}

// Put mocks base method
func (m *MockPointerDBClient) Put(arg0 context.Context, arg1 *pb.PutRequest, arg2 ...grpc.CallOption) (*pb.PutResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Put", varargs...)
	ret0, _ := ret[0].(*pb.PutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put
func (mr *MockPointerDBClientMockRecorder) Put(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockPointerDBClient)(nil).Put), varargs...)
}
