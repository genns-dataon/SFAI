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
