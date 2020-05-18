package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/datafeed"
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/internal/oracle"
	"p2pderivatives-oracle/test"
	mock_datafeed "p2pderivatives-oracle/test/mock/datafeed"
	mock_dlccrypto "p2pderivatives-oracle/test/mock/dlccrypto"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var TestAsset = &entity.Asset{
	AssetID:     "btcusd",
	Description: "Some test",
}

var TestDlcData = &entity.DLCData{
	PublishedDate: time.Now().Add(-2 * time.Hour).UTC().Round(time.Second),
	AssetID:       TestAsset.AssetID,
	Rvalue:        "inDB rvalue",
	Signature:     "inDB signature",
	Value:         "inDB value",
	Kvalue:        "inDB kvalue",
}

var TestAssetConfig = &api.AssetConfig{
	Frequency: time.Hour,
	RangeD:    time.Hour * 48,
}

type ResponseValue struct {
	AssetID   string
	Rvalue    string
	Kvalue    string
	Value     string
	Signature string
}

var TestResponseValues = &ResponseValue{
	AssetID:   TestAsset.AssetID,
	Kvalue:    "71957542da5ea5e6b58607a3e21f0eb4e871c785b5b3470fb5274213f050cd60",
	Rvalue:    "02d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a",
	Value:     "100",
	Signature: "26772fddf151463655823c906d1c01afb682b8c3533589fe71753a579bf59b22",
}

func SetupMockValues() (*dlccrypto.PrivateKey, *dlccrypto.PublicKey, *dlccrypto.Signature, uint64, error) {
	kvalue, err := dlccrypto.NewPrivateKey(TestResponseValues.Kvalue)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	rvalue, err := dlccrypto.NewPublicKey(TestResponseValues.Rvalue)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	signature, err := dlccrypto.NewSignature(TestResponseValues.Signature)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	value, err := strconv.ParseUint(TestResponseValues.Value, 10, 64)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	return kvalue, rvalue, signature, value, nil
}

func SetupAssetEngine(recorder *httptest.ResponseRecorder, o *oracle.Oracle, crypto dlccrypto.CryptoService, feed datafeed.DataFeed) (*gin.Context, *gin.Engine) {
	assetController := api.NewAssetController(TestAsset.AssetID, *TestAssetConfig)
	orm := test.NewOrm(&entity.Asset{}, &entity.DLCData{})
	orm.GetDB().Create(TestAsset)
	orm.GetDB().Create(TestDlcData)
	setup := func(c *gin.Context) {
		c.Set(api.ContextIDOracle, o)
		c.Set(api.ContextIDCryptoService, crypto)
		c.Set(api.ContextIDDataFeed, feed)
		c.Set(api.ContextIDOrm, orm)
	}
	c, r := SetupEngine(recorder, assetController, api.ErrorHandler(), setup)
	return c, r
}

func TestAssetController_GetConfiguration(t *testing.T) {
	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, nil, nil, nil)
	c.Request, _ = http.NewRequest(http.MethodGet, api.RouteGETAssetConfig, nil)
	r.ServeHTTP(resp, c.Request)
	if assert.Equal(t, http.StatusOK, resp.Code) {
		expected := &api.AssetConfigResponse{
			Frequency: "PT1H",
			RangeD:    "P2DT",
		}
		actual := &api.AssetConfigResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetRvalue_NotInConfigRange_ReturnsCorrectErrorResponse(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, nil, crypto, nil)
	date := time.Now().Add(TestAssetConfig.RangeD + time.Hour)
	route := GetRouteWithTimeParam(api.RouteGETAssetRvalue, date)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusBadRequest, resp.Code) {
		expectedCode := api.InvalidTimeTooLateBadRequestErrorCode
		actual := &api.ErrorResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedCode, actual.ErrorCode)
		}
	}
}

