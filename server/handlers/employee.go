package handlers

import (
        "net/http"
        "time"

        "hcm-backend/database"
        "hcm-backend/models"

        "github.com/gin-gonic/gin"
)

func GetEmployees(c *gin.Context) {
        var employees []models.Employee
        result := database.DB.Preload("Department").Preload("Manager").Find(&employees)
        if result.Error != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
                return
        }
        c.JSON(http.StatusOK, employees)
}

func GetEmployee(c *gin.Context) {
        id := c.Param("id")
        var employee models.Employee
        result := database.DB.Preload("Department").Preload("Manager").Preload("Reports").First(&employee, id)
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

        // Store original ID before binding
        originalID := employee.ID
        
        // Bind request into a separate struct to avoid overwriting ID
        var updateData struct {
                Name         string     `json:"name" binding:"required"`
                Email        string     `json:"email" binding:"required,email"`
                DepartmentID uint       `json:"department_id"`
                JobTitle     string     `json:"job_title"`
                HireDate     string     `json:"hire_date"`
                ManagerID    *uint      `json:"manager_id"`
        }

        if err := c.ShouldBindJSON(&updateData); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        // Prevent self-reporting using original ID
        if updateData.ManagerID != nil && *updateData.ManagerID == originalID {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Employee cannot be their own manager"})
                return
        }

        // Update only permitted fields
        employee.Name = updateData.Name
        employee.Email = updateData.Email
        employee.DepartmentID = updateData.DepartmentID
        employee.JobTitle = updateData.JobTitle
        employee.ManagerID = updateData.ManagerID
        
        // Parse hire date if provided
        if updateData.HireDate != "" {
                hireDate, err := time.Parse("2006-01-02", updateData.HireDate)
                if err == nil {
                        employee.HireDate = hireDate
                }
        }

        database.DB.Save(&employee)
        database.DB.Preload("Department").Preload("Manager").First(&employee, employee.ID)
        c.JSON(http.StatusOK, employee)
}
