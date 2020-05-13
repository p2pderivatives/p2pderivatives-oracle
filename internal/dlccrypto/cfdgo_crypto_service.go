package dlccrypto

import (
	"crypto/sha256"
	"encoding/hex"

	cfdgo "github.com/cryptogarageinc/cfd-go"
	"github.com/pkg/errors"
)

// NewCfdgoCryptoService returns a CryptoService implemented using cfd-go library
func NewCfdgoCryptoService() CryptoService {
	return &CfdgoCryptoService{}
}

// CfdgoCryptoService crypto service implementing the schnorr api of cfd-go library
type CfdgoCryptoService struct{}

// ComputeRvalue computes a public nonce from a kvalue private key
func (o *CfdgoCryptoService) ComputeRvalue(nonce *PrivateKey) (*PublicKey, error) {
	rvalue, err := cfdgo.CfdGoGetSchnorrPublicNonce(nonce.EncodeToString())
	if err != nil {
		return nil, errors.WithMessage(err, "Error while computing rvalue")
	}
	pubkey, err := NewPublicKey(rvalue)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

// GenerateKvalue generate a new unique K value to be used with schnorr signature
func (o *CfdgoCryptoService) GenerateKvalue() (*PrivateKey, error) {
	_, key, _, err := cfdgo.CfdGoCreateKeyPair(true, 0)
	if err != nil {
		return nil, errors.WithMessage(err, "Error while generating one time signing K")
	}
	kvalue, err := NewPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return kvalue, nil
}

// PublicKeyFromPrivateKey computes public key from private key
func (o *CfdgoCryptoService) PublicKeyFromPrivateKey(privateKey *PrivateKey) (*PublicKey, error) {
	bs, err := cfdgo.CfdGoGetPubkeyFromPrivkey(privateKey.EncodeToString(), "", true)
	if err != nil {
		return nil, errors.WithMessage(err, "Error while calculating public key from private key")
	}
	pubkey, err := NewPublicKey(bs)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

// ComputeSchnorrSignature computes a schnorr signature on the given message (will be hashed by sha256)
func (o *CfdgoCryptoService) ComputeSchnorrSignature(privateKey *PrivateKey, kvalue *PrivateKey, message string) (*Signature, error) {
	hash32 := sha256.Sum256([]byte(message))
	hash := hex.EncodeToString(hash32[:])
	bs, err := cfdgo.CfdGoCalculateSchnorrSignatureWithNonce(privateKey.EncodeToString(), kvalue.EncodeToString(), hash)
	if err != nil {
		return nil, errors.WithMessage(err, "Error while computing schnorr signature")
	}
	sig, err := NewSignature(bs)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// VerifySchnorrSignature verifies the schnorr signature against public key and rvalue on the given message (will be hashed by sha256)
func (o *CfdgoCryptoService) VerifySchnorrSignature(publicKey *PublicKey, rvalue *PublicKey, signature *Signature, message string) (bool, error) {
	hash32 := sha256.Sum256([]byte(message))
	hash := hex.EncodeToString(hash32[:])
	ok, err := cfdgo.CfdGoVerifySchnorrSignatureWithNonce(publicKey.EncodeToString(), rvalue.EncodeToString(), signature.EncodeToString(), hash)
	if err != nil {
		return false, errors.WithMessage(err, "Error while verifying schnorr signature")
	}
	return ok, nil
}
