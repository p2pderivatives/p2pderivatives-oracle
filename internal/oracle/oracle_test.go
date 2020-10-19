package oracle_test

import (
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/internal/oracle"
	"p2pderivatives-oracle/test"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ExpectedKeyPair = struct {
	privateKey string
	publicKey  string
	keyPath    string
	password   string
	passPath   string
}{
	privateKey: "c85c333c73eb6daf3479d0236b261d7512cb5daf6955beea5d66a180d34260ae",
	publicKey:  "d557c15ea53c46245be38a062e33c22c5880f6b776e36695befc6e28a3934bef",
	password:   "XqH/LtBBrBJ/iSSxx5gejDAwA3nbiIGBV0w/SqGnPXc=",
	keyPath:    filepath.Join(test.VectorsDirectoryPath, "oracle/key.pem"),
	passPath:   filepath.Join(test.VectorsDirectoryPath, "oracle/pass.txt"),
}

func Test_New_WithInValidKey_ReturnsError(t *testing.T) {
	invalidPrivKey := &dlccrypto.PrivateKey{}
	_, err := oracle.New(invalidPrivKey)
	assert.NotNil(t, err)
}

func Test_FromConfig_WithPass_ReturnsOracle(t *testing.T) {
	config := &oracle.Config{
		KeyFile: ExpectedKeyPair.keyPath,
		KeyPass: ExpectedKeyPair.password,
	}

	oracleInstance, err := oracle.FromConfig(config)
	assert.NoError(t, err)
	if assert.NotNil(t, oracleInstance) {
		assert.Equal(t, ExpectedKeyPair.privateKey, oracleInstance.PrivateKey.EncodeToString())
		assert.Equal(t, ExpectedKeyPair.publicKey, oracleInstance.PublicKey.EncodeToString())
	}
}

func Test_FromConfig_WithPassFile_ReturnsOracle(t *testing.T) {
	config := &oracle.Config{
		KeyFile:     ExpectedKeyPair.keyPath,
		KeyPassFile: ExpectedKeyPair.passPath,
	}
	oracleInstance, err := oracle.FromConfig(config)
	assert.NoError(t, err)
	if assert.NotNil(t, oracleInstance) {
		assert.Equal(t, ExpectedKeyPair.privateKey, oracleInstance.PrivateKey.EncodeToString())
		assert.Equal(t, ExpectedKeyPair.publicKey, oracleInstance.PublicKey.EncodeToString())
	}
}

func Test_FromConfig_WithNoPass_ReturnsError(t *testing.T) {
	config := &oracle.Config{
		KeyFile: ExpectedKeyPair.keyPath,
	}
	_, err := oracle.FromConfig(config)
	assert.NotNil(t, err)
}
