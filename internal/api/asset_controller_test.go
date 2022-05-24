package api_test

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/internal/cfddlccrypto"
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
	Unit: "usd/btc",
}

var InDbDLCData = &entity.EventData{
	PublishedDate: TestAssetConfig.StartDate.Add(10 * TestAssetConfig.Frequency),
	AssetID:       TestAsset.AssetID,
	Nonces:        []string{"inDB rvalue 1", "inDB rvalue 2", "inDB rvalue 3"},
	Signatures:    []string{"inDB signature", "inDB signature", "inDB signature"},
	Values:        []string{"inDB value", "inDB value", "inDB value"},
	Base:          10,
	Kvalues:       []string{"inDB kvalue", "inDB kvalue", "inDB kvalue"},
}

type ResponseValue struct {
	AssetID               string
	Rvalues               []string
	Kvalues               []string
	Values                []string
	Signatures            []string
	AnnouncementSignature string
}

const datafeedValue = 100.06

var TestResponseValues = &ResponseValue{
	AssetID: TestAsset.AssetID,
	Kvalues: []string{"af1e8c793ee16165ff653310a83964f2dac9bc2f831b2c687f0463b6c6f6ae38", "cf9c8213904b85f3002b1c4f34a788e707f08e6e387c29510b63ca15dbb7955f", "ca828e462a51a9d45baa0ed9c54b920ac19f7d9b2d7b3331bcf29044872996a6"},
	Rvalues: []string{"5b08032667f80bf22cf2f3445f0b5b7abf653d4cdd5d0fcd43af277f01cae31a", "5b0210b2f4ece4c0fabf400dc4bbaa2bacf44542812b7e8dabd1c8f1501f3622", "2bf59128d9302d7a668fafd4d1511232dc991665973567ecd5156ec02f81d9a8"},
	Values:  decompose.Value(int(math.Round(datafeedValue)), 10, 3),
	Signatures: []string{
		"5b08032667f80bf22cf2f3445f0b5b7abf653d4cdd5d0fcd43af277f01cae31afb039d703165c37579b0e15557e8cdeababbb53100e095c132682c0ac39a99d2",
		"5b0210b2f4ece4c0fabf400dc4bbaa2bacf44542812b7e8dabd1c8f1501f3622101a0502ef8e71d0c5ebf403b87dd152076f2cbd4d380fff768911c50a396e7a",
		"2bf59128d9302d7a668fafd4d1511232dc991665973567ecd5156ec02f81d9a80772068d6faa2a026e722b367d48027dc1ad8a010d359afe11f9092cb252d017",
	},
	AnnouncementSignature: "246869c18a3f07f767cbc61f326d4447a2ee052430f87ec8fa659435849aecc10ef8afe468459a687d65992d819c19314f42cdf35ee4411d1408ae3559e02241",
}

