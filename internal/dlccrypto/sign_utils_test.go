package dlccrypto_test

import (
	"encoding/hex"
	"math/rand"
	"p2pderivatives-oracle/internal/cfddlccrypto"
	"p2pderivatives-oracle/internal/dlccrypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validEventSignature = "319dfb9ced3c34242aad5920e1f2862346accfa19f013726beb1d7d1678737805eccecedea7abbe7296ff94a394043a219d3087d2d47fc3cff95d0b9f5595d92"
	validSerialization  = "0002abf8f63630a0b1dec98ce8db50e9680f89f3390105454510420048d050aaa05df4a731b0d25a291f7bbc33f391003e87dcfae98a7484e37646453725405f7f3160bf0bb0fdd80a1200020008736174732f73656300000000000a0454657374"
)

func TestEventSerialization_ReturnsExpectedByteArray(t *testing.T) {
	nonce0, _ := dlccrypto.NewSchnorrPublicKey("abf8f63630a0b1dec98ce8db50e9680f89f3390105454510420048d050aaa05d")
	nonce1, _ := dlccrypto.NewSchnorrPublicKey("f4a731b0d25a291f7bbc33f391003e87dcfae98a7484e37646453725405f7f31")
	nonces := []dlccrypto.SchnorrPublicKey{*nonce0, *nonce1}

	ser := dlccrypto.SerializeEvent(nonces, 1623133104, 2, false, "sats/sec", 0, 10, "Test")

	assert.Equal(t, validSerialization, hex.EncodeToString(ser))
}

func TestGenerateEventSignature_ReturnsExpectedSignature(t *testing.T) {
	rand.Seed(1)
	cryptoService := cfddlccrypto.NewCfdgoCryptoService()
	privKey, _ := dlccrypto.NewPrivateKey("c251ebf21fcf41e4875ddfc0a02e5ae849e847b3f528ae0413363f47d2c02e66")
	nonce0, _ := dlccrypto.NewSchnorrPublicKey("abf8f63630a0b1dec98ce8db50e9680f89f3390105454510420048d050aaa05d")
	nonce1, _ := dlccrypto.NewSchnorrPublicKey("f4a731b0d25a291f7bbc33f391003e87dcfae98a7484e37646453725405f7f31")
	nonces := []dlccrypto.SchnorrPublicKey{*nonce0, *nonce1}

	bs, err := dlccrypto.GenerateEventSignature(privKey, nonces, 1623133104, 2, false, "sats/sec", 0, 10, "Test", cryptoService)
	assert.NoError(t, err)
	assert.Equal(t, validEventSignature, bs)
}
