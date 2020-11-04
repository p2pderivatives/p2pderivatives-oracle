package api_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/datafeed"
	"p2pderivatives-oracle/internal/decompose"
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
	SignConfig: api.SigningConfig{
		Base:     10,
		NbDigits: 3,
	},
}

var InDbDLCData = &entity.DLCData{
	PublishedDate: TestAssetConfig.StartDate.Add(10 * TestAssetConfig.Frequency),
	AssetID:       TestAsset.AssetID,
	Rvalues:       []string{"inDB rvalue 1", "inDB rvalue 2", "inDB rvalue 3"},
	Signatures:    []string{"inDB signature", "inDB signature", "inDB signature"},
	Values:        []string{"inDB value", "inDB value", "inDB value"},
	Kvalues:       []string{"inDB kvalue", "inDB kvalue", "inDB kvalue"},
}

type ResponseValue struct {
	AssetID    string
	Rvalues    []string
	Kvalues    []string
	Values     []string
	Signatures []string
}

const datafeedValue = 100.06

var TestResponseValues = &ResponseValue{
	AssetID: TestAsset.AssetID,
	Kvalues: []string{"af1e8c793ee16165ff653310a83964f2dac9bc2f831b2c687f0463b6c6f6ae38", "cf9c8213904b85f3002b1c4f34a788e707f08e6e387c29510b63ca15dbb7955f", "ca828e462a51a9d45baa0ed9c54b920ac19f7d9b2d7b3331bcf29044872996a6"},
	Rvalues: []string{"5b08032667f80bf22cf2f3445f0b5b7abf653d4cdd5d0fcd43af277f01cae31a", "5b0210b2f4ece4c0fabf400dc4bbaa2bacf44542812b7e8dabd1c8f1501f3622", "2bf59128d9302d7a668fafd4d1511232dc991665973567ecd5156ec02f81d9a8"},
	Values:  decompose.DecomposeValue(int(math.Round(datafeedValue)), 10, 3),
	Signatures: []string{
		"5b08032667f80bf22cf2f3445f0b5b7abf653d4cdd5d0fcd43af277f01cae31afb039d703165c37579b0e15557e8cdeababbb53100e095c132682c0ac39a99d2",
		"5b0210b2f4ece4c0fabf400dc4bbaa2bacf44542812b7e8dabd1c8f1501f3622101a0502ef8e71d0c5ebf403b87dd152076f2cbd4d380fff768911c50a396e7a",
		"2bf59128d9302d7a668fafd4d1511232dc991665973567ecd5156ec02f81d9a80772068d6faa2a026e722b367d48027dc1ad8a010d359afe11f9092cb252d017",
	},
}

func SetupMockValues() ([]*dlccrypto.PrivateKey, []*dlccrypto.SchnorrPublicKey, []*dlccrypto.Signature, *float64, error) {
	nb := len(TestResponseValues.Kvalues)
	kvalues := make([]*dlccrypto.PrivateKey, nb)
	rvalues := make([]*dlccrypto.SchnorrPublicKey, nb)
	signatures := make([]*dlccrypto.Signature, nb)
	for i := 0; i < len(TestResponseValues.Kvalues); i++ {
		var err error
		kvalues[i], err = dlccrypto.NewPrivateKey(TestResponseValues.Kvalues[i])
		if err != nil {
			return nil, nil, nil, nil, err
		}
		rvalues[i], err = dlccrypto.NewSchnorrPublicKey(TestResponseValues.Rvalues[i])
		if err != nil {
			return nil, nil, nil, nil, err
		}
		signatures[i], err = dlccrypto.NewSignature(TestResponseValues.Signatures[i])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	value := datafeedValue

	return kvalues, rvalues, signatures, &value, nil
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

func TestAssetController_GetAssetRvalues_NotInDB_ReturnsCorrectValue(t *testing.T) {
	// parameters
	date := InDbDLCData.PublishedDate.Add(30 * time.Minute)

	// expected
	expected := &api.DLCDataResponse{
		OraclePublicKey: OraclePublicKey,
		PublishedDate:   InDbDLCData.PublishedDate.Add(TestAssetConfig.Frequency),
		AssetID:         InDbDLCData.AssetID,
		Rvalues:         TestResponseValues.Rvalues,
	}

	// setup mocks
	ctrl := gomock.NewController(t)
	kvalue, rvalue, _, _, err := SetupMockValues()
	if !assert.NoError(t, err) {
		t.Fail()
	}

	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	for i := 0; i < len(kvalue); i++ {
		crypto.EXPECT().GenerateSchnorrKeyPair().Return(kvalue[i], rvalue[i], nil)
	}

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
		Rvalues:         TestResponseValues.Rvalues,
		Signatures:      TestResponseValues.Signatures,
		Values:          TestResponseValues.Values,
	}

	oracleInstance, err := NewTestOracleService()
	if assert.NoError(t, err) {
		// setup mocks
		ctrl := gomock.NewController(t)
		kvalues, rvalues, sigs, sigValue, err := SetupMockValues()
		if err != nil {
			t.Error(err)
		}
		// mock datafeed
		feed := mock_datafeed.NewMockDataFeed(ctrl)
		feed.EXPECT().FindPastAssetPrice("btc", "usd", expectedDate).Return(sigValue, nil)
		// mock crypto
		crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
		for i := 0; i < len(kvalues); i++ {
			crypto.EXPECT().GenerateSchnorrKeyPair().Return(kvalues[i], rvalues[i], nil)
			crypto.EXPECT().ComputeSchnorrSignature(
				oracleInstance.PrivateKey,
				kvalues[i],
				TestResponseValues.Values[i]).Return(sigs[i], nil)
		}

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
