package main

import (
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/yogesh4952/ebookstore/docs"
	authhandler "github.com/yogesh4952/ebookstore/internal/auth/handler"
	authrepo "github.com/yogesh4952/ebookstore/internal/auth/repository"
	authservice "github.com/yogesh4952/ebookstore/internal/auth/service"
	"github.com/yogesh4952/ebookstore/internal/initializers"
	"github.com/yogesh4952/ebookstore/internal/middleware"
	"github.com/yogesh4952/ebookstore/pkg/logger"
	userhandler "github.com/yogesh4952/ebookstore/internal/user/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	userrepo "github.com/yogesh4952/ebookstore/internal/user/repository"
	userservice "github.com/yogesh4952/ebookstore/internal/user/service"
	"github.com/yogesh4952/ebookstore/pkg/utils"

	bookhandler "github.com/yogesh4952/ebookstore/internal/book/handler"
	bookrepo "github.com/yogesh4952/ebookstore/internal/book/repository"
	bookservice "github.com/yogesh4952/ebookstore/internal/book/service"
	sellerrepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
)

// @title           EBook Store API
// @version         1.0
// @description     REST API for the EBook Store platform
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func init() {

	logger.Init()
	initializers.LoadEnv()
	initializers.InitDb()
	initializers.InitRedis()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ZerologMiddleware())

	//userDependency injection
	userRepo := userrepo.NewUserRepository(initializers.DB)
	userService := userservice.NewUserService(userRepo)
	userHandler := userhandler.NewUserHandler(userService)

	sellerRepo := sellerrepo.NewSellerRepository(initializers.DB)
	bookRepo := bookrepo.NewBookRepo(initializers.DB)
	bookService := bookservice.NewBookService(bookRepo, sellerRepo)
	bookHandler := bookhandler.NewBookHandler(bookService)

	jwtManager := utils.NewJwtManager()
	emailService := utils.NewEmailService()
	authRepo := authrepo.NewAuthRepository(initializers.DB, initializers.RDB)
	authService := authservice.NewAuthService(authRepo, userRepo, jwtManager, emailService)
	authHandler := authhandler.NewAuthHandler(authService)

	apiRoutes := router.Group("/api")
	{

		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.GET("", userHandler.ListUser)
		}

		authRoutes := apiRoutes.Group("/auth")
		{
			authRoutes.POST("/send-otp", authHandler.SendOTP)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/register", authHandler.Register)

		}

		bookRoutes := apiRoutes.Group("/book")
		{
			bookRoutes.POST("/publish-book", middleware.AuthRequired(), middleware.AuthorizeRoles("seller", "admin"), bookHandler.PublishBook)
			bookRoutes.PATCH("/update-book", middleware.AuthRequired(), middleware.AuthorizeRoles("seller", "admin"), bookHandler.UpdateBook)
			bookRoutes.GET("/list-books", bookHandler.ListBooks)
		}
		apiRoutes.GET("/", middleware.AuthRequired())
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + port)
}
