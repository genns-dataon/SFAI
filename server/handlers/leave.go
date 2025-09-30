package handlers

import (
	"net/http"

	"hcm-backend/database"
	"hcm-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateLeaveRequest(c *gin.Context) {
	var leave models.LeaveRequest
	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave.Status = "pending"
	result := database.DB.Create(&leave)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	database.DB.Preload("Employee").First(&leave, leave.ID)
	c.JSON(http.StatusCreated, leave)
}

func GetLeaveRequests(c *gin.Context) {
	var leaves []models.LeaveRequest
	result := database.DB.Preload("Employee").Order("created_at desc").Find(&leaves)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, leaves)
}
