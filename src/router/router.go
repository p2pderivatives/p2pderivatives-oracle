package router

import (
	"p2pderivatives-oracle/src/api"
	"p2pderivatives-oracle/src/database/orm"
	"p2pderivatives-oracle/src/log"
	"p2pderivatives-oracle/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Router represent an Router instance
type Router struct {
	ginRouter   *gin.Engine
	orm         *orm.ORM
	config      *Config
	log         *log.Log
	logger      *logrus.Logger
	initialized bool
}

// NewRouter creates a new Router structure with the given parameters.
func NewRouter(config *Config, orm *orm.ORM, l *log.Log) *Router {
	return &Router{
		orm:         orm,
		config:      config,
		log:         l,
		initialized: false,
	}
}

// Initialize initializes the Router structure.
func (r *Router) Initialize() error {

	if r.initialized {
		return nil
	}

	r.log.Logger.Info("Router initialization starts")
	defer r.log.Logger.Info("Router initialization end")

	if !r.orm.IsInitialized() {
		err := errors.New("ORM is not initialized")
		r.log.Logger.Error(err)
		return err
	}

	if r.config.APIConfig == nil {
		err := errors.New("Router's Api configuration is not set")
		r.log.Logger.Error(err)
		return err
	}

	router := gin.New()

	// middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.GinLogrus(r.log.Logger))
	router.Use(middleware.Orm(r.orm))
	router.Use(middleware.RequestID())

	// routes
	api.BindAllRoutes(router, r.config.APIConfig)

	r.ginRouter = router
	r.initialized = true
	return nil
}

// IsInitialized returns whether the router is initialized.
func (r *Router) IsInitialized() bool {
	return r.initialized
}

// Finalize releases the resources held by the router.
func (r *Router) Finalize() error {
	err := r.orm.Finalize()
	if err != nil {
		return errors.Errorf("failed to close http router")
	}
	return nil
}

// GetEngine returns the gin engine instance associated with the Router object. Panics if the
// object is not initialized.
func (r *Router) GetEngine() *gin.Engine {
	if !r.IsInitialized() {
		panic("Trying to access uninitialized Router object.")
	}

	return r.ginRouter
}
