package handlers

import (
        "net/http"
        "time"

        "hcm-backend/database"
        "hcm-backend/models"

        "github.com/gin-gonic/gin"
)

func ClockIn(c *gin.Context) {
        var input struct {
                EmployeeID uint   `json:"employee_id" binding:"required"`
                Location   string `json:"location"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        attendance := models.Attendance{
                EmployeeID: input.EmployeeID,
                Date:       time.Now(),
                ClockIn:    time.Now(),
                Location:   input.Location,
        }

        result := database.DB.Create(&attendance)
        if result.Error != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
                return
        }

        database.DB.Preload("Employee").First(&attendance, attendance.ID)
        c.JSON(http.StatusCreated, attendance)
}

func ClockOut(c *gin.Context) {
        var input struct {
                EmployeeID uint   `json:"employee_id" binding:"required"`
                Location   string `json:"location"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        // Find today's attendance record without clock out
        var attendance models.Attendance
        err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE AND clock_out IS NULL", input.EmployeeID).First(&attendance).Error
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "No active clock-in found for today"})
                return
        }

        // Update with clock out time
        now := time.Now()
        attendance.ClockOut = &now
        if input.Location != "" {
                attendance.Location = input.Location
        }

        result := database.DB.Save(&attendance)
        if result.Error != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
                return
        }

        database.DB.Preload("Employee").First(&attendance, attendance.ID)
        c.JSON(http.StatusOK, attendance)
}

func GetAttendance(c *gin.Context) {
        var attendances []models.Attendance
        result := database.DB.Preload("Employee").Order("date desc").Find(&attendances)
        if result.Error != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
                return
        }
        c.JSON(http.StatusOK, attendances)
}
