package dlccrypto

import (
	"bytes"
	"encoding/binary"
	"math"
	"p2pderivatives-oracle/internal/decompose"

	"github.com/sirupsen/logrus"
)

type bigSize struct {
	inner uint64
}

func (b *bigSize) write(buf *bytes.Buffer) {
	if b.inner <= 0xFC {
		val := uint8(b.inner)
		binary.Write(buf, binary.BigEndian, val)
	} else if b.inner <= 0xffff {
		binary.Write(buf, binary.BigEndian, uint8(0xFD))
		val := uint16(b.inner)
		binary.Write(buf, binary.BigEndian, val)
	} else if b.inner <= 0xFFFFFFFF {
		binary.Write(buf, binary.BigEndian, uint8(0xFE))
		val := uint32(b.inner)
		binary.Write(buf, binary.BigEndian, val)
	} else {
		binary.Write(buf, binary.BigEndian, uint8(0xFF))
		binary.Write(buf, binary.BigEndian, b.inner)
	}
}

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
		sig, err := cryptoService.ComputeSchnorrSignatureFixedK(privKey, kvalue, digit)
		if err != nil {
			return nil, nil, err
		}

		sigs[i] = sig.EncodeToString()
	}

	return sigs, decomposedValue, nil
}

// SerializeEvent serializes the given data as an OracleEvent
func SerializeEvent(
	nonces []SchnorrPublicKey, eventMaturity uint32, base uint16, isSigned bool, unit string, precision int32, nbDigits uint16, eventId string,
) []byte {
	buf := new(bytes.Buffer)
	nbNonces := uint16(len(nonces))
	binary.Write(buf, binary.BigEndian, nbNonces)

	for _, nonce := range nonces {
		buf.Write(nonce.bytes)
	}

	binary.Write(buf, binary.BigEndian, eventMaturity)
	digitDecompositionPrefix := bigSize{inner: 55306}
	digitDecompositionPrefix.write(buf)
	subBuf := new(bytes.Buffer)
	binary.Write(subBuf, binary.BigEndian, base)
	binary.Write(subBuf, binary.BigEndian, isSigned)
	unitLen := bigSize{inner: uint64(len(unit))}
	unitLen.write(subBuf)
	subBuf.Write([]byte(unit))
	binary.Write(subBuf, binary.BigEndian, precision)
	binary.Write(subBuf, binary.BigEndian, nbDigits)
	subBufLen := bigSize{inner: uint64(subBuf.Len())}
	subBufLen.write(buf)
	buf.Write(subBuf.Bytes())
	eventIdLen := bigSize{inner: uint64(len(eventId))}
	eventIdLen.write(buf)
	buf.Write([]byte(eventId))
	return buf.Bytes()
}

// GenerateEventSignature serializes the given data to the appropriate format
// and returns a Schnorr signature over the resulting data
func GenerateEventSignature(
	privKey *PrivateKey, nonces []SchnorrPublicKey, eventMaturity uint32, base uint16, isSigned bool, unit string, precision int32, nbDigits uint16, eventId string, cryptoService CryptoService,
) (string, error) {

	ser := SerializeEvent(nonces, eventMaturity, base, isSigned, unit, precision, nbDigits, eventId)

	sig, err := cryptoService.ComputeSchnorrSignature(privKey, ser)

	if err != nil {
		return "", err
	}

	return sig.EncodeToString(), nil
}
