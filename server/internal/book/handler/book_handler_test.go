package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/pkg/utils"
)

type mockBookService struct {
	publishBookFn func(ctx context.Context, userId uint, data *models.PublishBookPayload) error
	updateBookFn  func(ctx context.Context, userId uint, data *models.UpdateBookPayload) (models.Book, error)
	listBooksFn   func(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error)
}

func (m *mockBookService) PublishBook(ctx context.Context, userId uint, data *models.PublishBookPayload) error {
	return m.publishBookFn(ctx, userId, data)
}

func (m *mockBookService) UpdateBook(ctx context.Context, userId uint, data *models.UpdateBookPayload) (models.Book, error) {
	return m.updateBookFn(ctx, userId, data)
}

func (m *mockBookService) ListBooks(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error) {
	return m.listBooksFn(ctx, p)
}

func (m *mockBookService) BatchBookSeed(data []*models.PublishBookPayload) error {
	return nil
}

func setupRouter(handler *BookHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userId", uint(1))
		c.Next()
	})
	return r
}

func TestPublishBook_Success(t *testing.T) {
	svc := &mockBookService{
		publishBookFn: func(ctx context.Context, userId uint, data *models.PublishBookPayload) error {
			if userId != 1 {
				t.Errorf("expected userId 1, got %d", userId)
			}
			return nil
		},
	}

	h := NewBookHandler(svc)
	r := setupRouter(h)
	r.POST("/books", h.PublishBook)

	body := models.PublishBookPayload{
		Title:       "Test Book",
		AuthorName:  "Author",
		Genre:       "Fiction",
		Category:    "Novel",
		Pages:       100,
		Publication: "Publisher",
		Price:       9.99,
		Units:       5,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Errorf("expected success true, got %v", resp["success"])
	}
}

func TestPublishBook_Unauthorized(t *testing.T) {
	svc := &mockBookService{}
	h := NewBookHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/books", h.PublishBook)

	body := models.PublishBookPayload{
		Title:       "Test",
		AuthorName:  "A",
		Genre:       "F",
		Category:    "C",
		Pages:       1,
		Publication: "P",
		Price:       1,
		Units:       1,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestPublishBook_BadRequest(t *testing.T) {
	svc := &mockBookService{
		publishBookFn: func(ctx context.Context, userId uint, data *models.PublishBookPayload) error {
			return nil
		},
	}

	h := NewBookHandler(svc)
	r := setupRouter(h)
	r.POST("/books", h.PublishBook)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPublishBook_ServiceError(t *testing.T) {
	svc := &mockBookService{
		publishBookFn: func(ctx context.Context, userId uint, data *models.PublishBookPayload) error {
			return &testError{"service failure"}
		},
	}

	h := NewBookHandler(svc)
	r := setupRouter(h)
	r.POST("/books", h.PublishBook)

	body := models.PublishBookPayload{
		Title:       "Test",
		AuthorName:  "A",
		Genre:       "F",
		Category:    "C",
		Pages:       1,
		Publication: "P",
		Price:       1,
		Units:       1,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("expected status %d, got %d", http.StatusBadGateway, rec.Code)
	}
}

func TestUpdateBook_Success(t *testing.T) {
	title := "Updated Title"
	author := "Updated Author"
	genre := "Genre"
	category := "Cat"
	pages := uint(200)
	pub := "Pub"
	price := float32(19.99)
	units := 10
	cover := "http://example.com/cover.jpg"

	svc := &mockBookService{
		updateBookFn: func(ctx context.Context, userId uint, data *models.UpdateBookPayload) (models.Book, error) {
			return models.Book{
				Title:       title,
				AuthorName:  author,
				Genre:       genre,
				Category:    category,
				Pages:       pages,
				Publication: pub,
				Price:       price,
				Units:       units,
			}, nil
		},
	}

	h := NewBookHandler(svc)
	r := setupRouter(h)
	r.PUT("/books", h.UpdateBook)

	body := models.UpdateBookPayload{
		Title:        &title,
		AuthorName:   &author,
		Genre:        &genre,
		Category:     &category,
		Pages:        &pages,
		Publication:  &pub,
		Price:        &price,
		Units:        &units,
		CoverPageUrl: &cover,
		BookId:       1,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Errorf("expected success true, got %v", resp["success"])
	}
}

func TestUpdateBook_MissingUserId(t *testing.T) {
	svc := &mockBookService{}
	h := NewBookHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/books", h.UpdateBook)

	body := models.UpdateBookPayload{BookId: 1}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestListBooks_Success(t *testing.T) {
	books := []models.Book{
		{Title: "Book 1"},
		{Title: "Book 2"},
	}

	svc := &mockBookService{
		listBooksFn: func(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error) {
			return books, 2, nil
		},
	}

	h := NewBookHandler(svc)
	r := setupRouter(h)
	r.GET("/books", h.ListBooks)

	req := httptest.NewRequest(http.MethodGet, "/books?page=1&limit=10", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Errorf("expected success true, got %v", resp["success"])
	}
}

func TestListBooks_ServiceError(t *testing.T) {
	svc := &mockBookService{
		listBooksFn: func(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error) {
			return nil, 0, &testError{"not found"}
		},
	}

	h := NewBookHandler(svc)
	r := setupRouter(h)
	r.GET("/books", h.ListBooks)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestListBooks_DefaultPagination(t *testing.T) {
	svc := &mockBookService{
		listBooksFn: func(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error) {
			if p.GetPage() != 1 {
				t.Errorf("expected default page 1, got %d", p.GetPage())
			}
			if p.Getlimit() != 10 {
				t.Errorf("expected default limit 10, got %d", p.Getlimit())
			}
			return nil, 0, nil
		},
	}

	h := NewBookHandler(svc)
	r := setupRouter(h)
	r.GET("/books", h.ListBooks)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// Ensure mockBookService satisfies IBookService at compile time.
var _ interface {
	PublishBook(ctx context.Context, userId uint, data *models.PublishBookPayload) error
	UpdateBook(ctx context.Context, userId uint, data *models.UpdateBookPayload) (models.Book, error)
	ListBooks(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error)
} = (*mockBookService)(nil)

// Helper to convert int to string for query params.
func itoa(i int) string {
	return strconv.Itoa(i)
}
