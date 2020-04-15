package hello

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller represent the helloApi Controller
type Controller struct {
	config *Config
}

// Config represent the helloApi Configuration
type Config struct {
	Count int
}

// NewController creates a new Controller structure with the given parameters.
func NewController(config *Config) *Controller {
	return &Controller{
		config: config,
	}
}

// Routes list and binds all routes to the router group provided
func (ct *Controller) Routes(route *gin.RouterGroup) {
	route.GET("", ct.GetHello)
}

// GetHello handler hello world count
func (ct *Controller) GetHello(c *gin.Context) {
	var world []string
	for i := 0; i < ct.config.Count; i++ {
		world = append(world, fmt.Sprintf("world %d", (i+1)))
	}
	c.JSON(http.StatusOK, gin.H{
		"hello": world,
	})
}
