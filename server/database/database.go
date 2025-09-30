package database

import (
        "log"
        "os"
        "strings"
        "time"

        "hcm-backend/models"

        "golang.org/x/crypto/bcrypt"
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

        // Hash the default password "password" for all users
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
        if err != nil {
                log.Fatal("Failed to hash password:", err)
        }

        departments := []models.Department{
                {Name: "Engineering"},
                {Name: "Human Resources"},
                {Name: "Sales"},
        }

        for i := range departments {
                DB.Create(&departments[i])
        }

        // Create users and employees with linked accounts
        employeeData := []struct {
                Name       string
                Email      string
                Department uint
                JobTitle   string
                HireDate   time.Time
        }{
                {"Alice Johnson", "alice.johnson@company.com", 1, "Senior Software Engineer", time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)},
                {"Bob Smith", "bob.smith@company.com", 1, "Frontend Developer", time.Date(2021, 3, 20, 0, 0, 0, 0, time.UTC)},
                {"Carol White", "carol.white@company.com", 2, "HR Manager", time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC)},
                {"David Brown", "david.brown@company.com", 2, "Recruiter", time.Date(2022, 2, 5, 0, 0, 0, 0, time.UTC)},
                {"Emma Davis", "emma.davis@company.com", 3, "Sales Director", time.Date(2018, 9, 1, 0, 0, 0, 0, time.UTC)},
                {"Frank Wilson", "frank.wilson@company.com", 3, "Account Executive", time.Date(2021, 11, 15, 0, 0, 0, 0, time.UTC)},
                {"Grace Lee", "grace.lee@company.com", 1, "DevOps Engineer", time.Date(2020, 7, 22, 0, 0, 0, 0, time.UTC)},
                {"Henry Martinez", "henry.martinez@company.com", 3, "Sales Representative", time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)},
                {"Iris Taylor", "iris.taylor@company.com", 2, "HR Coordinator", time.Date(2022, 8, 30, 0, 0, 0, 0, time.UTC)},
                {"Jack Anderson", "jack.anderson@company.com", 1, "Backend Developer", time.Date(2021, 5, 18, 0, 0, 0, 0, time.UTC)},
        }

        for _, data := range employeeData {
                // Generate username from first name (lowercase)
                firstName := strings.Split(data.Name, " ")[0]
                username := strings.ToLower(firstName)

                // Create user account
                user := models.User{
                        Username: username,
                        Email:    data.Email,
                        Password: string(hashedPassword),
                }
                DB.Create(&user)

                // Create employee linked to user
                employee := models.Employee{
                        Name:         data.Name,
                        Email:        data.Email,
                        DepartmentID: data.Department,
                        JobTitle:     data.JobTitle,
                        HireDate:     data.HireDate,
                        UserID:       &user.ID,
                }
                DB.Create(&employee)
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

        log.Println("Database seeded with 10 employees (with user accounts) and 3 departments")
        log.Println("All user accounts have username = first name (lowercase) and password = 'password'")
}
