package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/internal/oracle"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

const OraclePrivateKey = "18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"
const OraclePublicKey = "0250863ad64a87ae8a2fe83c1af1a8403cb53f53e486d8511dad8a04887e5b2352"

func NewTestOracleService() (*oracle.Oracle, error) {
	priv, err := dlccrypto.NewPrivateKey(OraclePrivateKey)
	if err != nil {
		return nil, err
	}
	return oracle.New(priv)
}

func SetupOracleEngine(recorder *httptest.ResponseRecorder, o *oracle.Oracle) (*gin.Context, *gin.Engine) {
	oracleController := api.NewOracleController()
	setup := func(c *gin.Context) {
		c.Set(api.ContextIDOracle, o)
	}
	c, r := SetupEngine(recorder, oracleController, setup)
	return c, r
}

func TestOracleController_GetPublicKey_ReturnsCorrectValue(t *testing.T) {
	oracleService, err := NewTestOracleService()
	if assert.NoError(t, err) {
		resp := httptest.NewRecorder()
		c, r := SetupOracleEngine(resp, oracleService)
		c.Request, _ = http.NewRequest(http.MethodGet, api.RouteGETOraclePublicKey, nil)
		r.ServeHTTP(resp, c.Request)
		if assert.Equal(t, http.StatusOK, resp.Code) {
			expected := &api.OraclePublicKeyResponse{
				PublicKey: OraclePublicKey,
			}
			actual := &api.OraclePublicKeyResponse{}
			err := json.Unmarshal([]byte(resp.Body.String()), &actual)
			if assert.NoError(t, err) {
				assert.Equal(t, expected, actual)
			}
		}
	}
}
