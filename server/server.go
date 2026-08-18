package main

import (
	"os"

	"github.com/gin-gonic/gin"

	authhandler "github.com/yogesh4952/ebookstore/auth/handler"
	authrepo "github.com/yogesh4952/ebookstore/auth/repository"
	authservice "github.com/yogesh4952/ebookstore/auth/service"
	"github.com/yogesh4952/ebookstore/initializers"
	userhandler "github.com/yogesh4952/ebookstore/user/handler"
	userrepo "github.com/yogesh4952/ebookstore/user/repository"
	userservice "github.com/yogesh4952/ebookstore/user/service"
)

func init() {

	initializers.LoadEnv()
	initializers.InitDb()
	initializers.InitRedis()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	//userDependency injection
	userRepo := userrepo.NewUserRepository(initializers.DB)
	userService := userservice.NewUserService(userRepo)
	userHandler := userhandler.NewUserHandler(userService)

	authRepo := authrepo.NewAuthRepository(initializers.DB, initializers.RDB)
	authService := authservice.NewAuthService(authRepo, userRepo)
	authHandler := authhandler.NewAuthHandler(authService)

	apiRoutes := router.Group("/api")
	{
		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.GET("", userHandler.ListUser)
		}

		auhtRoutes := apiRoutes.Group("/auth")
		{
			auhtRoutes.POST("/send-otp", authHandler.SendOTP)
			auhtRoutes.POST("/login", authHandler.Login)

		}
	}

	router.Run(":" + port)
}
