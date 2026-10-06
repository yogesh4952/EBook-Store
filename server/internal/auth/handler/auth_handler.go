package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/auth"
	authmodels "github.com/yogesh4952/ebookstore/internal/auth/models"
	"github.com/yogesh4952/ebookstore/internal/auth/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginPayload struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// SendOTP godoc
// @Summary      Send OTP to email
// @Description  Send a one-time password to the user's email for verification. The code is valid for five minutes and can only be used once. Note that this endpoint uses two different response envelopes: a 400 returns `{"error": "..."}` while a 500 returns `{"success": false, "message": "..."}`.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body SendOTPRequest true "Email address"
// @Success      200 {object} map[string]interface{}  "OTP sent, containing success and message keys"
// @Failure      400 {object} map[string]interface{}  "Invalid request body, containing an error key"
// @Failure      500 {object} map[string]interface{}  "Unknown email, or the code could not be generated, stored, or emailed"
// @Router       /auth/send-otp [post]
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SendOTP(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP verification code sent to your email",
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isVerified, err := h.service.VerifyOTP(c.Request.Context(), req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"verified": isVerified,
	})
}

// Login godoc
// @Summary      Login with email and OTP
// @Description  Authenticate a user with the email and the one-time code sent by `/auth/send-otp`, and return a JWT access token. The OTP is consumed on success and cannot be reused. Note that the success response contains `message` and `token` keys only, with no `success` key, whereas the error responses use an `error` key.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body LoginPayload true "Email and OTP"
// @Success      200 {object} map[string]interface{}  "Login successful, containing message and token keys"
// @Failure      400 {object} map[string]interface{}  "Invalid request body, containing an error key"
// @Failure      401 {object} map[string]interface{}  "Incorrect or expired OTP, containing an error key"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account as either a customer or a seller. Registering with an email that already exists returns 409. Use a role of `admin` only for internal tooling; normal signups should use `customer` or `seller`.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body authmodels.RegisterPayload true "Registration data"
// @Success      200 {object} map[string]interface{}  "Account created, containing success and message keys"
// @Failure      400 {object} map[string]interface{}  "Invalid request body, or an invalid role"
// @Failure      409 {object} map[string]interface{}  "Email is already registered"
// @Failure      500 {object} map[string]interface{}  "Account could not be created"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req authmodels.RegisterPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := h.service.Register(c.Request.Context(), &req)

	if err != nil {
		if errors.Is(err, auth.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": fmt.Errorf("Duplicate email: %v", err.Error()),
			})
			return
		}

		if errors.Is(err, auth.ErrInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Errorf("Internal server error: %v", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": result,
	})

}

func (h *AuthHandler) Google(c *gin.Context) {
	var code string
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("CODE:%s", code)

}
