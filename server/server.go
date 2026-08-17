package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/yogesh4952/ebookstore/initializers"
	userhandler "github.com/yogesh4952/ebookstore/user/handler"
	userrepo "github.com/yogesh4952/ebookstore/user/repository"
	"github.com/yogesh4952/ebookstore/user/service"
)

func init() {

	initializers.LoadEnv()
	initializers.InitDb()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	//userDependency injection
	userRepo := userrepo.NewUserRepository(initializers.DB)
	userService := service.NewUserService(userRepo)
	userHandler := userhandler.NewUserHandler(userService)

	apiRoutes := router.Group("/api")
	{
		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.GET("", userHandler.ListUser)
		}
	}

	// 4. Run server on dynamic port
	router.Run(":" + port)
}
