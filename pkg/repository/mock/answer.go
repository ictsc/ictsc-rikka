// Code generated by MockGen. DO NOT EDIT.
// Source: answer.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	entity "github.com/ictsc/ictsc-rikka/pkg/entity"
)

// MockAnswerRepository is a mock of AnswerRepository interface.
type MockAnswerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAnswerRepositoryMockRecorder
}

// MockAnswerRepositoryMockRecorder is the mock recorder for MockAnswerRepository.
type MockAnswerRepositoryMockRecorder struct {
	mock *MockAnswerRepository
}

// NewMockAnswerRepository creates a new mock instance.
func NewMockAnswerRepository(ctrl *gomock.Controller) *MockAnswerRepository {
	mock := &MockAnswerRepository{ctrl: ctrl}
	mock.recorder = &MockAnswerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAnswerRepository) EXPECT() *MockAnswerRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAnswerRepository) Create(answer *entity.Answer) (*entity.Answer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", answer)
	ret0, _ := ret[0].(*entity.Answer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAnswerRepositoryMockRecorder) Create(answer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAnswerRepository)(nil).Create), answer)
}

// Delete mocks base method.
func (m *MockAnswerRepository) Delete(answer *entity.Answer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", answer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAnswerRepositoryMockRecorder) Delete(answer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAnswerRepository)(nil).Delete), answer)
}

// FindByID mocks base method.
func (m *MockAnswerRepository) FindByID(id uuid.UUID) (*entity.Answer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*entity.Answer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockAnswerRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockAnswerRepository)(nil).FindByID), id)
}

// FindByProblem mocks base method.
func (m *MockAnswerRepository) FindByProblem(probid uuid.UUID, groupid *uuid.UUID) ([]*entity.Answer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByProblem", probid, groupid)
	ret0, _ := ret[0].([]*entity.Answer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByProblem indicates an expected call of FindByProblem.
func (mr *MockAnswerRepositoryMockRecorder) FindByProblem(probid, groupid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByProblem", reflect.TypeOf((*MockAnswerRepository)(nil).FindByProblem), probid, groupid)
}

// FindByProblemAndUserGroup mocks base method.
func (m *MockAnswerRepository) FindByProblemAndUserGroup(probid, groupid uuid.UUID) ([]*entity.Answer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByProblemAndUserGroup", probid, groupid)
	ret0, _ := ret[0].([]*entity.Answer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByProblemAndUserGroup indicates an expected call of FindByProblemAndUserGroup.
func (mr *MockAnswerRepositoryMockRecorder) FindByProblemAndUserGroup(probid, groupid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByProblemAndUserGroup", reflect.TypeOf((*MockAnswerRepository)(nil).FindByProblemAndUserGroup), probid, groupid)
}

// FindByUserGroup mocks base method.
func (m *MockAnswerRepository) FindByUserGroup(id uuid.UUID) ([]*entity.Answer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserGroup", id)
	ret0, _ := ret[0].([]*entity.Answer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserGroup indicates an expected call of FindByUserGroup.
func (mr *MockAnswerRepositoryMockRecorder) FindByUserGroup(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserGroup", reflect.TypeOf((*MockAnswerRepository)(nil).FindByUserGroup), id)
}

// GetAll mocks base method.
func (m *MockAnswerRepository) GetAll() ([]*entity.Answer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*entity.Answer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockAnswerRepositoryMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockAnswerRepository)(nil).GetAll))
}

// Update mocks base method.
func (m *MockAnswerRepository) Update(answer *entity.Answer) (*entity.Answer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", answer)
	ret0, _ := ret[0].(*entity.Answer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockAnswerRepositoryMockRecorder) Update(answer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAnswerRepository)(nil).Update), answer)
}
