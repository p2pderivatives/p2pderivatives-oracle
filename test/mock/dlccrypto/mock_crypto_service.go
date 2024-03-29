// Code generated by MockGen. DO NOT EDIT.
// Source: internal/dlccrypto/crypto_service.go

// Package mock_dlccrypto is a generated GoMock package.
package mock_dlccrypto

import (
	dlccrypto "p2pderivatives-oracle/internal/dlccrypto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCryptoService is a mock of CryptoService interface.
type MockCryptoService struct {
	ctrl     *gomock.Controller
	recorder *MockCryptoServiceMockRecorder
}

// MockCryptoServiceMockRecorder is the mock recorder for MockCryptoService.
type MockCryptoServiceMockRecorder struct {
	mock *MockCryptoService
}

// NewMockCryptoService creates a new mock instance.
func NewMockCryptoService(ctrl *gomock.Controller) *MockCryptoService {
	mock := &MockCryptoService{ctrl: ctrl}
	mock.recorder = &MockCryptoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCryptoService) EXPECT() *MockCryptoServiceMockRecorder {
	return m.recorder
}

// ComputeSchnorrSignature mocks base method.
func (m *MockCryptoService) ComputeSchnorrSignature(privateKey *dlccrypto.PrivateKey, message []byte) (*dlccrypto.Signature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ComputeSchnorrSignature", privateKey, message)
	ret0, _ := ret[0].(*dlccrypto.Signature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ComputeSchnorrSignature indicates an expected call of ComputeSchnorrSignature.
func (mr *MockCryptoServiceMockRecorder) ComputeSchnorrSignature(privateKey, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComputeSchnorrSignature", reflect.TypeOf((*MockCryptoService)(nil).ComputeSchnorrSignature), privateKey, message)
}

// ComputeSchnorrSignatureFixedK mocks base method.
func (m *MockCryptoService) ComputeSchnorrSignatureFixedK(privateKey, oneTimeSigningK *dlccrypto.PrivateKey, message string) (*dlccrypto.Signature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ComputeSchnorrSignatureFixedK", privateKey, oneTimeSigningK, message)
	ret0, _ := ret[0].(*dlccrypto.Signature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ComputeSchnorrSignatureFixedK indicates an expected call of ComputeSchnorrSignatureFixedK.
func (mr *MockCryptoServiceMockRecorder) ComputeSchnorrSignatureFixedK(privateKey, oneTimeSigningK, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComputeSchnorrSignatureFixedK", reflect.TypeOf((*MockCryptoService)(nil).ComputeSchnorrSignatureFixedK), privateKey, oneTimeSigningK, message)
}

// GenerateSchnorrKeyPair mocks base method.
func (m *MockCryptoService) GenerateSchnorrKeyPair() (*dlccrypto.PrivateKey, *dlccrypto.SchnorrPublicKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSchnorrKeyPair")
	ret0, _ := ret[0].(*dlccrypto.PrivateKey)
	ret1, _ := ret[1].(*dlccrypto.SchnorrPublicKey)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateSchnorrKeyPair indicates an expected call of GenerateSchnorrKeyPair.
func (mr *MockCryptoServiceMockRecorder) GenerateSchnorrKeyPair() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSchnorrKeyPair", reflect.TypeOf((*MockCryptoService)(nil).GenerateSchnorrKeyPair))
}

// SchnorrPublicKeyFromPrivateKey mocks base method.
func (m *MockCryptoService) SchnorrPublicKeyFromPrivateKey(privateKey *dlccrypto.PrivateKey) (*dlccrypto.SchnorrPublicKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SchnorrPublicKeyFromPrivateKey", privateKey)
	ret0, _ := ret[0].(*dlccrypto.SchnorrPublicKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SchnorrPublicKeyFromPrivateKey indicates an expected call of SchnorrPublicKeyFromPrivateKey.
func (mr *MockCryptoServiceMockRecorder) SchnorrPublicKeyFromPrivateKey(privateKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SchnorrPublicKeyFromPrivateKey", reflect.TypeOf((*MockCryptoService)(nil).SchnorrPublicKeyFromPrivateKey), privateKey)
}

// VerifySchnorrSignature mocks base method.
func (m *MockCryptoService) VerifySchnorrSignature(publicKey *dlccrypto.SchnorrPublicKey, signature *dlccrypto.Signature, message string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySchnorrSignature", publicKey, signature, message)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifySchnorrSignature indicates an expected call of VerifySchnorrSignature.
func (mr *MockCryptoServiceMockRecorder) VerifySchnorrSignature(publicKey, signature, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySchnorrSignature", reflect.TypeOf((*MockCryptoService)(nil).VerifySchnorrSignature), publicKey, signature, message)
}

// VerifySchnorrSignatureRaw mocks base method.
func (m *MockCryptoService) VerifySchnorrSignatureRaw(publicKey *dlccrypto.SchnorrPublicKey, signature *dlccrypto.Signature, message []byte) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySchnorrSignatureRaw", publicKey, signature, message)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifySchnorrSignatureRaw indicates an expected call of VerifySchnorrSignatureRaw.
func (mr *MockCryptoServiceMockRecorder) VerifySchnorrSignatureRaw(publicKey, signature, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySchnorrSignatureRaw", reflect.TypeOf((*MockCryptoService)(nil).VerifySchnorrSignatureRaw), publicKey, signature, message)
}
