package cfddlccrypto

import (
	"crypto/sha256"
	"math/rand"
	"p2pderivatives-oracle/internal/dlccrypto"

	cfdgo "github.com/cryptogarageinc/cfd-go"
	"github.com/pkg/errors"
)

// NewCfdgoCryptoService returns a CryptoService implemented using cfd-go library
func NewCfdgoCryptoService() dlccrypto.CryptoService {
	schnorrUtil := cfdgo.NewSchnorrUtil()
	return &CfdgoCryptoService{schnorrUtil}
}

// CfdgoCryptoService crypto service implementing the schnorr api of cfd-go library
type CfdgoCryptoService struct {
	schnorrUtil *cfdgo.SchnorrUtil
}

// GenerateSchnorrKeyPair returns a freshly generated Schnorr public/private key pair
func (o *CfdgoCryptoService) GenerateSchnorrKeyPair() (*dlccrypto.PrivateKey, *dlccrypto.SchnorrPublicKey, error) {
	_, seckey, _, err := cfdgo.CfdGoCreateKeyPair(true, 0)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Error while generating cfd go key pair")
	}

	privkey, err := dlccrypto.NewPrivateKey(seckey)
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
func (o *CfdgoCryptoService) SchnorrPublicKeyFromPrivateKey(privateKey *dlccrypto.PrivateKey) (*dlccrypto.SchnorrPublicKey, error) {
	bs, err := o.schnorrUtil.GetPubkeyFromPrivkey(*cfdgo.NewByteDataFromHexIgnoreError(privateKey.EncodeToString()))
	if err != nil {
		return nil, errors.WithMessage(err, "Error while calculating public key from private key")
	}
	pubkey, err := dlccrypto.NewSchnorrPublicKey(bs.ToHex())
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

// ComputeSchnorrSignature computes a schnorr signature on the given message (will be hashed by sha256)
func (o *CfdgoCryptoService) ComputeSchnorrSignatureFixedK(privateKey *dlccrypto.PrivateKey, kvalue *dlccrypto.PrivateKey, message string) (*dlccrypto.Signature, error) {
	hash := sha256.Sum256([]byte(message))

	bs, err := o.schnorrUtil.SignWithNonce(cfdgo.NewByteData(hash[:]),
		*cfdgo.NewByteDataFromHexIgnoreError(privateKey.EncodeToString()),
		*cfdgo.NewByteDataFromHexIgnoreError(kvalue.EncodeToString()))
	if err != nil {
		return nil, errors.WithMessage(err, "Error while computing schnorr signature")
	}
	sig, err := dlccrypto.NewSignature(bs.ToHex())
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// ComputeSchnorrSignature computes a schnorr signature on the given byte buffer message (will be hashed by sha256)
func (o *CfdgoCryptoService) ComputeSchnorrSignature(privateKey *dlccrypto.PrivateKey, message []byte) (*dlccrypto.Signature, error) {
	hash := sha256.Sum256(message)

	auxRand := make([]byte, 32)
	rand.Read(auxRand)

	bs, err := o.schnorrUtil.Sign(cfdgo.NewByteData(hash[:]),
		*cfdgo.NewByteDataFromHexIgnoreError(privateKey.EncodeToString()), cfdgo.NewByteData(auxRand))
	if err != nil {
		return nil, errors.WithMessage(err, "Error while computing schnorr signature")
	}
	sig, err := dlccrypto.NewSignature(bs.ToHex())
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// VerifySchnorrSignature verifies the schnorr signature against a given public key on the given message (will be hashed with sha256)
func (o *CfdgoCryptoService) VerifySchnorrSignature(publicKey *dlccrypto.SchnorrPublicKey, signature *dlccrypto.Signature, message string) (bool, error) {
	hash := sha256.Sum256([]byte(message))
	ok, err := o.schnorrUtil.Verify(*cfdgo.NewByteDataFromHexIgnoreError(signature.EncodeToString()), cfdgo.NewByteData(hash[:]), *cfdgo.NewByteDataFromHexIgnoreError(publicKey.EncodeToString()))
	if err != nil {
		return false, errors.WithMessage(err, "Error while verifying schnorr signature")
	}
	return ok, nil
}
