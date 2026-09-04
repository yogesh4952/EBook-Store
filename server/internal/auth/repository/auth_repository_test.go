package repository

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	bookModel "github.com/yogesh4952/ebookstore/internal/book/models"
	sellerModel "github.com/yogesh4952/ebookstore/internal/sellers/models"
	usermodels "github.com/yogesh4952/ebookstore/internal/user/models"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Use the pure-Go in-memory SQLite driver so tests run without a Postgres server.
	// A unique name per test keeps each DB isolated from the others.
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(
		&usermodels.User{},
		&sellerModel.Seller{},
		&bookModel.Book{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func setupRedisMock() *redis.Client {
	// The RegisterUser method under test never touches redis,
	// so a non-connected client is enough to construct the repo.
	return redis.NewClient(&redis.Options{})
}

func TestRegisterUser_CreatesSellerProfileForSellerRole(t *testing.T) {
	tests := []struct {
		name          string
		role          usermodels.Role
		wantSellerRow bool
	}{
		{
			name:          "seller role auto-creates a seller row",
			role:          usermodels.Roleseller,
			wantSellerRow: true,
		},
		{
			name:          "customer role does not create a seller row",
			role:          usermodels.RoleCustomer,
			wantSellerRow: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewAuthRepository(db, setupRedisMock())

			user := &usermodels.User{
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "john@example.com",
				Role:      tc.role,
			}

			err := repo.RegisterUser(context.Background(), user)
			if err != nil {
				t.Fatalf("RegisterUser returned error: %v", err)
			}

			// User must have been persisted with a generated ID.
			if user.ID == 0 {
				t.Fatal("expected user to be persisted with a non-zero ID")
			}

			var sellers []sellerModel.Seller
			if err := db.Find(&sellers).Error; err != nil {
				t.Fatalf("failed to query sellers: %v", err)
			}

			if tc.wantSellerRow {
				if len(sellers) != 1 {
					t.Fatalf("expected exactly 1 seller row, got %d", len(sellers))
				}
				if sellers[0].UserID != user.ID {
					t.Errorf("expected seller.UserID=%d to match user.ID=%d", sellers[0].UserID, user.ID)
				}
			} else {
				if len(sellers) != 0 {
					t.Errorf("expected 0 seller rows for role %q, got %d", tc.role, len(sellers))
				}
			}
		})
	}
}

func TestRegisterUser_RollbacksSellerOnTransactionError(t *testing.T) {
	// A unique email makes the duplicate-key error path roll back the
	// whole transaction, so no orphaned seller rows should remain.
	db := setupTestDB(t)
	repo := NewAuthRepository(db, setupRedisMock())

	seller := &usermodels.User{
		Firstname: "Jane",
		Lastname:  "Doe",
		Email:     "jane@example.com",
		Role:      usermodels.Roleseller,
	}
	if err := repo.RegisterUser(context.Background(), seller); err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	duplicate := &usermodels.User{
		Firstname: "Jane",
		Lastname:  "Doe",
		Email:     "jane@example.com", // same email -> duplicate key
		Role:      usermodels.Roleseller,
	}
	if err := repo.RegisterUser(context.Background(), duplicate); err == nil {
		t.Fatal("expected duplicate email to fail")
	}

	var sellers []sellerModel.Seller
	if err := db.Find(&sellers).Error; err != nil {
		t.Fatalf("failed to query sellers: %v", err)
	}
	if len(sellers) != 1 {
		t.Errorf("expected 1 seller row after rollback, got %d", len(sellers))
	}
}
