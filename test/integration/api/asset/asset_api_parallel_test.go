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

const RepeatRequestCounter = 3

func TestGetAssetRvalue_Concurrent_WithValidTime_ReturnsCorrectValue(t *testing.T) {
	t.Parallel()

	apiResults := map[string]*api.DLCDataResponse{}

	for asset := range helper.APIConfig.AssetConfigs {
		t.Run(fmt.Sprintf("asset %s", asset), func(subT *testing.T) {
			assertSub := assert.New(subT)
			subT.Parallel()

			// arrange
			handler := helper.NewHttpConcurrentHandler()
			for _, dur := range TestDuration {
				requestedDate := Now.Add(dur)
				route := GetRouteAssetRvalue(asset, requestedDate)
				for i := 0; i < RepeatRequestCounter; i++ {
					client := helper.CreateDefaultClient()
					req := client.R().SetResult(&api.DLCDataResponse{})
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
						actual := resp.Result().(*api.DLCDataResponse)
						if val, ok := apiResults[resp.Request.URL]; ok {
							assertSub.Equal(val, actual)
						} else {
							assertSub.Equal(asset, actual.AssetID)
							apiResults[resp.Request.URL] = actual
						}
					}
				}
			}
		})
	}
}

func TestGetAssetSignature_Concurrent_WithValidTime_ReturnsCorrectValue(t *testing.T) {
	t.Parallel()

	apiResults := map[string]*api.DLCDataResponse{}

	for asset, config := range helper.APIConfig.AssetConfigs {
		t.Run(fmt.Sprintf("asset %s", asset), func(subT *testing.T) {
			assertSub := assert.New(subT)
			subT.Parallel()

			// arrange
			handler := helper.NewHttpConcurrentHandler()
			for _, dur := range TestDuration {
				// past date for signature so negative duration
				requestedDate := Now.Add(-(config.Frequency + dur))
				route := GetRouteAssetSignature(asset, requestedDate)
				for i := 0; i < RepeatRequestCounter; i++ {
					client := helper.CreateDefaultClient()
					req := client.R().SetResult(&api.DLCDataResponse{})
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
						actual := resp.Result().(*api.DLCDataResponse)
						if val, ok := apiResults[resp.Request.URL]; ok {
							assertSub.Equal(val, actual)
						} else {
							assertSub.Equal(asset, actual.AssetID)
							assertSignature(assertSub, actual.Rvalue, actual.Signature, actual.Value)
							apiResults[resp.Request.URL] = actual
						}
					}
				}
			}
		})
	}
}
