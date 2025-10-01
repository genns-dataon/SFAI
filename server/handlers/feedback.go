package handlers

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "hcm-backend/database"
        "hcm-backend/models"
)

func CreateFeedback(c *gin.Context) {
        var input struct {
                Question string `json:"question" binding:"required"`
                Response string `json:"response" binding:"required"`
                Rating   string `json:"rating" binding:"required"`
                Comment  string `json:"comment"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        userID, exists := c.Get("userID")
        var userIDPtr *uint
        if exists {
                uid := userID.(uint)
                userIDPtr = &uid
        }

        feedback := models.ChatFeedback{
                UserID:   userIDPtr,
                Question: input.Question,
                Response: input.Response,
                Rating:   input.Rating,
                Comment:  input.Comment,
        }

        if err := database.DB.Create(&feedback).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save feedback"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Feedback saved successfully", "feedback": feedback})
}

func GetAllFeedback(c *gin.Context) {
        var feedbacks []models.ChatFeedback
        if err := database.DB.Preload("User").Order("created_at DESC").Find(&feedbacks).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feedback"})
                return
        }

        c.JSON(http.StatusOK, feedbacks)
}

func UpdateFeedback(c *gin.Context) {
        id := c.Param("id")

        var input struct {
                Rating  string `json:"rating" binding:"required"`
                Comment string `json:"comment"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        var feedback models.ChatFeedback
        if err := database.DB.First(&feedback, id).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
                return
        }

        feedback.Rating = input.Rating
        if input.Comment != "" {
                feedback.Comment = input.Comment
        }

        if err := database.DB.Save(&feedback).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feedback"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Feedback updated successfully", "feedback": feedback})
}
