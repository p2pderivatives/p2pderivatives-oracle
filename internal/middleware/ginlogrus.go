package middleware

import (
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLogrus returns a Handler middleware which defines Logrus as router logger
func GinLogrus(logger *logrus.Logger) gin.HandlerFunc {
	useBanner := false
	useUTC := true
	return ginlogrus.WithTracing(logger,
		useBanner,
		time.RFC3339,
		useUTC,
		RequestIDHeaderTag,
		[]byte("trace-id"),
		[]byte(RequestIDHeaderTag),
		ginlogrus.WithAggregateLogging(true))
}
