package dlccrypto

import (
	"crypto/sha256"

	cfdgo "github.com/cryptogarageinc/cfd-go"
	"github.com/pkg/errors"
)

// NewCfdgoCryptoService returns a CryptoService implemented using cfd-go library
func NewCfdgoCryptoService() CryptoService {
	schnorrUtil := cfdgo.NewSchnorrUtil()
	return &CfdgoCryptoService{schnorrUtil}
}

// CfdgoCryptoService crypto service implementing the schnorr api of cfd-go library
type CfdgoCryptoService struct {
	schnorrUtil *cfdgo.SchnorrUtil
}

// GenerateSchnorrKeyPair returns a freshly generated Schnorr public/private key pair
func (o *CfdgoCryptoService) GenerateSchnorrKeyPair() (*PrivateKey, *SchnorrPublicKey, error) {
	_, seckey, _, err := cfdgo.CfdGoCreateKeyPair(true, 0)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Error while generating cfd go key pair")
	}

	privkey, err := NewPrivateKey(seckey)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Error while generating private key")
	}

	pubkey, err := o.SchnorrPublicKeyFromPrivateKey(privkey)

	if err != nil {
		return nil, nil, err
	}

	return privkey, pubkey, nil
}

// SchnorrPublicKeyFromPrivateKey computes a Schnorr public key from a private key
func (o *CfdgoCryptoService) SchnorrPublicKeyFromPrivateKey(privateKey *PrivateKey) (*SchnorrPublicKey, error) {
	bs, err := o.schnorrUtil.GetPubkeyFromPrivkey(cfdgo.NewByteData(privateKey.bytes))
	if err != nil {
		return nil, errors.WithMessage(err, "Error while calculating public key from private key")
	}
	pubkey, err := NewSchnorrPublicKey(bs.ToHex())
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

// ComputeSchnorrSignature computes a schnorr signature on the given message (will be hashed by sha256)
func (o *CfdgoCryptoService) ComputeSchnorrSignature(privateKey *PrivateKey, kvalue *PrivateKey, message string) (*Signature, error) {
	hash := sha256.Sum256([]byte(message))

	bs, err := o.schnorrUtil.SignWithNonce(cfdgo.NewByteData(hash[:]), cfdgo.NewByteData(privateKey.bytes), cfdgo.NewByteData(kvalue.bytes))
	if err != nil {
		return nil, errors.WithMessage(err, "Error while computing schnorr signature")
	}
	sig, err := NewSignature(bs.ToHex())
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// VerifySchnorrSignature verifies the schnorr signature against a given public key on the given message (will be hashed with sha256)
func (o *CfdgoCryptoService) VerifySchnorrSignature(publicKey *SchnorrPublicKey, signature *Signature, message string) (bool, error) {
	hash := sha256.Sum256([]byte(message))
	ok, err := o.schnorrUtil.Verify(cfdgo.NewByteData(signature.bytes), cfdgo.NewByteData(hash[:]), cfdgo.NewByteData(publicKey.bytes))
	if err != nil {
		return false, errors.WithMessage(err, "Error while verifying schnorr signature")
	}
	return ok, nil
}
