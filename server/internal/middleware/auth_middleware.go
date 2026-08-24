package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/pkg/utils"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing token"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJwt(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid or expired token"})
			c.Abort()
			return
		}

		// store the identity so handlers can read it
		c.Set("userId", claims.UserId)
		c.Set("email", claims.Email)
		c.Next()
	}
}
