// +build cryptocompare

package cryptocompare_test

import (
	"p2pderivatives-oracle/internal/cryptocompare"
	"p2pderivatives-oracle/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	now            = time.Now()
	testAssets     = []string{"btc"}
	testCurrencies = []string{"usd", "jpy", "eur"}
	testPastTimes  = []time.Time{
		now.Add(-time.Hour),
		now.Add(-20 * time.Hour),
		now.Add(-48 * time.Hour),
		now.Add(-70 * time.Hour),
		now.Add(-100 * time.Hour),
	}
)

func TestNewClient_WithValidConfig_NoPanics(t *testing.T) {
	config := &cryptocompare.Config{}
	err := test.InitializeConfig(config)
	assert.NoError(t, err)
	assert.NotPanics(t, func() { cryptocompare.NewClient(config) })
}

func NewTestClient() *cryptocompare.Client {
	config := &cryptocompare.Config{}
	err := test.InitializeConfig(config)
	if err != nil {
		panic(err)
	}
	return cryptocompare.NewClient(config)
}

func TestClient_FindCurrentAssetPrice_NotInitialized_ReturnsError(t *testing.T) {
	client := NewTestClient()
	val, err := client.FindCurrentAssetPrice(testAssets[0], testCurrencies[0])
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestClient_FindPastAssetPrice_NotInitialized_ReturnsError(t *testing.T) {
	client := NewTestClient()
	val, err := client.FindPastAssetPrice(testAssets[0], testCurrencies[0], testPastTimes[0])
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestClient_IsInitialized_ReturnsCorrectValue(t *testing.T) {
	client := NewTestClient()
	assert.False(t, client.IsInitialized())
	client.Initialize()
	assert.True(t, client.IsInitialized())

}

func TestClient_FindCurrentAssetPrice_WithValidParameters_ReturnsCorrectValue(t *testing.T) {
	client := NewTestClient()
	client.Initialize()
	for _, asset := range testAssets {
		for _, cur := range testCurrencies {
			val, err := client.FindCurrentAssetPrice(asset, cur)
			assert.NoError(t, err)
			assert.IsType(t, float64(0), *val)
		}
	}
}

func TestClient_FindPastAssetPrice_WithFutureDate_ReturnsError(t *testing.T) {
	client := NewTestClient()
	client.Initialize()
	val, err := client.FindPastAssetPrice(testAssets[0], testCurrencies[0], time.Now().Add(time.Hour))
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestClient_FindPastAssetPrice_WithValidParameters_ReturnsCorrectValue(t *testing.T) {
	client := NewTestClient()
	client.Initialize()
	for _, asset := range testAssets {
		for _, cur := range testCurrencies {
			for _, date := range testPastTimes {
				val, err := client.FindPastAssetPrice(asset, cur, date)
				assert.NoError(t, err)
				assert.IsType(t, float64(0), *val)
			}
		}
	}
}
