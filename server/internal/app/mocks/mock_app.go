// Code generated by MockGen. DO NOT EDIT.
// Source: app.go

// Package mock_app is a generated GoMock package.
package mock_app

import (
	context "context"
	multipart "mime/multipart"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	app "github.com/saifahn/expenseus/internal/app"
	oauth2 "golang.org/x/oauth2"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateSharedTxn mocks base method.
func (m *MockStore) CreateSharedTxn(txn app.SharedTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSharedTxn", txn)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSharedTxn indicates an expected call of CreateSharedTxn.
func (mr *MockStoreMockRecorder) CreateSharedTxn(txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSharedTxn", reflect.TypeOf((*MockStore)(nil).CreateSharedTxn), txn)
}

// CreateTracker mocks base method.
func (m *MockStore) CreateTracker(tracker app.Tracker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTracker", tracker)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTracker indicates an expected call of CreateTracker.
func (mr *MockStoreMockRecorder) CreateTracker(tracker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTracker", reflect.TypeOf((*MockStore)(nil).CreateTracker), tracker)
}

// CreateTransaction mocks base method.
func (m *MockStore) CreateTransaction(transactionDetails app.TransactionDetails) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", transactionDetails)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockStoreMockRecorder) CreateTransaction(transactionDetails interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockStore)(nil).CreateTransaction), transactionDetails)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(user app.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), user)
}

// GetAllTransactions mocks base method.
func (m *MockStore) GetAllTransactions() ([]app.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTransactions")
	ret0, _ := ret[0].([]app.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTransactions indicates an expected call of GetAllTransactions.
func (mr *MockStoreMockRecorder) GetAllTransactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTransactions", reflect.TypeOf((*MockStore)(nil).GetAllTransactions))
}

// GetAllUsers mocks base method.
func (m *MockStore) GetAllUsers() ([]app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers")
	ret0, _ := ret[0].([]app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockStoreMockRecorder) GetAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockStore)(nil).GetAllUsers))
}

