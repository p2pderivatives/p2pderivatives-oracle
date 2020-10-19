package api_test

import (
	"encoding/json"
	"fmt"
	"math"
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

var TestAssetConfig = &api.AssetConfig{
	StartDate: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
	Frequency: time.Hour,
	RangeD:    time.Hour * 48,
}

var InDbDLCData = &entity.DLCData{
	PublishedDate: TestAssetConfig.StartDate.Add(10 * TestAssetConfig.Frequency),
	AssetID:       TestAsset.AssetID,
	Rvalue:        "inDB rvalue",
	Signature:     "inDB signature",
	Value:         "inDB value",
	Kvalue:        "inDB kvalue",
}

type ResponseValue struct {
	AssetID   string
	Rvalue    string
	Kvalue    string
	Value     string
	Signature string
}

const datafeedValue = 100.06

var TestResponseValues = &ResponseValue{
	AssetID:   TestAsset.AssetID,
	Kvalue:    "d1e3dcdda619833ec2b91d4fd304e9be0ad85326c9f524dfbc53f443ab54063e",
	Rvalue:    "44b1350439fc9a098db6edd5bd417eb1aeaa17ec60f9e5a799605feebd5c19eb",
	Value:     fmt.Sprintf("%d", int(math.Round(datafeedValue))),
	Signature: "44b1350439fc9a098db6edd5bd417eb1aeaa17ec60f9e5a799605feebd5c19ebf742bea67ff64f738c0426d04a22b30fe61258e074c0c90a1b13ce29d11f4b67",
}

func SetupMockValues() (*dlccrypto.PrivateKey, *dlccrypto.SchnorrPublicKey, *dlccrypto.Signature, *float64, error) {
	kvalue, err := dlccrypto.NewPrivateKey(TestResponseValues.Kvalue)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	rvalue, err := dlccrypto.NewSchnorrPublicKey(TestResponseValues.Rvalue)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	signature, err := dlccrypto.NewSignature(TestResponseValues.Signature)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	value := datafeedValue

	return kvalue, rvalue, signature, &value, nil
}

func SetupAssetEngine(recorder *httptest.ResponseRecorder, o *oracle.Oracle, crypto dlccrypto.CryptoService, feed datafeed.DataFeed) (*gin.Context, *gin.Engine) {
	assetController := api.NewAssetController(TestAsset.AssetID, *TestAssetConfig)
	orm := test.NewOrm(&entity.Asset{}, &entity.DLCData{})
	orm.GetDB().Create(TestAsset)
	orm.GetDB().Create(InDbDLCData)
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
			StartDate: TestAssetConfig.StartDate,
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
	date := time.Now().UTC().Add(TestAssetConfig.RangeD + time.Hour)
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
	oracleService, err := NewTestOracleService()
	if err != nil {
		t.Error(err)
	}
	c, r := SetupAssetEngine(resp, oracleService, crypto, nil)
	route := GetRouteWithTimeParam(api.RouteGETAssetRvalue, InDbDLCData.PublishedDate)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code) {
		expected := api.NewDLCDataResponse(oracleService.PublicKey, InDbDLCData)
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
	oracleService, err := NewTestOracleService()
	if err != nil {
		t.Error(err)
	}
	c, r := SetupAssetEngine(resp, oracleService, crypto, nil)

	// 30 minutes before (< Frequency)
	date := InDbDLCData.PublishedDate.Add((30 * time.Minute) - TestAssetConfig.Frequency)
	route := GetRouteWithTimeParam(api.RouteGETAssetRvalue, date)

	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code) {
		expected := api.NewDLCDataResponse(oracleService.PublicKey, InDbDLCData)
		actual := &api.DLCDataResponse{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetRvalue_NotInDB_ReturnsCorrectValue(t *testing.T) {
	// parameters
	date := InDbDLCData.PublishedDate.Add(30 * time.Minute)

	// expected
	expected := &api.DLCDataResponse{
		OraclePublicKey: OraclePublicKey,
		PublishedDate:   InDbDLCData.PublishedDate.Add(TestAssetConfig.Frequency),
		AssetID:         InDbDLCData.AssetID,
		Rvalue:          TestResponseValues.Rvalue,
	}

	// setup mocks
	ctrl := gomock.NewController(t)
	kvalue, rvalue, _, _, err := SetupMockValues()
	if !assert.NoError(t, err) {
		t.Fail()
	}

	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	crypto.EXPECT().GenerateSchnorrKeyPair().Return(kvalue, rvalue, nil)

	oracleService, err := NewTestOracleService()
	if err != nil {
		t.Error(err)
	}

	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, oracleService, crypto, nil)
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
	date := InDbDLCData.PublishedDate.Add(time.Minute * 30)

	// expected
	expectedDate := InDbDLCData.PublishedDate.Add(TestAssetConfig.Frequency)
	expected := &api.DLCDataResponse{
		OraclePublicKey: OraclePublicKey,
		PublishedDate:   expectedDate,
		AssetID:         TestAsset.AssetID,
		Rvalue:          TestResponseValues.Rvalue,
		Signature:       TestResponseValues.Signature,
		Value:           TestResponseValues.Value,
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
		crypto.EXPECT().GenerateSchnorrKeyPair().Return(kvalue, rvalue, nil)
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
	oracleService, err := NewTestOracleService()
	if err != nil {
		t.Error(err)
	}
	c, r := SetupAssetEngine(resp, oracleService, crypto, nil)

	// 30 minutes before (< Frequency)
	date := InDbDLCData.PublishedDate.Add(-(TestAssetConfig.Frequency / 2))
	route := GetRouteWithTimeParam(api.RouteGETAssetSignature, date)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String()) {
		expected := api.NewDLCDataResponse(oracleService.PublicKey, InDbDLCData)
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

	date := time.Now().UTC().Add(30 * time.Minute)
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
