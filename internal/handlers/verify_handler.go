package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifyHandler обрабатывает запрос на подтверждение email.
// GET /verify?email=...&code=...
func (h *UserHandler) VerifyHandler(c *gin.Context) {
	email := c.Query("email")
	code := c.Query("code")

	if email == "" || code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and code are required"})
		return
	}

	err := h.userService.VerifyEmail(c.Request.Context(), email, code)
	if err != nil {
		// Обработка ошибок (пользователь не найден, код не совпадает и т.д.)
		log.Printf("Error verifying email: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}