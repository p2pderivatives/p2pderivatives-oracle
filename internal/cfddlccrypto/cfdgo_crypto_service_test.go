package cfddlccrypto_test

import (
	"math/rand"
	"p2pderivatives-oracle/internal/cfddlccrypto"
	"p2pderivatives-oracle/internal/dlccrypto"
	"testing"
	"time"

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
	TestOracleKeyPair = &KeyPair{PrivateKey: "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725", PublicKey: "50863ad64a87ae8a2fe83c1af1a8403cb53f53e486d8511dad8a04887e5b2352"}
	TestKeyPairs      = [...]KeyPair{
		{PrivateKey: "dcfbf4fc4f357ac42e038c00dc8a9d4f51f04b5ef7b31ce413b6daad5b9efb69", PublicKey: "2b4ec6a8dff179be54a40e68ba48c2f81a239aaf1dae5b625a371c52dca7e649"},
		{PrivateKey: "b3b2b54604efa25ad90f04c51a70a415fdb246253cb74f7cce3918ed54f24df1", PublicKey: "fcc2004734a187853e98b62a7715bcf63c44bd8bc3ade6e9c518fe34f1602a22"},
		{PrivateKey: "4cb24cab7de4b95e1d37164dd1583a720f366b1801549a0e90aceea988932ce6", PublicKey: "c900116a6219b8f24dde0e7828e73b2258cd7fc9475192d013f30cd69c1666c6"},
		{PrivateKey: "9607eaacb30a1b2e84ea3b15130b7448be2cfcaf882be229da254d3dcd46b7f0", PublicKey: "50f88b6346375cc847526c57962368cc9ba2f0d26404a0f38c09b122b31488ef"},
		{PrivateKey: "cf7d4800d89f4673d87a97f4a16b02f1890e673ab57c5c3c31c1b01ff0df13d0", PublicKey: "c38ddcd58cede91f22488be1b5ab07b4ca3c84bf8dc94c0712d911aa5030bd20"},
		{PrivateKey: "abf4c668af6c9fa8c15edc9a996e35d8fbb3feed69259f0e875946e8e37621c9", PublicKey: "9599d5293b3ef1eed597234a0af94190e11e692d02ff91cc3648dd9debfe8637"},
	}
	TestKRValues = [...]RKPair{
		{k: "d8667a07d8a66cbdeda3a8da8c8ce802bf22493abea287df37f92ee0d7725fb0", rvalue: "5a00f102a9a2c789046da82a900b4b1b34fcf73dce5ac1063a653c2bf9b3f5f0"},
		{k: "2644083242f5cf7ff89331f219cb064ee81f6279face75d96ea5a22b2180fa72", rvalue: "d34fba30e1d6f8e82e37ae34ded1f16aac0e1527257201bbb21025ea49c9bdea"},
		{k: "7db3f7091798d2d426205bdb194f74401014755e3e58f9303390c1ff0e4bd44a", rvalue: "1fc82267e136cb89bd2ac0c2b0c7e3ef202d895193f071556bd91769bd45c752"},
		{k: "82b12720e1ee86ef2c2fb29726026357438b509e00cada4bd368c110e30f0540", rvalue: "8edc285addfcf2b3d3868a419cca4aeb89dff72e966be63dd57dac59a9c7c6df"},
		{k: "b80916acc3c31ddb05b72526fff1b7d813a55dd84616635fe2907e1a36755f75", rvalue: "6b02fb3fdaf0a9e6d05e4d1973a56ce7a4f2e44e02a6c4a0574f372d7771a9f3"},
		{k: "f8f6f58f646202749081862693763b48b740964597c2ea641b267aeb9fbecf4d", rvalue: "1ef3179d5bb5ded119de7c63d3e2abf6aea94b5fd4458ad9b46dd8f59df2ccf5"},
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
		{krPair: TestKRValues[0], message: TestMessage[0], signature: "5a00f102a9a2c789046da82a900b4b1b34fcf73dce5ac1063a653c2bf9b3f5f0c50e5b475b7fefc5deac176a91fde56e1fa8c522661af2e6a6ea60b3f3e4d82e"},
		{krPair: TestKRValues[1], message: TestMessage[1], signature: "d34fba30e1d6f8e82e37ae34ded1f16aac0e1527257201bbb21025ea49c9bdea536c4b9fceeed308bf06d6b4e006555adc481b1e93995a599b1b98f1ac66b195"},
		{krPair: TestKRValues[2], message: TestMessage[2], signature: "1fc82267e136cb89bd2ac0c2b0c7e3ef202d895193f071556bd91769bd45c752d69f9d53fdb1a7fb280d718a8bf083cbb5a796c0021fdde34e3280b63d109e6a"},
		{krPair: TestKRValues[3], message: TestMessage[3], signature: "8edc285addfcf2b3d3868a419cca4aeb89dff72e966be63dd57dac59a9c7c6dfa56233133551e0f1bb7693490eba8e2b41caa0d33ace97bc7598d26b659980f6"},
		{krPair: TestKRValues[4], message: TestMessage[4], signature: "6b02fb3fdaf0a9e6d05e4d1973a56ce7a4f2e44e02a6c4a0574f372d7771a9f365af59c67f1460b0090bf2357c6c387fc67f45f7b7e07cb8e65172d7898dee4a"},
		{krPair: TestKRValues[5], message: TestMessage[5], signature: "1ef3179d5bb5ded119de7c63d3e2abf6aea94b5fd4458ad9b46dd8f59df2ccf548c13c2c6e9d2aab66ea8f393d29bc9f658ff56c629528ab6e80370c92b5afec"},
	}
)

