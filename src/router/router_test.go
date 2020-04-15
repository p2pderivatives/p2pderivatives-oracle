package router

import (
	"p2pderivatives-oracle/src/api"
	"p2pderivatives-oracle/src/database/orm"
	"p2pderivatives-oracle/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestRouter() *Router {
	config := test.GetTestConfig()
	ormConfig := orm.Config{}
	config.InitializeComponentConfig(&ormConfig)
	l := test.GetTestLogger(config)
	ormInstance := orm.NewORM(&ormConfig, l)
	ormInstance.Initialize()
	apiConfig := &api.Config{}
	config.InitializeComponentConfig(apiConfig)
	routerConfig := &Config{APIConfig: apiConfig}
	config.InitializeComponentConfig(routerConfig)
	return NewRouter(routerConfig, ormInstance, l)
}

func TestRouterInitializeFinalize_NoError(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	router := getTestRouter()
	// Act
	err := router.Initialize()
	err2 := router.Finalize()

	// Assert
	assert.NoError(err)
	assert.NoError(err2)
}

func TestRouterInitialize_IsInitialized_True(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	router := getTestRouter()
	router.Initialize()
	defer router.Finalize()

	// Assert
	assert.True(router.IsInitialized())
}

func TestRouterGetEngine_Initialized_Succeeds(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	router := getTestRouter()
	router.Initialize()
	// Act
	routerEngine := router.GetEngine()

	// Assert
	assert.NotNil(routerEngine)
}

func TestRouterGetEngine_NotInitialized_Panics(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	router := Router{}

	// Act
	act := func() { router.GetEngine() }

	// Assert
	assert.Panics(act)
}
