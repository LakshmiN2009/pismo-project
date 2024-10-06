package middleware

import (
	"src/drivers"

	"github.com/gin-gonic/gin"
)

func DbMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := drivers.InitDB()
		c.Set("db", db)
		c.Next()
	}
}
