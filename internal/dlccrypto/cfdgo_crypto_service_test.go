package dlccrypto_test

import (
	"p2pderivatives-oracle/internal/dlccrypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

type RKPair struct {
	rvalue string
	k      string
}

type SignaturePair struct {
	signature string
	message   string
	krPair    RKPair
}

var (
	TestOracleKeyPair = &KeyPair{PrivateKey: "18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725", PublicKey: "0250863ad64a87ae8a2fe83c1af1a8403cb53f53e486d8511dad8a04887e5b2352"}
	TestKeyPairs      = [...]KeyPair{
		{PrivateKey: "71957542da5ea5e6b58607a3e21f0eb4e871c785b5b3470fb5274213f050cd60", PublicKey: "02d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a"},
		{PrivateKey: "a11ad4c96610096f956769cbd060e745750932c03fa30d1ea3b080428e12a336", PublicKey: "0226721aa761d4483df8105fb8ed2e1c6aa7cf309f5e1f0b48dcb615407697e193"},
		{PrivateKey: "f912b67c2796be8fa70d55c5f240be799efa9908484af9bd43344bd95a62614b", PublicKey: "03bb3186cfae1a48d74a7ed0d6428226fbfac45bd5c6fc0d3f377f8c474e804735"},
		{PrivateKey: "1ce5743c1d546ea501a421cd979581d62fd81fbd27e79b652ad2a273d45e92b5", PublicKey: "03b6e8ff7065ff02fb415542713b55b9c596cd7db908c2ecf782fcf29025b26943"},
		{PrivateKey: "574ba4e483bce2581b11552fb301b0db7d549aa54d0a139f49ce14a3bb9fa107", PublicKey: "038b80e907786c8dc18921ac3ef7fc82e0e23a115da59dc630150bf49b5c970e7b"},
		{PrivateKey: "224a8c4ce0fa3903b1d483625521cf1f46993bb8afcbf81a5b84ee94ce7d65de", PublicKey: "02325cfa2600ea6487e2467ce2c1efbc497e04b03f1fa02b8e9bd90146e8049381"},
		{PrivateKey: "50617a8fc7d133e8c9ab31ff5eab79312c0f0a7b280d70c0e562eee3f05cc9c8", PublicKey: "03364591c8a6e01a79bc635af1fe331844c301e460e74b77cdab44227a57105a60"},
		{PrivateKey: "ab0c1bd576e85d0c7817ef1f7c25684769ac10e71fadf618ee97f003aa1b14ff", PublicKey: "02d83ca02cdce6d8e934fdcd8a3bc5f874182de950db732184ccde01d77994a58b"},
	}
	TestKRValues = [...]RKPair{
		{k: "71957542da5ea5e6b58607a3e21f0eb4e871c785b5b3470fb5274213f050cd60", rvalue: "02d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a"},
		{k: "a11ad4c96610096f956769cbd060e745750932c03fa30d1ea3b080428e12a336", rvalue: "0226721aa761d4483df8105fb8ed2e1c6aa7cf309f5e1f0b48dcb615407697e193"},
		{k: "f912b67c2796be8fa70d55c5f240be799efa9908484af9bd43344bd95a62614b", rvalue: "02bb3186cfae1a48d74a7ed0d6428226fbfac45bd5c6fc0d3f377f8c474e804735"},
		{k: "1ce5743c1d546ea501a421cd979581d62fd81fbd27e79b652ad2a273d45e92b5", rvalue: "02b6e8ff7065ff02fb415542713b55b9c596cd7db908c2ecf782fcf29025b26943"},
		{k: "574ba4e483bce2581b11552fb301b0db7d549aa54d0a139f49ce14a3bb9fa107", rvalue: "028b80e907786c8dc18921ac3ef7fc82e0e23a115da59dc630150bf49b5c970e7b"},
		{k: "224a8c4ce0fa3903b1d483625521cf1f46993bb8afcbf81a5b84ee94ce7d65de", rvalue: "02325cfa2600ea6487e2467ce2c1efbc497e04b03f1fa02b8e9bd90146e8049381"},
		{k: "50617a8fc7d133e8c9ab31ff5eab79312c0f0a7b280d70c0e562eee3f05cc9c8", rvalue: "02364591c8a6e01a79bc635af1fe331844c301e460e74b77cdab44227a57105a60"},
		{k: "ab0c1bd576e85d0c7817ef1f7c25684769ac10e71fadf618ee97f003aa1b14ff", rvalue: "02d83ca02cdce6d8e934fdcd8a3bc5f874182de950db732184ccde01d77994a58b"},
	}
	TestMessage = [...]string{
		"1200",
		"3000",
		"0",
		"30",
		"5475673",
		"76428764",
	}

	TestSignature = [...]SignaturePair{
		{krPair: TestKRValues[0], message: TestMessage[0], signature: "26772fddf151463655823c906d1c01afb682b8c3533589fe71753a579bf59b22"},
		{krPair: TestKRValues[1], message: TestMessage[0], signature: "61724003ecdf295d220c481f87a7905399af05671a0c4acd2979efaad778cf5a"},
		{krPair: TestKRValues[2], message: TestMessage[1], signature: "60433d967d606c7ebd204d0f7d269bc435ea75d0fe3bb8b88bb102cde7972f00"},
		{krPair: TestKRValues[3], message: TestMessage[1], signature: "0ad1c72b5abbdcb311bd878b960847e53808cfb8cb2bf3ae03eff74d3058028c"},
		{krPair: TestKRValues[4], message: TestMessage[2], signature: "5f7e28241fbdc768f8a541f523c88887c133751acab2ef7a580e636abff0f85b"},
	}
)

func NewTestCfdgoCryptoService() dlccrypto.CryptoService {
	return dlccrypto.NewCfdgoCryptoService()
}

func Test_CfdgoCryptoService_PublicKeyFromPrivateKey(t *testing.T) {
	crypto := NewTestCfdgoCryptoService()
	for _, keypair := range TestKeyPairs {
		privKey, err := dlccrypto.NewPrivateKey(keypair.PrivateKey)
		assert.NoError(t, err)
		pubkey, err := crypto.PublicKeyFromPrivateKey(privKey)
		assert.NoError(t, err)
		assert.Equal(t, keypair.PublicKey, pubkey.EncodeToString())
	}
}

func Test_CfdgoCryptoService_GenerateKvalue(t *testing.T) {
	crypto := dlccrypto.NewCfdgoCryptoService()
	k1, err1 := crypto.GenerateKvalue()
	assert.NoError(t, err1)
	k2, err2 := crypto.GenerateKvalue()
	assert.NoError(t, err2)
	assert.NotEqual(t, k1.EncodeToString(), k2.EncodeToString())
}

func Test_CfdgoCryptoService_ComputeRvalue(t *testing.T) {
	crypto := dlccrypto.NewCfdgoCryptoService()
	for _, rk := range TestKRValues {
		privKey, err := dlccrypto.NewPrivateKey(rk.k)
		assert.NoError(t, err)
		rvalue, err := crypto.ComputeRvalue(privKey)
		assert.NoError(t, err)
		assert.Equal(t, rk.rvalue, rvalue.EncodeToString())
	}
}

func Test_CfdgoCryptoService_ComputeSchnorrSignature(t *testing.T) {
	crypto := dlccrypto.NewCfdgoCryptoService()
	oracleKey, err := dlccrypto.NewPrivateKey(TestOracleKeyPair.PrivateKey)
	assert.NoError(t, err)
	for _, sigpair := range TestSignature {
		kvalue, err := dlccrypto.NewPrivateKey(sigpair.krPair.k)
		assert.NoError(t, err)
		sig, err := crypto.ComputeSchnorrSignature(oracleKey, kvalue, sigpair.message)
		assert.NoError(t, err)
		assert.Equal(t, sigpair.signature, sig.EncodeToString())
	}
}

func Test_CfdgoCryptoService_VerifySignature(t *testing.T) {
	crypto := dlccrypto.NewCfdgoCryptoService()
	oraclePub, err := dlccrypto.NewPublicKey(TestOracleKeyPair.PublicKey)
	assert.NoError(t, err)
	for _, sigpair := range TestSignature {
		rvalue, err := dlccrypto.NewPublicKey(sigpair.krPair.rvalue)
		assert.NoError(t, err)
		sig, err := dlccrypto.NewSignature(sigpair.signature)
		assert.NoError(t, err)
		check, err := crypto.VerifySchnorrSignature(oraclePub, rvalue, sig, sigpair.message)
		assert.NoError(t, err)
		assert.True(t, check)
	}
}
