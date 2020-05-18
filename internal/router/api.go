package router

import "github.com/gin-gonic/gin"

// API should implement the routes logic and services that will be associated with the router
// as well as eventual default middleware (the order matters)
type API interface {
	Routes(route *gin.RouterGroup)
	GlobalMiddlewares() []gin.HandlerFunc
	InitializeServices() error
	AreServicesInitialized() bool
	FinalizeServices() error
}
