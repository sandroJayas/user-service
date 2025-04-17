package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireEmployeeRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("account_type")
		if !exists || userType != "employee" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}
		c.Next()
	}
}
