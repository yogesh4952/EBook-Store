package handler

import (
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
// @Description  Create a new book listing. Restricted to the `seller` and `admin` roles; any other authenticated user receives 403. The authenticated user id is taken from the JWT and used as the book owner, so it cannot be spoofed in the request body.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        payload body models.PublishBookPayload true "Book data"
// @Success      201 {object} map[string]interface{}  "Book published, containing success and message keys"
// @Failure      400 {object} map[string]interface{}  "Invalid request body, or the book could not be saved"
// @Failure      401 {object} map[string]interface{}  "Missing or invalid token, or no user identity on the request context"
// @Failure      403 {object} map[string]interface{}  "Authenticated user is not a seller or admin"
// @Failure      500 {object} map[string]interface{}  "User id on the request context had an unexpected type"
// @Failure      502 {object} map[string]interface{}  "The book service failed to process the request"
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
// @Summary      Update an existing book
// @Description  Update the details of a book the authenticated user owns. Only the fields present in the request body are modified. Restricted to the `seller` and `admin` roles; any other authenticated user receives 403.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        payload body models.UpdateBookPayload true "Updated book data"
// @Success      202 {object} map[string]interface{}  "Book updated, containing success, message, and the updated book under data"
// @Failure      400 {object} map[string]interface{}  "Missing user identity, invalid request body, or the update failed"
// @Failure      401 {object} map[string]interface{}  "Missing or invalid token"
// @Failure      403 {object} map[string]interface{}  "Authenticated user is not a seller or admin"
// @Failure      500 {object} map[string]interface{}  "User id on the request context had an unexpected type"
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
// @Description  Get a paginated list of books. The response echoes the `page` and `limit` actually applied by the server, which may differ from the requested values when out of range. Any service or database error is reported as 404.
// @Tags         books
// @Produce      json
// @Param        page  query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Success      200 {object} map[string]interface{}  "Books fetched, containing success, message, the books under data, and the total, page, and limit counts"
// @Failure      404 {object} map[string]interface{}  "Books could not be fetched"
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

func (h *BookHandler) Genre(c *gin.Context) {

	genre, err := h.service.ListGenre(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Book Genre  fetches successfully",
		"data":    genre,
	})
}

func (h *BookHandler) Category(c *gin.Context) {

	category, err := h.service.ListCategory(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Book Category  fetches successfully",
		"data":    category,
	})
}
