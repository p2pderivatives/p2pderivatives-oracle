package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"p2pderivatives-oracle/internal/router"
	"p2pderivatives-oracle/test"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func NewTestRouter(api router.API) *router.Router {
	gin.SetMode(gin.TestMode)
	l := test.NewLogger()
	return router.NewRouter(l, api)
}

func TestRouterInitializeFinalize_NoError(t *testing.T) {
	router := NewTestRouter(NewMockAPI())
	err := router.Initialize()
	err2 := router.Finalize()
	assert.NoError(t, err)
	assert.NoError(t, err2)
}

func TestRouterInitialize_IsInitialized_True(t *testing.T) {
	router := NewTestRouter(NewMockAPI())
	router.Initialize()
	defer router.Finalize()

	assert.True(t, router.IsInitialized())
}

func TestRouterGetEngine_Initialized_Succeeds(t *testing.T) {
	router := NewTestRouter(NewMockAPI())
	router.Initialize()
	routerEngine := router.GetEngine()

	assert.NotNil(t, routerEngine)
}

func TestRouterGetEngine_NotInitialized_Panics(t *testing.T) {
	router := NewTestRouter(NewMockAPI())

	// Act
	act := func() { router.GetEngine() }

	// Assert
	assert.Panics(t, act)
}

func TestRouterInitialized_API_Succeeds(t *testing.T) {
	router := NewTestRouter(NewMockAPI())
	err := router.Initialize()
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", GetRoute, nil)
	router.GetEngine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := &testAPIResponse{}
	err = json.Unmarshal([]byte(w.Body.String()), response)
	assert.Nil(t, err)
	assert.EqualValues(t, &testAPIResponse{
		HasMiddleware: true,
		HasService:    true,
	}, response)
}
