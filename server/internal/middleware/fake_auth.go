package middleware

import "github.com/gin-gonic/gin"

func FakeAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("email", "yogeshshah2063@gmail.com")
		c.Set("userId", "2")
		c.Set("role", "user")
		c.Next()
	}

}
