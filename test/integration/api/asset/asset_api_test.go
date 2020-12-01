// +build integration

package assetapi_test

import (
	"fmt"
	"net/http"
	"os"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/internal/cfddlccrypto"
	"p2pderivatives-oracle/internal/dlccrypto"
	helper "p2pderivatives-oracle/test/integration"
	"testing"
	"time"

	"github.com/cryptogarageinc/server-common-go/pkg/utils/iso8601"

	"github.com/stretchr/testify/assert"
)

var (
	// now in UTC truncated to second to avoid precision mismatch
	Now = time.Now().UTC().Truncate(time.Second)
)

func TestMain(m *testing.M) {
	helper.InitHelper()
	os.Exit(m.Run())
}

func TestGetAvailableAssets_ReturnsCorrectValue(t *testing.T) {
	// arrange
	client := helper.CreateDefaultClient()
	req := client.R().SetResult([]string{})

	// no Set in go using a map with bool instead check equality
	expected := map[string]bool{}
	for key := range helper.APIConfig.AssetConfigs {
		expected[key] = true
	}

	// act
	resp, err := req.Get(api.AssetBaseRoute)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	actual := resp.Result().(*[]string)
	actualMap := map[string]bool{}
	for _, key := range *actual {
		actualMap[key] = true
	}
	assert.Equal(t, expected, actualMap)
}

func TestGetAssetConfig_WithValidAssets_ReturnsCorrectValue(t *testing.T) {
	for asset, expectedConfig := range helper.APIConfig.AssetConfigs {
		t.Run(fmt.Sprintf("asset %s", asset), func(t *testing.T) {
			// arrange
			client := helper.CreateDefaultClient()
			req := client.R().SetResult(&api.AssetConfigResponse{})
			// act

			resp, err := req.Get(api.AssetBaseRoute + "/" + asset + api.RouteGETAssetConfig)
			assert.NoError(t, err)
			actual := resp.Result().(*api.AssetConfigResponse)
			assert.Equal(t, expectedConfig.StartDate, actual.StartDate)
			assert.Equal(t, iso8601.EncodeDuration(expectedConfig.Frequency), actual.Frequency)
			assert.Equal(t, iso8601.EncodeDuration(expectedConfig.RangeD), actual.RangeD)
		})
	}
}

func TestGetAssetConfig_WithInValidAssets_Returns404NotFound(t *testing.T) {
	// arrange
	client := helper.CreateDefaultClient()
	req := client.R()

	// act
	resp, err := req.Get(api.AssetBaseRoute + "/" + "unknownasset" + api.RouteGETAssetConfig)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())
}

func TestGetAssetRvalue_WithTimeNotInRange_ReturnsCorrectErrorResponse(t *testing.T) {
	for asset, config := range helper.APIConfig.AssetConfigs {
		t.Run(fmt.Sprintf("asset %s", asset), func(t *testing.T) {
			// arrange
			client := helper.CreateDefaultClient()
			req := client.R().SetError(&api.ErrorResponse{})
			requestedDate := Now.Add(config.RangeD + 30*time.Minute)

			// act
			resp, err := req.Get(GetRouteAssetRvalue(asset, requestedDate))

			// assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
			actual := resp.Error().(*api.ErrorResponse)
			assert.Equal(t, api.InvalidTimeTooLateBadRequestErrorCode, actual.ErrorCode)
		})
	}
}

func TestGetAssetRvalue_WithValidTime_ReturnsCorrectValue(t *testing.T) {
	for asset, config := range helper.APIConfig.AssetConfigs {
		t.Run(fmt.Sprintf("asset %s", asset), func(subT *testing.T) {
			assertSub := assert.New(subT)
			// arrange
			client := helper.CreateDefaultClient()
			req := client.R().SetResult(&api.DLCDataResponse{})
			requestedDate := Now.Add(30 * time.Minute)

			// act
			resp, err := req.Get(GetRouteAssetRvalue(asset, requestedDate))

			// assert
			assertSub.NoError(err)
			assertSub.Equal(http.StatusOK, resp.StatusCode())
			actual := resp.Result().(*api.DLCDataResponse)

			assertSub.Equal(asset, actual.AssetID)
			assertPublishedDate(assertSub, requestedDate, actual.PublishedDate, config.Frequency)
			_, err = dlccrypto.NewSchnorrPublicKey(actual.Rvalue)
			assertSub.NoError(err)
		})
	}
}

func TestGetAssetSignature_ReturnsCorrectValue(t *testing.T) {
	for asset, config := range helper.APIConfig.AssetConfigs {
		t.Run(fmt.Sprintf("asset %s", asset), func(subT *testing.T) {
			assertSub := assert.New(subT)
			// arrange
			client := helper.CreateDefaultClient()
			req := client.R().SetResult(&api.DLCDataResponse{})
			requestedDate := Now.Add(-(config.Frequency + config.Frequency/2))

			// act
			resp, err := req.Get(GetRouteAssetSignature(asset, requestedDate))

			// assert
			assertSub.NoError(err)
			if assertSub.Equal(http.StatusOK, resp.StatusCode(), resp.String()) {
				actual := resp.Result().(*api.DLCDataResponse)

				assertSub.Equal(asset, actual.AssetID)
				assertPublishedDate(assertSub, requestedDate, actual.PublishedDate, config.Frequency)
				assertSignature(assertSub, actual)
			}
		})
	}
}

func TestGetAssetSignature_WithFutureTime_ReturnsError(t *testing.T) {
	for asset := range helper.APIConfig.AssetConfigs {
		t.Run(fmt.Sprintf("asset %s", asset), func(t *testing.T) {
			// arrange
			client := helper.CreateDefaultClient()
			req := client.R().SetError(&api.ErrorResponse{})
			requestedDate := Now.Add(time.Hour)

			// act
			resp, err := req.Get(GetRouteAssetSignature(asset, requestedDate))

			// assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
			actual := resp.Error().(*api.ErrorResponse)
			assert.Equal(t, api.InvalidTimeTooEarlyBadRequestErrorCode, actual.ErrorCode)
		})
	}
}

// IsValidPublishedDate checks published date is after requested date and in frequency range
func IsValidPublishedDate(requestedDate time.Time, publishedDate time.Time, frequency time.Duration) bool {
	diff := publishedDate.Sub(requestedDate)
	return 0 <= diff && diff < frequency
}

func assertPublishedDate(
	assertSub *assert.Assertions,
	requestedDate time.Time,
	publishedDate time.Time,
	frequency time.Duration) bool {
	return assertSub.True(
		IsValidPublishedDate(requestedDate, publishedDate, frequency),
		"Invalid Published date, requested date: %v | actual published date: %v",
		requestedDate,
		publishedDate)
}

func assertSignature(assertSub *assert.Assertions, resp *api.DLCDataResponse) bool {
	sig, err := dlccrypto.NewSignature(resp.Signature)
	assertSub.NoError(err)
	isValidSignature, err := cfddlccrypto.NewCfdgoCryptoService().VerifySchnorrSignature(
		helper.ExpectedOracle.PublicKey,
		sig,
		resp.Value)
	assertSub.NoError(err)
	return assertSub.True(
		isValidSignature,
		"Signature %v does not match using oracle public key: %s rvalue:　%s message: %s",
		sig,
		helper.ExpectedOracle.PublicKey.EncodeToString(),
		resp.Rvalue,
		resp.Signature)
}
