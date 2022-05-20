package dlccrypto

// CryptoService interface for an utility crypto service
type CryptoService interface {
	GenerateSchnorrKeyPair() (*PrivateKey, *SchnorrPublicKey, error)
	SchnorrPublicKeyFromPrivateKey(privateKey *PrivateKey) (*SchnorrPublicKey, error)
	ComputeSchnorrSignatureFixedK(privateKey *PrivateKey, oneTimeSigningK *PrivateKey, message string) (*Signature, error)
	ComputeSchnorrSignature(privateKey *PrivateKey, message []byte) (*Signature, error)
	VerifySchnorrSignature(publicKey *SchnorrPublicKey, signature *Signature, message string) (bool, error)
}
