package hello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestHelloRouter() *gin.Engine {
	return gin.Default()
}

func setupTestController() *Controller {

	return NewController(&Config{Count: 3})
}

func TestGetHello(t *testing.T) {

	router := setupTestHelloRouter()
	controller := setupTestController()
	router.GET("/hello", controller.GetHello)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var world []string
	for i := 0; i < controller.config.Count; i++ {
		world = append(world, fmt.Sprintf("world %d", (i+1)))
	}
	body := gin.H{
		"hello": world,
	}
	// Convert the JSON response to a map
	var response map[string][]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["hello"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["hello"], value)
}
