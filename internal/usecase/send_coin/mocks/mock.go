// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "AvitoTask/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v5"
)

// Mockuser is a mock of user interface.
type Mockuser struct {
	ctrl     *gomock.Controller
	recorder *MockuserMockRecorder
}

// MockuserMockRecorder is the mock recorder for Mockuser.
type MockuserMockRecorder struct {
	mock *Mockuser
}

// NewMockuser creates a new mock instance.
func NewMockuser(ctrl *gomock.Controller) *Mockuser {
	mock := &Mockuser{ctrl: ctrl}
	mock.recorder = &MockuserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockuser) EXPECT() *MockuserMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *Mockuser) BeginTx(ctx context.Context) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockuserMockRecorder) BeginTx(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*Mockuser)(nil).BeginTx), ctx)
}

// GetUserById mocks base method.
func (m *Mockuser) GetUserById(ctx context.Context, tx pgx.Tx, userID string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", ctx, tx, userID)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockuserMockRecorder) GetUserById(ctx, tx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*Mockuser)(nil).GetUserById), ctx, tx, userID)
}

// GetUserByLoginWithTx mocks base method.
func (m *Mockuser) GetUserByLoginWithTx(ctx context.Context, tx pgx.Tx, login string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLoginWithTx", ctx, tx, login)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLoginWithTx indicates an expected call of GetUserByLoginWithTx.
func (mr *MockuserMockRecorder) GetUserByLoginWithTx(ctx, tx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLoginWithTx", reflect.TypeOf((*Mockuser)(nil).GetUserByLoginWithTx), ctx, tx, login)
}

// IsUserExists mocks base method.
func (m *Mockuser) IsUserExists(ctx context.Context, user models.User) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserExists", ctx, user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserExists indicates an expected call of IsUserExists.
func (mr *MockuserMockRecorder) IsUserExists(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserExists", reflect.TypeOf((*Mockuser)(nil).IsUserExists), ctx, user)
}

// UpdateUserCoins mocks base method.
func (m *Mockuser) UpdateUserCoins(ctx context.Context, tx pgx.Tx, userID string, newCoins int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserCoins", ctx, tx, userID, newCoins)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserCoins indicates an expected call of UpdateUserCoins.
func (mr *MockuserMockRecorder) UpdateUserCoins(ctx, tx, userID, newCoins interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserCoins", reflect.TypeOf((*Mockuser)(nil).UpdateUserCoins), ctx, tx, userID, newCoins)
}

// Mocktransaction is a mock of transaction interface.
type Mocktransaction struct {
	ctrl     *gomock.Controller
	recorder *MocktransactionMockRecorder
}

// MocktransactionMockRecorder is the mock recorder for Mocktransaction.
type MocktransactionMockRecorder struct {
	mock *Mocktransaction
}

// NewMocktransaction creates a new mock instance.
func NewMocktransaction(ctrl *gomock.Controller) *Mocktransaction {
	mock := &Mocktransaction{ctrl: ctrl}
	mock.recorder = &MocktransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktransaction) EXPECT() *MocktransactionMockRecorder {
	return m.recorder
}

// InsertTransaction mocks base method.
func (m *Mocktransaction) InsertTransaction(ctx context.Context, tx pgx.Tx, id, fromUserID, toUserID string, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransaction", ctx, tx, id, fromUserID, toUserID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertTransaction indicates an expected call of InsertTransaction.
func (mr *MocktransactionMockRecorder) InsertTransaction(ctx, tx, id, fromUserID, toUserID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransaction", reflect.TypeOf((*Mocktransaction)(nil).InsertTransaction), ctx, tx, id, fromUserID, toUserID, amount)
}
