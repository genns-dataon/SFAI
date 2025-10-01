# HCM API Documentation

## Base URL
```
http://your-domain:8080/api
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

---

## Authentication Endpoints

### Sign Up
Create a new user account.

**Endpoint:** `POST /api/auth/signup`

**Request Body:**
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response (200):**
```json
{
  "message": "User created successfully"
}
```

**Error Responses:**
- `400` - Invalid request data
- `500` - User already exists or server error

---

### Login
Authenticate and receive a JWT token.

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com"
  }
}
```

**Error Responses:**
- `400` - Missing credentials
- `401` - Invalid credentials
- `500` - Server error

---

## User Endpoints

### Get Current User
Get information about the currently logged-in user.

**Endpoint:** `GET /api/me`

**Headers:** Requires authentication

**Response (200):**
```json
{
  "id": 1,
  "username": "alice",
  "email": "alice@example.com"
}
```

**Error Responses:**
- `401` - Unauthorized
- `404` - User not found

---

## Employee Endpoints

### Get All Employees
Retrieve a list of all employees.

**Endpoint:** `GET /api/employees`

**Headers:** Requires authentication

**Response (200):**
```json
[
  {
    "id": 1,
    "name": "Alice Johnson",
    "email": "alice@company.com",
    "phone": "+1234567890",
    "job_title": "Software Engineer",
    "department": {
      "id": 1,
      "name": "Engineering"
    },
    "manager_id": 5,
    "hire_date": "2023-01-15T00:00:00Z",
    "employment_type": "full-time",
    "work_location": "San Francisco Office",
    "base_salary": 120000.00,
    "currency": "USD",
    "pay_frequency": "annually"
  }
]
```

---

### Get Single Employee
Retrieve details of a specific employee.

**Endpoint:** `GET /api/employees/:id`

**Headers:** Requires authentication

**URL Parameters:**
- `id` (integer) - Employee ID

**Response (200):**
```json
{
  "id": 1,
  "name": "Alice Johnson",
  "email": "alice@company.com",
  "phone": "+1234567890",
  "job_title": "Software Engineer",
  "department": {
    "id": 1,
    "name": "Engineering"
  },
  "manager": {
    "id": 5,
    "name": "Bob Manager"
  },
  "hire_date": "2023-01-15T00:00:00Z",
  "employment_type": "full-time",
  "work_location": "San Francisco Office",
  "base_salary": 120000.00,
  "currency": "USD",
  "pay_frequency": "annually"
}
```

**Error Responses:**
- `404` - Employee not found

---

### Create Employee
Create a new employee record.

**Endpoint:** `POST /api/employees`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@company.com",
  "phone": "+1234567890",
  "job_title": "Product Manager",
  "department_id": 2,
  "manager_id": 5,
  "hire_date": "2024-01-01",
  "employment_type": "full-time",
  "work_location": "New York Office",
  "base_salary": 130000.00,
  "currency": "USD",
  "pay_frequency": "annually"
}
```

**Response (201):**
```json
{
  "id": 15,
  "name": "John Doe",
  "email": "john@company.com",
  ...
}
```

**Error Responses:**
- `400` - Invalid request data
- `500` - Server error

---

### Update Employee
Update an existing employee record.

**Endpoint:** `PUT /api/employees/:id`

**Headers:** Requires authentication

**URL Parameters:**
- `id` (integer) - Employee ID

**Request Body:** (All fields optional)
```json
{
  "name": "John Doe Updated",
  "job_title": "Senior Product Manager",
  "base_salary": 150000.00
}
```

**Response (200):**
```json
{
  "id": 15,
  "name": "John Doe Updated",
  "job_title": "Senior Product Manager",
  ...
}
```

---

## Attendance Endpoints

### Clock In
Record employee arrival/start of work.

**Endpoint:** `POST /api/attendance/clockin`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "employee_id": 1
}
```

**Response (201):**
```json
{
  "id": 100,
  "employee_id": 1,
  "date": "2024-10-01T08:00:00Z",
  "clock_in": "2024-10-01T08:00:00Z",
  "clock_out": null
}
```

**Error Responses:**
- `400` - Already clocked in today
- `404` - Employee not found

---

### Clock Out
Record employee departure/end of work.

