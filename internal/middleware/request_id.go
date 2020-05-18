package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeaderTag Tag used in request header to recover the request Id
var RequestIDHeaderTag = "Request-Id"

// RequestID returns a gin middleware function which assign (or recover)
// a unique uuid to the request header and to the gin context
func RequestID(contextID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(RequestIDHeaderTag)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(contextID, requestID)

		c.Writer.Header().Set(RequestIDHeaderTag, requestID)
		c.Next()
	}
}
