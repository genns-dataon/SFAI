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
        var createData struct {
                Name               string  `json:"name" binding:"required"`
                Email              string  `json:"email" binding:"required,email"`
                DepartmentID       uint    `json:"department_id"`
                JobTitle           string  `json:"job_title"`
                HireDate           string  `json:"hire_date"`
                ManagerID          *uint   `json:"manager_id"`
                EmployeeNumber     string  `json:"employee_number"`
                DateOfBirth        string  `json:"date_of_birth"`
                NationalID         string  `json:"national_id"`
                TaxID              string  `json:"tax_id"`
                MaritalStatus      string  `json:"marital_status"`
                EmploymentType     string  `json:"employment_type"`
                EmploymentStatus   string  `json:"employment_status"`
                JobLevel           string  `json:"job_level"`
                WorkLocation       string  `json:"work_location"`
                WorkArrangement    string  `json:"work_arrangement"`
                BaseSalary         float64 `json:"base_salary"`
                PayFrequency       string  `json:"pay_frequency"`
                Currency           string  `json:"currency"`
                BankAccount        string  `json:"bank_account"`
                BenefitEligibility string  `json:"benefit_eligibility"`
                ProbationEndDate   string  `json:"probation_end_date"`
                PerformanceRating  string  `json:"performance_rating"`
                Skills             string  `json:"skills"`
                TrainingCompleted  string  `json:"training_completed"`
                CareerNotes        string  `json:"career_notes"`
        }

        if err := c.ShouldBindJSON(&createData); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        // Prevent self-reporting
        if createData.ManagerID != nil && *createData.ManagerID == 0 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID"})
                return
        }

        employee := models.Employee{
                Name:               createData.Name,
                Email:              createData.Email,
                DepartmentID:       createData.DepartmentID,
                JobTitle:           createData.JobTitle,
                ManagerID:          createData.ManagerID,
                EmployeeNumber:     createData.EmployeeNumber,
                NationalID:         createData.NationalID,
                TaxID:              createData.TaxID,
                MaritalStatus:      createData.MaritalStatus,
                EmploymentType:     createData.EmploymentType,
                EmploymentStatus:   createData.EmploymentStatus,
                JobLevel:           createData.JobLevel,
                WorkLocation:       createData.WorkLocation,
                WorkArrangement:    createData.WorkArrangement,
                BaseSalary:         createData.BaseSalary,
                PayFrequency:       createData.PayFrequency,
                Currency:           createData.Currency,
                BankAccount:        createData.BankAccount,
                BenefitEligibility: createData.BenefitEligibility,
                PerformanceRating:  createData.PerformanceRating,
                Skills:             createData.Skills,
                TrainingCompleted:  createData.TrainingCompleted,
                CareerNotes:        createData.CareerNotes,
        }

        // Parse hire date if provided
        if createData.HireDate != "" {
                hireDate, err := time.Parse("2006-01-02", createData.HireDate)
                if err == nil {
                        employee.HireDate = hireDate
                }
        }

        // Parse date of birth if provided
        if createData.DateOfBirth != "" {
                dob, err := time.Parse("2006-01-02", createData.DateOfBirth)
                if err == nil {
                        employee.DateOfBirth = &dob
                }
        }

        // Parse probation end date if provided
        if createData.ProbationEndDate != "" {
                probationEnd, err := time.Parse("2006-01-02", createData.ProbationEndDate)
                if err == nil {
                        employee.ProbationEndDate = &probationEnd
                }
        }

        result := database.DB.Create(&employee)
        if result.Error != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
                return
        }

        database.DB.Preload("Department").Preload("Manager").First(&employee, employee.ID)
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
                Name               string  `json:"name" binding:"required"`
                Email              string  `json:"email" binding:"required,email"`
                DepartmentID       uint    `json:"department_id"`
                JobTitle           string  `json:"job_title"`
                HireDate           string  `json:"hire_date"`
                ManagerID          *uint   `json:"manager_id"`
                EmployeeNumber     string  `json:"employee_number"`
                DateOfBirth        string  `json:"date_of_birth"`
                NationalID         string  `json:"national_id"`
                TaxID              string  `json:"tax_id"`
                MaritalStatus      string  `json:"marital_status"`
                EmploymentType     string  `json:"employment_type"`
                EmploymentStatus   string  `json:"employment_status"`
                JobLevel           string  `json:"job_level"`
                WorkLocation       string  `json:"work_location"`
                WorkArrangement    string  `json:"work_arrangement"`
                BaseSalary         float64 `json:"base_salary"`
                PayFrequency       string  `json:"pay_frequency"`
                Currency           string  `json:"currency"`
                BankAccount        string  `json:"bank_account"`
                BenefitEligibility string  `json:"benefit_eligibility"`
                ProbationEndDate   string  `json:"probation_end_date"`
                PerformanceRating  string  `json:"performance_rating"`
                Skills             string  `json:"skills"`
                TrainingCompleted  string  `json:"training_completed"`
                CareerNotes        string  `json:"career_notes"`
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
        employee.EmployeeNumber = updateData.EmployeeNumber
        employee.NationalID = updateData.NationalID
        employee.TaxID = updateData.TaxID
        employee.MaritalStatus = updateData.MaritalStatus
        employee.EmploymentType = updateData.EmploymentType
        employee.EmploymentStatus = updateData.EmploymentStatus
        employee.JobLevel = updateData.JobLevel
        employee.WorkLocation = updateData.WorkLocation
        employee.WorkArrangement = updateData.WorkArrangement
        employee.BaseSalary = updateData.BaseSalary
        employee.PayFrequency = updateData.PayFrequency
        employee.Currency = updateData.Currency
        employee.BankAccount = updateData.BankAccount
        employee.BenefitEligibility = updateData.BenefitEligibility
        employee.PerformanceRating = updateData.PerformanceRating
        employee.Skills = updateData.Skills
        employee.TrainingCompleted = updateData.TrainingCompleted
        employee.CareerNotes = updateData.CareerNotes
        
        // Parse hire date if provided
        if updateData.HireDate != "" {
                hireDate, err := time.Parse("2006-01-02", updateData.HireDate)
                if err == nil {
                        employee.HireDate = hireDate
                }
        }

        // Parse date of birth if provided
        if updateData.DateOfBirth != "" {
                dob, err := time.Parse("2006-01-02", updateData.DateOfBirth)
                if err == nil {
                        employee.DateOfBirth = &dob
                }
        }

        // Parse probation end date if provided
        if updateData.ProbationEndDate != "" {
                probationEnd, err := time.Parse("2006-01-02", updateData.ProbationEndDate)
                if err == nil {
                        employee.ProbationEndDate = &probationEnd
                }
        }

        database.DB.Save(&employee)
        database.DB.Preload("Department").Preload("Manager").First(&employee, employee.ID)
        c.JSON(http.StatusOK, employee)
}