**Endpoint:** `POST /api/attendance/clockout`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "employee_id": 1
}
```

**Response (200):**
```json
{
  "id": 100,
  "employee_id": 1,
  "date": "2024-10-01T08:00:00Z",
  "clock_in": "2024-10-01T08:00:00Z",
  "clock_out": "2024-10-01T17:00:00Z"
}
```

**Error Responses:**
- `400` - No active clock-in found
- `404` - Employee not found

---

### Get Attendance Records
Retrieve attendance records.

**Endpoint:** `GET /api/attendance`

**Headers:** Requires authentication

**Response (200):**
```json
[
  {
    "id": 100,
    "employee_id": 1,
    "employee": {
      "id": 1,
      "name": "Alice Johnson"
    },
    "date": "2024-10-01T00:00:00Z",
    "clock_in": "2024-10-01T08:00:00Z",
    "clock_out": "2024-10-01T17:00:00Z"
  }
]
```

---

## Leave Request Endpoints

### Create Leave Request
Submit a new leave request.

**Endpoint:** `POST /api/leave`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "employee_id": 1,
  "leave_type": "Vacation",
  "start_date": "2024-12-20",
  "end_date": "2024-12-27",
  "reason": "Holiday vacation",
  "status": "pending"
}
```

**Leave Types:**
- `Vacation`
- `Sick Leave`
- `Personal`
- `Emergency`

**Response (201):**
```json
{
  "id": 25,
  "employee_id": 1,
  "leave_type": "Vacation",
  "start_date": "2024-12-20T00:00:00Z",
  "end_date": "2024-12-27T00:00:00Z",
  "reason": "Holiday vacation",
  "status": "pending"
}
```

---

### Get Leave Requests
Retrieve all leave requests.

**Endpoint:** `GET /api/leave`

**Headers:** Requires authentication

**Response (200):**
```json
[
  {
    "id": 25,
    "employee_id": 1,
    "employee": {
      "id": 1,
      "name": "Alice Johnson"
    },
    "leave_type": "Vacation",
    "start_date": "2024-12-20T00:00:00Z",
    "end_date": "2024-12-27T00:00:00Z",
    "reason": "Holiday vacation",
    "status": "pending"
  }
]
```

---

### Update Leave Status
Approve or reject a leave request.

**Endpoint:** `PUT /api/leave/:id`

**Headers:** Requires authentication

**URL Parameters:**
- `id` (integer) - Leave request ID

**Request Body:**
```json
{
  "status": "approved"
}
```

**Status Options:**
- `pending`
- `approved`
- `rejected`

**Response (200):**
```json
{
  "id": 25,
  "status": "approved",
  ...
}
```

---

## Salary & Payroll Endpoints

### Export Salary Data
Export salary information for all employees.

**Endpoint:** `GET /api/salary/export`

**Headers:** Requires authentication

**Response (200):**
```json
[
  {
    "id": 1,
    "employee_id": 1,
    "employee_name": "Alice Johnson",
    "base_salary": 120000.00,
    "currency": "USD",
    "pay_frequency": "annually",
    "bonus": 10000.00,
    "deductions": 5000.00
  }
]
```

---

### Generate Payslip
Generate a payslip for a specific employee.

**Endpoint:** `POST /api/salary/payslip`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "employee_id": 1
}
```

**Response (201):**
```json
{
  "id": 50,
  "employee_id": 1,
  "employee_name": "Alice Johnson",
  "period_start": "2024-10-01T00:00:00Z",
  "period_end": "2024-10-31T00:00:00Z",
  "gross_pay": 10000.00,
  "deductions": 2000.00,
  "net_pay": 8000.00
}
```

---

## AI Chatbot Endpoints

### Chat with AI Assistant
Send a message to the AI chatbot and receive a response.

**Endpoint:** `POST /api/chat`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "message": "How many employees are in Engineering?",
  "history": [
    {
      "role": "user",
      "content": "Hello"
    },
    {
      "role": "assistant",
      "content": "Hi! How can I help you today?"
    }
  ]
}
```

**Notes:**
- `history` is optional but recommended for multi-turn conversations
- The AI can handle queries about employees, attendance, leave requests, and salaries
- It supports natural language date parsing for leave requests

**Response (200):**
```json
{
  "response": "There are 15 employees in the Engineering department."
}
```

**Chatbot Capabilities:**
- Employee queries (list, search by department, get details)
- Attendance management (clock in/out, view records)
- Leave requests (create with natural language dates, view requests)
- Salary information (view your own salary, list all salaries)
- Employee analytics (count by type, tenure, location)

---

## Feedback Endpoints

### Submit Feedback
Submit feedback on AI chatbot responses.

