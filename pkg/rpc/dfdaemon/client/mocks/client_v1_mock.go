// Code generated by MockGen. DO NOT EDIT.
// Source: client_v1.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	common "d7y.io/api/v2/pkg/apis/common/v1"
	dfdaemon "d7y.io/api/v2/pkg/apis/dfdaemon/v1"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockV1 is a mock of V1 interface.
type MockV1 struct {
	ctrl     *gomock.Controller
	recorder *MockV1MockRecorder
}

// MockV1MockRecorder is the mock recorder for MockV1.
type MockV1MockRecorder struct {
	mock *MockV1
}

// NewMockV1 creates a new mock instance.
func NewMockV1(ctrl *gomock.Controller) *MockV1 {
	mock := &MockV1{ctrl: ctrl}
	mock.recorder = &MockV1MockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockV1) EXPECT() *MockV1MockRecorder {
	return m.recorder
}

// CheckHealth mocks base method.
func (m *MockV1) CheckHealth(arg0 context.Context, arg1 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckHealth", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckHealth indicates an expected call of CheckHealth.
func (mr *MockV1MockRecorder) CheckHealth(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckHealth", reflect.TypeOf((*MockV1)(nil).CheckHealth), varargs...)
}

// Close mocks base method.
func (m *MockV1) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockV1MockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockV1)(nil).Close))
}

// DeleteTask mocks base method.
func (m *MockV1) DeleteTask(arg0 context.Context, arg1 *dfdaemon.DeleteTaskRequest, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteTask", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockV1MockRecorder) DeleteTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockV1)(nil).DeleteTask), varargs...)
}

// Download mocks base method.
func (m *MockV1) Download(arg0 context.Context, arg1 *dfdaemon.DownRequest, arg2 ...grpc.CallOption) (dfdaemon.Daemon_DownloadClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Download", varargs...)
	ret0, _ := ret[0].(dfdaemon.Daemon_DownloadClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockV1MockRecorder) Download(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockV1)(nil).Download), varargs...)
}

// ExportTask mocks base method.
func (m *MockV1) ExportTask(arg0 context.Context, arg1 *dfdaemon.ExportTaskRequest, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExportTask", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExportTask indicates an expected call of ExportTask.
func (mr *MockV1MockRecorder) ExportTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportTask", reflect.TypeOf((*MockV1)(nil).ExportTask), varargs...)
}

// GetPieceTasks mocks base method.
func (m *MockV1) GetPieceTasks(arg0 context.Context, arg1 *common.PieceTaskRequest, arg2 ...grpc.CallOption) (*common.PiecePacket, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPieceTasks", varargs...)
	ret0, _ := ret[0].(*common.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieceTasks indicates an expected call of GetPieceTasks.
func (mr *MockV1MockRecorder) GetPieceTasks(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieceTasks", reflect.TypeOf((*MockV1)(nil).GetPieceTasks), varargs...)
}

// ImportTask mocks base method.
func (m *MockV1) ImportTask(arg0 context.Context, arg1 *dfdaemon.ImportTaskRequest, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ImportTask", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ImportTask indicates an expected call of ImportTask.
func (mr *MockV1MockRecorder) ImportTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportTask", reflect.TypeOf((*MockV1)(nil).ImportTask), varargs...)
}

// StatTask mocks base method.
func (m *MockV1) StatTask(arg0 context.Context, arg1 *dfdaemon.StatTaskRequest, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StatTask", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// StatTask indicates an expected call of StatTask.
func (mr *MockV1MockRecorder) StatTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatTask", reflect.TypeOf((*MockV1)(nil).StatTask), varargs...)
}

// SyncPieceTasks mocks base method.
func (m *MockV1) SyncPieceTasks(arg0 context.Context, arg1 *common.PieceTaskRequest, arg2 ...grpc.CallOption) (dfdaemon.Daemon_SyncPieceTasksClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SyncPieceTasks", varargs...)
	ret0, _ := ret[0].(dfdaemon.Daemon_SyncPieceTasksClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncPieceTasks indicates an expected call of SyncPieceTasks.
func (mr *MockV1MockRecorder) SyncPieceTasks(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPieceTasks", reflect.TypeOf((*MockV1)(nil).SyncPieceTasks), varargs...)
}
