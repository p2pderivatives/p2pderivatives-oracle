package oracle

import (
	"p2pderivatives-oracle/internal/dlccrypto"

	"github.com/cryptogarageinc/server-common-go/pkg/utils/file"

	"github.com/pkg/errors"
)

// Oracle represents an oracle with private key, public key pair
type Oracle struct {
	PrivateKey *dlccrypto.PrivateKey
	PublicKey  *dlccrypto.SchnorrPublicKey
}

// New returns a new Oracle instance
// the public key will be calculated from the private key
func New(privateKey *dlccrypto.PrivateKey, publicKey *dlccrypto.SchnorrPublicKey) *Oracle {
	return &Oracle{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

// FromConfig returns an oracle from configuration
// password has to be defined either from a file or directly in configuration (environment variable)
// in case of using a txt file as password, the first line will be considered as password
func FromConfig(config *Config, cryptoService dlccrypto.CryptoService) (*Oracle, error) {
	var pass string
	var err error
	if config.KeyPass == "" && config.KeyPassFile == "" {
		return nil, errors.Errorf("No password or password file provided for key %s", config.KeyFile)
	}
	if config.KeyPass != "" {
		pass = config.KeyPass
	}
	if config.KeyPassFile != "" {
		pass, err = file.ReadFirstLineFromFile(config.KeyPassFile)
		if err != nil {
			return nil, err
		}
	}
	privKey, err := dlccrypto.ReadPemKeyFile(config.KeyFile, []byte(pass))
	if err != nil {
		return nil, errors.WithMessage(err, "Could not recover Oracle Private Key")
	}
	publicKey, err := cryptoService.SchnorrPublicKeyFromPrivateKey(privKey)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not get public key from private key")
	}
	return New(privKey, publicKey), nil
}
