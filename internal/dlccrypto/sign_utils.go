package dlccrypto

import (
	"math"
	"p2pderivatives-oracle/internal/decompose"

	"github.com/sirupsen/logrus"
)

// GetRoundedDecomposedSignaturesForValue rounds and decompose a given value and
// produces signatures over its digits using the provided private key and nonces.
func GetRoundedDecomposedSignaturesForValue(
	value float64, base int, nbDigits int, privKey *PrivateKey, kValues []string, cryptoService CryptoService) ([]string, []string, error) {
	// round datafeed price to neareast integer
	roundedValue := int(math.Round(value))
	// decompose value
	maxValue := int(math.Pow(float64(base), float64(nbDigits)) - 1)
	if roundedValue > maxValue {
		roundedValue = maxValue
	}
	decomposedValue := decompose.Value(roundedValue, base, nbDigits)
	if len(decomposedValue) != len(kValues) {
		logrus.Panic("Incompatible lengths for decomposed value")
	}
	sigs := make([]string, len(decomposedValue))
	for i, digit := range decomposedValue {
		kvalue, err := NewPrivateKey(kValues[i])
		if err != nil {
			return nil, nil, err
		}
		sig, err := cryptoService.ComputeSchnorrSignature(privKey, kvalue, digit)
		if err != nil {
			return nil, nil, err
		}

		sigs[i] = sig.EncodeToString()
	}

	return sigs, decomposedValue, nil
}
