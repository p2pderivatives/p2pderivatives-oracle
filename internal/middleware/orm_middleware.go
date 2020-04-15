package middleware

import (
	"p2pderivatives-oracle/internal/database/orm"

	"github.com/gin-gonic/gin"
)

// Orm returns a Handler middleware which add the orm into the router context
func Orm(o *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("orm", o)
		c.Next()
	}
}
