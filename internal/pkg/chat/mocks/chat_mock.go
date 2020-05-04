// Code generated by MockGen. DO NOT EDIT.
// Source: failless/internal/pkg/chat (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	forms "failless/internal/pkg/forms"
	models "failless/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddMessageToChat mocks base method
func (m *MockRepository) AddMessageToChat(arg0 *forms.Message, arg1 []int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMessageToChat", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMessageToChat indicates an expected call of AddMessageToChat
func (mr *MockRepositoryMockRecorder) AddMessageToChat(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessageToChat", reflect.TypeOf((*MockRepository)(nil).AddMessageToChat), arg0, arg1)
}

// CheckRoom mocks base method
func (m *MockRepository) CheckRoom(arg0, arg1 int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRoom", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRoom indicates an expected call of CheckRoom
func (mr *MockRepositoryMockRecorder) CheckRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRoom", reflect.TypeOf((*MockRepository)(nil).CheckRoom), arg0, arg1)
}

// GetRoomMessages mocks base method
func (m *MockRepository) GetRoomMessages(arg0, arg1 int64, arg2, arg3 int) ([]forms.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoomMessages", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]forms.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoomMessages indicates an expected call of GetRoomMessages
func (mr *MockRepositoryMockRecorder) GetRoomMessages(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoomMessages", reflect.TypeOf((*MockRepository)(nil).GetRoomMessages), arg0, arg1, arg2, arg3)
}

// GetUserTopMessages mocks base method
func (m *MockRepository) GetUserTopMessages(arg0 int64, arg1, arg2 int) ([]models.ChatMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserTopMessages", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.ChatMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserTopMessages indicates an expected call of GetUserTopMessages
func (mr *MockRepositoryMockRecorder) GetUserTopMessages(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserTopMessages", reflect.TypeOf((*MockRepository)(nil).GetUserTopMessages), arg0, arg1, arg2)
}

// GetUsersRooms mocks base method
func (m *MockRepository) GetUsersRooms(arg0 int64) ([]models.ChatRoom, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersRooms", arg0)
	ret0, _ := ret[0].([]models.ChatRoom)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersRooms indicates an expected call of GetUsersRooms
func (mr *MockRepositoryMockRecorder) GetUsersRooms(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersRooms", reflect.TypeOf((*MockRepository)(nil).GetUsersRooms), arg0)
}

// InsertDialogue mocks base method
func (m *MockRepository) InsertDialogue(arg0, arg1, arg2 int, arg3 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertDialogue", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertDialogue indicates an expected call of InsertDialogue
func (mr *MockRepositoryMockRecorder) InsertDialogue(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertDialogue", reflect.TypeOf((*MockRepository)(nil).InsertDialogue), arg0, arg1, arg2, arg3)
}
