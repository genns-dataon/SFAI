package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Chat(c *gin.Context) {
	var input struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Feature coming soon. This chatbot will be powered by AI to answer your HR queries.",
		"message":  input.Message,
	})
}