func TestAssetController_GetAssetRvalue_WithExactValidDateInDB_ReturnsCorrectValue(t *testing.T) {
	// arrange
	resp := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	c, r := SetupAssetEngine(resp, nil, crypto, nil)
	route := GetRouteWithTimeParam(api.RouteGETAssetRvalue, TestDlcData.PublishedDate)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code) {
		expected := api.NewResponse(TestDlcData)
		actual := &api.DLCDataResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetRvalue_WithNearValidDateInDB_ReturnsCorrectValue(t *testing.T) {
	// arrange
	resp := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	c, r := SetupAssetEngine(resp, nil, crypto, nil)

	// 30 minutes before (< Frequency)
	date := TestDlcData.PublishedDate.Add((30 * time.Minute) - TestAssetConfig.Frequency)
	route := GetRouteWithTimeParam(api.RouteGETAssetRvalue, date)

	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code) {
		expected := api.NewResponse(TestDlcData)
		actual := &api.DLCDataResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetRvalue_NotInDB_ReturnsCorrectValue(t *testing.T) {
	// parameters
	date := TestDlcData.PublishedDate.Add(30 * time.Minute)

	// expected
	expected := &api.DLCDataResponse{
		PublishedDate: TestDlcData.PublishedDate.Add(TestAssetConfig.Frequency),
		AssetID:       TestDlcData.AssetID,
		Rvalue:        TestResponseValues.Rvalue,
	}

	// setup mocks
	ctrl := gomock.NewController(t)
	kvalue, rvalue, _, _, err := SetupMockValues()
	if !assert.NoError(t, err) {
		t.Fail()
	}

	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	crypto.EXPECT().GenerateKvalue().Return(kvalue, nil)
	crypto.EXPECT().ComputeRvalue(kvalue).Return(rvalue, nil)

	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, nil, crypto, nil)
	route := GetRouteWithTimeParam(api.RouteGETAssetRvalue, date)

	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)
	r.ServeHTTP(resp, c.Request)

	if assert.Equal(t, http.StatusOK, resp.Code) {
		actual := &api.DLCDataResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetSignature_NotInDB_ReturnsCorrectValue(t *testing.T) {
	// params
	date := TestDlcData.PublishedDate.Add(time.Minute * 30)

	// expected
	expectedDate := TestDlcData.PublishedDate.Add(TestAssetConfig.Frequency)
	expected := &api.DLCDataResponse{
		PublishedDate: expectedDate,
		AssetID:       TestAsset.AssetID,
		Rvalue:        TestResponseValues.Rvalue,
		Signature:     TestResponseValues.Signature,
		Value:         TestResponseValues.Value,
	}

	oracleInstance, err := NewTestOracleService()
	if assert.NoError(t, err) {
		// setup mocks
		ctrl := gomock.NewController(t)
		kvalue, rvalue, sig, sigValue, err := SetupMockValues()
		if err != nil {
			t.Error(err)
		}
		// mock datafeed
		feed := mock_datafeed.NewMockDataFeed(ctrl)
		feed.EXPECT().FindPastAssetPrice("btc", "usd", expectedDate).Return(sigValue, nil)
		// mock crypto
		crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
		crypto.EXPECT().GenerateKvalue().Return(kvalue, nil)
		crypto.EXPECT().ComputeRvalue(kvalue).Return(rvalue, nil)
		crypto.EXPECT().ComputeSchnorrSignature(
			oracleInstance.PrivateKey,
			kvalue,
			TestResponseValues.Value).Return(sig, nil)

		resp := httptest.NewRecorder()
		c, r := SetupAssetEngine(resp, oracleInstance, crypto, feed)
		route := GetRouteWithTimeParam(api.RouteGETAssetSignature, date)
		c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

		// act
		r.ServeHTTP(resp, c.Request)

		// assert
		if assert.Equal(t, http.StatusOK, resp.Code) {
			actual := &api.DLCDataResponse{}
			err := json.Unmarshal([]byte(resp.Body.String()), actual)
			if assert.NoError(t, err) {
				assert.Equal(t, expected, actual, resp.Body.String())
			}
		}
	}
}

func TestAssetController_GetAssetSignature_WithNearValidDateInDB_ReturnsCorrectValue(t *testing.T) {
	resp := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	c, r := SetupAssetEngine(resp, nil, crypto, nil)

	// 30 minutes before (< Frequency)
	date := TestDlcData.PublishedDate.Add(-(TestAssetConfig.Frequency / 2))
	route := GetRouteWithTimeParam(api.RouteGETAssetSignature, date)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String()) {
		expected := api.NewResponse(TestDlcData)
		actual := &api.DLCDataResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetSignature_WithFutureDate_ReturnsBadRequestValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, nil, crypto, nil)

	date := time.Now().Add(30 * time.Minute)
	route := GetRouteWithTimeParam(api.RouteGETAssetSignature, date)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusBadRequest, resp.Code) {
		expectedErrorCode := api.InvalidTimeTooEarlyBadRequestErrorCode
		actual := &api.ErrorResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedErrorCode, actual.ErrorCode, actual)
		}
	}
}

func GetRouteWithTimeParam(route string, date time.Time) string {
	return strings.Replace(
		route,
		":"+api.URLParamTagTime,
		date.Format(api.TimeFormatISO8601),
		1)
}
