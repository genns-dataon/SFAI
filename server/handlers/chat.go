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
                "who is", "who works", "who's in", "people in", "who hasn't", "who didn't",
                "engineering", "sales", "hr", "human resources",
                "developer", "manager", "engineer", "director", "reports to", "reports",
                "email", "contact", "hired", "hire date", "hired this", "hired in",
                "team", "department", "list all", "show me",
                "attendance", "clocked in", "clock in", "at work", "came last", "came first",
                "on leave", "leave", "vacation", "absent",
                "organization chart", "org chart", "edit employee", "how to", "how can i",
        }

        isEmployeeQuery := false
        for _, keyword := range employeeKeywords {
                if strings.Contains(messageLower, keyword) {
                        isEmployeeQuery = true
                        break
                }
        }
        
        // Also check for navigation/help queries
        navigationKeywords := []string{
                "where", "how do i", "how can i", "how to", "where can i",
                "organization chart", "org chart", "edit", "update", "modify",
        }
        for _, keyword := range navigationKeywords {
                if strings.Contains(messageLower, keyword) {
                        isEmployeeQuery = true
                        break
                }
        }

        // Handle all employee data requests locally (without sending PII to OpenAI)
        if isEmployeeQuery {
                var employees []models.Employee
                if err := database.DB.Preload("Department").Preload("Manager").Find(&employees).Error; err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee data"})
                        return
                }

                var response strings.Builder
                
                // Handle navigation/help queries first
                if strings.Contains(messageLower, "organization chart") || strings.Contains(messageLower, "org chart") {
                        response.WriteString("📊 To view the Organization Chart:\n\n")
                        response.WriteString("1. Click on 'Organization Chart' in the left sidebar menu (has a tree/hierarchy icon)\n")
                        response.WriteString("2. You'll see a visual representation of the company's reporting structure\n")
                        response.WriteString("3. The chart shows managers and their direct reports in a tree view\n\n")
                        response.WriteString("The organization chart automatically updates when you assign managers to employees.")
                } else if (strings.Contains(messageLower, "edit") || strings.Contains(messageLower, "update") || strings.Contains(messageLower, "modify")) && 
                           strings.Contains(messageLower, "employee") {
                        response.WriteString("✏️ To edit employee information:\n\n")
                        response.WriteString("1. Go to the 'Employees' page from the sidebar\n")
                        response.WriteString("2. Find the employee in the table\n")
                        response.WriteString("3. Click the 'Edit' button (pencil icon) in the Actions column\n")
                        response.WriteString("4. Update the information in the form\n")
                        response.WriteString("5. Click 'Save' to apply changes\n\n")
                        response.WriteString("You can edit: name, email, department, job title, hire date, and manager assignment.")
                } else if strings.Contains(messageLower, "how") && (strings.Contains(messageLower, "add") || strings.Contains(messageLower, "create")) {
                        response.WriteString("➕ To add a new employee:\n\n")
                        response.WriteString("1. Navigate to the 'Employees' page\n")
                        response.WriteString("2. Click the '+ Add Employee' button (top right)\n")
                        response.WriteString("3. Fill in the employee information:\n")
                        response.WriteString("   - Name, Email, Department, Job Title\n")
                        response.WriteString("   - Hire Date (optional)\n")
                        response.WriteString("   - Manager (optional - for org chart hierarchy)\n")
                        response.WriteString("4. Click 'Save' to create the employee record")
                } else if strings.Contains(messageLower, "manager") && !strings.Contains(messageLower, "who reports") {
                        // Find who someone's manager is (handles "manager of X", "X's manager", "who is X's manager")
                        foundMatch := false
                        for _, emp := range employees {
                                nameParts := strings.Fields(strings.ToLower(emp.Name))
                                nameMatch := false
                                for _, part := range nameParts {
                                        if strings.Contains(messageLower, part) && len(part) > 2 {
                                                nameMatch = true
                                                break
                                        }
                                }
                                if nameMatch {
                                        if emp.Manager != nil {
                                                response.WriteString(fmt.Sprintf("%s's manager is %s (%s)\n", 
                                                        emp.Name, emp.Manager.Name, emp.Manager.JobTitle))
                                        } else {
                                                response.WriteString(fmt.Sprintf("%s has no manager assigned (top-level position)\n", emp.Name))
                                        }
                                        foundMatch = true
                                        break
                                }
                        }
                        if !foundMatch {
                                response.WriteString("I couldn't find that employee. Please check the name and try again.\n")
                        }
                } else if strings.Contains(messageLower, "who reports to") || strings.Contains(messageLower, "direct reports") || strings.Contains(messageLower, "team members") {
                        // Find who reports to someone
                        foundMatch := false
                        for _, emp := range employees {
                                nameParts := strings.Fields(strings.ToLower(emp.Name))
                                nameMatch := false
                                for _, part := range nameParts {
                                        if strings.Contains(messageLower, part) && len(part) > 2 {
                                                nameMatch = true
                                                break
                                        }
                                }
                                if nameMatch {
                                        var reports []models.Employee
                                        database.DB.Preload("Department").Where("manager_id = ?", emp.ID).Find(&reports)
                                        if len(reports) > 0 {
                                                response.WriteString(fmt.Sprintf("%s has %d direct report(s):\n\n", emp.Name, len(reports)))
                                                for _, report := range reports {
                                                        deptName := "N/A"
                                                        if report.Department != nil {
                                                                deptName = report.Department.Name
                                                        }
                                                        response.WriteString(fmt.Sprintf("• %s (%s) - %s\n", report.Name, report.JobTitle, deptName))
                                                }
                                        } else {
                                                response.WriteString(fmt.Sprintf("%s has no direct reports.\n", emp.Name))
                                        }
                                        foundMatch = true
                                        break
                                }
                        }
                        if !foundMatch {
                                response.WriteString("I couldn't find that employee. Please check the name and try again.\n")
                        }
                } else if strings.Contains(messageLower, "who hasn't clocked in") || strings.Contains(messageLower, "who didn't clock in") || 
                           strings.Contains(messageLower, "not clocked in") || strings.Contains(messageLower, "haven't clocked in") {
                        // Find who hasn't clocked in today
                        var attendances []models.Attendance
                        database.DB.Where("DATE(date) = CURRENT_DATE").Find(&attendances)
                        
                        clockedInIDs := make(map[uint]bool)
                        for _, att := range attendances {
                                clockedInIDs[att.EmployeeID] = true
                        }
                        
                        var notClockedIn []models.Employee
                        for _, emp := range employees {
                                if !clockedInIDs[emp.ID] {
                                        notClockedIn = append(notClockedIn, emp)
                                }
                        }
                        
                        if len(notClockedIn) > 0 {
                                response.WriteString(fmt.Sprintf("Employees who haven't clocked in today (%d):\n\n", len(notClockedIn)))
                                for _, emp := range notClockedIn {
                                        deptName := "N/A"
                                        if emp.Department != nil {
                                                deptName = emp.Department.Name
                                        }
                                        response.WriteString(fmt.Sprintf("• %s (%s) - %s\n", emp.Name, emp.JobTitle, deptName))
                                }
                        } else {
                                response.WriteString("Great! All employees have clocked in today. 🎉\n")
                        }
                } else if strings.Contains(messageLower, "came last") || strings.Contains(messageLower, "came in last") || 
                           strings.Contains(messageLower, "latest arrival") || strings.Contains(messageLower, "last to arrive") {
                        // Find who came in last today
                        var attendances []models.Attendance
                        database.DB.Preload("Employee.Department").Where("DATE(date) = CURRENT_DATE").Order("clock_in DESC").Limit(5).Find(&attendances)
                        
                        if len(attendances) > 0 {
                                response.WriteString("Latest arrivals today:\n\n")
                                for i, att := range attendances {
                                        if att.Employee != nil {
                                                deptName := "N/A"
                                                if att.Employee.Department != nil {
                                                        deptName = att.Employee.Department.Name
                                                }
                                                response.WriteString(fmt.Sprintf("%d. %s (%s) - %s | Clocked in at %s\n", 
                                                        i+1, att.Employee.Name, att.Employee.JobTitle, deptName, att.ClockIn.Format("3:04 PM")))
                                        }
                                }
                        } else {
                                response.WriteString("No one has clocked in today yet.\n")
                        }
                } else if strings.Contains(messageLower, "came first") || strings.Contains(messageLower, "came in first") || 
                           strings.Contains(messageLower, "earliest arrival") || strings.Contains(messageLower, "first to arrive") {
                        // Find who came in first today
                        var attendances []models.Attendance
                        database.DB.Preload("Employee.Department").Where("DATE(date) = CURRENT_DATE").Order("clock_in ASC").Limit(5).Find(&attendances)
                        
                        if len(attendances) > 0 {
                                response.WriteString("Earliest arrivals today:\n\n")
                                for i, att := range attendances {
                                        if att.Employee != nil {
                                                deptName := "N/A"
                                                if att.Employee.Department != nil {
                                                        deptName = att.Employee.Department.Name
                                                }
                                                response.WriteString(fmt.Sprintf("%d. %s (%s) - %s | Clocked in at %s\n", 
                                                        i+1, att.Employee.Name, att.Employee.JobTitle, deptName, att.ClockIn.Format("3:04 PM")))
                                        }
                                }
                        } else {
                                response.WriteString("No one has clocked in today yet.\n")
                        }
                } else if strings.Contains(messageLower, "on leave") || strings.Contains(messageLower, "vacation") || 
                           (strings.Contains(messageLower, "leave") && (strings.Contains(messageLower, "this month") || 
                            strings.Contains(messageLower, "today") || strings.Contains(messageLower, "this week"))) {
                        // Find who is on leave
                        var leaveRequests []models.LeaveRequest
                        database.DB.Preload("Employee.Department").Where("status = ?", "approved").Find(&leaveRequests)
                        
                        var currentLeave []models.LeaveRequest
                        now := database.DB.NowFunc()
                        
                        for _, leave := range leaveRequests {
                                // Use inclusive bounds: Start <= now <= End (handles same-day leave)
                                startMatches := leave.StartDate.Before(now) || leave.StartDate.Format("2006-01-02") == now.Format("2006-01-02")
                                endMatches := leave.EndDate.After(now) || leave.EndDate.Format("2006-01-02") == now.Format("2006-01-02")
                                
                                if startMatches && endMatches {
                                        currentLeave = append(currentLeave, leave)
                                } else if strings.Contains(messageLower, "this month") {
                                        // Include any leave that overlaps with this month
                                        leaveInMonth := (leave.StartDate.Month() == now.Month() && leave.StartDate.Year() == now.Year()) ||
                                                        (leave.EndDate.Month() == now.Month() && leave.EndDate.Year() == now.Year()) ||
                                                        (leave.StartDate.Before(now) && leave.EndDate.After(now))
                                        if leaveInMonth {
                                                currentLeave = append(currentLeave, leave)
                                        }
                                }
                        }
                        
                        if len(currentLeave) > 0 {
                                response.WriteString("Employees currently on leave:\n\n")
                                for _, leave := range currentLeave {
                                        if leave.Employee != nil {
                                                deptName := "N/A"
                                                if leave.Employee.Department != nil {
                                                        deptName = leave.Employee.Department.Name
                                                }
                                                response.WriteString(fmt.Sprintf("• %s (%s) - %s | %s: %s to %s\n",
                                                        leave.Employee.Name, leave.Employee.JobTitle, deptName,
                                                        leave.LeaveType, leave.StartDate.Format("Jan 2"), leave.EndDate.Format("Jan 2")))
                                        }
                                }
                        } else {
                                response.WriteString("No employees are currently on approved leave.\n")
                        }
                } else if (strings.Contains(messageLower, "hired") || strings.Contains(messageLower, "joined")) && 
                           (strings.Contains(messageLower, "this year") || strings.Contains(messageLower, "this month") || 
                            strings.Contains(messageLower, "2025") || strings.Contains(messageLower, "2024")) {
                        // Find employees hired in specific time period
                        now := database.DB.NowFunc()
                        var filtered []models.Employee
                        
                        for _, emp := range employees {
                                match := false
                                if strings.Contains(messageLower, "this year") {
                                        if emp.HireDate.Year() == now.Year() {
                                                match = true
                                        }
                                } else if strings.Contains(messageLower, "this month") {
                                        if emp.HireDate.Year() == now.Year() && emp.HireDate.Month() == now.Month() {
                                                match = true
                                        }
                                } else if strings.Contains(messageLower, "2025") {
                                        if emp.HireDate.Year() == 2025 {
                                                match = true
                                        }
                                } else if strings.Contains(messageLower, "2024") {
                                        if emp.HireDate.Year() == 2024 {
                                                match = true
                                        }
                                }
                                if match {
                                        filtered = append(filtered, emp)
                                }
                        }
                        
                        if len(filtered) > 0 {
                                response.WriteString(fmt.Sprintf("Found %d employee(s) matching your criteria:\n\n", len(filtered)))
                                for _, emp := range filtered {
                                        deptName := "N/A"
                                        if emp.Department != nil {
                                                deptName = emp.Department.Name
                                        }
                                        response.WriteString(fmt.Sprintf("• %s (%s) - %s | Hired: %s\n",
                                                emp.Name, emp.JobTitle, deptName, emp.HireDate.Format("Jan 2, 2006")))
                                }
                        } else {
                                response.WriteString("No employees were hired during that time period.\n")
                        }
                } else

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
                } else if strings.Contains(messageLower, "at work") || strings.Contains(messageLower, "working today") || strings.Contains(messageLower, "clocked in") {
                        // Show who is at work today
                        var attendances []models.Attendance
                        today := strings.Split(fmt.Sprintf("%v", database.DB.NowFunc()), " ")[0]
                        database.DB.Preload("Employee").Where("DATE(date) = ?", today).Find(&attendances)

                        if len(attendances) > 0 {
                                response.WriteString(fmt.Sprintf("Employees at work today (%s):\n\n", today))
                                for _, att := range attendances {
                                        if att.Employee != nil {
                                                status := "Currently at work"
                                                if att.ClockOut != nil {
                                                        status = fmt.Sprintf("Clocked out at %s", att.ClockOut.Format("3:04 PM"))
                                                }
                                                response.WriteString(fmt.Sprintf(
                                                        "• %s (%s) - Clocked in: %s | %s\n",
                                                        att.Employee.Name,
                                                        att.Employee.JobTitle,
                                                        att.ClockIn.Format("3:04 PM"),
                                                        status,
                                                ))
                                        }
                                }
                        } else {
                                response.WriteString("No employees have clocked in today yet.\n")
                        }
                } else if strings.Contains(messageLower, "how many") && (strings.Contains(messageLower, "developer") || strings.Contains(messageLower, "engineer") || strings.Contains(messageLower, "manager") || strings.Contains(messageLower, "director")) {
                        // Count employees by job title/role
                        roleKeywords := map[string]string{
                                "developer": "developer",
                                "engineer": "engineer",
                                "manager": "manager",
                                "director": "director",
                                "executive": "executive",
                                "coordinator": "coordinator",
                                "representative": "representative",
                        }
                        
                        for keyword, role := range roleKeywords {
                                if strings.Contains(messageLower, keyword) {
                                        count := 0
                                        var matchedEmployees []string
                                        for _, emp := range employees {
                                                if strings.Contains(strings.ToLower(emp.JobTitle), role) {
                                                        count++
                                                        matchedEmployees = append(matchedEmployees, fmt.Sprintf("%s (%s)", emp.Name, emp.JobTitle))
                                                }
                                        }
                                        if count > 0 {
                                                response.WriteString(fmt.Sprintf("We have %d %s(s):\n\n", count, role))
                                                for _, emp := range matchedEmployees {
                                                        response.WriteString(fmt.Sprintf("• %s\n", emp))
                                                }
                                        } else {
                                                response.WriteString(fmt.Sprintf("We don't have any employees with '%s' in their job title.\n", role))
                                        }
                                        break
                                }
                        }
                } else {
                        // Try to find a specific employee by name (including partial name/first name)
                        foundEmployee := false
                        for _, emp := range employees {
                                nameParts := strings.Fields(strings.ToLower(emp.Name))
                                // Check full name or any part of the name
                                nameMatch := strings.Contains(messageLower, strings.ToLower(emp.Name))
                                for _, part := range nameParts {
                                        if strings.Contains(messageLower, part) && len(part) > 2 {
                                                nameMatch = true
                                                break
                                        }
                                }
                                
                                if nameMatch {
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
                                response.WriteString("• Get employee contact information by name (e.g., 'What is Emma's email?')\n")
                                response.WriteString("• View attendance records for any employee\n")
                                response.WriteString("• Check who is at work today\n")
                                response.WriteString("• Count employees by role (e.g., 'How many developers do we have?')\n")
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
