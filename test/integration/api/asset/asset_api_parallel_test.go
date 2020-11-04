// +build integration

package assetapi_test

import (
	"fmt"
	"net/http"
	"p2pderivatives-oracle/internal/api"
	helper "p2pderivatives-oracle/test/integration"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/stretchr/testify/assert"
)

type TestTime struct {
	FrequencyCoefficient int
	DurationToAdd        time.Duration
}

var TestDuration = []time.Duration{
	0,
	30 * time.Second,
	30 * time.Minute,
	3 * time.Hour,
	10 * time.Hour,
	10*time.Hour + 30*time.Second,
	13*time.Hour + 3*time.Second,
	17*time.Hour + 30*time.Minute,
	24*time.Hour + 30*time.Minute,
}

const RepeatRequestCounter = 5

func concurrentTestBase(
	t *testing.T,
	durations []time.Duration,
	repeatRequestCounter int,
	assets map[string]api.AssetConfig,
	routeBuilder func(string, time.Time) string,
	extraCheck func(*assert.Assertions, interface{}) bool,
	computeDuration func(dur time.Duration, config api.AssetConfig) time.Duration,
	result interface{}) {
	apiResults := map[string]interface{}{}

	for asset, config := range assets {
		t.Run(fmt.Sprintf("asset %s", asset), func(subT *testing.T) {
			assertSub := assert.New(subT)
			subT.Parallel()

			// arrange
			handler := helper.NewHttpConcurrentHandler()
			for _, dur := range durations {
				requestedDate := Now.Add(computeDuration(dur, config))
				route := routeBuilder(asset, requestedDate)
				for i := 0; i < repeatRequestCounter; i++ {
					client := helper.CreateDefaultClient()
					req := client.R().SetResult(result)
					handler.RegisterRequest(req, resty.MethodGet, route)
				}
			}

			// act
			results := handler.RunAndWait()

			// assert
			for _, r := range results {
				if assertSub.NoError(r.Error, "Error while sending the request") {
					resp := r.Response
					if assertSub.Equal(http.StatusOK, resp.StatusCode(), resp.String()) {
						actual := resp.Result()
						if val, ok := apiResults[resp.Request.URL]; ok {
							assertSub.Equal(val, actual)
						} else {
							extraCheck(assertSub, actual)
							apiResults[resp.Request.URL] = actual
						}
					}
				}
			}
		})
	}

}

func dummyCheck(t *assert.Assertions, resp interface{}) bool {
	return true
}

func dummyDuration(dur time.Duration, config api.AssetConfig) time.Duration {
	return dur
}

func signatureDuration(dur time.Duration, config api.AssetConfig) time.Duration {
	return -(config.Frequency + dur)
}

func signatureCheck(t *assert.Assertions, resp interface{}) bool {
	attestation, ok := resp.(*api.OracleAttestation)
	t.True(ok)
	return assertSignature(t, attestation.Signatures, attestation.Values)
}

func TestGetAssetRvalue_Concurrent_WithValidTime_ReturnsCorrectValue(t *testing.T) {
	t.Parallel()

	concurrentTestBase(t, TestDuration, RepeatRequestCounter, helper.APIConfig.AssetConfigs, GetRouteAssetAnnouncement, dummyCheck, dummyDuration, &api.OracleAnnouncement{})
}

func TestGetAssetSignature_Concurrent_WithValidTime_ReturnsCorrectValue(t *testing.T) {
	t.Parallel()

	concurrentTestBase(t, TestDuration, RepeatRequestCounter, helper.APIConfig.AssetConfigs, GetRouteAssetAttestation, signatureCheck, signatureDuration, &api.OracleAttestation{})
}
