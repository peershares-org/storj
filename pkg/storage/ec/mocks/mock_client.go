// Code generated by MockGen. DO NOT EDIT.
// Source: czarcoin.org/czarcoin/pkg/storage/ec (interfaces: Client)

// Package mock_ecclient is a generated GoMock package.
package mock_ecclient

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	eestream "czarcoin.org/czarcoin/pkg/eestream"
	pb "czarcoin.org/czarcoin/pkg/pb"
	client "czarcoin.org/czarcoin/pkg/piecestore/psclient"
	ranger "czarcoin.org/czarcoin/pkg/ranger"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockClient) Delete(arg0 context.Context, arg1 []*pb.Node, arg2 client.PieceID, arg3 *pb.SignedMessage) error {
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockClientMockRecorder) Delete(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockClient)(nil).Delete), arg0, arg1, arg2, arg3)
}

// Get mocks base method
func (m *MockClient) Get(arg0 context.Context, arg1 []*pb.Node, arg2 eestream.ErasureScheme, arg3 client.PieceID, arg4 int64, arg5 *pb.PayerBandwidthAllocation, arg6 *pb.SignedMessage) (ranger.Ranger, error) {
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(ranger.Ranger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockClientMockRecorder) Get(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// Put mocks base method
func (m *MockClient) Put(arg0 context.Context, arg1 []*pb.Node, arg2 eestream.RedundancyStrategy, arg3 client.PieceID, arg4 io.Reader, arg5 time.Time, arg6 *pb.PayerBandwidthAllocation, arg7 *pb.SignedMessage) ([]*pb.Node, error) {
	ret := m.ctrl.Call(m, "Put", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	ret0, _ := ret[0].([]*pb.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put
func (mr *MockClientMockRecorder) Put(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockClient)(nil).Put), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}
