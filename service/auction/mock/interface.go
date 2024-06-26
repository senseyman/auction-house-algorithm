// Code generated by MockGen. DO NOT EDIT.
// Source: service/auction/interface.go
//
// Generated by this command:
//
//	mockgen -package mock -source service/auction/interface.go -destination service/auction/mock/interface.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/senseyman/auction-house/model"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// BidOrder mocks base method.
func (m *MockStorage) BidOrder(ctx context.Context, bid model.BidCommand) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BidOrder", ctx, bid)
	ret0, _ := ret[0].(error)
	return ret0
}

// BidOrder indicates an expected call of BidOrder.
func (mr *MockStorageMockRecorder) BidOrder(ctx, bid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BidOrder", reflect.TypeOf((*MockStorage)(nil).BidOrder), ctx, bid)
}

// CreateOrder mocks base method.
func (m *MockStorage) CreateOrder(ctx context.Context, order model.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockStorageMockRecorder) CreateOrder(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockStorage)(nil).CreateOrder), ctx, order)
}

// FinishAllAuctions mocks base method.
func (m *MockStorage) FinishAllAuctions(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishAllAuctions", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishAllAuctions indicates an expected call of FinishAllAuctions.
func (mr *MockStorageMockRecorder) FinishAllAuctions(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishAllAuctions", reflect.TypeOf((*MockStorage)(nil).FinishAllAuctions), arg0)
}

// FinishExpiredAuctions mocks base method.
func (m *MockStorage) FinishExpiredAuctions(ctx context.Context, timestamp int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishExpiredAuctions", ctx, timestamp)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishExpiredAuctions indicates an expected call of FinishExpiredAuctions.
func (mr *MockStorageMockRecorder) FinishExpiredAuctions(ctx, timestamp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishExpiredAuctions", reflect.TypeOf((*MockStorage)(nil).FinishExpiredAuctions), ctx, timestamp)
}

// GetAuctionResults mocks base method.
func (m *MockStorage) GetAuctionResults(ctx context.Context) ([]model.ActionResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuctionResults", ctx)
	ret0, _ := ret[0].([]model.ActionResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuctionResults indicates an expected call of GetAuctionResults.
func (mr *MockStorageMockRecorder) GetAuctionResults(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuctionResults", reflect.TypeOf((*MockStorage)(nil).GetAuctionResults), ctx)
}

// MockReadService is a mock of ReadService interface.
type MockReadService struct {
	ctrl     *gomock.Controller
	recorder *MockReadServiceMockRecorder
}

// MockReadServiceMockRecorder is the mock recorder for MockReadService.
type MockReadServiceMockRecorder struct {
	mock *MockReadService
}

// NewMockReadService creates a new mock instance.
func NewMockReadService(ctrl *gomock.Controller) *MockReadService {
	mock := &MockReadService{ctrl: ctrl}
	mock.recorder = &MockReadServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadService) EXPECT() *MockReadServiceMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockReadService) Read(filename string, outputCh chan model.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", filename, outputCh)
	ret0, _ := ret[0].(error)
	return ret0
}

// Read indicates an expected call of Read.
func (mr *MockReadServiceMockRecorder) Read(filename, outputCh any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReadService)(nil).Read), filename, outputCh)
}

// MockReportService is a mock of ReportService interface.
type MockReportService struct {
	ctrl     *gomock.Controller
	recorder *MockReportServiceMockRecorder
}

// MockReportServiceMockRecorder is the mock recorder for MockReportService.
type MockReportServiceMockRecorder struct {
	mock *MockReportService
}

// NewMockReportService creates a new mock instance.
func NewMockReportService(ctrl *gomock.Controller) *MockReportService {
	mock := &MockReportService{ctrl: ctrl}
	mock.recorder = &MockReportServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReportService) EXPECT() *MockReportServiceMockRecorder {
	return m.recorder
}

// Report mocks base method.
func (m *MockReportService) Report(fos []model.ActionResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Report", fos)
	ret0, _ := ret[0].(error)
	return ret0
}

// Report indicates an expected call of Report.
func (mr *MockReportServiceMockRecorder) Report(fos any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Report", reflect.TypeOf((*MockReportService)(nil).Report), fos)
}
