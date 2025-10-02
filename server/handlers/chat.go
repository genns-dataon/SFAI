package handlers

import (
        "context"
        "encoding/json"
        "fmt"
        "net/http"
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

// handleChatWithAI uses OpenAI function calling to intelligently handle all chatbot operations
func handleChatWithAI(userMessage string, history []map[string]string, verbose bool, userID interface{}) (string, []string, error) {
        client := getOpenAIClient()
        ctx := context.Background()
        var verboseSteps []string
        
        // Define available functions for the AI to call
        tools := []openai.ChatCompletionToolUnionParam{
                // Employee Information Functions
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "list_all_employees",
                        Description: openai.String("Get a list of all employees with their basic information including ID, name, email, job title, and department"),
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "get_employees_by_department",
                        Description: openai.String("Get employees from a specific department"),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "department": map[string]interface{}{
                                                "type":        "string",
                                                "description": "The department name (e.g., Engineering, Sales, HR, Human Resources)",
                                        },
                                },
                                "required": []string{"department"},
                        },
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "get_employee_reporting_structure",
                        Description: openai.String("Get the reporting structure for an employee - who they report to and who reports to them"),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "employee_name": map[string]interface{}{
                                                "type":        "string",
                                                "description": "The employee's name (first name, last name, or full name)",
                                        },
                                },
                                "required": []string{"employee_name"},
                        },
                }),
                // Attendance Functions
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "clock_in",
                        Description: openai.String("Record when the current user starts work (clock in). Use when user wants to mark their arrival, start work, or check in."),
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "clock_out",
                        Description: openai.String("Record when the current user ends work (clock out). Use when user wants to mark their departure, end work, or leave."),
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "record_attendance_for_employee",
                        Description: openai.String("Record attendance (clock in or clock out) for another employee. Managers use this to record attendance for their team members."),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "employee_name": map[string]interface{}{
                                                "type":        "string",
                                                "description": "The employee's name to record attendance for",
                                        },
                                        "action": map[string]interface{}{
                                                "type":        "string",
                                                "enum":        []string{"clock_in", "clock_out"},
                                                "description": "Whether to clock in or clock out",
                                        },
                                },
                                "required": []string{"employee_name", "action"},
                        },
                }),
                // Leave Request Functions
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "create_leave_request",
                        Description: openai.String("Create a leave request for the current user. Use when user wants to request time off, vacation, sick leave, etc."),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "start_date": map[string]interface{}{
                                                "type":        "string",
                                                "description": "Start date in YYYY-MM-DD format. Can derive from 'tomorrow', 'today', 'next week', etc.",
                                        },
                                        "end_date": map[string]interface{}{
                                                "type":        "string",
                                                "description": "End date in YYYY-MM-DD format. Can derive from duration like '2 days', '1 week', etc.",
                                        },
                                        "leave_type": map[string]interface{}{
                                                "type":        "string",
                                                "enum":        []string{"Vacation", "Sick Leave", "Personal", "Emergency"},
                                                "description": "Type of leave being requested",
                                        },
                                        "reason": map[string]interface{}{
                                                "type":        "string",
                                                "description": "Optional reason for the leave request",
                                        },
                                },
                                "required": []string{"start_date", "end_date", "leave_type"},
                        },
                }),
                // Additional Employee Information Functions
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "get_employee_details",
                        Description: openai.String("Get detailed information about a specific employee including email, phone, job title, department, manager, etc."),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "employee_name": map[string]interface{}{
                                                "type":        "string",
                                                "description": "The employee's name (first name, last name, or full name)",
                                        },
                                },
                                "required": []string{"employee_name"},
                        },
                }),
                // Leave Request Query Functions
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "list_leave_requests",
                        Description: openai.String("Get a list of leave requests, optionally filtered by month or date range"),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "month": map[string]interface{}{
                                                "type":        "string",
                                                "description": "Optional month filter in format 'YYYY-MM' (e.g., '2025-09' for September 2025) or 'this month', 'current month'",
                                        },
                                },
                        },
                }),
                // Attendance Query Functions
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "list_todays_attendance",
                        Description: openai.String("Get a list of employees who have clocked in today (who came to work today)"),
                }),
                // Additional Analytics & Filtering Functions
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "get_employees_by_work_location",
                        Description: openai.String("Get employees filtered by work location (e.g., San Francisco Office, New York Office, etc.)"),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "location": map[string]interface{}{
                                                "type":        "string",
                                                "description": "The work location to filter by",
                                        },
                                },
                                "required": []string{"location"},
                        },
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "get_employees_with_tenure",
                        Description: openai.String("Get list of employees with their years of service/tenure. Useful for 'how long have they worked here' questions."),
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "get_employee_salaries",
                        Description: openai.String("Get salary information for all employees (base salary, currency, pay frequency)"),
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "get_my_salary",
                        Description: openai.String("Get the current logged-in user's salary information. Use when user asks 'my salary', 'how much do I make', 'what's my pay', etc."),
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "count_employees_by_type",
                        Description: openai.String("Count employees by employment type (full-time, part-time, contract, etc.)"),
                        Parameters: openai.FunctionParameters{
                                "type": "object",
                                "properties": map[string]interface{}{
                                        "employment_type": map[string]interface{}{
                                                "type":        "string",
                                                "description": "The employment type to count (full-time, part-time, contract, intern, or 'all' for breakdown)",
                                        },
                                },
                                "required": []string{"employment_type"},
                        },
                }),
                openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
                        Name:        "count_total_employees",
                        Description: openai.String("Get the total count of employees. Use this when user asks 'how many employees' without listing them."),
                }),
        }
        
        // Load chatbot settings and build system prompt
        var settings []models.ChatbotSettings
        database.DB.Find(&settings)
        
        systemPrompt := "You are an intelligent HR assistant that helps employees with various tasks including: viewing employee information, managing attendance (clock in/out), and submitting leave requests. Use the available functions to help users. Never make up information. Always be helpful and professional."
        
        // Append settings to system prompt
        if len(settings) > 0 {
                systemPrompt += "\n\n📋 Additional Configuration & Helpers:\n"
                for _, setting := range settings {
                        systemPrompt += fmt.Sprintf("\n• %s: %s", setting.Key, setting.Value)
                        if setting.Description != "" {
                                systemPrompt += fmt.Sprintf(" (%s)", setting.Description)
                        }
                }
        }
        
        // Build messages array with conversation history
        messages := []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage(systemPrompt),
        }
        
        // Add conversation history
        for _, msg := range history {
                if msg["role"] == "user" {
                        messages = append(messages, openai.UserMessage(msg["content"]))
                } else if msg["role"] == "assistant" {
                        messages = append(messages, openai.AssistantMessage(msg["content"]))
                }
        }
        
        // Add current user message
        messages = append(messages, openai.UserMessage(userMessage))
        
        // Initial API call
        params := openai.ChatCompletionNewParams{
                Messages: messages,
                Tools:    tools,
                Model:    openai.ChatModelGPT4oMini,
        }
        
        if verbose {
                verboseSteps = append(verboseSteps, "🔍 Analyzing your question...")
        }
        
        response, err := client.Chat.Completions.New(ctx, params)
        if err != nil {
                return "", verboseSteps, err
        }
        
        // Check if AI wants to call a function
        toolCalls := response.Choices[0].Message.ToolCalls
        if len(toolCalls) == 0 {
                // No function call, return direct response
                if verbose {
                        verboseSteps = append(verboseSteps, "💬 Responding directly without database access")
                }
                if len(response.Choices) > 0 {
                        return response.Choices[0].Message.Content, verboseSteps, nil
                }
                return "", verboseSteps, fmt.Errorf("no response from AI")
        }
        
        // Add assistant's message to conversation
        params.Messages = append(params.Messages, response.Choices[0].Message.ToParam())
        
        // Execute each function call
        for _, toolCall := range toolCalls {
                functionName := toolCall.Function.Name
                argumentsJSON := toolCall.Function.Arguments
                
                switch functionName {
                case "list_all_employees":
                        var employees []models.Employee
                        if err := database.DB.Preload("Department").Find(&employees).Error; err != nil {
                                return "", verboseSteps, fmt.Errorf("database error: %v", err)
                        }
                        
                        result := "📋 Employee List:\n\n"
                        for _, emp := range employees {
                                deptName := "N/A"
                                if emp.Department != nil {
                                        deptName = emp.Department.Name
                                }
                                result += fmt.Sprintf("• ID: %d | %s (%s) - %s | Email: %s\n", 
                                        emp.ID, emp.Name, emp.JobTitle, deptName, emp.Email)
                        }
                        return result, verboseSteps, nil
                        
                case "get_employees_by_department":
                        var args struct {
                                Department string `json:"department"`
                        }
                        if err := json.Unmarshal([]byte(argumentsJSON), &args); err != nil {
                                return "", verboseSteps, fmt.Errorf("invalid arguments: %v", err)
                        }
                        
                        var employees []models.Employee
                        if err := database.DB.Preload("Department").Joins("JOIN departments ON departments.id = employees.department_id").
                                Where("LOWER(departments.name) LIKE ?", "%"+strings.ToLower(args.Department)+"%").
                                Find(&employees).Error; err != nil {
                                return "", verboseSteps, fmt.Errorf("database error: %v", err)
                        }
                        
                        if len(employees) == 0 {
                                return fmt.Sprintf("No employees found in %s department", args.Department), verboseSteps, nil
                        }
                        
                        result := fmt.Sprintf("👥 Employees in %s:\n\n", args.Department)
                        for _, emp := range employees {
                                result += fmt.Sprintf("• ID: %d | %s (%s) | Email: %s\n", 
                                        emp.ID, emp.Name, emp.JobTitle, emp.Email)
                        }
                        return result, verboseSteps, nil
                        
                case "get_employee_reporting_structure":
                        var args struct {
                                EmployeeName string `json:"employee_name"`
                        }
                        if err := json.Unmarshal([]byte(argumentsJSON), &args); err != nil {
                                return "", verboseSteps, fmt.Errorf("invalid arguments: %v", err)
                        }
                        
                        // Find the employee
                        var employee models.Employee
                        if err := database.DB.Preload("Manager").Preload("Department").
                                Where("LOWER(name) LIKE ?", "%"+strings.ToLower(args.EmployeeName)+"%").
                                First(&employee).Error; err != nil {
                                return fmt.Sprintf("Employee '%s' not found. Please check the spelling and try again.", args.EmployeeName), verboseSteps, nil
                        }
                        
                        result := fmt.Sprintf("📊 Reporting Structure for %s:\n\n", employee.Name)
                        result += fmt.Sprintf("• Job Title: %s\n", employee.JobTitle)
                        if employee.Department != nil {
                                result += fmt.Sprintf("• Department: %s\n", employee.Department.Name)
                        }
                        
                        if employee.Manager != nil {
                                result += fmt.Sprintf("• Reports to: %s (%s)\n", employee.Manager.Name, employee.Manager.JobTitle)
                        } else {
                                result += "• Reports to: No one (Top level)\n"
                        }
                        
                        // Find direct reports
                        var directReports []models.Employee
                        if err := database.DB.Where("manager_id = ?", employee.ID).Find(&directReports).Error; err != nil {
                                return "", verboseSteps, fmt.Errorf("database error: %v", err)
                        }
                        
                        if len(directReports) > 0 {
                                result += fmt.Sprintf("• Direct reports (%d):\n", len(directReports))
                                for _, report := range directReports {
                                        result += fmt.Sprintf("  - %s (%s)\n", report.Name, report.JobTitle)
                                }
                        } else {
                                result += "• Direct reports: None\n"
                        }
                        
                        return result, verboseSteps, nil
                        
        case "clock_in":
                if userID == nil {
                        return "⚠️ You need to be logged in to clock in.", verboseSteps, nil
                }
                
                var employee models.Employee
                if err := database.DB.Where("user_id = ?", userID).First(&employee).Error; err != nil {
                        return "❌ I couldn't find your employee record. Please contact HR.", verboseSteps, nil
                }
                
                var existingAttendance models.Attendance
                err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE", employee.ID).First(&existingAttendance).Error
                if err == nil {
                        clockOutStatus := ""
                        if existingAttendance.ClockOut == nil {
                                clockOutStatus = " (not yet clocked out)"
                        } else {
                                clockOutStatus = fmt.Sprintf(" (clocked out at %s)", existingAttendance.ClockOut.Format("3:04 PM"))
                        }
                        return fmt.Sprintf("✅ You already clocked in today at %s%s.", 
                                existingAttendance.ClockIn.Format("3:04 PM"), clockOutStatus), verboseSteps, nil
                }
                
                attendance := models.Attendance{
                        EmployeeID: employee.ID,
                        Date:       time.Now(),
                        ClockIn:    time.Now(),
                }
                
                if err := database.DB.Create(&attendance).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("failed to record attendance: %v", err)
                }
                
                result := fmt.Sprintf("✅ Attendance recorded successfully!\n\n👋 Welcome, %s!\n⏰ Clock-in time: %s\n\nHave a productive day!", 
                        employee.Name, attendance.ClockIn.Format("3:04 PM"))
                return result, verboseSteps, nil
                
        case "clock_out":
                if userID == nil {
                        return "⚠️ You need to be logged in to clock out.", verboseSteps, nil
                }
                
                var employee models.Employee
                if err := database.DB.Where("user_id = ?", userID).First(&employee).Error; err != nil {
                        return "❌ I couldn't find your employee record. Please contact HR.", verboseSteps, nil
                }
                
                var attendance models.Attendance
                err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE AND clock_out IS NULL", employee.ID).First(&attendance).Error
                if err != nil {
                        return "❌ You don't have an active clock-in for today. Please clock in first!", verboseSteps, nil
                }
                
                now := time.Now()
                attendance.ClockOut = &now
                
                if err := database.DB.Save(&attendance).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("failed to update attendance: %v", err)
                }
                
                duration := now.Sub(attendance.ClockIn)
                hours := int(duration.Hours())
                minutes := int(duration.Minutes()) % 60
                
                result := fmt.Sprintf("✅ Successfully clocked out!\n\n👋 See you later, %s!\n⏰ Clock-out time: %s\n📊 Total time: %dh %dm\n\nHave a great evening!", 
                        employee.Name, now.Format("3:04 PM"), hours, minutes)
                return result, verboseSteps, nil
                
        case "record_attendance_for_employee":
                var args struct {
                        EmployeeName string `json:"employee_name"`
                        Action       string `json:"action"`
                }
                if err := json.Unmarshal([]byte(argumentsJSON), &args); err != nil {
                        return "", verboseSteps, fmt.Errorf("invalid arguments: %v", err)
                }
                
                // Authorization: Only managers can record attendance for others
                if userID == nil {
                        return "⚠️ You need to be logged in to perform this action.", verboseSteps, nil
                }
                
                var requestingEmployee models.Employee
                if err := database.DB.Preload("Reports").Where("user_id = ?", userID).First(&requestingEmployee).Error; err != nil {
                        return "❌ I couldn't find your employee record. Please contact HR.", verboseSteps, nil
                }
                
                // Check if user is a manager (has direct reports)
                if len(requestingEmployee.Reports) == 0 {
                        return "⚠️ Only managers can record attendance for other employees. This action requires manager privileges.", verboseSteps, nil
                }
                
                var targetEmployee models.Employee
                if err := database.DB.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(args.EmployeeName)+"%").
                        First(&targetEmployee).Error; err != nil {
                        return fmt.Sprintf("Employee '%s' not found. Please check the spelling and try again.", args.EmployeeName), verboseSteps, nil
                }
                
                // Verify the target employee reports to the requesting manager
                if targetEmployee.ManagerID == nil || *targetEmployee.ManagerID != requestingEmployee.ID {
                        return fmt.Sprintf("⚠️ You can only record attendance for your direct reports. %s does not report to you.", targetEmployee.Name), verboseSteps, nil
                }
                
                if args.Action == "clock_out" {
                        var attendance models.Attendance
                        err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE AND clock_out IS NULL", targetEmployee.ID).First(&attendance).Error
                        if err != nil {
                                return fmt.Sprintf("❌ %s doesn't have an active clock-in for today.", targetEmployee.Name), verboseSteps, nil
                        }
                        
                        now := time.Now()
                        attendance.ClockOut = &now
                        
                        if err := database.DB.Save(&attendance).Error; err != nil {
                                return "", verboseSteps, fmt.Errorf("failed to update attendance: %v", err)
                        }
                        
                        duration := now.Sub(attendance.ClockIn)
                        hours := int(duration.Hours())
                        minutes := int(duration.Minutes()) % 60
                        
                        result := fmt.Sprintf("✅ Clock-out recorded for %s!\n\n⏰ Clock-out time: %s\n📊 Total time: %dh %dm", 
                                targetEmployee.Name, now.Format("3:04 PM"), hours, minutes)
                        return result, verboseSteps, nil
                } else {
                        var existingAttendance models.Attendance
                        err := database.DB.Where("employee_id = ? AND DATE(date) = CURRENT_DATE", targetEmployee.ID).First(&existingAttendance).Error
                        if err == nil {
                                return fmt.Sprintf("✅ %s already clocked in today at %s.", 
                                        targetEmployee.Name, existingAttendance.ClockIn.Format("3:04 PM")), verboseSteps, nil
                        }
                        
                        attendance := models.Attendance{
                                EmployeeID: targetEmployee.ID,
                                Date:       time.Now(),
                                ClockIn:    time.Now(),
                        }
                        
                        if err := database.DB.Create(&attendance).Error; err != nil {
                                return "", verboseSteps, fmt.Errorf("failed to record attendance: %v", err)
                        }
                        
                        result := fmt.Sprintf("✅ Attendance recorded for %s!\n\n⏰ Clock-in time: %s", 
                                targetEmployee.Name, attendance.ClockIn.Format("3:04 PM"))
                        return result, verboseSteps, nil
                }
                
        case "create_leave_request":
                var args struct {
                        StartDate string `json:"start_date"`
                        EndDate   string `json:"end_date"`
                        LeaveType string `json:"leave_type"`
                }
                if err := json.Unmarshal([]byte(argumentsJSON), &args); err != nil {
                        return "", verboseSteps, fmt.Errorf("invalid arguments: %v", err)
                }
                
                if userID == nil {
                        return "⚠️ You need to be logged in to create a leave request.", verboseSteps, nil
                }
                
                var employee models.Employee
                if err := database.DB.Where("user_id = ?", userID).First(&employee).Error; err != nil {
                        return "❌ I couldn't find your employee record. Please contact HR.", verboseSteps, nil
                }
                
                startDate, err := time.Parse("2006-01-02", args.StartDate)
                if err != nil {
                        return "❌ Invalid start date format. Please use YYYY-MM-DD format (e.g., 2025-10-15).", verboseSteps, nil
                }
                
                endDate, err := time.Parse("2006-01-02", args.EndDate)
                if err != nil {
                        return "❌ Invalid end date format. Please use YYYY-MM-DD format (e.g., 2025-10-20).", verboseSteps, nil
                }
                
                if endDate.Before(startDate) {
                        return "❌ End date cannot be before start date.", verboseSteps, nil
                }
                
                days := int(endDate.Sub(startDate).Hours()/24) + 1
                
                leaveRequest := models.LeaveRequest{
                        EmployeeID: employee.ID,
                        StartDate:  startDate,
                        EndDate:    endDate,
                        LeaveType:  args.LeaveType,
                        Status:     "Pending",
                }
                
                if err := database.DB.Create(&leaveRequest).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("failed to create leave request: %v", err)
                }
                
                result := fmt.Sprintf("✅ Leave request submitted successfully!\n\n"+
                        "📋 Request Details:\n"+
                        "• Employee: %s\n"+
                        "• Type: %s\n"+
                        "• Start Date: %s\n"+
                        "• End Date: %s\n"+
                        "• Duration: %d day(s)\n"+
                        "• Status: Pending\n\n"+
                        "Your manager will review your request soon.",
                        employee.Name, args.LeaveType, 
                        startDate.Format("Jan 02, 2006"), 
                        endDate.Format("Jan 02, 2006"), 
                        days)
                return result, verboseSteps, nil
                
        case "get_employee_details":
                var args struct {
                        EmployeeName string `json:"employee_name"`
                }
                if err := json.Unmarshal([]byte(argumentsJSON), &args); err != nil {
                        return "", verboseSteps, fmt.Errorf("invalid arguments: %v", err)
                }
                
                var employee models.Employee
                if err := database.DB.Preload("Department").Preload("Manager").
                        Where("LOWER(name) LIKE ?", "%"+strings.ToLower(args.EmployeeName)+"%").
                        First(&employee).Error; err != nil {
                        return fmt.Sprintf("❌ Employee '%s' not found. Please check the spelling and try again.", args.EmployeeName), verboseSteps, nil
                }
                
                result := fmt.Sprintf("👤 Employee Details for %s:\n\n", employee.Name)
                result += fmt.Sprintf("• ID: %d\n", employee.ID)
                result += fmt.Sprintf("• Email: %s\n", employee.Email)
                if employee.JobTitle != "" {
                        result += fmt.Sprintf("• Job Title: %s\n", employee.JobTitle)
                }
                if employee.Department != nil {
                        result += fmt.Sprintf("• Department: %s\n", employee.Department.Name)
                }
                if employee.Manager != nil {
                        result += fmt.Sprintf("• Reports to: %s\n", employee.Manager.Name)
                }
                if !employee.HireDate.IsZero() {
                        result += fmt.Sprintf("• Hire Date: %s\n", employee.HireDate.Format("Jan 02, 2006"))
                }
                if employee.WorkLocation != "" {
                        result += fmt.Sprintf("• Work Location: %s\n", employee.WorkLocation)
                }
                return result, verboseSteps, nil
                
        case "list_leave_requests":
                var args struct {
                        Month string `json:"month"`
                }
                json.Unmarshal([]byte(argumentsJSON), &args)
                
                var leaveRequests []models.LeaveRequest
                query := database.DB.Preload("Employee")
                
                // Parse month filter if provided
                if args.Month != "" {
                        var year, month int
                        if args.Month == "this month" || args.Month == "current month" {
                                now := time.Now()
                                year = now.Year()
                                month = int(now.Month())
                        } else {
                                parsedTime, err := time.Parse("2006-01", args.Month)
                                if err == nil {
                                        year = parsedTime.Year()
                                        month = int(parsedTime.Month())
                                }
                        }
                        
                        if year > 0 && month > 0 {
                                startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
                                endOfMonth := startOfMonth.AddDate(0, 1, -1)
                                query = query.Where("start_date >= ? AND start_date <= ?", startOfMonth, endOfMonth)
                        }
                }
                
                if err := query.Find(&leaveRequests).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("database error: %v", err)
                }
                
                if len(leaveRequests) == 0 {
                        if args.Month != "" {
                                return fmt.Sprintf("📅 No leave requests found for %s.", args.Month), verboseSteps, nil
                        }
                        return "📅 No leave requests found.", verboseSteps, nil
                }
                
                result := "📅 Leave Requests:\n\n"
                for i, lr := range leaveRequests {
                        result += fmt.Sprintf("%d. %s\n", i+1, lr.Employee.Name)
                        result += fmt.Sprintf("   • Type: %s\n", lr.LeaveType)
                        result += fmt.Sprintf("   • Dates: %s to %s\n", 
                                lr.StartDate.Format("Jan 02"), lr.EndDate.Format("Jan 02, 2006"))
                        result += fmt.Sprintf("   • Status: %s\n", lr.Status)
                        if i < len(leaveRequests)-1 {
                                result += "\n"
                        }
                }
                return result, verboseSteps, nil
                
        case "list_todays_attendance":
                var attendances []models.Attendance
                if err := database.DB.Preload("Employee").
                        Where("DATE(date) = CURRENT_DATE").
                        Find(&attendances).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("database error: %v", err)
                }
                
                if len(attendances) == 0 {
                        return "📊 No one has clocked in today yet.", verboseSteps, nil
                }
                
                result := fmt.Sprintf("📊 Today's Attendance (%s):\n\n", time.Now().Format("Jan 02, 2006"))
                for i, att := range attendances {
                        status := "Clocked In"
                        timeInfo := fmt.Sprintf("at %s", att.ClockIn.Format("3:04 PM"))
                        
                        if att.ClockOut != nil {
                                status = "Clocked Out"
                                duration := att.ClockOut.Sub(att.ClockIn)
                                hours := int(duration.Hours())
                                minutes := int(duration.Minutes()) % 60
                                timeInfo = fmt.Sprintf("(%dh %dm worked)", hours, minutes)
                        }
                        
                        result += fmt.Sprintf("%d. %s - %s %s\n", i+1, att.Employee.Name, status, timeInfo)
                }
                return result, verboseSteps, nil
                
        case "get_employees_by_work_location":
                var args struct {
                        Location string `json:"location"`
                }
                if err := json.Unmarshal([]byte(argumentsJSON), &args); err != nil {
                        return "", verboseSteps, fmt.Errorf("invalid arguments: %v", err)
                }
                
                var employees []models.Employee
                if err := database.DB.Preload("Department").
                        Where("LOWER(work_location) LIKE ?", "%"+strings.ToLower(args.Location)+"%").
                        Find(&employees).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("database error: %v", err)
                }
                
                if len(employees) == 0 {
                        return fmt.Sprintf("No employees found at %s", args.Location), verboseSteps, nil
                }
                
                result := fmt.Sprintf("📍 Employees at %s:\n\n", args.Location)
                for _, emp := range employees {
                        deptName := "N/A"
                        if emp.Department != nil {
                                deptName = emp.Department.Name
                        }
                        result += fmt.Sprintf("• %s (%s) - %s | Location: %s\n", 
                                emp.Name, emp.JobTitle, deptName, emp.WorkLocation)
                }
                return result, verboseSteps, nil
                
        case "get_employees_with_tenure":
                var employees []models.Employee
                if err := database.DB.Preload("Department").Find(&employees).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("database error: %v", err)
                }
                
                type EmployeeTenure struct {
                        Employee      models.Employee
                        YearsOfService float64
                }
                
                var tenures []EmployeeTenure
                now := time.Now()
                
                for _, emp := range employees {
                        years := now.Sub(emp.HireDate).Hours() / (24 * 365.25)
                        tenures = append(tenures, EmployeeTenure{
                                Employee:      emp,
                                YearsOfService: years,
                        })
                }
                
                // Sort by tenure descending (longest first)
                for i := 0; i < len(tenures); i++ {
                        for j := i + 1; j < len(tenures); j++ {
                                if tenures[j].YearsOfService > tenures[i].YearsOfService {
                                        tenures[i], tenures[j] = tenures[j], tenures[i]
                                }
                        }
                }
                
                result := "📅 Employee Tenure (Longest to Newest):\n\n"
                for _, t := range tenures {
                        deptName := "N/A"
                        if t.Employee.Department != nil {
                                deptName = t.Employee.Department.Name
                        }
                        result += fmt.Sprintf("• %s (%s) - %s | %.1f years of service\n", 
                                t.Employee.Name, t.Employee.JobTitle, deptName, t.YearsOfService)
                }
                return result, verboseSteps, nil
                
        case "get_employee_salaries":
                var employees []models.Employee
                if err := database.DB.Preload("Department").Find(&employees).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("database error: %v", err)
                }
                
                result := "💰 Employee Salaries:\n\n"
                for _, emp := range employees {
                        deptName := "N/A"
                        if emp.Department != nil {
                                deptName = emp.Department.Name
                        }
                        
                        salaryStr := "Not specified"
                        if emp.BaseSalary > 0 {
                                currency := emp.Currency
                                if currency == "" {
                                        currency = "USD"
                                }
                                payFreq := emp.PayFrequency
                                if payFreq == "" {
                                        payFreq = "annually"
                                }
                                salaryStr = fmt.Sprintf("%.2f %s (%s)", emp.BaseSalary, currency, payFreq)
                        }
                        
                        result += fmt.Sprintf("• %s (%s) - %s | Salary: %s\n", 
                                emp.Name, emp.JobTitle, deptName, salaryStr)
                }
                return result, verboseSteps, nil
                
        case "get_my_salary":
                if userID == nil {
                        return "⚠️ You need to be logged in to view your salary.", verboseSteps, nil
                }
                
                var employee models.Employee
                if err := database.DB.Preload("Department").Where("user_id = ?", userID).First(&employee).Error; err != nil {
                        return "❌ I couldn't find your employee record. Please contact HR.", verboseSteps, nil
                }
                
                salaryStr := "Not specified"
                if employee.BaseSalary > 0 {
                        currency := employee.Currency
                        if currency == "" {
                                currency = "USD"
                        }
                        payFreq := employee.PayFrequency
                        if payFreq == "" {
                                payFreq = "annually"
                        }
                        salaryStr = fmt.Sprintf("%.2f %s (%s)", employee.BaseSalary, currency, payFreq)
                }
                
                deptName := "N/A"
                if employee.Department != nil {
                        deptName = employee.Department.Name
                }
                
                result := fmt.Sprintf("💰 Your Salary Information:\n\n• Name: %s\n• Job Title: %s\n• Department: %s\n• Salary: %s", 
                        employee.Name, employee.JobTitle, deptName, salaryStr)
                return result, verboseSteps, nil
                
        case "count_employees_by_type":
                var args struct {
                        EmploymentType string `json:"employment_type"`
                }
                if err := json.Unmarshal([]byte(argumentsJSON), &args); err != nil {
                        return "", verboseSteps, fmt.Errorf("invalid arguments: %v", err)
                }
                
                if strings.ToLower(args.EmploymentType) == "all" {
                        // Get breakdown of all types
                        type TypeCount struct {
                                Type  string
                                Count int
                        }
                        
                        rows, err := database.DB.Raw(`
                                SELECT employment_type as type, COUNT(*) as count 
                                FROM employees 
                                WHERE deleted_at IS NULL 
                                GROUP BY employment_type
                                ORDER BY count DESC
                        `).Rows()
                        if err != nil {
                                return "", verboseSteps, fmt.Errorf("database error: %v", err)
                        }
                        defer rows.Close()
                        
                        result := "📊 Employee Count by Employment Type:\n\n"
                        totalCount := 0
                        for rows.Next() {
                                var tc TypeCount
                                if err := rows.Scan(&tc.Type, &tc.Count); err != nil {
                                        return "", verboseSteps, fmt.Errorf("scan error: %v", err)
                                }
                                if tc.Type == "" {
                                        tc.Type = "Not specified"
                                }
                                result += fmt.Sprintf("• %s: %d employees\n", tc.Type, tc.Count)
                                totalCount += tc.Count
                        }
                        result += fmt.Sprintf("\n📈 Total: %d employees", totalCount)
                        return result, verboseSteps, nil
                } else {
                        // Count specific type
                        var count int64
                        if err := database.DB.Model(&models.Employee{}).
                                Where("LOWER(employment_type) = ?", strings.ToLower(args.EmploymentType)).
                                Count(&count).Error; err != nil {
                                return "", verboseSteps, fmt.Errorf("database error: %v", err)
                        }
                        return fmt.Sprintf("📊 There are %d %s employees.", count, args.EmploymentType), verboseSteps, nil
                }
                
        case "count_total_employees":
                var count int64
                if err := database.DB.Model(&models.Employee{}).Count(&count).Error; err != nil {
                        return "", verboseSteps, fmt.Errorf("database error: %v", err)
                }
                return fmt.Sprintf("📊 Total number of employees: %d", count), verboseSteps, nil
                
                default:
                        return "", verboseSteps, fmt.Errorf("unknown function: %s", functionName)
                }
        }
        
        return "", verboseSteps, nil
}

func Chat(c *gin.Context) {
        var input struct {
                Message  string `json:"message" binding:"required"`
                History  []map[string]string `json:"history"`
                Verbose  bool `json:"verbose"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        userID, _ := c.Get("userID")

        aiResponse, verboseSteps, err := handleChatWithAI(input.Message, input.History, input.Verbose, userID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                        "error": fmt.Sprintf("Failed to process request: %v", err),
                })
                return
        }

        response := gin.H{
                "response": aiResponse,
        }
        
        if input.Verbose && len(verboseSteps) > 0 {
                response["verbose_steps"] = verboseSteps
        }

        c.JSON(http.StatusOK, response)
}
