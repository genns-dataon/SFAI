package handlers

import (
        "fmt"
        "net/http"
        "os"
        "strings"
        "sync"

        "github.com/gin-gonic/gin"
        "github.com/openai/openai-go/v2"
        "hcm-backend/database"
        "hcm-backend/models"
)

var (
        openaiClient     openai.Client
        openaiClientOnce sync.Once
)

func getOpenAIClient() openai.Client {
        openaiClientOnce.Do(func() {
                openaiClient = openai.NewClient()
        })
        return openaiClient
}

func getDepartments(employees []models.Employee) map[string]bool {
        depts := make(map[string]bool)
        for _, emp := range employees {
                if emp.Department != nil {
                        depts[emp.Department.Name] = true
                }
        }
        return depts
}

func Chat(c *gin.Context) {
        var input struct {
                Message string `json:"message" binding:"required"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        messageLower := strings.ToLower(input.Message)

        // Check if this is an employee-related query (handle locally to protect PII)
        employeeKeywords := []string{
                "employee", "employees", "staff", "worker", "workers",
                "who is", "who works", "who's in", "people in",
                "engineering", "sales", "hr", "human resources",
                "developer", "manager", "engineer", "director",
                "email", "contact", "hired", "hire date",
                "team", "department", "list all", "show me",
        }

        isEmployeeQuery := false
        for _, keyword := range employeeKeywords {
                if strings.Contains(messageLower, keyword) {
                        isEmployeeQuery = true
                        break
                }
        }

        // Handle all employee data requests locally (without sending PII to OpenAI)
        if isEmployeeQuery {
                var employees []models.Employee
                if err := database.DB.Preload("Department").Find(&employees).Error; err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee data"})
                        return
                }

                var response strings.Builder

                // Check for specific query types (check department-specific queries first)
                if strings.Contains(messageLower, "engineering") {
                        response.WriteString("Engineering Department Employees:\n\n")
                        for _, emp := range employees {
                                if emp.Department != nil && strings.Contains(strings.ToLower(emp.Department.Name), "engineering") {
                                        response.WriteString(fmt.Sprintf(
                                                "• %s (%s) | Email: %s | Hired: %s\n",
                                                emp.Name, emp.JobTitle, emp.Email, emp.HireDate.Format("Jan 2, 2006"),
                                        ))
                                }
                        }
                } else if strings.Contains(messageLower, "sales") {
                        response.WriteString("Sales Department Employees:\n\n")
                        for _, emp := range employees {
                                if emp.Department != nil && strings.Contains(strings.ToLower(emp.Department.Name), "sales") {
                                        response.WriteString(fmt.Sprintf(
                                                "• %s (%s) | Email: %s | Hired: %s\n",
                                                emp.Name, emp.JobTitle, emp.Email, emp.HireDate.Format("Jan 2, 2006"),
                                        ))
                                }
                        }
                } else if strings.Contains(messageLower, "hr") || strings.Contains(messageLower, "human resources") {
                        response.WriteString("Human Resources Department Employees:\n\n")
                        for _, emp := range employees {
                                if emp.Department != nil && strings.Contains(strings.ToLower(emp.Department.Name), "human resources") {
                                        response.WriteString(fmt.Sprintf(
                                                "• %s (%s) | Email: %s | Hired: %s\n",
                                                emp.Name, emp.JobTitle, emp.Email, emp.HireDate.Format("Jan 2, 2006"),
                                        ))
                                }
                        }
                } else if strings.Contains(messageLower, "list") || strings.Contains(messageLower, "show") || strings.Contains(messageLower, "all") {
                        response.WriteString(fmt.Sprintf("Here are all %d employees in the system:\n\n", len(employees)))
                        for _, emp := range employees {
                                deptName := "N/A"
                                if emp.Department != nil {
                                        deptName = emp.Department.Name
                                }
                                response.WriteString(fmt.Sprintf(
                                        "• %s (%s) - %s | Email: %s | Hired: %s\n",
                                        emp.Name, emp.JobTitle, deptName, emp.Email, emp.HireDate.Format("Jan 2, 2006"),
                                ))
                        }
                } else if strings.Contains(messageLower, "attendance") {
                        // Handle attendance queries
                        var attendances []models.Attendance
                        database.DB.Preload("Employee.Department").Order("date desc").Find(&attendances)

                        // Check if asking for specific employee's attendance
                        foundAttendance := false
                        for _, emp := range employees {
                                if strings.Contains(messageLower, strings.ToLower(emp.Name)) {
                                        response.WriteString(fmt.Sprintf("Attendance records for %s:\n\n", emp.Name))
                                        recordCount := 0
                                        for _, att := range attendances {
                                                if att.EmployeeID == emp.ID {
                                                        recordCount++
                                                        clockOut := "Still clocked in"
                                                        if att.ClockOut != nil {
                                                                clockOut = att.ClockOut.Format("3:04 PM")
                                                        }
                                                        response.WriteString(fmt.Sprintf(
                                                                "• %s | Clock In: %s | Clock Out: %s | Location: %s\n",
                                                                att.Date.Format("Jan 2, 2006"),
                                                                att.ClockIn.Format("3:04 PM"),
                                                                clockOut,
                                                                att.Location,
                                                        ))
                                                        if recordCount >= 5 {
                                                                break
                                                        }
                                                }
                                        }
                                        if recordCount == 0 {
                                                response.WriteString("No attendance records found for this employee.\n")
                                        }
                                        foundAttendance = true
                                        break
                                }
                        }

                        if !foundAttendance {
                                // Show all recent attendance
                                response.WriteString("Recent attendance records:\n\n")
                                count := 0
                                for _, att := range attendances {
                                        if att.Employee != nil {
                                                clockOut := "Still clocked in"
                                                if att.ClockOut != nil {
                                                        clockOut = att.ClockOut.Format("3:04 PM")
                                                }
                                                response.WriteString(fmt.Sprintf(
                                                        "• %s - %s | %s | In: %s | Out: %s\n",
                                                        att.Date.Format("Jan 2, 2006"),
                                                        att.Employee.Name,
                                                        att.Location,
                                                        att.ClockIn.Format("3:04 PM"),
                                                        clockOut,
                                                ))
                                                count++
                                                if count >= 10 {
                                                        break
                                                }
                                        }
                                }
                                if count == 0 {
                                        response.WriteString("No attendance records found in the system.\n")
                                }
                        }
                } else {
                        // Try to find a specific employee by name
                        foundEmployee := false
                        for _, emp := range employees {
                                // Check if employee name is mentioned in the query
                                if strings.Contains(messageLower, strings.ToLower(emp.Name)) {
                                        deptName := "N/A"
                                        if emp.Department != nil {
                                                deptName = emp.Department.Name
                                        }
                                        response.WriteString(fmt.Sprintf("Here's the information for %s:\n\n", emp.Name))
                                        response.WriteString(fmt.Sprintf("• Name: %s\n", emp.Name))
                                        response.WriteString(fmt.Sprintf("• Job Title: %s\n", emp.JobTitle))
                                        response.WriteString(fmt.Sprintf("• Department: %s\n", deptName))
                                        response.WriteString(fmt.Sprintf("• Email: %s\n", emp.Email))
                                        response.WriteString(fmt.Sprintf("• Hire Date: %s\n", emp.HireDate.Format("Jan 2, 2006")))
                                        foundEmployee = true
                                        break
                                }
                        }
                        
                        if !foundEmployee {
                                // General employee info request
                                response.WriteString(fmt.Sprintf("I can help with employee information! We have %d employees across %d departments.\n\n", 
                                        len(employees), len(getDepartments(employees))))
                                response.WriteString("You can ask me to:\n")
                                response.WriteString("• List all employees\n")
                                response.WriteString("• Show employees in a specific department (Engineering, Sales, HR)\n")
                                response.WriteString("• Get employee contact information by name\n")
                                response.WriteString("• View attendance records for any employee\n")
                        }
                }

                c.JSON(http.StatusOK, gin.H{
                        "response": response.String(),
                        "message":  input.Message,
                })
                return
        }

        // For general HR questions, use OpenAI with aggregated statistics only (no PII)
        apiKey := os.Getenv("OPENAI_API_KEY")
        if apiKey == "" {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI API key not configured"})
                return
        }

        // Fetch aggregated statistics only (no PII)
        var totalEmployees int64
        database.DB.Model(&models.Employee{}).Count(&totalEmployees)

        var departments []models.Department
        database.DB.Find(&departments)

        var statsInfo strings.Builder
        statsInfo.WriteString("\n\nCOMPANY STATISTICS:\n")
        statsInfo.WriteString(fmt.Sprintf("- Total Employees: %d\n", totalEmployees))
        statsInfo.WriteString("- Departments: ")
        for i, dept := range departments {
                if i > 0 {
                        statsInfo.WriteString(", ")
                }
                statsInfo.WriteString(dept.Name)
        }
        statsInfo.WriteString("\n")

        client := getOpenAIClient()

        systemPrompt := `You are an AI assistant for an HR Management System. Your role is to help answer questions about:
- General HR policies and best practices
- Attendance tracking policies
- Leave request procedures
- Salary and payroll general information
- HR guidelines and recommendations

For specific employee data requests (like "list all employees" or "who works in engineering"), inform users that you've provided the information.

Be professional, helpful, and concise in your responses.` + statsInfo.String()

        chatCompletion, err := client.Chat.Completions.New(c.Request.Context(), openai.ChatCompletionNewParams{
                Messages: []openai.ChatCompletionMessageParamUnion{
                        openai.SystemMessage(systemPrompt),
                        openai.UserMessage(input.Message),
                },
                Model: openai.ChatModelGPT4oMini,
        })

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI response", "details": err.Error()})
                return
        }

        if len(chatCompletion.Choices) == 0 {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "No response from AI"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "response": chatCompletion.Choices[0].Message.Content,
                "message":  input.Message,
        })
}
