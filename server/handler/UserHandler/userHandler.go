package userhandler

import "github.com/gin-gonic/gin"

func ListUser(c *gin.Context) {
	// req herne
	// redirect to servuce+
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
