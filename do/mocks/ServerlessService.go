// Code generated by MockGen. DO NOT EDIT.
// Source: serverless.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	exec "os/exec"
	reflect "reflect"

	whisk "github.com/apache/openwhisk-client-go/whisk"
	do "github.com/digitalocean/doctl/do"
	gomock "github.com/golang/mock/gomock"
)

// MockServerlessService is a mock of ServerlessService interface.
type MockServerlessService struct {
	ctrl     *gomock.Controller
	recorder *MockServerlessServiceMockRecorder
}

// MockServerlessServiceMockRecorder is the mock recorder for MockServerlessService.
type MockServerlessServiceMockRecorder struct {
	mock *MockServerlessService
}

// NewMockServerlessService creates a new mock instance.
func NewMockServerlessService(ctrl *gomock.Controller) *MockServerlessService {
	mock := &MockServerlessService{ctrl: ctrl}
	mock.recorder = &MockServerlessServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerlessService) EXPECT() *MockServerlessServiceMockRecorder {
	return m.recorder
}

// CheckServerlessStatus mocks base method.
func (m *MockServerlessService) CheckServerlessStatus(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckServerlessStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckServerlessStatus indicates an expected call of CheckServerlessStatus.
func (mr *MockServerlessServiceMockRecorder) CheckServerlessStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckServerlessStatus", reflect.TypeOf((*MockServerlessService)(nil).CheckServerlessStatus), arg0)
}

