package dlccrypto_test

import (
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/test"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var KeyFileDirectory = filepath.Join(test.VectorsDirectoryPath, "keys")

var TestKeyFile = []struct {
	path     string
	password string
	key      string
}{
	{path: "key_0.pem", password: "wKeEhq0DP/rNtcD8u/NxLyJYKmyKqOzklgOamGJlbSA=", key: "83e03f14bd6ae801ff21430ec4f745c8ca749945c9c0ba147216c2fd4a6e6df6"},
	{path: "key_1.pem", password: "z6Re1aGzRaVewoIX+3HHsR6dELtLL8aR4LFKLLCypJc=", key: "ceb79d2836048151d6b579e82c08554949422844a45b982b75b27ad76e840cb1"},
	{path: "key_2.pem", password: "om9fTkErrpZfF7Q85sC6AM+hjeYh9zxN7yi2iD7IWRQ=", key: "572a422a93577e6e89c17324600d50e618190a7a4870cdf135c392f88914ede1"},
	{path: "key_3.pem", password: "svMKJcSUkpXVrrwhiznzs7fWm7i3h9WwsdzCazB3NkA=", key: "b636e3cb3bd804e78cec51d27243b3fd370c4933f1bc318aed08fc6a7500691c"},
}

func Test_ReadPemKeyFile_WithValidPassword_Success(t *testing.T) {
	for _, v := range TestKeyFile {
		path := filepath.Join(KeyFileDirectory, v.path)
		actual, err := dlccrypto.ReadPemKeyFile(path, []byte(v.password))
		assert.NoError(t, err)
		assert.Equal(t, v.key, actual.EncodeToString())
	}
}

func Test_ReadPemKeyFile_WithInValidPassword_Error(t *testing.T) {
	path := filepath.Join(KeyFileDirectory, TestKeyFile[0].path)
	_, err := dlccrypto.ReadPemKeyFile(path, []byte("invalid pass"))
	assert.Error(t, err)
}

func Test_ReadPemKeyFile_WithInvalidPemFile_Error(t *testing.T) {
	path := filepath.Join(KeyFileDirectory)
	_, err := dlccrypto.ReadPemKeyFile(path, []byte("irrelevant"))
	assert.Error(t, err)
}
