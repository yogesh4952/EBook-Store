package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("role")

		if exists == false {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized. User identity missing",
				"success": "false",
			})
			c.Abort()
			return
		}

		userRole, ok := roleInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: invalid role format"})
			c.Abort()
			return
		}

		roleAllowed := false

		for _, role := range allowedRoles {
			if userRole == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden role",
				"success": false,
			})
			c.Abort()
			return
		}

		c.Next()

	}
}