// Cmd mocks base method.
func (m *MockServerlessService) Cmd(arg0 string, arg1 []string) (*exec.Cmd, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cmd", arg0, arg1)
	ret0, _ := ret[0].(*exec.Cmd)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cmd indicates an expected call of Cmd.
func (mr *MockServerlessServiceMockRecorder) Cmd(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cmd", reflect.TypeOf((*MockServerlessService)(nil).Cmd), arg0, arg1)
}

// CreateNamespace mocks base method.
func (m *MockServerlessService) CreateNamespace(arg0 context.Context, arg1, arg2 string) (do.ServerlessCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNamespace", arg0, arg1, arg2)
	ret0, _ := ret[0].(do.ServerlessCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNamespace indicates an expected call of CreateNamespace.
func (mr *MockServerlessServiceMockRecorder) CreateNamespace(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNamespace", reflect.TypeOf((*MockServerlessService)(nil).CreateNamespace), arg0, arg1, arg2)
}

// DeleteNamespace mocks base method.
func (m *MockServerlessService) DeleteNamespace(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNamespace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNamespace indicates an expected call of DeleteNamespace.
func (mr *MockServerlessServiceMockRecorder) DeleteNamespace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNamespace", reflect.TypeOf((*MockServerlessService)(nil).DeleteNamespace), arg0, arg1)
}

// DeleteTrigger mocks base method.
func (m *MockServerlessService) DeleteTrigger(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrigger", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrigger indicates an expected call of DeleteTrigger.
func (mr *MockServerlessServiceMockRecorder) DeleteTrigger(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrigger", reflect.TypeOf((*MockServerlessService)(nil).DeleteTrigger), arg0, arg1)
}

// Exec mocks base method.
func (m *MockServerlessService) Exec(arg0 *exec.Cmd) (do.ServerlessOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", arg0)
	ret0, _ := ret[0].(do.ServerlessOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockServerlessServiceMockRecorder) Exec(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockServerlessService)(nil).Exec), arg0)
}

// FireTrigger mocks base method.
func (m *MockServerlessService) FireTrigger(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FireTrigger", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FireTrigger indicates an expected call of FireTrigger.
func (mr *MockServerlessServiceMockRecorder) FireTrigger(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FireTrigger", reflect.TypeOf((*MockServerlessService)(nil).FireTrigger), arg0, arg1)
}

// GetConnectedAPIHost mocks base method.
func (m *MockServerlessService) GetConnectedAPIHost() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectedAPIHost")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectedAPIHost indicates an expected call of GetConnectedAPIHost.
func (mr *MockServerlessServiceMockRecorder) GetConnectedAPIHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectedAPIHost", reflect.TypeOf((*MockServerlessService)(nil).GetConnectedAPIHost))
}

// GetFunction mocks base method.
func (m *MockServerlessService) GetFunction(arg0 string, arg1 bool) (whisk.Action, []do.FunctionParameter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFunction", arg0, arg1)
	ret0, _ := ret[0].(whisk.Action)
	ret1, _ := ret[1].([]do.FunctionParameter)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFunction indicates an expected call of GetFunction.
func (mr *MockServerlessServiceMockRecorder) GetFunction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFunction", reflect.TypeOf((*MockServerlessService)(nil).GetFunction), arg0, arg1)
}

// GetHostInfo mocks base method.
func (m *MockServerlessService) GetHostInfo(arg0 string) (do.ServerlessHostInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHostInfo", arg0)
	ret0, _ := ret[0].(do.ServerlessHostInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHostInfo indicates an expected call of GetHostInfo.
func (mr *MockServerlessServiceMockRecorder) GetHostInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHostInfo", reflect.TypeOf((*MockServerlessService)(nil).GetHostInfo), arg0)
}

// GetNamespace mocks base method.
func (m *MockServerlessService) GetNamespace(arg0 context.Context, arg1 string) (do.ServerlessCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespace", arg0, arg1)
	ret0, _ := ret[0].(do.ServerlessCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespace indicates an expected call of GetNamespace.
func (mr *MockServerlessServiceMockRecorder) GetNamespace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespace", reflect.TypeOf((*MockServerlessService)(nil).GetNamespace), arg0, arg1)
}

// GetNamespaceFromCluster mocks base method.
func (m *MockServerlessService) GetNamespaceFromCluster(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespaceFromCluster", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespaceFromCluster indicates an expected call of GetNamespaceFromCluster.
func (mr *MockServerlessServiceMockRecorder) GetNamespaceFromCluster(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespaceFromCluster", reflect.TypeOf((*MockServerlessService)(nil).GetNamespaceFromCluster), arg0, arg1)
}

// GetServerlessNamespace mocks base method.
func (m *MockServerlessService) GetServerlessNamespace(arg0 context.Context) (do.ServerlessCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServerlessNamespace", arg0)
	ret0, _ := ret[0].(do.ServerlessCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServerlessNamespace indicates an expected call of GetServerlessNamespace.
func (mr *MockServerlessServiceMockRecorder) GetServerlessNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServerlessNamespace", reflect.TypeOf((*MockServerlessService)(nil).GetServerlessNamespace), arg0)
}

// GetTrigger mocks base method.
func (m *MockServerlessService) GetTrigger(arg0 context.Context, arg1 string) (do.ServerlessTrigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrigger", arg0, arg1)
	ret0, _ := ret[0].(do.ServerlessTrigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrigger indicates an expected call of GetTrigger.
func (mr *MockServerlessServiceMockRecorder) GetTrigger(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrigger", reflect.TypeOf((*MockServerlessService)(nil).GetTrigger), arg0, arg1)
}

// InstallServerless mocks base method.
func (m *MockServerlessService) InstallServerless(arg0 string, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallServerless", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InstallServerless indicates an expected call of InstallServerless.
func (mr *MockServerlessServiceMockRecorder) InstallServerless(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallServerless", reflect.TypeOf((*MockServerlessService)(nil).InstallServerless), arg0, arg1)
}

// InvokeFunction mocks base method.
func (m *MockServerlessService) InvokeFunction(arg0 string, arg1 interface{}, arg2, arg3 bool) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeFunction", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InvokeFunction indicates an expected call of InvokeFunction.
func (mr *MockServerlessServiceMockRecorder) InvokeFunction(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeFunction", reflect.TypeOf((*MockServerlessService)(nil).InvokeFunction), arg0, arg1, arg2, arg3)
}

// InvokeFunctionViaWeb mocks base method.
func (m *MockServerlessService) InvokeFunctionViaWeb(arg0 string, arg1 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeFunctionViaWeb", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvokeFunctionViaWeb indicates an expected call of InvokeFunctionViaWeb.
func (mr *MockServerlessServiceMockRecorder) InvokeFunctionViaWeb(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeFunctionViaWeb", reflect.TypeOf((*MockServerlessService)(nil).InvokeFunctionViaWeb), arg0, arg1)
}

// ListFunctions mocks base method.
func (m *MockServerlessService) ListFunctions(arg0 string, arg1, arg2 int) ([]whisk.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFunctions", arg0, arg1, arg2)
	ret0, _ := ret[0].([]whisk.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFunctions indicates an expected call of ListFunctions.
func (mr *MockServerlessServiceMockRecorder) ListFunctions(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFunctions", reflect.TypeOf((*MockServerlessService)(nil).ListFunctions), arg0, arg1, arg2)
}

// ListNamespaces mocks base method.
func (m *MockServerlessService) ListNamespaces(arg0 context.Context) (do.NamespaceListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNamespaces", arg0)
	ret0, _ := ret[0].(do.NamespaceListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNamespaces indicates an expected call of ListNamespaces.
func (mr *MockServerlessServiceMockRecorder) ListNamespaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNamespaces", reflect.TypeOf((*MockServerlessService)(nil).ListNamespaces), arg0)
}

// ListTriggers mocks base method.
func (m *MockServerlessService) ListTriggers(arg0 context.Context, arg1 string) ([]do.ServerlessTrigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTriggers", arg0, arg1)
	ret0, _ := ret[0].([]do.ServerlessTrigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTriggers indicates an expected call of ListTriggers.
func (mr *MockServerlessServiceMockRecorder) ListTriggers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTriggers", reflect.TypeOf((*MockServerlessService)(nil).ListTriggers), arg0, arg1)
}

// ReadCredentials mocks base method.
func (m *MockServerlessService) ReadCredentials() (do.ServerlessCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCredentials")
	ret0, _ := ret[0].(do.ServerlessCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCredentials indicates an expected call of ReadCredentials.
func (mr *MockServerlessServiceMockRecorder) ReadCredentials() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCredentials", reflect.TypeOf((*MockServerlessService)(nil).ReadCredentials))
}

// ReadProject mocks base method.
func (m *MockServerlessService) ReadProject(arg0 *do.ServerlessProject, arg1 []string) (do.ServerlessOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadProject", arg0, arg1)
	ret0, _ := ret[0].(do.ServerlessOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadProject indicates an expected call of ReadProject.
func (mr *MockServerlessServiceMockRecorder) ReadProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadProject", reflect.TypeOf((*MockServerlessService)(nil).ReadProject), arg0, arg1)
}

// SetTriggerEnablement mocks base method.
func (m *MockServerlessService) SetTriggerEnablement(arg0 context.Context, arg1 string, arg2 bool) (do.ServerlessTrigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTriggerEnablement", arg0, arg1, arg2)
	ret0, _ := ret[0].(do.ServerlessTrigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetTriggerEnablement indicates an expected call of SetTriggerEnablement.
func (mr *MockServerlessServiceMockRecorder) SetTriggerEnablement(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTriggerEnablement", reflect.TypeOf((*MockServerlessService)(nil).SetTriggerEnablement), arg0, arg1, arg2)
}

// Stream mocks base method.
func (m *MockServerlessService) Stream(arg0 *exec.Cmd) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stream", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stream indicates an expected call of Stream.
func (mr *MockServerlessServiceMockRecorder) Stream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stream", reflect.TypeOf((*MockServerlessService)(nil).Stream), arg0)
}

// WriteCredentials mocks base method.
func (m *MockServerlessService) WriteCredentials(arg0 do.ServerlessCredentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteCredentials", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteCredentials indicates an expected call of WriteCredentials.
func (mr *MockServerlessServiceMockRecorder) WriteCredentials(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteCredentials", reflect.TypeOf((*MockServerlessService)(nil).WriteCredentials), arg0)
}

// WriteProject mocks base method.
func (m *MockServerlessService) WriteProject(arg0 do.ServerlessProject) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteProject", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteProject indicates an expected call of WriteProject.
func (mr *MockServerlessServiceMockRecorder) WriteProject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteProject", reflect.TypeOf((*MockServerlessService)(nil).WriteProject), arg0)
}
