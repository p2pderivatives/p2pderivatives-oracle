package dlccrypto_test

import (
	"p2pderivatives-oracle/internal/dlccrypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validBytestring = "abcdef0100"
	validPrivateKey = "71957542da5ea5e6b58607a3e21f0eb4e871c785b5b3470fb5274213f050cd60"
	validPublicKey  = "d557c15ea53c46245be38a062e33c22c5880f6b776e36695befc6e28a3934bef"
	validSignature  = "26ba8bc2e81388a1c75b3fa6de8a90ee3c45d35793c4e0327496c3d70b49f5a182931af82dd50ea6807cc4395715edc4c2af39c6c695c5a027d1fada0570bd6c"

	invalidBytestring = "abcdefghrt65749"
	invalidPrivateKey = "957542da5ea5e6b58607a3e21f0eb4e871c785b5b3470fb5274213f050cd60"
	invalidPublicKey  = "03d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a"
	invalidSignature  = "26ba8bc2e81388a1c75b3fa6de8a90ee3c45d35793c4e0327496c3d70b49f5a182931af82dd50ea6807cc4395715edc4c2af39c6c695c5a027d1fada0570bd6ca"
)

func TestByteString_EncodeToString_ReturnsCorrectValue(t *testing.T) {
	bs, err := dlccrypto.NewByteString(validBytestring)
	assert.NoError(t, err)
	assert.Equal(t, validBytestring, bs.EncodeToString())
}

func TestNewByteString_WithValidBytestring_ReturnsNoError(t *testing.T) {
	_, err := dlccrypto.NewByteString(validBytestring)
	assert.NoError(t, err)
}

func TestNewByteString_WithInvalidBytestring_ReturnsError(t *testing.T) {
	bs, err := dlccrypto.NewByteString(invalidBytestring)
	assert.Error(t, err)
	assert.Nil(t, bs)
}

func TestNewPrivateKey_WithInvalidSizeBytestring_ReturnsError(t *testing.T) {
	bs, err := dlccrypto.NewPrivateKey(invalidPrivateKey)
	assert.Error(t, err)
	assert.Nil(t, bs)
}

func TestNewPrivateKey_WithValidSizeBytestring_ReturnsNoError(t *testing.T) {
	_, err := dlccrypto.NewPrivateKey(validPrivateKey)
	assert.NoError(t, err)
}

func TestNewPublicKey_WithInvalidSizeBytestring_ReturnsError(t *testing.T) {
	bs, err := dlccrypto.NewSchnorrPublicKey(invalidPublicKey)
	assert.Error(t, err)
	assert.Nil(t, bs)
}

func TestNewPublicKey_WithValidSizeBytestring_ReturnsNoError(t *testing.T) {
	_, err := dlccrypto.NewSchnorrPublicKey(validPublicKey)
	assert.NoError(t, err)
}

func TestNewSignature_WithInvalidSizeBytestring_ReturnsError(t *testing.T) {
	_, err := dlccrypto.NewSignature(invalidSignature)
	assert.Error(t, err)
}

func TestNewSignature_WithValidSizeBytestring_ReturnsNoError(t *testing.T) {
	_, err := dlccrypto.NewSignature(validSignature)
	assert.NoError(t, err)
}
