package api

import (
	"net/http"
	"p2pderivatives-oracle/internal/oracle"

	ginlogrus "github.com/Bose/go-gin-logrus"

	"github.com/gin-gonic/gin"
)

// RouteGETOraclePublicKey route for the GET oracle public key from OracleController
const RouteGETOraclePublicKey = "/publickey"

// OracleController represents the oracle api Controller
type OracleController struct {
}

// NewOracleController creates a new Controller structure with the given parameters.
func NewOracleController() Controller {
	return &OracleController{}
}

// Routes list and binds all routes to the router group provided
func (ct *OracleController) Routes(route *gin.RouterGroup) {
	route.GET(RouteGETOraclePublicKey, ct.GetPublicKey)
}

// GetPublicKey handler returns the Oracle public key
func (ct *OracleController) GetPublicKey(c *gin.Context) {
	ginlogrus.SetCtxLoggerHeader(c, "request-header", "Get Oracle Public Key")
	logger := ginlogrus.GetCtxLogger(c)
	oracleInstance := c.MustGet(ContextIDOracle).(*oracle.Oracle)
	logger.Info("Accessing Oracle instance")
	c.JSON(http.StatusOK, &OraclePublicKeyResponse{
		PublicKey: oracleInstance.PublicKey.EncodeToString(),
	})
}
