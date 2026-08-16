package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/initializers"
)

func main() {
	// 1. Load environment variables first
	initializers.LoadEnv()
	initializers.InitDb()

	// 2. Read PORT after LoadEnv() has executed
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Sensible fallback default
	}

	// 3. Initialize Gin router
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 4. Run server on dynamic port
	router.Run(":" + port)
}
