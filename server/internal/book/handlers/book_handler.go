package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/book/services"
	"github.com/yogesh4952/ebookstore/pkg/utils"
)

type BookHandler struct {
	service services.IBookService
}

func NewBookHandler(svc services.IBookService) *BookHandler {
	return &BookHandler{service: svc}
}

func (h *BookHandler) PublishBook(c *gin.Context) {

	var req models.BookPayload
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
