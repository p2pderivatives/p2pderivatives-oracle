package dlccrypto

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// ReadPemKeyFile returns an secp256k1 private key from a (encrypted) pem file.
// If no password necessary, use nil
func ReadPemKeyFile(filePath string, pass []byte) (*PrivateKey, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "Error while checking for file at %s", filePath)
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(content)
	var keyDer []byte
	if pemBlock == nil {
		return nil, errors.New("The file is not of PEM format")
	}
	if pass != nil {
		keyDer, err = x509.DecryptPEMBlock(pemBlock, pass)
		if err != nil {
			return nil, err
		}
	} else {
		keyDer = pemBlock.Bytes
	}

	priv := make([]byte, 32)
	copy(priv, keyDer[7:39])
	bs := hex.EncodeToString(priv)

	return NewPrivateKey(bs)
}
