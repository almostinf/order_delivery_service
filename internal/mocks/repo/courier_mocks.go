// Code generated by MockGen. DO NOT EDIT.
// Source: internal/infrastructure/interfaces/courier.go

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	context "context"
	reflect "reflect"
	time "time"

	entity "github.com/almostinf/order_delivery_service/internal/entity"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockCourier is a mock of Courier interface.
type MockCourier struct {
	ctrl     *gomock.Controller
	recorder *MockCourierMockRecorder
}

// MockCourierMockRecorder is the mock recorder for MockCourier.
type MockCourierMockRecorder struct {
	mock *MockCourier
}

// NewMockCourier creates a new mock instance.
func NewMockCourier(ctrl *gomock.Controller) *MockCourier {
	mock := &MockCourier{ctrl: ctrl}
	mock.recorder = &MockCourierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCourier) EXPECT() *MockCourierMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCourier) Create(ctx context.Context, courier *entity.Courier) (*entity.CourierResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, courier)
	ret0, _ := ret[0].(*entity.CourierResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCourierMockRecorder) Create(ctx, courier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCourier)(nil).Create), ctx, courier)
}

// Get mocks base method.
func (m *MockCourier) Get(ctx context.Context, id uuid.UUID) (*entity.CourierResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.CourierResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCourierMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCourier)(nil).Get), ctx, id)
}

// GetAll mocks base method.
func (m *MockCourier) GetAll(ctx context.Context, limit, offset int) ([]*entity.CourierResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*entity.CourierResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockCourierMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockCourier)(nil).GetAll), ctx, limit, offset)
}

// GetAssignments mocks base method.
func (m *MockCourier) GetAssignments(ctx context.Context, date time.Time, courierID uuid.UUID, isAllCouriers bool) ([]*entity.CourierAssignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssignments", ctx, date, courierID, isAllCouriers)
	ret0, _ := ret[0].([]*entity.CourierAssignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssignments indicates an expected call of GetAssignments.
func (mr *MockCourierMockRecorder) GetAssignments(ctx, date, courierID, isAllCouriers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssignments", reflect.TypeOf((*MockCourier)(nil).GetAssignments), ctx, date, courierID, isAllCouriers)
}

// GetMetaInfo mocks base method.
func (m *MockCourier) GetMetaInfo(ctx context.Context, courierID uuid.UUID, startDate, endDate time.Time) (*entity.CourierMetaInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetaInfo", ctx, courierID, startDate, endDate)
	ret0, _ := ret[0].(*entity.CourierMetaInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetaInfo indicates an expected call of GetMetaInfo.
func (mr *MockCourierMockRecorder) GetMetaInfo(ctx, courierID, startDate, endDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetaInfo", reflect.TypeOf((*MockCourier)(nil).GetMetaInfo), ctx, courierID, startDate, endDate)
}
