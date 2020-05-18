package api

import (
	"fmt"
	"net/http"
	"p2pderivatives-oracle/internal/database/orm"
	"p2pderivatives-oracle/internal/datafeed"
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/internal/log"
	"p2pderivatives-oracle/internal/middleware"
	"p2pderivatives-oracle/internal/oracle"
	"p2pderivatives-oracle/internal/router"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Controller interface
type Controller interface {
	Routes(route *gin.RouterGroup)
}

const (
	// TimeFormatISO8601 time format of the api using ISO8601
	TimeFormatISO8601 = "2006-01-02T15:04:05Z"
)

const (
	// AssetBaseRoute base route of asset api
	AssetBaseRoute = "/asset"
	// OracleBaseRoute base route of oracle api
	OracleBaseRoute = "/oracle"
)

// NewOracleAPI returns a new oracle api instance
func NewOracleAPI(config *Config, log *log.Log, oracle *oracle.Oracle, orm *orm.ORM, cryptoService dlccrypto.CryptoService, feed datafeed.DataFeed) router.API {
	return &OracleAPI{
		logger:        log,
		config:        config,
		oracle:        oracle,
		orm:           orm,
		cryptoService: cryptoService,
		feed:          feed,
	}
}

// OracleAPI represents an oracle api containing the related services (like crypto service)
type OracleAPI struct {
	Controller
	logger        *log.Log
	config        *Config
	oracle        *oracle.Oracle
	orm           *orm.ORM
	cryptoService dlccrypto.CryptoService
	feed          datafeed.DataFeed
}

// Routes defines (and attached to a gin.routerGroup) the routes of the api
func (a *OracleAPI) Routes(route *gin.RouterGroup) {
	NewOracleController().Routes(route.Group(OracleBaseRoute))
	assetRoutes := []string{}
	for assetID, config := range a.config.AssetConfigs {
		assetRoute := fmt.Sprintf("%s/%s", AssetBaseRoute, assetID)
		assetRoutes = append(assetRoutes, assetID)
		NewAssetController(assetID, config).Routes(route.Group(assetRoute))
	}

	route.Group(AssetBaseRoute).GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, assetRoutes)
	})
}

// GlobalMiddlewares returns the global middlewares that the api should use
func (a *OracleAPI) GlobalMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.GinLogrus(a.logger.Logger),
		middleware.RequestID(ContextIDRequestID),
		ErrorHandler(),
		middleware.AddToContext(ContextIDOracle, a.oracle),
		middleware.AddToContext(ContextIDOrm, a.orm),
		middleware.AddToContext(ContextIDCryptoService, a.cryptoService),
		middleware.AddToContext(ContextIDDataFeed, a.feed),
	}
}

// InitializeServices initializes the api services
func (a *OracleAPI) InitializeServices() error {
	if !a.orm.IsInitialized() {
		if err := a.orm.Initialize(); err != nil {
			return err
		}
	}

	if a.cryptoService == nil {
		err := errors.New("Crypto Service is not set")
		return err
	}

	if a.feed == nil {
		err := errors.New("DataFeed Service is not set")
		return err
	}

	return nil
}

// AreServicesInitialized returns a boolean to check if the services are initialized
func (a *OracleAPI) AreServicesInitialized() bool {
	return a.orm.IsInitialized()
}

// FinalizeServices releases the resources held by the api services
func (a *OracleAPI) FinalizeServices() error {
	return a.orm.Finalize()
}
