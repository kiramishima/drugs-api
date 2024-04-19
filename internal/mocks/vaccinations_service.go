// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\interfaces\vaccinations_service.go
//
// Generated by this command:
//
//	mockgen -source .\internal\interfaces\vaccinations_service.go -destination .\internal\mocks\vaccinations_service.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "kiramishima/ionix/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockVaccinationService is a mock of VaccinationService interface.
type MockVaccinationService struct {
	ctrl     *gomock.Controller
	recorder *MockVaccinationServiceMockRecorder
}

// MockVaccinationServiceMockRecorder is the mock recorder for MockVaccinationService.
type MockVaccinationServiceMockRecorder struct {
	mock *MockVaccinationService
}

// NewMockVaccinationService creates a new mock instance.
func NewMockVaccinationService(ctrl *gomock.Controller) *MockVaccinationService {
	mock := &MockVaccinationService{ctrl: ctrl}
	mock.recorder = &MockVaccinationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVaccinationService) EXPECT() *MockVaccinationServiceMockRecorder {
	return m.recorder
}

// DeleteVaccination mocks base method.
func (m *MockVaccinationService) DeleteVaccination(ctx context.Context, vaccinationId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVaccination", ctx, vaccinationId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVaccination indicates an expected call of DeleteVaccination.
func (mr *MockVaccinationServiceMockRecorder) DeleteVaccination(ctx, vaccinationId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVaccination", reflect.TypeOf((*MockVaccinationService)(nil).DeleteVaccination), ctx, vaccinationId)
}

// GetListVaccinations mocks base method.
func (m *MockVaccinationService) GetListVaccinations(ctx context.Context) ([]*models.Vaccination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListVaccinations", ctx)
	ret0, _ := ret[0].([]*models.Vaccination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListVaccinations indicates an expected call of GetListVaccinations.
func (mr *MockVaccinationServiceMockRecorder) GetListVaccinations(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListVaccinations", reflect.TypeOf((*MockVaccinationService)(nil).GetListVaccinations), ctx)
}

// NewVaccination mocks base method.
func (m *MockVaccinationService) NewVaccination(ctx context.Context, form *models.VaccinationForm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewVaccination", ctx, form)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewVaccination indicates an expected call of NewVaccination.
func (mr *MockVaccinationServiceMockRecorder) NewVaccination(ctx, form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewVaccination", reflect.TypeOf((*MockVaccinationService)(nil).NewVaccination), ctx, form)
}

// UpdateVaccination mocks base method.
func (m *MockVaccinationService) UpdateVaccination(ctx context.Context, vaccinationId int, form *models.VaccinationForm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVaccination", ctx, vaccinationId, form)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVaccination indicates an expected call of UpdateVaccination.
func (mr *MockVaccinationServiceMockRecorder) UpdateVaccination(ctx, vaccinationId, form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVaccination", reflect.TypeOf((*MockVaccinationService)(nil).UpdateVaccination), ctx, vaccinationId, form)
}
