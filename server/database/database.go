package database

import (
        "log"
        "os"
        "time"

        "hcm-backend/models"

        "gorm.io/driver/postgres"
        "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
        var err error
        dsn := os.Getenv("DATABASE_URL")
        if dsn == "" {
                log.Fatal("DATABASE_URL environment variable is required")
        }

        DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
                log.Fatal("Failed to connect to database:", err)
        }

        log.Println("Database connected successfully")
}

func Migrate() {
        err := DB.AutoMigrate(
                &models.User{},
                &models.Department{},
                &models.Employee{},
                &models.Attendance{},
                &models.LeaveRequest{},
                &models.SalaryComponent{},
                &models.Document{},
                &models.PayrollExport{},
        )
        if err != nil {
                log.Fatal("Failed to migrate database:", err)
        }
        log.Println("Database migrated successfully")
}

func SeedData() {
        var count int64
        DB.Model(&models.Department{}).Count(&count)
        if count > 0 {
                log.Println("Database already seeded, skipping...")
                return
        }

        departments := []models.Department{
                {Name: "Engineering"},
                {Name: "Human Resources"},
                {Name: "Sales"},
        }

        for i := range departments {
                DB.Create(&departments[i])
        }

        employees := []models.Employee{
                {Name: "Alice Johnson", Email: "alice.johnson@company.com", DepartmentID: 1, JobTitle: "Senior Software Engineer", HireDate: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)},
                {Name: "Bob Smith", Email: "bob.smith@company.com", DepartmentID: 1, JobTitle: "Frontend Developer", HireDate: time.Date(2021, 3, 20, 0, 0, 0, 0, time.UTC)},
                {Name: "Carol White", Email: "carol.white@company.com", DepartmentID: 2, JobTitle: "HR Manager", HireDate: time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC)},
                {Name: "David Brown", Email: "david.brown@company.com", DepartmentID: 2, JobTitle: "Recruiter", HireDate: time.Date(2022, 2, 5, 0, 0, 0, 0, time.UTC)},
                {Name: "Emma Davis", Email: "emma.davis@company.com", DepartmentID: 3, JobTitle: "Sales Director", HireDate: time.Date(2018, 9, 1, 0, 0, 0, 0, time.UTC)},
                {Name: "Frank Wilson", Email: "frank.wilson@company.com", DepartmentID: 3, JobTitle: "Account Executive", HireDate: time.Date(2021, 11, 15, 0, 0, 0, 0, time.UTC)},
                {Name: "Grace Lee", Email: "grace.lee@company.com", DepartmentID: 1, JobTitle: "DevOps Engineer", HireDate: time.Date(2020, 7, 22, 0, 0, 0, 0, time.UTC)},
                {Name: "Henry Martinez", Email: "henry.martinez@company.com", DepartmentID: 3, JobTitle: "Sales Representative", HireDate: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)},
                {Name: "Iris Taylor", Email: "iris.taylor@company.com", DepartmentID: 2, JobTitle: "HR Coordinator", HireDate: time.Date(2022, 8, 30, 0, 0, 0, 0, time.UTC)},
                {Name: "Jack Anderson", Email: "jack.anderson@company.com", DepartmentID: 1, JobTitle: "Backend Developer", HireDate: time.Date(2021, 5, 18, 0, 0, 0, 0, time.UTC)},
        }

        for i := range employees {
                DB.Create(&employees[i])
        }

        salaries := []models.SalaryComponent{
                {EmployeeID: 1, Type: "Base Salary", Amount: 95000, EffectiveDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
                {EmployeeID: 2, Type: "Base Salary", Amount: 75000, EffectiveDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
                {EmployeeID: 3, Type: "Base Salary", Amount: 85000, EffectiveDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
                {EmployeeID: 4, Type: "Base Salary", Amount: 55000, EffectiveDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
                {EmployeeID: 5, Type: "Base Salary", Amount: 110000, EffectiveDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
        }

        for i := range salaries {
                DB.Create(&salaries[i])
        }

        log.Println("Database seeded with 10 employees and 3 departments")
}
