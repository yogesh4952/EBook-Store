package userhandler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/user/models"
	"github.com/yogesh4952/ebookstore/internal/user/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ListUser godoc
// @Summary      List all users
// @Description  Get all registered users. Responds with a bare JSON array of user objects; there is no response envelope on this endpoint. Note that a healthy but empty database also returns 500, because the service reports an empty result set as an error.
// @Tags         users
// @Produce      json
// @Success      200 {array} github_com_yogesh4952_ebookstore_internal_user_models.User  "List of registered users"
// @Failure      500 {object} map[string]interface{}  "Database error, or no users registered yet"
// @Router       /users [get]
// @Security     BearerAuth
func (uh *UserHandler) ListUser(c *gin.Context) {

	users, err := uh.service.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) EditUserData(c *gin.Context) {
	var req models.EditUserPayload

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sucess":  false,
			"message": fmt.Errorf("%v: error request", err.Error()),
		})
		return
	}

	userId, exist := c.Get("userId")
	log.Printf("%v: USERID", userId)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sucess":  false,
			"message": "UserId doesn't exist!",
		})
		return
	}
	req.ID = userId.(uint)

	err := uh.service.EditUserData(c, req)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"sucess":  false,
			"message": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"sucess":  false,
		"message": "User data edited sucessfully",
	})

}
