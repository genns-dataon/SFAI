package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"hcm-backend/database"
	"hcm-backend/models"

	"github.com/gin-gonic/gin"
)

func ExportSalary(c *gin.Context) {
	var salaries []models.SalaryComponent
	result := database.DB.Preload("Employee").Find(&salaries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	dataJSON, _ := json.Marshal(salaries)
	export := models.PayrollExport{
		Period:     time.Now().Format("2006-01"),
		DataJSON:   string(dataJSON),
		ExportedAt: time.Now(),
	}

	database.DB.Create(&export)
	c.JSON(http.StatusOK, gin.H{
		"export_id": export.ID,
		"period":    export.Period,
		"data":      salaries,
	})
}

func GeneratePayslip(c *gin.Context) {
	var input struct {
		EmployeeID uint `json:"employee_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var employee models.Employee
	if err := database.DB.First(&employee, input.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var salaries []models.SalaryComponent
	database.DB.Where("employee_id = ?", input.EmployeeID).Find(&salaries)

	c.JSON(http.StatusOK, gin.H{
		"employee": employee,
		"salaries": salaries,
		"period":   time.Now().Format("2006-01"),
	})
}
