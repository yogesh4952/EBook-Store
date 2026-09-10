package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	bookRepo "github.com/yogesh4952/ebookstore/internal/book/repository"
	"github.com/yogesh4952/ebookstore/internal/book/service"
	"github.com/yogesh4952/ebookstore/internal/initializers"
	sellersModels "github.com/yogesh4952/ebookstore/internal/sellers/models"
	sellerRepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
	userModels "github.com/yogesh4952/ebookstore/internal/user/models"
	userRepo "github.com/yogesh4952/ebookstore/internal/user/repository"
	"github.com/yogesh4952/ebookstore/pkg/logger"
	"gorm.io/gorm"
)

func main() {
	logger.Init()
	initializers.LoadEnv()

	db, err := initializers.InitDb()
	if err != nil {
		logger.Fatal("Failed to initialize DB: %v", err)
	}

	bookFile, err := os.Open("/home/yst/code/EBook-Store/server/data/book.json")

	if err != nil {
		logger.Error("%v", err)
		return
	}

	defer bookFile.Close()
	var books []*models.PublishBookPayload

	err = json.NewDecoder(bookFile).Decode(&books)

	if err != nil {
		logger.Error("%v", err)
		return
	}

	logger.Info("Loaded %d books", len(books))

	sellerID, err := ensureDemoSeller(db)
	if err != nil {
		logger.Fatal("Failed to ensure demo seller: %v", err)
	}

	if err := db.Where("seller_id = ? AND NOT EXISTS (SELECT 1 FROM order_items oi WHERE oi.book_id = books.id)", sellerID).Delete(&models.Book{}).Error; err != nil {
		logger.Fatal("Failed to clear existing seeded books: %v", err)
	}
	for _, b := range books {
		b.SellerId = &sellerID
	}

	bRepo := bookRepo.NewBookRepo(db)
	sRepo := sellerRepo.NewSellerRepository(db)
	svc := service.NewBookService(bRepo, sRepo)

	err = svc.BatchBookSeed(books)
	if err != nil {
		logger.Fatal("Failed to seed books: %v", err)
	}
}

// ensureDemoSeller returns the ID of a seller to attach seeded books to. It
// reuses the first existing seller so re-seeding never collides on the unique
// seller_number constraint, and only creates a demo user + seller on a fresh
// database where no seller exists yet.
func ensureDemoSeller(db *gorm.DB) (uint, error) {
	ctx := context.Background()

	var seller sellersModels.Seller
	err := db.WithContext(ctx).Order("id ASC").First(&seller).Error
	if err == nil {
		return seller.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	uRepo := userRepo.NewUserRepository(db)
	user, err := uRepo.FindByEmail(ctx, "seller@demo.com")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &userModels.User{
			Firstname:   "Demo",
			Lastname:    "Seller",
			Email:       "seller@demo.com",
			Role:        userModels.Roleseller,
			PhoneNumber: "9800000000",
		}
		if err := db.WithContext(ctx).Create(user).Error; err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	seller = sellersModels.Seller{UserID: user.ID, SellerNumber: 1}
	if err := db.WithContext(ctx).Create(&seller).Error; err != nil {
		return 0, err
	}

	return seller.ID, nil
}
