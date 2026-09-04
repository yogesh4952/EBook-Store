package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthRequired_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.Use(AuthRequired())

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource",
		})
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer invalid-token",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			recorder.Code,
		)
	}
}
