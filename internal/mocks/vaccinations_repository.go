// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\interfaces\vaccinations_repository.go
//
// Generated by this command:
//
//	mockgen -source .\internal\interfaces\vaccinations_repository.go -destination .\internal\mocks\vaccinations_repository.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "kiramishima/ionix/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockVaccinationRepository is a mock of VaccinationRepository interface.
type MockVaccinationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVaccinationRepositoryMockRecorder
}

// MockVaccinationRepositoryMockRecorder is the mock recorder for MockVaccinationRepository.
type MockVaccinationRepositoryMockRecorder struct {
	mock *MockVaccinationRepository
}

// NewMockVaccinationRepository creates a new mock instance.
func NewMockVaccinationRepository(ctrl *gomock.Controller) *MockVaccinationRepository {
	mock := &MockVaccinationRepository{ctrl: ctrl}
	mock.recorder = &MockVaccinationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVaccinationRepository) EXPECT() *MockVaccinationRepositoryMockRecorder {
	return m.recorder
}

// CreateNewVaccinationItem mocks base method.
func (m *MockVaccinationRepository) CreateNewVaccinationItem(ctx context.Context, form *models.VaccinationForm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewVaccinationItem", ctx, form)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewVaccinationItem indicates an expected call of CreateNewVaccinationItem.
func (mr *MockVaccinationRepositoryMockRecorder) CreateNewVaccinationItem(ctx, form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewVaccinationItem", reflect.TypeOf((*MockVaccinationRepository)(nil).CreateNewVaccinationItem), ctx, form)
}

// DeleteVaccinationItem mocks base method.
func (m *MockVaccinationRepository) DeleteVaccinationItem(ctx context.Context, vaccinationId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVaccinationItem", ctx, vaccinationId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVaccinationItem indicates an expected call of DeleteVaccinationItem.
func (mr *MockVaccinationRepositoryMockRecorder) DeleteVaccinationItem(ctx, vaccinationId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVaccinationItem", reflect.TypeOf((*MockVaccinationRepository)(nil).DeleteVaccinationItem), ctx, vaccinationId)
}

// GetVaccinationItemByID mocks base method.
func (m *MockVaccinationRepository) GetVaccinationItemByID(ctx context.Context, vaccinationId int) (*models.Vaccination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaccinationItemByID", ctx, vaccinationId)
	ret0, _ := ret[0].(*models.Vaccination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaccinationItemByID indicates an expected call of GetVaccinationItemByID.
func (mr *MockVaccinationRepositoryMockRecorder) GetVaccinationItemByID(ctx, vaccinationId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaccinationItemByID", reflect.TypeOf((*MockVaccinationRepository)(nil).GetVaccinationItemByID), ctx, vaccinationId)
}

// GetVaccinationsData mocks base method.
func (m *MockVaccinationRepository) GetVaccinationsData(ctx context.Context) ([]*models.Vaccination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaccinationsData", ctx)
	ret0, _ := ret[0].([]*models.Vaccination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaccinationsData indicates an expected call of GetVaccinationsData.
func (mr *MockVaccinationRepositoryMockRecorder) GetVaccinationsData(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaccinationsData", reflect.TypeOf((*MockVaccinationRepository)(nil).GetVaccinationsData), ctx)
}

// UpdateVaccinationItem mocks base method.
func (m *MockVaccinationRepository) UpdateVaccinationItem(ctx context.Context, vaccinationId int, form *models.Vaccination) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVaccinationItem", ctx, vaccinationId, form)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVaccinationItem indicates an expected call of UpdateVaccinationItem.
func (mr *MockVaccinationRepositoryMockRecorder) UpdateVaccinationItem(ctx, vaccinationId, form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVaccinationItem", reflect.TypeOf((*MockVaccinationRepository)(nil).UpdateVaccinationItem), ctx, vaccinationId, form)
}