// GetTracker mocks base method.
func (m *MockStore) GetTracker(trackerID string) (app.Tracker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTracker", trackerID)
	ret0, _ := ret[0].(app.Tracker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTracker indicates an expected call of GetTracker.
func (mr *MockStoreMockRecorder) GetTracker(trackerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTracker", reflect.TypeOf((*MockStore)(nil).GetTracker), trackerID)
}

// GetTrackersByUser mocks base method.
func (m *MockStore) GetTrackersByUser(userID string) ([]app.Tracker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrackersByUser", userID)
	ret0, _ := ret[0].([]app.Tracker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrackersByUser indicates an expected call of GetTrackersByUser.
func (mr *MockStoreMockRecorder) GetTrackersByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrackersByUser", reflect.TypeOf((*MockStore)(nil).GetTrackersByUser), userID)
}

// GetTransaction mocks base method.
func (m *MockStore) GetTransaction(txnID string) (app.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", txnID)
	ret0, _ := ret[0].(app.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockStoreMockRecorder) GetTransaction(txnID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockStore)(nil).GetTransaction), txnID)
}

// GetTransactionsByUser mocks base method.
func (m *MockStore) GetTransactionsByUser(userID string) ([]app.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsByUser", userID)
	ret0, _ := ret[0].([]app.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsByUser indicates an expected call of GetTransactionsByUser.
func (mr *MockStoreMockRecorder) GetTransactionsByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsByUser", reflect.TypeOf((*MockStore)(nil).GetTransactionsByUser), userID)
}

// GetTxnsByTracker mocks base method.
func (m *MockStore) GetTxnsByTracker(trackerID string) ([]app.SharedTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxnsByTracker", trackerID)
	ret0, _ := ret[0].([]app.SharedTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxnsByTracker indicates an expected call of GetTxnsByTracker.
func (mr *MockStoreMockRecorder) GetTxnsByTracker(trackerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxnsByTracker", reflect.TypeOf((*MockStore)(nil).GetTxnsByTracker), trackerID)
}

// GetUnsettledTxnsByTracker mocks base method.
func (m *MockStore) GetUnsettledTxnsByTracker(trackerID string) ([]app.SharedTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnsettledTxnsByTracker", trackerID)
	ret0, _ := ret[0].([]app.SharedTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnsettledTxnsByTracker indicates an expected call of GetUnsettledTxnsByTracker.
func (mr *MockStoreMockRecorder) GetUnsettledTxnsByTracker(trackerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnsettledTxnsByTracker", reflect.TypeOf((*MockStore)(nil).GetUnsettledTxnsByTracker), trackerID)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(id string) (app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), id)
}

// SettleTxns mocks base method.
func (m *MockStore) SettleTxns(txns []app.SharedTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SettleTxns", txns)
	ret0, _ := ret[0].(error)
	return ret0
}

// SettleTxns indicates an expected call of SettleTxns.
func (mr *MockStoreMockRecorder) SettleTxns(txns interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SettleTxns", reflect.TypeOf((*MockStore)(nil).SettleTxns), txns)
}

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// AuthCodeURL mocks base method.
func (m *MockAuth) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	m.ctrl.T.Helper()
	varargs := []interface{}{state}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AuthCodeURL", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// AuthCodeURL indicates an expected call of AuthCodeURL.
func (mr *MockAuthMockRecorder) AuthCodeURL(state interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{state}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthCodeURL", reflect.TypeOf((*MockAuth)(nil).AuthCodeURL), varargs...)
}

// Exchange mocks base method.
func (m *MockAuth) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, code}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exchange", varargs...)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exchange indicates an expected call of Exchange.
func (mr *MockAuthMockRecorder) Exchange(ctx, code interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, code}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exchange", reflect.TypeOf((*MockAuth)(nil).Exchange), varargs...)
}

// GetInfoAndGenerateUser mocks base method.
func (m *MockAuth) GetInfoAndGenerateUser(state, code string) (app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfoAndGenerateUser", state, code)
	ret0, _ := ret[0].(app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfoAndGenerateUser indicates an expected call of GetInfoAndGenerateUser.
func (mr *MockAuthMockRecorder) GetInfoAndGenerateUser(state, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfoAndGenerateUser", reflect.TypeOf((*MockAuth)(nil).GetInfoAndGenerateUser), state, code)
}

// MockSessionManager is a mock of SessionManager interface.
type MockSessionManager struct {
	ctrl     *gomock.Controller
	recorder *MockSessionManagerMockRecorder
}

// MockSessionManagerMockRecorder is the mock recorder for MockSessionManager.
type MockSessionManagerMockRecorder struct {
	mock *MockSessionManager
}

// NewMockSessionManager creates a new mock instance.
func NewMockSessionManager(ctrl *gomock.Controller) *MockSessionManager {
	mock := &MockSessionManager{ctrl: ctrl}
	mock.recorder = &MockSessionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionManager) EXPECT() *MockSessionManagerMockRecorder {
	return m.recorder
}

// GetUserID mocks base method.
func (m *MockSessionManager) GetUserID(r *http.Request) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserID", r)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserID indicates an expected call of GetUserID.
func (mr *MockSessionManagerMockRecorder) GetUserID(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserID", reflect.TypeOf((*MockSessionManager)(nil).GetUserID), r)
}

// Remove mocks base method.
func (m *MockSessionManager) Remove(rw http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Remove", rw, r)
}

// Remove indicates an expected call of Remove.
func (mr *MockSessionManagerMockRecorder) Remove(rw, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockSessionManager)(nil).Remove), rw, r)
}

// Save mocks base method.
func (m *MockSessionManager) Save(rw http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Save", rw, r)
}

// Save indicates an expected call of Save.
func (mr *MockSessionManagerMockRecorder) Save(rw, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSessionManager)(nil).Save), rw, r)
}

// Validate mocks base method.
func (m *MockSessionManager) Validate(r *http.Request) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", r)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockSessionManagerMockRecorder) Validate(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockSessionManager)(nil).Validate), r)
}

// MockImageStore is a mock of ImageStore interface.
type MockImageStore struct {
	ctrl     *gomock.Controller
	recorder *MockImageStoreMockRecorder
}

// MockImageStoreMockRecorder is the mock recorder for MockImageStore.
type MockImageStoreMockRecorder struct {
	mock *MockImageStore
}

// NewMockImageStore creates a new mock instance.
func NewMockImageStore(ctrl *gomock.Controller) *MockImageStore {
	mock := &MockImageStore{ctrl: ctrl}
	mock.recorder = &MockImageStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageStore) EXPECT() *MockImageStoreMockRecorder {
	return m.recorder
}

// AddImageToTransaction mocks base method.
func (m *MockImageStore) AddImageToTransaction(transaction app.Transaction) (app.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImageToTransaction", transaction)
	ret0, _ := ret[0].(app.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddImageToTransaction indicates an expected call of AddImageToTransaction.
func (mr *MockImageStoreMockRecorder) AddImageToTransaction(transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImageToTransaction", reflect.TypeOf((*MockImageStore)(nil).AddImageToTransaction), transaction)
}

// Upload mocks base method.
func (m *MockImageStore) Upload(file multipart.File, header multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", file, header)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockImageStoreMockRecorder) Upload(file, header interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockImageStore)(nil).Upload), file, header)
}

// Validate mocks base method.
func (m *MockImageStore) Validate(file multipart.File) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", file)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockImageStoreMockRecorder) Validate(file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockImageStore)(nil).Validate), file)
}
