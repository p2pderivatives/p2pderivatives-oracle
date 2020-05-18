package router_test

import (
	"net/http"
	"p2pderivatives-oracle/internal/router"

	"github.com/gin-gonic/gin"
)

const (
	ServiceValue           = "some value"
	ObjMiddleware          = "my middleware Test value"
	ObjMiddlewareContextID = "test-context-obj"
	GetRoute               = "/test"
)

func NewMockAPI() router.API {
	return &MockAPI{mockService: &MockService{
		initialized: false,
		finalized:   false,
		value:       ServiceValue,
	}}
}

type MockService struct {
	initialized bool
	finalized   bool
	value       string
}

type testAPIResponse struct {
	HasMiddleware bool
	HasService    bool
}

type MockAPI struct {
	mockService *MockService
}

func (m MockAPI) Routes(route *gin.RouterGroup) {
	route.GET(GetRoute, func(c *gin.Context) {
		testObj := c.GetString(ObjMiddlewareContextID)

		c.JSON(http.StatusOK, &testAPIResponse{
			HasMiddleware: testObj == ObjMiddleware,
			HasService:    ServiceValue == m.mockService.value,
		})
	})
}

func (m MockAPI) GlobalMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			c.Set(ObjMiddlewareContextID, ObjMiddleware)
			c.Next()
		},
	}
}

func (m MockAPI) InitializeServices() error {
	m.mockService.initialized = true
	return nil
}

func (m MockAPI) AreServicesInitialized() bool {
	return m.mockService.initialized
}

func (m MockAPI) FinalizeServices() error {
	m.mockService.finalized = true
	return nil
}
