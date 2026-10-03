package userhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