func NewTestCfdgoCryptoService() dlccrypto.CryptoService {
	return cfddlccrypto.NewCfdgoCryptoService()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Test_CfdgoCryptoService_PublicKeyFromPrivateKey(t *testing.T) {
	crypto := NewTestCfdgoCryptoService()
	for _, keypair := range TestKeyPairs {
		privKey, err := dlccrypto.NewPrivateKey(keypair.PrivateKey)
		assert.NoError(t, err)
		pubkey, err := crypto.SchnorrPublicKeyFromPrivateKey(privKey)
		assert.NoError(t, err)
		assert.Equal(t, keypair.PublicKey, pubkey.EncodeToString())
	}
}

func Test_SignAndVerify(t *testing.T) {
	crypto := NewTestCfdgoCryptoService()
	privkey, pubkey, _ := crypto.GenerateSchnorrKeyPair()
	seed := time.Now().UnixNano()
	t.Log("Seed: ", seed)
	rand.Seed(seed)

	for i := 0; i < 100; i++ {
		msg := RandString(10)
		kvalue, _, _ := crypto.GenerateSchnorrKeyPair()
		sig, err := crypto.ComputeSchnorrSignatureFixedK(privkey, kvalue, msg)
		assert.NoError(t, err)
		valid, err := crypto.VerifySchnorrSignature(pubkey, sig, msg)
		assert.NoError(t, err)
		assert.True(t, valid)
	}
}

func Test_CfdgoCryptoService_ComputeSchnorrSignature(t *testing.T) {
	crypto := cfddlccrypto.NewCfdgoCryptoService()
	oracleKey, err := dlccrypto.NewPrivateKey(TestOracleKeyPair.PrivateKey)
	assert.NoError(t, err)
	for _, sigpair := range TestSignature {
		kvalue, err := dlccrypto.NewPrivateKey(sigpair.krPair.k)
		assert.NoError(t, err)
		sig, err := crypto.ComputeSchnorrSignatureFixedK(oracleKey, kvalue, sigpair.message)
		assert.NoError(t, err)
		assert.Equal(t, sigpair.signature, sig.EncodeToString())
	}
}

func Test_CfdgoCryptoService_VerifySignature(t *testing.T) {
	crypto := cfddlccrypto.NewCfdgoCryptoService()
	oraclePub, err := dlccrypto.NewSchnorrPublicKey(TestOracleKeyPair.PublicKey)
	assert.NoError(t, err)
	for _, sigpair := range TestSignature {
		assert.NoError(t, err)
		sig, err := dlccrypto.NewSignature(sigpair.signature)
		assert.NoError(t, err)
		check, err := crypto.VerifySchnorrSignature(oraclePub, sig, sigpair.message)
		assert.NoError(t, err)
		assert.True(t, check)
	}
}
