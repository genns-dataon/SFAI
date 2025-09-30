package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hcm-backend/database"
	"hcm-backend/models"
)

func GetSettings(c *gin.Context) {
	var settings []models.ChatbotSettings
	if err := database.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func GetSetting(c *gin.Context) {
	key := c.Param("key")
	var setting models.ChatbotSettings
	
	if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

func UpsertSetting(c *gin.Context) {
	var input struct {
		Key         string `json:"key" binding:"required"`
		Value       string `json:"value"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var setting models.ChatbotSettings
	
	// Check if setting exists
	err := database.DB.Where("key = ?", input.Key).First(&setting).Error
	if err != nil {
		// Create new setting
		setting = models.ChatbotSettings{
			Key:         input.Key,
			Value:       input.Value,
			Description: input.Description,
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create setting"})
			return
		}
	} else {
		// Update existing setting
		setting.Value = input.Value
		setting.Description = input.Description
		if err := database.DB.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
			return
		}
	}

	c.JSON(http.StatusOK, setting)
}

func DeleteSetting(c *gin.Context) {
	key := c.Param("key")
	
	if err := database.DB.Where("key = ?", key).Delete(&models.ChatbotSettings{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting deleted successfully"})
}
