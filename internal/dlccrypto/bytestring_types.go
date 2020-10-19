package dlccrypto

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

const (
	sizePrivateKey = 32
	sizePublicKey  = 32
	sizeSignature  = 64
)

// ErrInvalidBytestringSize represents a bytestring of wrong length for the struct type used
var ErrInvalidBytestringSize = errors.New("Invalid ByteString size")

// NewByteString returns a new ByteString instance
func NewByteString(bytestring string) (*ByteString, error) {
	bytes, err := hex.DecodeString(bytestring)
	if err != nil {
		return nil, err
	}
	return &ByteString{bytes: bytes}, nil
}

// ByteString represents a bytestring
type ByteString struct {
	bytes []byte
}

// EncodeToString encodes the bytestring to a hexadecimal string
func (b *ByteString) EncodeToString() string {
	return hex.EncodeToString(b.bytes)
}

// NewPrivateKey returns a new PrivateKey instance
func NewPrivateKey(bytestring string) (*PrivateKey, error) {
	bt, err := NewByteString(bytestring)
	if err != nil {
		return nil, err
	}
	if len(bt.bytes) != sizePrivateKey {
		return nil, invalidSizeError("PrivateKey", sizePrivateKey)
	}
	return &PrivateKey{*bt}, nil
}

// PrivateKey represents a private key
type PrivateKey struct {
	ByteString
}

// NewSchnorrPublicKey returns a new PublicKey instance
func NewSchnorrPublicKey(bytestring string) (*SchnorrPublicKey, error) {
	bt, err := NewByteString(bytestring)
	if err != nil {
		return nil, err
	}
	if len(bt.bytes) != sizePublicKey {
		return nil, invalidSizeError("PublicKey", sizePublicKey)
	}
	return &SchnorrPublicKey{*bt}, nil
}

// SchnorrPublicKey represents a compressed public key
type SchnorrPublicKey struct {
	ByteString
}

// NewSignature returns a new Signature instance
func NewSignature(bytestring string) (*Signature, error) {
	bt, err := NewByteString(bytestring)
	if err != nil {
		return nil, err
	}
	if len(bt.bytes) != sizeSignature {
		return nil, invalidSizeError("Signature", sizeSignature)
	}
	return &Signature{*bt}, nil
}

// Signature represents the the s value of a Schnorr signature.
// As the r value is already known to the users in DLC we do not include it here.
type Signature struct {
	ByteString
}

func invalidSizeError(name string, size int) error {
	return errors.WithMessagef(
		ErrInvalidBytestringSize,
		"Bytestring of %s should be of %d byte size",
		name,
		size)
}
