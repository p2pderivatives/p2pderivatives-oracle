package api_test

import (
	"net/http"
	"net/http/httptest"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/test"
	mock_datafeed "p2pderivatives-oracle/test/mock/datafeed"
	mock_dlccrypto "p2pderivatives-oracle/test/mock/dlccrypto"
	"path/filepath"
	"testing"

	conf "github.com/cryptogarageinc/server-common-go/pkg/configuration"

	"github.com/cryptogarageinc/server-common-go/pkg/rest/router"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	"github.com/gin-gonic/gin"
)

func SetupEngine(recorder *httptest.ResponseRecorder, ct api.Controller, middlewares ...gin.HandlerFunc) (*gin.Context, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	r.Use(middlewares...)
	ct.Routes(r.Group(""))
	return c, r
}

func SetupTestOracleAPI(mockct *gomock.Controller) (router.API, error) {
	apiConfig := &api.Config{}
	err := test.InitializeConfig(apiConfig)
	if err != nil {
		return nil, err
	}
	oracleService, err := NewTestOracleService()
	if err != nil {
		return nil, err
	}

	crypto := mock_dlccrypto.NewMockCryptoService(mockct)
	feed := mock_datafeed.NewMockDataFeed(mockct)

	return api.NewOracleAPI(
		apiConfig,
		test.NewLogger(),
		oracleService,
		test.NewOrm(),
		crypto,
		feed), nil
}

func TestOracleAPI_WithEngine_RoutesAccessible(t *testing.T) {
	ctrl := gomock.NewController(t)
	oracleApi, err := SetupTestOracleAPI(ctrl)
	if !assert.NoError(t, err) {
		t.Fail()
	}
	resp := httptest.NewRecorder()
	c, r := SetupEngine(resp, oracleApi, oracleApi.GlobalMiddlewares()...)

	route := api.OracleBaseRoute + api.RouteGETOraclePublicKey
	c.Request, _ = http.NewRequest(http.MethodGet, route, nil)
	r.ServeHTTP(resp, c.Request)

	// assert
	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestOracleAPI_InitializeServices_InitializedWithNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	oracleApi, err := SetupTestOracleAPI(ctrl)
	if err != nil {
		t.Fail()
	}
	err = oracleApi.InitializeServices()
	assert.NoError(t, err)
	assert.True(t, oracleApi.AreServicesInitialized())
}

func TestOracleAPI_FinalizeServices_NoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	oracleApi, err := SetupTestOracleAPI(ctrl)
	if err != nil {
		t.Fail()
	}
	err = oracleApi.InitializeServices()
	err = oracleApi.FinalizeServices()
	assert.NoError(t, err)
}

func TestAPIConfig_Initialize_NoError(t *testing.T) {
	path := filepath.Join("..", "..", "test", "config")
	c := conf.NewConfiguration("core", "integration", []string{path})
	err := c.Initialize()
	assert.NoError(t, err)
	apiConfig := &api.Config{}
	err = c.InitializeComponentConfig(apiConfig)
	assert.NoError(t, err)
}
