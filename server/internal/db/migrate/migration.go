package main

import (
	userAddressesModel "github.com/yogesh4952/ebookstore/internal/address/models"
	bookModels "github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/initializers"
	order_model "github.com/yogesh4952/ebookstore/internal/order/models"
	sellerModels "github.com/yogesh4952/ebookstore/internal/sellers/models"
	userModels "github.com/yogesh4952/ebookstore/internal/user/models"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

func init() {
	logger.Init()
	initializers.LoadEnv()
}

func main() {
	db, err := initializers.InitDb()
	if err != nil {
		logger.Fatal("Failed to initialize DB: %v", err)
	}

	err = db.AutoMigrate(&userModels.User{}, &sellerModels.Seller{}, &bookModels.Book{}, &userAddressesModel.UserAddress{}, &order_model.Order{}, &order_model.OrderItem{})
	if err != nil {
		logger.Fatal("Error during automigrations: %v", err)
	} else {
		logger.Success("Succesfully migrated")
	}
}
