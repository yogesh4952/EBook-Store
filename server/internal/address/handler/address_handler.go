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
// @Description  Add a shipping address for the authenticated user
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        payload body models.UserAddress true "Address data"
// @Success      202 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
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
	})

}
