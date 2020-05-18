package router

import (
	"p2pderivatives-oracle/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Router represent a Router instance
// registered middleware :
//	- gin.Recovery
type Router struct {
	api         API
	ginRouter   *gin.Engine
	log         *log.Log
	logger      *logrus.Logger
	initialized bool
}

// NewRouter creates a new Router structure with the given parameters.
func NewRouter(l *log.Log, api API) *Router {
	return &Router{
		api:         api,
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

	if !r.api.AreServicesInitialized() {
		if err := r.api.InitializeServices(); err != nil {
			r.log.Logger.Errorf("Error while initializing api services")
			return err
		}

	}
	router := gin.New()

	// middlewares
	router.Use(gin.Recovery())
	router.Use(r.api.GlobalMiddlewares()...)

	// routes
	baseRoute := router.Group("")
	r.api.Routes(baseRoute)

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
	err := r.api.FinalizeServices()
	if err != nil {
		return errors.Errorf("failed to gracefully close http router")
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
