package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/book/service"
	"github.com/yogesh4952/ebookstore/pkg/utils"
)

type BookHandler struct {
	service service.IBookService
}

func NewBookHandler(svc service.IBookService) *BookHandler {
	return &BookHandler{service: svc}
}

// PublishBook godoc
// @Summary      Publish a new book
// @Description  Seller creates a new book listing
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        payload body models.PublishBookPayload true "Book data"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      502 {object} map[string]interface{}
// @Router       /book/publish-book [post]
// @Security     BearerAuth
func (h *BookHandler) PublishBook(c *gin.Context) {

	var req models.PublishBookPayload
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "invalid user ID",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err := h.service.PublishBook(c.Request.Context(), userID, &req)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Book Published Succesfully",
	})

}

// UpdateBook godoc
// @Summary      Update an existing book
// @Description  Update book details (owner only)
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        payload body models.UpdateBookPayload true "Updated book data"
// @Success      202 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /book/update-book [patch]
// @Security     BearerAuth
func (h *BookHandler) UpdateBook(c *gin.Context) {
	var req models.UpdateBookPayload

	userIdValue, exist := c.Get("userId")

	if exist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Missing userid",
		})
		return
	}

	userID, ok := userIdValue.(uint)
	log.Printf("UserId: %v", userID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "invalid user ID",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	data, err := h.service.UpdateBook(c.Request.Context(), userID, &req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Book data updated succesfully",
		"data":    data,
	})

}

// ListBooks godoc
// @Summary      List all books
// @Description  Get a paginated list of books
// @Tags         books
// @Produce      json
// @Param        page  query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /book/list-books [get]
func (h *BookHandler) ListBooks(c *gin.Context) {

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	p := utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	books, total, err := h.service.ListBooks(c.Request.Context(), p)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Book data fetches successfully",
		"data":    books,
		"total":   total,
		"page":    p.GetPage(),
		"limit":   p.Getlimit(),
	})
}