func SetupMockValues() ([]*dlccrypto.PrivateKey, []*dlccrypto.SchnorrPublicKey, []*dlccrypto.Signature, *float64, error) {
	rand.Seed(1)
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
	orm := test.NewOrm(&entity.Asset{}, &entity.EventData{})
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

func TestAssetController_GetAssetAnnouncement_NotInConfigRange_ReturnsCorrectErrorResponse(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, nil, crypto, nil)
	date := time.Now().UTC().Add(TestAssetConfig.RangeD + time.Hour)
	route := GetRouteWithTimeParam(api.RouteGETAssetAnnouncement, date)
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

func TestAssetController_GetAssetAnnouncement_WithExactValidDateInDB_ReturnsCorrectValue(t *testing.T) {
	// arrange
	resp := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	oracleService, err := NewTestOracleService()
	if err != nil {
		t.Error(err)
	}
	c, r := SetupAssetEngine(resp, oracleService, crypto, nil)
	route := GetRouteWithTimeParam(api.RouteGETAssetAnnouncement, InDbDLCData.PublishedDate)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code) {
		expected := api.NewOracleAnnouncement(oracleService.PublicKey, InDbDLCData)
		actual := &api.OracleAnnouncement{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetAnnouncement_WithNearValidDateInDB_ReturnsCorrectValue(t *testing.T) {
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
	route := GetRouteWithTimeParam(api.RouteGETAssetAnnouncement, date)

	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code) {
		expected := api.NewOracleAnnouncement(oracleService.PublicKey, InDbDLCData)
		actual := &api.OracleAnnouncement{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetAnnouncement_NotInDB_ReturnsCorrectValue(t *testing.T) {
	// parameters
	date := InDbDLCData.PublishedDate.Add((30 * time.Minute) + (2 * time.Second))
	oracleService, err := NewTestOracleService()
	if err != nil {
		t.Error(err)
	}

	// expected
	updatedDlcData := entity.EventData{
		Timestamp:             InDbDLCData.Timestamp,
		PublishedDate:         InDbDLCData.PublishedDate.Add(TestAssetConfig.Frequency),
		AssetID:               InDbDLCData.AssetID,
		Nonces:                TestResponseValues.Rvalues,
		Asset:                 InDbDLCData.Asset,
		AnnouncementSignature: TestResponseValues.AnnouncementSignature,
		Base:                  InDbDLCData.Base,
		Precision:             InDbDLCData.Precision,
		Unit:                  TestAssetConfig.Unit,
	}

	expected := api.NewOracleAnnouncement(oracleService.PublicKey, &updatedDlcData)
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

	expectedSig, _ := dlccrypto.NewSignature(TestResponseValues.AnnouncementSignature)
	crypto.EXPECT().ComputeSchnorrSignature(oracleService.PrivateKey, gomock.Any()).Return(expectedSig, nil)

	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, oracleService, crypto, nil)
	route := GetRouteWithTimeParam(api.RouteGETAssetAnnouncement, date)

	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)
	r.ServeHTTP(resp, c.Request)

	if assert.Equal(t, http.StatusOK, resp.Code) {
		actual := &api.OracleAnnouncement{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetAttestation_NotInDB_ReturnsCorrectValue(t *testing.T) {
	// params
	date := InDbDLCData.PublishedDate.Add(time.Minute * 30)

	// expected
	expectedDate := InDbDLCData.PublishedDate.Add(TestAssetConfig.Frequency)
	updatedDlcData := &entity.EventData{
		Timestamp:             InDbDLCData.Timestamp,
		PublishedDate:         expectedDate,
		AssetID:               TestAsset.AssetID,
		Nonces:                TestResponseValues.Rvalues,
		Signatures:            TestResponseValues.Signatures,
		Values:                TestResponseValues.Values,
		Asset:                 InDbDLCData.Asset,
		AnnouncementSignature: TestResponseValues.AnnouncementSignature,
	}
	expected := api.NewOracleAttestation(updatedDlcData)

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
		feed.EXPECT().FindPastAssetPrice("btcusd", expectedDate).Return(sigValue, nil)
		// mock crypto
		crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
		for i := 0; i < len(kvalues); i++ {
			crypto.EXPECT().GenerateSchnorrKeyPair().Return(kvalues[i], rvalues[i], nil)
			crypto.EXPECT().ComputeSchnorrSignatureFixedK(
				oracleInstance.PrivateKey,
				kvalues[i],
				TestResponseValues.Values[i]).Return(sigs[i], nil)
		}

		expectedSig, _ := dlccrypto.NewSignature(TestResponseValues.AnnouncementSignature)
		crypto.EXPECT().ComputeSchnorrSignature(oracleInstance.PrivateKey, gomock.Any()).Return(expectedSig, nil)

		resp := httptest.NewRecorder()
		c, r := SetupAssetEngine(resp, oracleInstance, crypto, feed)
		route := GetRouteWithTimeParam(api.RouteGETAssetAttestation, date)
		c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

		// act
		r.ServeHTTP(resp, c.Request)

		// assert
		if assert.Equal(t, http.StatusOK, resp.Code) {
			actual := &api.OracleAttestation{}
			err := json.Unmarshal([]byte(resp.Body.String()), actual)
			if assert.NoError(t, err) {
				assert.Equal(t, expected, actual, resp.Body.String())
			}
		}
	}
}

func TestAssetController_GetAssetAttestation_WithNearValidDateInDB_ReturnsCorrectValue(t *testing.T) {
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
	route := GetRouteWithTimeParam(api.RouteGETAssetAttestation, date)
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)

	// act
	r.ServeHTTP(resp, c.Request)

	// assert
	if assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String()) {
		expected := api.NewOracleAttestation(InDbDLCData)
		actual := &api.OracleAttestation{}
		err := json.Unmarshal([]byte(resp.Body.String()), actual)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestAssetController_GetAssetAttestation_WithFutureDate_ReturnsBadRequestValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	crypto := mock_dlccrypto.NewMockCryptoService(ctrl)
	resp := httptest.NewRecorder()
	c, r := SetupAssetEngine(resp, nil, crypto, nil)

	date := time.Now().UTC().Add(30 * time.Minute)
	route := GetRouteWithTimeParam(api.RouteGETAssetAttestation, date)
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

func TestAssetController_GetAssetAttestation_HasValidAnnouncementSignature(t *testing.T) {
	oracleInstance, _ := NewTestOracleService()
	resp := httptest.NewRecorder()
	crypto := cfddlccrypto.NewCfdgoCryptoService()

	c, r := SetupAssetEngine(resp, oracleInstance, crypto, nil)
	route := GetRouteWithTimeParam(api.RouteGETAssetAnnouncement, time.Now().Add(1*time.Hour))
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)
	r.ServeHTTP(resp, c.Request)

	input := resp.Body.String()
	println(input)
	var announcement api.OracleAnnouncement
	err := json.Unmarshal([]byte(input), &announcement)
	assert.NoError(t, err)

	nonces := make([]dlccrypto.SchnorrPublicKey, 0)

	for _, s := range announcement.OracleEvent.Nonces {
		k, _ := dlccrypto.NewSchnorrPublicKey(s)
		nonces = append(nonces, *k)
	}

	ser := dlccrypto.SerializeEvent(
		nonces,
		uint32(announcement.OracleEvent.EventMaturityEpoch),
		uint16(announcement.OracleEvent.EventDescriptor.DigitDecompositionDescriptor.Base),
		announcement.OracleEvent.EventDescriptor.DigitDecompositionDescriptor.IsSigned,
		announcement.OracleEvent.EventDescriptor.DigitDecompositionDescriptor.Unit,
		int32(announcement.OracleEvent.EventDescriptor.DigitDecompositionDescriptor.Precision),
		uint16(announcement.OracleEvent.EventDescriptor.DigitDecompositionDescriptor.NbDigits),
		announcement.OracleEvent.EventID,
	)

	pubkey, err := dlccrypto.NewSchnorrPublicKey(announcement.OraclePublicKey)
	assert.NoError(t, err)

	sig, err := dlccrypto.NewSignature(announcement.AnnouncementSignature)
	assert.NoError(t, err)

	valid, _ := crypto.VerifySchnorrSignatureRaw(pubkey, sig, ser)
	assert.True(t, valid)
}
