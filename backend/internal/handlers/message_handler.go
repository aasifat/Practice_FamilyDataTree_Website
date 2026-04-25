package handlers

import (
	"net/http"

	"family-tree-api/internal/models"
	"family-tree-api/internal/repository"

	"github.com/gin-gonic/gin"
)

func CreateMessage(c *gin.Context) {
	var req models.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &models.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	}

	if err := repository.CreateContactMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "message sent successfully",
		"id":      msg.ID,
	})
}

func GetAllMessages(c *gin.Context) {
	messages, err := repository.GetAllMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
