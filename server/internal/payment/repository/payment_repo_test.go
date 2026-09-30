package repository

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	bookModel "github.com/yogesh4952/ebookstore/internal/book/models"
	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(
		&orderModel.Order{},
		&orderModel.OrderItem{},
		&bookModel.Book{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func seedOrder(t *testing.T, db *gorm.DB, status orderModel.PaymentStatus) *orderModel.Order {
	t.Helper()

	order := &orderModel.Order{
		OrderCode:       "ORD-TEST",
		TransactionUUID: "ORD-TEST",
		PaymentMethod:   orderModel.PayementEsewa,
		PaymentStatus:   status,
		OrderStatus:     orderModel.OrderPlaced,
		TotalPrice:      1750,
		ShippingCity:    "Kathmandu",
	}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("failed to seed order: %v", err)
	}
	return order
}

func TestConfirmPayment_PersistsStatusAndTransactionCode(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepo(db)

	book := &bookModel.Book{Title: "Go", Price: 1750, Units: 10, Status: bookModel.StatusAvailable}
	if err := db.Create(book).Error; err != nil {
		t.Fatalf("failed to seed book: %v", err)
	}

	order := seedOrder(t, db, orderModel.PayementPending)
	item := &orderModel.OrderItem{OrderId: order.ID, BookId: book.ID, Quantity: 2, UnitPrice: 875}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("failed to seed order item: %v", err)
	}

	if err := repo.ConfirmPayment(context.Background(), order.ID, "000H8NP"); err != nil {
		t.Fatalf("ConfirmPayment: %v", err)
	}

	var got orderModel.Order
	if err := db.First(&got, order.ID).Error; err != nil {
		t.Fatalf("failed to reload order: %v", err)
	}
	if got.PaymentStatus != orderModel.PayementPaid {
		t.Errorf("payment status = %q, want %q", got.PaymentStatus, orderModel.PayementPaid)
	}
	if got.TransactionCode != "000H8NP" {
		t.Errorf("transaction code = %q, want %q", got.TransactionCode, "000H8NP")
	}

	var gotBook bookModel.Book
	if err := db.First(&gotBook, book.ID).Error; err != nil {
		t.Fatalf("failed to reload book: %v", err)
	}
	if gotBook.Units != 8 {
		t.Errorf("units = %d, want 8", gotBook.Units)
	}
}

// A replayed redirect must not decrement stock a second time.
func TestConfirmPayment_RefusesToReconfirmPaidOrder(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepo(db)

	book := &bookModel.Book{Title: "Go", Price: 1750, Units: 10, Status: bookModel.StatusAvailable}
	if err := db.Create(book).Error; err != nil {
		t.Fatalf("failed to seed book: %v", err)
	}

	order := seedOrder(t, db, orderModel.PayementPending)
	item := &orderModel.OrderItem{OrderId: order.ID, BookId: book.ID, Quantity: 2, UnitPrice: 875}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("failed to seed order item: %v", err)
	}

	if err := repo.ConfirmPayment(context.Background(), order.ID, "000H8NP"); err != nil {
		t.Fatalf("first ConfirmPayment: %v", err)
	}
	if err := repo.ConfirmPayment(context.Background(), order.ID, "000H8NP"); err == nil {
		t.Fatal("expected the second ConfirmPayment to be refused")
	}

	var gotBook bookModel.Book
	if err := db.First(&gotBook, book.ID).Error; err != nil {
		t.Fatalf("failed to reload book: %v", err)
	}
	if gotBook.Units != 8 {
		t.Errorf("units = %d, want 8 (stock must not be decremented twice)", gotBook.Units)
	}
}

func TestFindOrderByTransactionUUID_Missing(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepo(db)

	_, err := repo.FindOrderByTransactionUUID(context.Background(), "ORD-NOPE")
	if err == nil {
		t.Fatal("expected an error for an unknown transaction uuid")
	}
}
