// Code generated by MockGen. DO NOT EDIT.
// Source: ../../../../../pkg/depends/kit/sqlx/database.go

// Package mock_sqlx is a generated GoMock package.
package mock_sqlx

import (
	context "context"
	sql "database/sql"
	driver "database/sql/driver"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	sqlx "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	builder "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

// MockDBExecutor is a mock of DBExecutor interface.
type MockDBExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockDBExecutorMockRecorder
}

// MockDBExecutorMockRecorder is the mock recorder for MockDBExecutor.
type MockDBExecutorMockRecorder struct {
	mock *MockDBExecutor
}

// NewMockDBExecutor creates a new mock instance.
func NewMockDBExecutor(ctrl *gomock.Controller) *MockDBExecutor {
	mock := &MockDBExecutor{ctrl: ctrl}
	mock.recorder = &MockDBExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBExecutor) EXPECT() *MockDBExecutorMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockDBExecutor) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDBExecutorMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDBExecutor)(nil).Context))
}

// D mocks base method.
func (m *MockDBExecutor) D() *sqlx.Database {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "D")
	ret0, _ := ret[0].(*sqlx.Database)
	return ret0
}

// D indicates an expected call of D.
func (mr *MockDBExecutorMockRecorder) D() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "D", reflect.TypeOf((*MockDBExecutor)(nil).D))
}

// Dialect mocks base method.
func (m *MockDBExecutor) Dialect() builder.Dialect {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dialect")
	ret0, _ := ret[0].(builder.Dialect)
	return ret0
}

// Dialect indicates an expected call of Dialect.
func (mr *MockDBExecutorMockRecorder) Dialect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dialect", reflect.TypeOf((*MockDBExecutor)(nil).Dialect))
}

// Exec mocks base method.
func (m *MockDBExecutor) Exec(arg0 builder.SqlExpr) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", arg0)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockDBExecutorMockRecorder) Exec(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDBExecutor)(nil).Exec), arg0)
}

// ExecContext mocks base method.
func (m *MockDBExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockDBExecutorMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockDBExecutor)(nil).ExecContext), varargs...)
}

// Query mocks base method.
func (m *MockDBExecutor) Query(arg0 builder.SqlExpr) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockDBExecutorMockRecorder) Query(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDBExecutor)(nil).Query), arg0)
}

// QueryAndScan mocks base method.
func (m *MockDBExecutor) QueryAndScan(arg0 builder.SqlExpr, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAndScan", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueryAndScan indicates an expected call of QueryAndScan.
func (mr *MockDBExecutorMockRecorder) QueryAndScan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAndScan", reflect.TypeOf((*MockDBExecutor)(nil).QueryAndScan), arg0, arg1)
}

// QueryContext mocks base method.
func (m *MockDBExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockDBExecutorMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockDBExecutor)(nil).QueryContext), varargs...)
}

// T mocks base method.
func (m *MockDBExecutor) T(model builder.Model) *builder.Table {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "T", model)
	ret0, _ := ret[0].(*builder.Table)
	return ret0
}

// T indicates an expected call of T.
func (mr *MockDBExecutorMockRecorder) T(model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "T", reflect.TypeOf((*MockDBExecutor)(nil).T), model)
}

// WithContext mocks base method.
func (m *MockDBExecutor) WithContext(ctx context.Context) sqlx.DBExecutor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithContext", ctx)
	ret0, _ := ret[0].(sqlx.DBExecutor)
	return ret0
}

// WithContext indicates an expected call of WithContext.
func (mr *MockDBExecutorMockRecorder) WithContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithContext", reflect.TypeOf((*MockDBExecutor)(nil).WithContext), ctx)
}

// WithSchema mocks base method.
func (m *MockDBExecutor) WithSchema(arg0 string) sqlx.DBExecutor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithSchema", arg0)
	ret0, _ := ret[0].(sqlx.DBExecutor)
	return ret0
}

// WithSchema indicates an expected call of WithSchema.
func (mr *MockDBExecutorMockRecorder) WithSchema(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithSchema", reflect.TypeOf((*MockDBExecutor)(nil).WithSchema), arg0)
}

// MockWithDBName is a mock of WithDBName interface.
type MockWithDBName struct {
	ctrl     *gomock.Controller
	recorder *MockWithDBNameMockRecorder
}

// MockWithDBNameMockRecorder is the mock recorder for MockWithDBName.
type MockWithDBNameMockRecorder struct {
	mock *MockWithDBName
}

// NewMockWithDBName creates a new mock instance.
func NewMockWithDBName(ctrl *gomock.Controller) *MockWithDBName {
	mock := &MockWithDBName{ctrl: ctrl}
	mock.recorder = &MockWithDBNameMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithDBName) EXPECT() *MockWithDBNameMockRecorder {
	return m.recorder
}

// WithDBName mocks base method.
func (m *MockWithDBName) WithDBName(arg0 string) driver.Connector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithDBName", arg0)
	ret0, _ := ret[0].(driver.Connector)
	return ret0
}

// WithDBName indicates an expected call of WithDBName.
func (mr *MockWithDBNameMockRecorder) WithDBName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithDBName", reflect.TypeOf((*MockWithDBName)(nil).WithDBName), arg0)
}

// MockSqlExecutor is a mock of SqlExecutor interface.
type MockSqlExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockSqlExecutorMockRecorder
}

// MockSqlExecutorMockRecorder is the mock recorder for MockSqlExecutor.
type MockSqlExecutorMockRecorder struct {
	mock *MockSqlExecutor
}

// NewMockSqlExecutor creates a new mock instance.
func NewMockSqlExecutor(ctrl *gomock.Controller) *MockSqlExecutor {
	mock := &MockSqlExecutor{ctrl: ctrl}
	mock.recorder = &MockSqlExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSqlExecutor) EXPECT() *MockSqlExecutorMockRecorder {
	return m.recorder
}

