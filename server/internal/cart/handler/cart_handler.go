package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/cart/model"
	"github.com/yogesh4952/ebookstore/internal/cart/service"
)

type CartHandler struct {
	cartServ service.ICartService
}

func NewCartHandler(serv service.ICartService) *CartHandler {
	return &CartHandler{cartServ: serv}
}

// AddToCart godoc
// @Summary      Add a book to the cart
// @Description  Add a quantity of a book to the authenticated user's cart. Quantities are additive: adding the same book twice increments the stored count rather than replacing it. The user id is taken from the JWT, so a caller can only ever modify their own cart.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        payload body model.AddToCartPayload true "Book and quantity to add"
// @Success      201 {object} map[string]interface{}  "Book added to cart, containing success and message keys"
// @Failure      400 {object} map[string]interface{}  "Invalid request body"
// @Failure      401 {object} map[string]interface{}  "Missing or invalid token"
// @Failure      404 {object} map[string]interface{}  "Book not found"
// @Failure      500 {object} map[string]interface{}  "User id on the request context had an unexpected type"
// @Router       /cart/add-to-cart [post]
// @Security     BearerAuth
func (ch *CartHandler) AddToCart(c *gin.Context) {
	var addToCartPayload model.AddToCartPayload
	if err := c.ShouldBindJSON(&addToCartPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userIdValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
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

	_, err := ch.cartServ.AddToCart(c.Request.Context(), userID, addToCartPayload)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Item added to cart successfully",
		"data": gin.H{
			"book_id":  addToCartPayload.BookId,
			"quantity": addToCartPayload.Quantity,
		},
	})
}

func (ch *CartHandler) CartItems(c *gin.Context) {
	userIdValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user Id",
		})
		return
	}

	userId := userIdValue.(uint)
	data, err := ch.cartServ.GetCartItem(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": true,
		"message": "Cart Item Fetched sucessfully",
		"data":    data,
	})

}
