package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler handler function that will handle the api errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		// no need to log if ginlogrus is used

		errorResponse, ok := err.Err.(*Error)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
		}
	}
}
