package handlers

import (
        "fmt"
        "net/http"
        "os"
        "strings"
        "sync"
        "time"

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

        // Check for attendance clock-in requests
        clockInKeywords := []string{
                "clock in", "clock-in", "check in", "start work", "start my shift",
                "mark attendance", "i'm here", "im here", "arrived", "present",
        }
        for _, keyword := range clockInKeywords {
                if strings.Contains(messageLower, keyword) {
                        // Get current user ID from context
                        userID, exists := c.Get("userID")
                        if !exists {
                                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                                return
                        }

                        // Find employee associated with this user
                        var employee models.Employee
                        if err := database.DB.Where("user_id = ?", userID).First(&employee).Error; err != nil {
                                c.JSON(http.StatusOK, gin.H{
                                        "response": "❌ I couldn't find your employee record. Please contact HR.",
                                })
                                return
                        }

                        // Check if already clocked in today
                        var existingAttendance models.Attendance
                        err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE", employee.ID).First(&existingAttendance).Error
                        if err == nil {
                                clockOutStatus := ""
                                if existingAttendance.ClockOut == nil {
                                        clockOutStatus = " (not yet clocked out)"
                                } else {
                                        clockOutStatus = fmt.Sprintf(" (clocked out at %s)", existingAttendance.ClockOut.Format("3:04 PM"))
                                }
                                c.JSON(http.StatusOK, gin.H{
                                        "response": fmt.Sprintf("✅ You already clocked in today at %s%s.", 
                                                existingAttendance.ClockIn.Format("3:04 PM"), clockOutStatus),
                                })
                                return
                        }

                        // Create attendance record
                        attendance := models.Attendance{
                                EmployeeID: employee.ID,
                                Date:       time.Now(),
                                ClockIn:    time.Now(),
                        }

                        if err := database.DB.Create(&attendance).Error; err != nil {
                                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record attendance"})
                                return
                        }

                        c.JSON(http.StatusOK, gin.H{
                                "response": fmt.Sprintf("✅ Attendance recorded successfully!\n\n👋 Welcome, %s!\n⏰ Clock-in time: %s\n\nHave a productive day!", 
                                        employee.Name, attendance.ClockIn.Format("3:04 PM")),
                        })
                        return
                }
        }

        // Check for attendance clock-out requests
        clockOutKeywords := []string{
                "clock out", "clock-out", "check out", "end work", "end my shift",
                "leaving", "done for the day", "finish work", "log out",
        }
        for _, keyword := range clockOutKeywords {
                if strings.Contains(messageLower, keyword) {
                        // Get current user ID from context
                        userID, exists := c.Get("userID")
                        if !exists {
                                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                                return
                        }

                        // Find employee associated with this user
                        var employee models.Employee
                        if err := database.DB.Where("user_id = ?", userID).First(&employee).Error; err != nil {
                                c.JSON(http.StatusOK, gin.H{
                                        "response": "❌ I couldn't find your employee record. Please contact HR.",
                                })
                                return
                        }

                        // Find today's attendance record without clock out
                        var attendance models.Attendance
                        err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE AND clock_out IS NULL", employee.ID).First(&attendance).Error
                        if err != nil {
                                c.JSON(http.StatusOK, gin.H{
                                        "response": "❌ You don't have an active clock-in for today. Please clock in first!",
                                })
                                return
                        }

                        // Update with clock out time
                        now := time.Now()
                        attendance.ClockOut = &now

                        if err := database.DB.Save(&attendance).Error; err != nil {
                                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
                                return
                        }

                        duration := now.Sub(attendance.ClockIn)
                        hours := int(duration.Hours())
                        minutes := int(duration.Minutes()) % 60

                        c.JSON(http.StatusOK, gin.H{
                                "response": fmt.Sprintf("✅ Successfully clocked out!\n\n👋 See you later, %s!\n⏰ Clock-out time: %s\n📊 Total time: %dh %dm\n\nHave a great evening!", 
                                        employee.Name, now.Format("3:04 PM"), hours, minutes),
                        })
                        return
                }
        }

        // Check for attendance recording requests (for managers recording for others)
        recordAttendanceKeywords := []string{
                "record attendance", "mark attendance for", "clock in for", "clock out for",
                "record start time", "record end time", "attendance for",
        }
        isRecordingForOther := false
        for _, keyword := range recordAttendanceKeywords {
                if strings.Contains(messageLower, keyword) {
                        isRecordingForOther = true
                        break
                }
        }

        if isRecordingForOther {
                // Load all employees to help match names
                var employees []models.Employee
                database.DB.Preload("Department").Find(&employees)

                // Try to extract employee name from the message
                var targetEmployee *models.Employee
                for _, emp := range employees {
                        nameLower := strings.ToLower(emp.Name)
                        if strings.Contains(messageLower, nameLower) {
                                targetEmployee = &emp
                                break
                        }
                }

                if targetEmployee == nil {
                        // Ask for employee name
                        var response strings.Builder
                        response.WriteString("📝 I can help you record attendance!\n\n")
                        response.WriteString("Please tell me whose attendance you want to record. Here are the employees:\n\n")
                        for _, emp := range employees {
                                deptName := "N/A"
                                if emp.Department != nil {
                                        deptName = emp.Department.Name
                                }
                                response.WriteString(fmt.Sprintf("• %s (%s) - %s\n", emp.Name, emp.JobTitle, deptName))
                        }
                        response.WriteString("\nExample: \"Record start time for John Smith\"")
                        c.JSON(http.StatusOK, gin.H{"response": response.String()})
                        return
                }

                // Determine if it's clock in or clock out
                isClockOut := strings.Contains(messageLower, "end") || strings.Contains(messageLower, "out") || 
                        strings.Contains(messageLower, "finish") || strings.Contains(messageLower, "leaving")

                if isClockOut {
                        // Record clock out
                        var attendance models.Attendance
                        err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE AND clock_out IS NULL", targetEmployee.ID).First(&attendance).Error
                        if err != nil {
                                c.JSON(http.StatusOK, gin.H{
                                        "response": fmt.Sprintf("❌ %s doesn't have an active clock-in for today.", targetEmployee.Name),
                                })
                                return
                        }

                        now := time.Now()
                        attendance.ClockOut = &now

                        if err := database.DB.Save(&attendance).Error; err != nil {
                                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
                                return
                        }

                        duration := now.Sub(attendance.ClockIn)
                        hours := int(duration.Hours())
                        minutes := int(duration.Minutes()) % 60

                        c.JSON(http.StatusOK, gin.H{
                                "response": fmt.Sprintf("✅ Clock-out recorded for %s!\n\n⏰ Clock-out time: %s\n📊 Total time: %dh %dm", 
                                        targetEmployee.Name, now.Format("3:04 PM"), hours, minutes),
                        })
                } else {
                        // Record clock in
                        var existingAttendance models.Attendance
                        err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE", targetEmployee.ID).First(&existingAttendance).Error
                        if err == nil {
                                c.JSON(http.StatusOK, gin.H{
                                        "response": fmt.Sprintf("✅ %s already clocked in today at %s.", 
                                                targetEmployee.Name, existingAttendance.ClockIn.Format("3:04 PM")),
                                })
                                return
                        }

                        attendance := models.Attendance{
                                EmployeeID: targetEmployee.ID,
                                Date:       time.Now(),
                                ClockIn:    time.Now(),
                        }

                        if err := database.DB.Create(&attendance).Error; err != nil {
                                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record attendance"})
                                return
                        }

                        c.JSON(http.StatusOK, gin.H{
                                "response": fmt.Sprintf("✅ Attendance recorded for %s!\n\n⏰ Clock-in time: %s", 
                                        targetEmployee.Name, attendance.ClockIn.Format("3:04 PM")),
                        })
                }
                return
        }

        // Check for leave request submissions
        leaveRequestKeywords := []string{
                "request leave", "apply for leave", "take leave", "need leave",
                "submit leave", "leave request", "time off", "vacation request",
                "i want leave", "i need time off", "book leave",
        }
        isLeaveRequest := false
        for _, keyword := range leaveRequestKeywords {
                if strings.Contains(messageLower, keyword) {
                        isLeaveRequest = true
                        break
                }
        }

        if isLeaveRequest {
                // Get current user ID from context
                userID, exists := c.Get("userID")
                if !exists {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                        return
                }

                // Find employee associated with this user
                var employee models.Employee
                if err := database.DB.Where("user_id = ?", userID).First(&employee).Error; err != nil {
                        c.JSON(http.StatusOK, gin.H{
                                "response": "❌ I couldn't find your employee record. Please contact HR.",
                        })
                        return
                }

                // Parse dates from the message
                var startDate, endDate time.Time
                var leaveType string
                dateFound := false

                // Check for "tomorrow"
                if strings.Contains(messageLower, "tomorrow") {
                        tomorrow := time.Now().AddDate(0, 0, 1)
                        startDate = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
                        endDate = startDate
                        dateFound = true
                }

                // Check for "today"
                if strings.Contains(messageLower, "today") {
                        today := time.Now()
                        startDate = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
                        endDate = startDate
                        dateFound = true
                }

                // Check for "next week"
                if strings.Contains(messageLower, "next week") {
                        nextWeek := time.Now().AddDate(0, 0, 7)
                        // Find Monday of next week
                        daysUntilMonday := (8 - int(nextWeek.Weekday())) % 7
                        if daysUntilMonday == 0 {
                                daysUntilMonday = 7
                        }
                        monday := nextWeek.AddDate(0, 0, daysUntilMonday)
                        friday := monday.AddDate(0, 0, 4)
                        startDate = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
                        endDate = time.Date(friday.Year(), friday.Month(), friday.Day(), 0, 0, 0, 0, time.UTC)
                        dateFound = true
                }

                // Determine leave type
                if strings.Contains(messageLower, "sick") {
                        leaveType = "Sick Leave"
                } else if strings.Contains(messageLower, "vacation") {
                        leaveType = "Vacation"
                } else if strings.Contains(messageLower, "personal") {
                        leaveType = "Personal"
                } else if strings.Contains(messageLower, "emergency") {
                        leaveType = "Emergency"
                } else {
                        leaveType = "Vacation" // Default
                }

                // If we found dates, create the leave request
                if dateFound {
                        leaveRequest := models.LeaveRequest{
                                EmployeeID: employee.ID,
                                LeaveType:  leaveType,
                                StartDate:  startDate,
                                EndDate:    endDate,
                                Status:     "pending",
                        }

                        if err := database.DB.Create(&leaveRequest).Error; err != nil {
                                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create leave request"})
                                return
                        }

                        days := int(endDate.Sub(startDate).Hours()/24) + 1

                        c.JSON(http.StatusOK, gin.H{
                                "response": fmt.Sprintf("✅ Leave request submitted successfully!\n\n📝 Details:\n"+
                                        "• Employee: %s\n"+
                                        "• Leave Type: %s\n"+
                                        "• Start Date: %s\n"+
                                        "• End Date: %s\n"+
                                        "• Duration: %d day(s)\n"+
                                        "• Status: Pending\n\n"+
                                        "Your manager will review and approve/reject the request.",
                                        employee.Name, leaveType, startDate.Format("Jan 02, 2006"),
                                        endDate.Format("Jan 02, 2006"), days),
                        })
                        return
                }

                // Ask for more information
                c.JSON(http.StatusOK, gin.H{
                        "response": "📝 I can help you submit a leave request!\n\n" +
                                "Please provide the following information:\n" +
                                "• When do you want to take leave? (e.g., 'tomorrow', 'today', 'next week')\n" +
                                "• What type of leave? (Vacation, Sick Leave, Personal, Emergency)\n\n" +
                                "Example: \"I want to request vacation leave for tomorrow\"\n" +
                                "Example: \"I need sick leave for next week\"",
                })
                return
        }

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
                "salary", "compensation", "pay", "benefits", "performance", "rating",
                "skills", "certifications", "training", "employment status", "work location",
                "employment type", "full-time", "part-time", "contract", "probation",
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
                
                // Handle queries about new employee fields
                if strings.Contains(messageLower, "salary") || strings.Contains(messageLower, "compensation") || strings.Contains(messageLower, "pay") {
                        response.WriteString("💰 Employee Salary Information:\n\n")
                        for _, emp := range employees {
                                if emp.BaseSalary > 0 {
                                        currency := emp.Currency
                                        if currency == "" {
                                                currency = "USD"
                                        }
                                        frequency := emp.PayFrequency
                                        if frequency == "" {
                                                frequency = "Monthly"
                                        }
                                        response.WriteString(fmt.Sprintf("• %s (%s): %s %.2f (%s)\n",
                                                emp.Name, emp.JobTitle, currency, emp.BaseSalary, frequency))
                                }
                        }
                        if response.Len() <= 40 {
                                response.WriteString("No salary information available in the system.")
                        }
                } else if strings.Contains(messageLower, "skills") || strings.Contains(messageLower, "certifications") {
                        response.WriteString("🎯 Employee Skills & Certifications:\n\n")
                        for _, emp := range employees {
                                if emp.Skills != "" {
                                        response.WriteString(fmt.Sprintf("• %s (%s):\n  Skills: %s\n",
                                                emp.Name, emp.JobTitle, emp.Skills))
                                        if emp.TrainingCompleted != "" {
                                                response.WriteString(fmt.Sprintf("  Training: %s\n", emp.TrainingCompleted))
                                        }
                                        response.WriteString("\n")
                                }
                        }
                        if response.Len() <= 40 {
                                response.WriteString("No skills information available in the system.")
                        }
                } else if strings.Contains(messageLower, "performance") || strings.Contains(messageLower, "rating") {
                        response.WriteString("⭐ Employee Performance Ratings:\n\n")
                        for _, emp := range employees {
                                if emp.PerformanceRating != "" {
                                        response.WriteString(fmt.Sprintf("• %s (%s): %s\n",
                                                emp.Name, emp.JobTitle, emp.PerformanceRating))
                                }
                        }
                        if response.Len() <= 40 {
                                response.WriteString("No performance ratings available in the system.")
                        }
                } else if strings.Contains(messageLower, "employment status") || strings.Contains(messageLower, "active") || 
                           strings.Contains(messageLower, "probation") || strings.Contains(messageLower, "resigned") {
                        response.WriteString("📊 Employment Status:\n\n")
                        statusCounts := make(map[string]int)
                        for _, emp := range employees {
                                if emp.EmploymentStatus != "" {
                                        statusCounts[emp.EmploymentStatus]++
                                }
                        }
                        for status, count := range statusCounts {
                                response.WriteString(fmt.Sprintf("• %s: %d employees\n", status, count))
                        }
                        response.WriteString("\nDetailed List:\n\n")
                        for _, emp := range employees {
                                if emp.EmploymentStatus != "" {
                                        deptName := "N/A"
                                        if emp.Department != nil {
                                                deptName = emp.Department.Name
                                        }
                                        response.WriteString(fmt.Sprintf("• %s (%s) - %s | Status: %s\n",
                                                emp.Name, emp.JobTitle, deptName, emp.EmploymentStatus))
                                }
                        }
                } else if strings.Contains(messageLower, "work location") || strings.Contains(messageLower, "office") || 
                           strings.Contains(messageLower, "remote") || strings.Contains(messageLower, "hybrid") {
                        response.WriteString("🏢 Work Arrangements:\n\n")
                        for _, emp := range employees {
                                if emp.WorkArrangement != "" || emp.WorkLocation != "" {
                                        location := emp.WorkLocation
                                        if location == "" {
                                                location = "Not specified"
                                        }
                                        arrangement := emp.WorkArrangement
                                        if arrangement == "" {
                                                arrangement = "Not specified"
                                        }
                                        response.WriteString(fmt.Sprintf("• %s (%s): %s | %s\n",
                                                emp.Name, emp.JobTitle, arrangement, location))
                                }
                        }
                } else if strings.Contains(messageLower, "employment type") || strings.Contains(messageLower, "full-time") || 
                           strings.Contains(messageLower, "part-time") || strings.Contains(messageLower, "contract") {
                        response.WriteString("👥 Employment Types:\n\n")
                        typeCounts := make(map[string]int)
                        for _, emp := range employees {
                                if emp.EmploymentType != "" {
                                        typeCounts[emp.EmploymentType]++
                                }
                        }
                        for empType, count := range typeCounts {
                                response.WriteString(fmt.Sprintf("• %s: %d employees\n", empType, count))
                        }
                } else if strings.Contains(messageLower, "benefits") {
                        response.WriteString("🎁 Employee Benefits Eligibility:\n\n")
                        for _, emp := range employees {
                                if emp.BenefitEligibility != "" {
                                        response.WriteString(fmt.Sprintf("• %s (%s): %s\n",
                                                emp.Name, emp.JobTitle, emp.BenefitEligibility))
                                }
                        }
                } else if strings.Contains(messageLower, "organization chart") || strings.Contains(messageLower, "org chart") {
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
                        response.WriteString("4. Update the information in the tabbed form:\n")
                        response.WriteString("   - Basic Info: name, email, department, job title, manager, hire date\n")
                        response.WriteString("   - Personal & ID: employee number, date of birth, IDs, marital status\n")
                        response.WriteString("   - Employment: type, status, level, location, work arrangement\n")
                        response.WriteString("   - Compensation: salary, pay frequency, currency, bank account, benefits\n")
                        response.WriteString("   - Performance: probation date, rating, skills, training, notes\n")
                        response.WriteString("5. Click 'Update Employee' to save changes")
                } else if strings.Contains(messageLower, "how") && (strings.Contains(messageLower, "add") || strings.Contains(messageLower, "create")) {
                        response.WriteString("➕ To add a new employee:\n\n")
                        response.WriteString("1. Navigate to the 'Employees' page\n")
                        response.WriteString("2. Click the '+ Add Employee' button (top right)\n")
                        response.WriteString("3. Fill in the employee information across 5 tabs:\n")
                        response.WriteString("   - Basic Info (required): Name, Email, Department, Job Title\n")
                        response.WriteString("   - Personal & ID: Employee number, Date of birth, IDs, Marital status\n")
                        response.WriteString("   - Employment: Type, Status, Level, Location, Work arrangement\n")
                        response.WriteString("   - Compensation: Salary, Pay frequency, Currency, Bank account, Benefits\n")
                        response.WriteString("   - Performance: Probation date, Rating, Skills, Training, Career notes\n")
                        response.WriteString("4. Click 'Add Employee' to create the record")
                } else if strings.Contains(messageLower, "manager") && !strings.Contains(messageLower, "who reports") {
                        // Find who someone's manager is
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
                                startMatches := leave.StartDate.Before(now) || leave.StartDate.Format("2006-01-02") == now.Format("2006-01-02")
                                endMatches := leave.EndDate.After(now) || leave.EndDate.Format("2006-01-02") == now.Format("2006-01-02")
                                
                                if startMatches && endMatches {
                                        currentLeave = append(currentLeave, leave)
                                } else if strings.Contains(messageLower, "this month") {
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
                } else if strings.Contains(messageLower, "engineering") {
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
                        database.DB.Preload("Employee.Department").Where("DATE(date) = CURRENT_DATE").Order("clock_in ASC").Find(&attendances)
                        
                        if len(attendances) > 0 {
                                response.WriteString(fmt.Sprintf("Today's Attendance (%d employees clocked in):\n\n", len(attendances)))
                                for i, att := range attendances {
                                        if att.Employee != nil {
                                                deptName := "N/A"
                                                if att.Employee.Department != nil {
                                                        deptName = att.Employee.Department.Name
                                                }
                                                response.WriteString(fmt.Sprintf("%d. %s (%s) - %s | %s\n", 
                                                        i+1, att.Employee.Name, att.Employee.JobTitle, deptName, att.ClockIn.Format("3:04 PM")))
                                        }
                                }
                        } else {
                                response.WriteString("No attendance records for today yet.\n")
                        }
                } else {
                        // Default response for unmatched employee queries
                        response.WriteString("I can help you with:\n\n")
                        response.WriteString("📋 Employee Information:\n")
                        response.WriteString("• List all employees\n")
                        response.WriteString("• Show employees by department (Engineering, Sales, HR)\n")
                        response.WriteString("• View employee details (salary, skills, performance, status)\n")
                        response.WriteString("• Find managers and direct reports\n\n")
                        response.WriteString("⏰ Attendance:\n")
                        response.WriteString("• Clock in (say 'clock in' or 'report attendance')\n")
                        response.WriteString("• View today's attendance\n")
                        response.WriteString("• Check who hasn't clocked in\n\n")
                        response.WriteString("🏖️ Leave Management:\n")
                        response.WriteString("• Request leave (I'll guide you to the Leave page)\n")
                        response.WriteString("• Check who's on leave\n\n")
                        response.WriteString("💡 Try asking me questions like:\n")
                        response.WriteString("• 'List all employees in Engineering'\n")
                        response.WriteString("• 'Who is on leave this month?'\n")
                        response.WriteString("• 'Show me salary information'\n")
                        response.WriteString("• 'Clock in'\n")
                        response.WriteString("• 'Request leave'")
                }

                c.JSON(http.StatusOK, gin.H{"response": response.String()})
                return
        }

        // For non-employee queries, use OpenAI (with aggregated stats only, no PII)
        apiKey := os.Getenv("OPENAI_API_KEY")
        if apiKey == "" {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI API key not configured"})
                return
        }

        // Get aggregated statistics (no PII)
        var employees []models.Employee
        database.DB.Preload("Department").Find(&employees)

        deptStats := getDepartments(employees)
        var deptList []string
        for dept := range deptStats {
                deptList = append(deptList, dept)
        }

        contextMessage := fmt.Sprintf(
                "You are a helpful HR assistant for an HCM (Human Capital Management) system. "+
                "The company has %d employees across these departments: %s. "+
                "You can answer general HR policy questions, explain features, and provide guidance. "+
                "Do not make up specific employee information. Keep responses professional and helpful.",
                len(employees), strings.Join(deptList, ", "),
        )

        client := getOpenAIClient()
        completion, err := client.Chat.Completions.New(c.Request.Context(), openai.ChatCompletionNewParams{
                Messages: []openai.ChatCompletionMessageParamUnion{
                        openai.SystemMessage(contextMessage),
                        openai.UserMessage(input.Message),
                },
                Model: openai.ChatModelGPT4oMini,
        })

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("OpenAI API error: %v", err)})
                return
        }

        if len(completion.Choices) == 0 {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "No response from OpenAI"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "response": completion.Choices[0].Message.Content,
        })
}