**Endpoint:** `POST /api/feedback`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "question": "How many employees are there?",
  "response": "There are 50 employees in total.",
  "rating": "positive",
  "comment": "Very helpful response!"
}
```

**Rating Options:**
- `positive`
- `negative`
- `resolved`

**Response (201):**
```json
{
  "id": 10,
  "question": "How many employees are there?",
  "response": "There are 50 employees in total.",
  "rating": "positive",
  "comment": "Very helpful response!"
}
```

---

### Get All Feedback
Retrieve all feedback entries.

**Endpoint:** `GET /api/feedback`

**Headers:** Requires authentication

**Response (200):**
```json
[
  {
    "id": 10,
    "question": "How many employees are there?",
    "response": "There are 50 employees in total.",
    "rating": "positive",
    "comment": "Very helpful response!",
    "created_at": "2024-10-01T10:00:00Z"
  }
]
```

---

### Update Feedback
Update an existing feedback entry (e.g., change rating from negative to resolved).

**Endpoint:** `PUT /api/feedback/:id`

**Headers:** Requires authentication

**URL Parameters:**
- `id` (integer) - Feedback ID

**Request Body:**
```json
{
  "rating": "resolved",
  "comment": "Issue has been fixed"
}
```

**Response (200):**
```json
{
  "id": 10,
  "rating": "resolved",
  "comment": "Issue has been fixed",
  ...
}
```

---

## Settings Endpoints

### Get All Settings
Retrieve all chatbot configuration settings.

**Endpoint:** `GET /api/settings`

**Headers:** Requires authentication

**Response (200):**
```json
[
  {
    "key": "greeting_message",
    "value": "Welcome to HR Assistant!",
    "description": "Initial greeting shown to users"
  }
]
```

---

### Get Single Setting
Retrieve a specific setting by key.

**Endpoint:** `GET /api/settings/:key`

**Headers:** Requires authentication

**URL Parameters:**
- `key` (string) - Setting key

**Response (200):**
```json
{
  "key": "greeting_message",
  "value": "Welcome to HR Assistant!",
  "description": "Initial greeting shown to users"
}
```

---

### Create or Update Setting
Create a new setting or update an existing one.

**Endpoint:** `POST /api/settings`

**Headers:** Requires authentication

**Request Body:**
```json
{
  "key": "max_leave_days",
  "value": "30",
  "description": "Maximum leave days per year"
}
```

**Response (200):**
```json
{
  "key": "max_leave_days",
  "value": "30",
  "description": "Maximum leave days per year"
}
```

---

### Delete Setting
Delete a setting by key.

**Endpoint:** `DELETE /api/settings/:key`

**Headers:** Requires authentication

**URL Parameters:**
- `key` (string) - Setting key

**Response (200):**
```json
{
  "message": "Setting deleted successfully"
}
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request data"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Rate Limiting

Currently, there are no rate limits implemented. Consider implementing rate limiting for production use.

---

## CORS

CORS is enabled for all origins in development. Configure appropriately for production.

---

## Example Mobile Implementation

### Swift (iOS)
```swift
struct LoginRequest: Codable {
    let username: String
    let password: String
}

func login(username: String, password: String) async throws -> String {
    let url = URL(string: "http://your-domain:8080/api/auth/login")!
    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    
    let body = LoginRequest(username: username, password: password)
    request.httpBody = try JSONEncoder().encode(body)
    
    let (data, _) = try await URLSession.shared.data(for: request)
    let response = try JSONDecoder().decode(LoginResponse.self, from: data)
    return response.token
}
```

### Kotlin (Android)
```kotlin
suspend fun login(username: String, password: String): String {
    val client = OkHttpClient()
    val json = JSONObject()
    json.put("username", username)
    json.put("password", password)
    
    val body = json.toString().toRequestBody("application/json".toMediaType())
    val request = Request.Builder()
        .url("http://your-domain:8080/api/auth/login")
        .post(body)
        .build()
    
    val response = client.newCall(request).execute()
    val responseData = JSONObject(response.body?.string() ?: "")
    return responseData.getString("token")
}
```

### React Native
```javascript
async function login(username, password) {
    const response = await fetch('http://your-domain:8080/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
    });
    
    const data = await response.json();
    return data.token;
}
```

---

## Additional Notes

- All timestamps are in ISO 8601 format (UTC)
- All monetary values are represented as floating-point numbers
- The AI chatbot uses OpenAI GPT-4o-mini and supports function calling for database operations
- JWT tokens expire after a configured period (check server configuration)
- All protected endpoints return 401 if the token is invalid or expired
