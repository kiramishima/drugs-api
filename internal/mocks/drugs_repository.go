// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\interfaces\drugs_repository.go
//
// Generated by this command:
//
//	mockgen -source .\internal\interfaces\drugs_repository.go -destination .\internal\mocks\drugs_repository.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "kiramishima/ionix/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDrugRepository is a mock of DrugRepository interface.
type MockDrugRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDrugRepositoryMockRecorder
}

// MockDrugRepositoryMockRecorder is the mock recorder for MockDrugRepository.
type MockDrugRepositoryMockRecorder struct {
	mock *MockDrugRepository
}

// NewMockDrugRepository creates a new mock instance.
func NewMockDrugRepository(ctrl *gomock.Controller) *MockDrugRepository {
	mock := &MockDrugRepository{ctrl: ctrl}
	mock.recorder = &MockDrugRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDrugRepository) EXPECT() *MockDrugRepositoryMockRecorder {
	return m.recorder
}

// CreateNewDrugItem mocks base method.
func (m *MockDrugRepository) CreateNewDrugItem(ctx context.Context, form *models.DrugForm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewDrugItem", ctx, form)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewDrugItem indicates an expected call of CreateNewDrugItem.
func (mr *MockDrugRepositoryMockRecorder) CreateNewDrugItem(ctx, form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewDrugItem", reflect.TypeOf((*MockDrugRepository)(nil).CreateNewDrugItem), ctx, form)
}

// DeleteDrugItem mocks base method.
func (m *MockDrugRepository) DeleteDrugItem(ctx context.Context, drugId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDrugItem", ctx, drugId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDrugItem indicates an expected call of DeleteDrugItem.
func (mr *MockDrugRepositoryMockRecorder) DeleteDrugItem(ctx, drugId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDrugItem", reflect.TypeOf((*MockDrugRepository)(nil).DeleteDrugItem), ctx, drugId)
}

// GetDrugItemByID mocks base method.
func (m *MockDrugRepository) GetDrugItemByID(ctx context.Context, drugId int) (*models.Drug, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDrugItemByID", ctx, drugId)
	ret0, _ := ret[0].(*models.Drug)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDrugItemByID indicates an expected call of GetDrugItemByID.
func (mr *MockDrugRepositoryMockRecorder) GetDrugItemByID(ctx, drugId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDrugItemByID", reflect.TypeOf((*MockDrugRepository)(nil).GetDrugItemByID), ctx, drugId)
}

// GetDrugsData mocks base method.
func (m *MockDrugRepository) GetDrugsData(ctx context.Context) ([]*models.Drug, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDrugsData", ctx)
	ret0, _ := ret[0].([]*models.Drug)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDrugsData indicates an expected call of GetDrugsData.
func (mr *MockDrugRepositoryMockRecorder) GetDrugsData(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDrugsData", reflect.TypeOf((*MockDrugRepository)(nil).GetDrugsData), ctx)
}

// UpdateDrugItem mocks base method.
func (m *MockDrugRepository) UpdateDrugItem(ctx context.Context, drugId int, form *models.Drug) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDrugItem", ctx, drugId, form)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDrugItem indicates an expected call of UpdateDrugItem.
func (mr *MockDrugRepositoryMockRecorder) UpdateDrugItem(ctx, drugId, form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDrugItem", reflect.TypeOf((*MockDrugRepository)(nil).UpdateDrugItem), ctx, drugId, form)
}