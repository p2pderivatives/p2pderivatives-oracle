package api

import (
	"github.com/gin-gonic/gin"

	"p2pderivatives-oracle/internal/api/hello"
)

// Controller interface
type Controller interface {
	Routes(route *gin.RouterGroup)
}

const (
	// HelloBaseRoute base route of hello api
	HelloBaseRoute = "/hello"
)

// BindAllRoutes binds all the routes to the router parameter
// using the configuration provided.
// Each api gin.RouterGroup should be unique
func BindAllRoutes(router *gin.Engine, config *Config) {
	hello.NewController(&hello.Config{Count: config.HelloCount}).Routes(router.Group(HelloBaseRoute))

	// TODO Add new Api
}
