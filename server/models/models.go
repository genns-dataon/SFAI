package models

import (
        "time"
        "gorm.io/gorm"
)

type Employee struct {
        ID           uint           `gorm:"primarykey" json:"id"`
        CreatedAt    time.Time      `json:"created_at"`
        UpdatedAt    time.Time      `json:"updated_at"`
        DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
        Name         string         `json:"name" binding:"required"`
        Email        string         `gorm:"unique" json:"email" binding:"required,email"`
        DepartmentID uint           `json:"department_id"`
        Department   *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
        JobTitle     string         `json:"job_title"`
        HireDate     time.Time      `json:"hire_date"`
        ManagerID    *uint          `json:"manager_id"`
        Manager      *Employee      `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
        Reports      []Employee     `gorm:"foreignKey:ManagerID" json:"reports,omitempty"`
        
        EmployeeNumber      string     `json:"employee_number"`
        DateOfBirth         *time.Time `json:"date_of_birth"`
        NationalID          string     `json:"national_id"`
        TaxID               string     `json:"tax_id"`
        MaritalStatus       string     `json:"marital_status"`
        
        EmploymentType      string     `json:"employment_type"`
        EmploymentStatus    string     `json:"employment_status" gorm:"default:'active'"`
        JobLevel            string     `json:"job_level"`
        WorkLocation        string     `json:"work_location"`
        WorkArrangement     string     `json:"work_arrangement"`
        
        BaseSalary          float64    `json:"base_salary"`
        PayFrequency        string     `json:"pay_frequency"`
        Currency            string     `json:"currency" gorm:"default:'USD'"`
        BankAccount         string     `json:"bank_account"`
        BenefitEligibility  string     `json:"benefit_eligibility"`
        
        ProbationEndDate    *time.Time `json:"probation_end_date"`
        PerformanceRating   string     `json:"performance_rating"`
        Skills              string     `gorm:"type:text" json:"skills"`
        TrainingCompleted   string     `gorm:"type:text" json:"training_completed"`
        CareerNotes         string     `gorm:"type:text" json:"career_notes"`
}

type Department struct {
        ID        uint           `gorm:"primarykey" json:"id"`
        CreatedAt time.Time      `json:"created_at"`
        UpdatedAt time.Time      `json:"updated_at"`
        DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
        Name      string         `json:"name" binding:"required"`
        ParentID  *uint          `json:"parent_id"`
        Parent    *Department    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

type Attendance struct {
        ID         uint           `gorm:"primarykey" json:"id"`
        CreatedAt  time.Time      `json:"created_at"`
        UpdatedAt  time.Time      `json:"updated_at"`
        DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
        EmployeeID uint           `json:"employee_id" binding:"required"`
        Employee   *Employee      `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
        Date       time.Time      `json:"date" binding:"required"`
        ClockIn    time.Time      `json:"clock_in"`
        ClockOut   *time.Time     `json:"clock_out"`
        Location   string         `json:"location"`
}

type LeaveRequest struct {
        ID         uint           `gorm:"primarykey" json:"id"`
        CreatedAt  time.Time      `json:"created_at"`
        UpdatedAt  time.Time      `json:"updated_at"`
        DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
        EmployeeID uint           `json:"employee_id" binding:"required"`
        Employee   *Employee      `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
        LeaveType  string         `json:"leave_type" binding:"required"`
        StartDate  time.Time      `json:"start_date" binding:"required"`
        EndDate    time.Time      `json:"end_date" binding:"required"`
        Status     string         `json:"status" gorm:"default:'pending'"`
}

type SalaryComponent struct {
        ID            uint           `gorm:"primarykey" json:"id"`
        CreatedAt     time.Time      `json:"created_at"`
        UpdatedAt     time.Time      `json:"updated_at"`
        DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
        EmployeeID    uint           `json:"employee_id" binding:"required"`
        Employee      *Employee      `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
        Type          string         `json:"type" binding:"required"`
        Amount        float64        `json:"amount" binding:"required"`
        EffectiveDate time.Time      `json:"effective_date"`
}

type Document struct {
        ID         uint           `gorm:"primarykey" json:"id"`
        CreatedAt  time.Time      `json:"created_at"`
        UpdatedAt  time.Time      `json:"updated_at"`
        DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
        EmployeeID uint           `json:"employee_id" binding:"required"`
        Employee   *Employee      `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
        DocType    string         `json:"doc_type" binding:"required"`
        FilePath   string         `json:"file_path"`
        ExpiryDate *time.Time     `json:"expiry_date"`
}

type PayrollExport struct {
        ID         uint           `gorm:"primarykey" json:"id"`
        CreatedAt  time.Time      `json:"created_at"`
        UpdatedAt  time.Time      `json:"updated_at"`
        DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
        Period     string         `json:"period" binding:"required"`
        DataJSON   string         `gorm:"type:text" json:"data_json"`
        ExportedAt time.Time      `json:"exported_at"`
}

type User struct {
        ID        uint           `gorm:"primarykey" json:"id"`
        CreatedAt time.Time      `json:"created_at"`
        UpdatedAt time.Time      `json:"updated_at"`
        DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
        Email     string         `gorm:"unique" json:"email" binding:"required,email"`
        Password  string         `json:"-"`
}