// ExecContext mocks base method.
func (m *MockSqlExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockSqlExecutorMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockSqlExecutor)(nil).ExecContext), varargs...)
}

// QueryContext mocks base method.
func (m *MockSqlExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockSqlExecutorMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockSqlExecutor)(nil).QueryContext), varargs...)
}

// MockExprExecutor is a mock of ExprExecutor interface.
type MockExprExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockExprExecutorMockRecorder
}

// MockExprExecutorMockRecorder is the mock recorder for MockExprExecutor.
type MockExprExecutorMockRecorder struct {
	mock *MockExprExecutor
}

// NewMockExprExecutor creates a new mock instance.
func NewMockExprExecutor(ctrl *gomock.Controller) *MockExprExecutor {
	mock := &MockExprExecutor{ctrl: ctrl}
	mock.recorder = &MockExprExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExprExecutor) EXPECT() *MockExprExecutorMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockExprExecutor) Exec(arg0 builder.SqlExpr) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", arg0)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockExprExecutorMockRecorder) Exec(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockExprExecutor)(nil).Exec), arg0)
}

// ExecContext mocks base method.
func (m *MockExprExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockExprExecutorMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockExprExecutor)(nil).ExecContext), varargs...)
}

// Query mocks base method.
func (m *MockExprExecutor) Query(arg0 builder.SqlExpr) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockExprExecutorMockRecorder) Query(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockExprExecutor)(nil).Query), arg0)
}

// QueryAndScan mocks base method.
func (m *MockExprExecutor) QueryAndScan(arg0 builder.SqlExpr, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAndScan", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueryAndScan indicates an expected call of QueryAndScan.
func (mr *MockExprExecutorMockRecorder) QueryAndScan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAndScan", reflect.TypeOf((*MockExprExecutor)(nil).QueryAndScan), arg0, arg1)
}

// QueryContext mocks base method.
func (m *MockExprExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockExprExecutorMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockExprExecutor)(nil).QueryContext), varargs...)
}

// MockTableResolver is a mock of TableResolver interface.
type MockTableResolver struct {
	ctrl     *gomock.Controller
	recorder *MockTableResolverMockRecorder
}

// MockTableResolverMockRecorder is the mock recorder for MockTableResolver.
type MockTableResolverMockRecorder struct {
	mock *MockTableResolver
}

// NewMockTableResolver creates a new mock instance.
func NewMockTableResolver(ctrl *gomock.Controller) *MockTableResolver {
	mock := &MockTableResolver{ctrl: ctrl}
	mock.recorder = &MockTableResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTableResolver) EXPECT() *MockTableResolverMockRecorder {
	return m.recorder
}

// T mocks base method.
func (m *MockTableResolver) T(model builder.Model) *builder.Table {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "T", model)
	ret0, _ := ret[0].(*builder.Table)
	return ret0
}

// T indicates an expected call of T.
func (mr *MockTableResolverMockRecorder) T(model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "T", reflect.TypeOf((*MockTableResolver)(nil).T), model)
}

// MockTxExecutor is a mock of TxExecutor interface.
type MockTxExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockTxExecutorMockRecorder
}

// MockTxExecutorMockRecorder is the mock recorder for MockTxExecutor.
type MockTxExecutorMockRecorder struct {
	mock *MockTxExecutor
}

// NewMockTxExecutor creates a new mock instance.
func NewMockTxExecutor(ctrl *gomock.Controller) *MockTxExecutor {
	mock := &MockTxExecutor{ctrl: ctrl}
	mock.recorder = &MockTxExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxExecutor) EXPECT() *MockTxExecutorMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockTxExecutor) Begin() (sqlx.DBExecutor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(sqlx.DBExecutor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockTxExecutorMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockTxExecutor)(nil).Begin))
}

// BeginTx mocks base method.
func (m *MockTxExecutor) BeginTx(arg0 *sql.TxOptions) (sqlx.DBExecutor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", arg0)
	ret0, _ := ret[0].(sqlx.DBExecutor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockTxExecutorMockRecorder) BeginTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockTxExecutor)(nil).BeginTx), arg0)
}

// Commit mocks base method.
func (m *MockTxExecutor) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTxExecutorMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTxExecutor)(nil).Commit))
}

// IsTx mocks base method.
func (m *MockTxExecutor) IsTx() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTx")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTx indicates an expected call of IsTx.
func (mr *MockTxExecutorMockRecorder) IsTx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTx", reflect.TypeOf((*MockTxExecutor)(nil).IsTx))
}

// Rollback mocks base method.
func (m *MockTxExecutor) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTxExecutorMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTxExecutor)(nil).Rollback))
}

// MockMigrator is a mock of Migrator interface.
type MockMigrator struct {
	ctrl     *gomock.Controller
	recorder *MockMigratorMockRecorder
}

// MockMigratorMockRecorder is the mock recorder for MockMigrator.
type MockMigratorMockRecorder struct {
	mock *MockMigrator
}

// NewMockMigrator creates a new mock instance.
func NewMockMigrator(ctrl *gomock.Controller) *MockMigrator {
	mock := &MockMigrator{ctrl: ctrl}
	mock.recorder = &MockMigratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMigrator) EXPECT() *MockMigratorMockRecorder {
	return m.recorder
}

// Migrate mocks base method.
func (m *MockMigrator) Migrate(ctx context.Context, db sqlx.DBExecutor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate", ctx, db)
	ret0, _ := ret[0].(error)
	return ret0
}

// Migrate indicates an expected call of Migrate.
func (mr *MockMigratorMockRecorder) Migrate(ctx, db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockMigrator)(nil).Migrate), ctx, db)
}