package handlers

import (
	"net/http"

	"hcm-backend/database"
	"hcm-backend/models"

	"github.com/gin-gonic/gin"
)

func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	result := database.DB.Preload("Department").Find(&employees)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	result := database.DB.Preload("Department").First(&employee, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&employee)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	database.DB.Preload("Department").First(&employee, employee.ID)
	c.JSON(http.StatusCreated, employee)
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	
	if err := database.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&employee)
	database.DB.Preload("Department").First(&employee, employee.ID)
	c.JSON(http.StatusOK, employee)
}
