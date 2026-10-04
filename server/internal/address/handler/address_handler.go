package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/address/models"
	"github.com/yogesh4952/ebookstore/internal/address/service"
)

type AddressHandler struct {
	addressService service.IAddressService
}

func NewAddressHandler(serv service.IAddressService) *AddressHandler {
	return &AddressHandler{addressService: serv}
}

// AddAddress godoc
// @Summary      Add a new address
// @Description  Add a shipping address owned by the authenticated caller. The owner is always taken from the JWT, so any `user_id` in the request body is ignored and a caller cannot attach an address to another account. Both the city and the delivery address are required. The created record is returned under `data` so the client can use its id.
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        payload body models.UserAddress true "Address data"
// @Success      202 {object} map[string]interface{}  "Address added, containing success and message keys, and the created record under data"
// @Failure      400 {object} map[string]interface{}  "Invalid request body, the user does not exist, or the address could not be saved"
// @Failure      401 {object} map[string]interface{}  "Missing or invalid token"
// @Failure      500 {object} map[string]interface{}  "User id on the request context had an unexpected type"
// @Router       /address/add-address [post]
// @Security     BearerAuth
func (ah *AddressHandler) AddAddress(ctx *gin.Context) {
	var data models.UserAddress
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, ok := ctx.Get("userId")
	if !ok {
	}
	userIDVal, ok := userID.(uint)
	if !ok {
	}
	data.UserId = userIDVal

	err = ah.addressService.AddAddress(ctx, &data)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "User address added succesfully",
		"data":    data,
	})

}

// ListUserAddress godoc
// @Summary      List the authenticated user's addresses
// @Description  Return every shipping address belonging to the caller. Ownership is taken from the JWT and enforced in the query, so one user can never receive another user's addresses. A user with no saved addresses gets an empty array and a 200, not a 404.
// @Tags         addresses
// @Produce      json
// @Success      200 {object} map[string]interface{}  "Addresses fetched, containing the list under data"
// @Failure      401 {object} map[string]interface{}  "Missing or invalid token"
// @Failure      404 {object} map[string]interface{}  "Addresses could not be fetched"
// @Failure      500 {object} map[string]interface{}  "User id on the request context had an unexpected type"
// @Router       /address/ [get]
// @Security     BearerAuth
func (ah *AddressHandler) ListUserAddress(ctx *gin.Context) {
	userId, exist := ctx.Get("userId")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Userid doesnt exist",
		})
		return
	}
	userIdVal, ok := userId.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "invalid user ID",
		})
		return
	}

	data, err := ah.addressService.ListUserAddress(ctx, userIdVal)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User address fetches sucesfully",
		"data":    data,
	})

}
