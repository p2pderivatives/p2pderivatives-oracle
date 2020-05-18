package middleware

import (
	"github.com/gin-gonic/gin"
)

// AddToContext returns a Handler middleware which add an object to the request context
func AddToContext(id string, obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(id, obj)
		c.Next()
	}
}
