package dlccrypto

// CryptoService interface for an utility crypto service
type CryptoService interface {
	GenerateKvalue() (*PrivateKey, error)
	ComputeRvalue(nonce *PrivateKey) (*PublicKey, error)
	PublicKeyFromPrivateKey(privateKey *PrivateKey) (*PublicKey, error)
	ComputeSchnorrSignature(privateKey *PrivateKey, oneTimeSigningK *PrivateKey, message string) (*Signature, error)
	VerifySchnorrSignature(publicKey *PublicKey, rvalue *PublicKey, signature *Signature, message string) (bool, error)
}
