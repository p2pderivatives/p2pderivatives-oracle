package dlccrypto_test

import (
	"p2pderivatives-oracle/internal/dlccrypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validBytestring = "abcdef0100"
	validPrivateKey = "71957542da5ea5e6b58607a3e21f0eb4e871c785b5b3470fb5274213f050cd60"
	validPublicKey  = "02d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a"
	validSignature  = "60433d967d606c7ebd204d0f7d269bc435ea75d0fe3bb8b88bb102cde7972f00"

	invalidBytestring = "abcdefghrt65749"
	invalidPrivateKey = "957542da5ea5e6b58607a3e21f0eb4e871c785b5b3470fb5274213f050cd60"
	invalidPublicKey  = "d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a"
	invalidSignature  = "433d967d606c7ebd204d0f7d269bc435ea75d0fe3bb8b88bb102cde7972f00"
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
	bs, err := dlccrypto.NewPublicKey(invalidPublicKey)
	assert.Error(t, err)
	assert.Nil(t, bs)
}

func TestNewPublicKey_WithValidSizeBytestring_ReturnsNoError(t *testing.T) {
	_, err := dlccrypto.NewPublicKey(validPublicKey)
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
