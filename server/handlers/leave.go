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

func UpdateLeaveStatus(c *gin.Context) {
        id := c.Param("id")
        
        var input struct {
                Status string `json:"status" binding:"required"`
        }
        
        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Validate status
        validStatuses := []string{"pending", "approved", "rejected"}
        isValid := false
        for _, status := range validStatuses {
                if input.Status == status {
                        isValid = true
                        break
                }
        }
        
        if !isValid {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be pending, approved, or rejected"})
                return
        }
        
        var leave models.LeaveRequest
        if err := database.DB.First(&leave, id).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
                return
        }
        
        leave.Status = input.Status
        
        if err := database.DB.Save(&leave).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave status"})
                return
        }
        
        database.DB.Preload("Employee").First(&leave, leave.ID)
        c.JSON(http.StatusOK, gin.H{"message": "Leave status updated successfully", "leave": leave})
}
