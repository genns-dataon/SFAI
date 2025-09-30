# Overview

This is an Employee Management System built with a React + Vite frontend and a backend API. The application handles core HR functions including employee records management with hierarchical reporting structure, attendance tracking, leave management, salary processing, organization chart visualization, and includes an AI chat assistant feature. The system uses a client-server architecture with REST API communication.

# User Preferences

Preferred communication style: Simple, everyday language.

# System Architecture

## Frontend Architecture

**Technology Stack:**
- **Framework**: React 19.1.1 with Vite 7.1.7 as the build tool
- **Routing**: React Router DOM 7.9.3 for client-side navigation
- **UI Library**: Ant Design (antd) for professional enterprise UI components
- **Styling**: Tailwind CSS 4.1.13 with PostCSS for utility-first styling
- **HTTP Client**: Axios 1.12.2 for API communication
- **Icons**: Ant Design Icons (@ant-design/icons) for consistent iconography
- **Date Utilities**: Day.js for Ant Design DatePicker and RangePicker components
- **Org Chart**: react-organizational-chart for hierarchical employee visualization

**Design Decisions:**
- Chose Vite over Create React App for faster build times and better development experience with HMR (Hot Module Replacement)
- Implemented Ant Design UI library for consistent, professional enterprise design system with pre-built components (Tables, Forms, Cards, Modals, etc.)
- Maintained Tailwind CSS for custom styling and layout utilities alongside Ant Design components
- Used Axios over fetch API for better error handling, request/response interceptors, and cleaner syntax

**Responsive Design:**
- **Collapsible Sidebar**: Sidebar automatically collapses on screens < 768px for mobile optimization
- **Manual Toggle**: Always-visible toggle button (hamburger icon) for user control
- **Icon-Only Mode**: Sidebar shows only icons when collapsed (80px on desktop, hidden on mobile)
- **Horizontal Scrolling**: Tables enable horizontal scroll on mobile to prevent data cutoff
- **Responsive Modals**: All modals use responsive width (90% with max-width constraints)
- **Adaptive Layouts**: Employee detail views display single-column layout on mobile devices
- **Mobile Padding**: Reduced padding (16px) on mobile for better space utilization

**Development Configuration:**
- ESLint configured with React Hooks and React Refresh plugins for code quality
- Server configured to run on port 5000 with host `0.0.0.0` for broader network access
- Supports Replit deployment with dynamic API URL resolution

## Backend Architecture

**Technology Stack:**
- **Framework**: Go with Gin web framework
- **ORM**: GORM for database operations
- **Authentication**: JWT with bcrypt password hashing
- **Database**: PostgreSQL with 7 tables (employees, departments, attendances, leave_requests, salaries, payslips, users)

**API Structure:**
The frontend communicates with a RESTful API at `/api` base path with the following endpoints:

1. **Authentication** (`/api/auth`)
   - POST /signup - Register new user
   - POST /login - Authenticate and get JWT token
   - GET /me - Get current user info (protected)

2. **Employee Management** (`/api/employees`)
   - GET all employees with manager relationships (protected)
   - GET employee by ID with manager and direct reports (protected)
   - POST create employee with 26 comprehensive fields (protected)
   - PUT update employee with manager hierarchy validation (protected)
   - **Manager Hierarchy**: Self-referential employee relationships with validation to prevent self-reporting
   - **Comprehensive Employee Data**: Supports 20 additional fields beyond basic info:
     - Personal & Identification (5 fields)
     - Employment & Job Details (5 fields)
     - Compensation & Benefits (5 fields)
     - Performance & Development (5 fields)
   - **UI Features**:
     - Clickable employee names open detailed view modal with all information
     - Tabbed edit form organizing fields into 5 sections for better UX
     - Detail view displays all fields categorized with visual organization

3. **Attendance Tracking** (`/api/attendance`)
   - GET all attendance records (protected)
   - POST clock-in records (protected)
   - POST clock-out records with duration calculation (protected)
   - **UI Features**: Clock In and Clock Out buttons on Attendance page with separate modals
   - **Duration Tracking**: Automatically calculates and displays total time worked (hours and minutes)

4. **Leave Management** (`/api/leave`)
   - GET all leave requests (protected)
   - POST create leave request (protected)

5. **Salary/Payroll** (`/api/salary`)
   - GET export salary data (protected)
   - POST generate payslip for employee (protected)

6. **AI Chat Assistant** (`/api/chat`)
   - POST send message and receive AI-powered response (protected, powered by OpenAI GPT-4o-mini)
   - **Security-First Design**: Employee data queries handled locally without sending PII to OpenAI
   - **Intelligent Query Routing**: 
     - Employee-related queries (list, department-specific, contact info) → Local database (no external API)
     - General HR policy questions → OpenAI with aggregated statistics only (no PII)
   - **Supported Employee Queries**: 
     - List all employees with basic information
     - Department-specific queries (Engineering, Sales, HR)
     - Employee contact information (email, phone)
     - New comprehensive employee fields:
       * Salary & compensation information (base salary, currency, pay frequency)
       * Skills & certifications
       * Performance ratings
       - Employment status statistics
       * Work location & arrangements (remote, hybrid, onsite)
       * Employment types (full-time, part-time, contract)
       * Benefits eligibility
   - **Attendance Recording via Chat**:
     - **Clock In**: Keywords like "clock in", "check in", "start work", "mark attendance"
     - **Clock Out**: Keywords like "clock out", "end work", "leaving", "done for the day"
     - Authenticates user and links to employee record via user_id
     - Prevents duplicate clock-ins for the same day
     - Calculates and displays total time worked when clocking out
     - **For Managers**: Record attendance for others with "record attendance for [name]", "clock in for [name]"
     - Conversational prompting asks for employee name if not provided
   - **Leave Request Creation via Chat**:
     - Creates actual leave requests directly from chat (not just guidance)
     - Keywords: "request leave", "apply for leave", "I want leave for tomorrow"
     - **Date Parsing**: Understands "tomorrow", "today", "next week" (Mon-Fri)
     - **Leave Types**: Detects "sick", "vacation", "personal", "emergency" from message
     - Conversational prompting asks for missing information (dates, leave type)
     - Returns confirmation with leave details and pending status

**API Client Design:**
- Centralized axios instance with base configuration
- Environment-aware API URL resolution (local development vs Replit deployment)
- Replit hostname detection replaces '-5000' with '-8080' for proper backend URL
- Modular API service layer with domain-specific exports (employeeAPI, attendanceAPI, etc.)
- Automatic JSON content-type headers
- JWT token stored in localStorage and sent in Authorization header

## Data Storage

The system uses PostgreSQL database with GORM to persist:
- **Employees**: Comprehensive employee records with 26 fields including:
  - Basic Information: name, email, job title, department, manager, hire date
  - Personal & Identification: employee number, date of birth, national ID, tax ID, marital status
  - Employment & Job Details: employment type, employment status, job level, work location, work arrangement
  - Compensation & Benefits: base salary, pay frequency, currency, bank account, benefit eligibility
  - Performance & Development: probation end date, performance rating, skills, training completed, career notes
- **Departments**: Organizational units
- **Attendances**: Clock-in/out records with timestamps
- **Leave Requests**: Employee leave applications with status tracking
- **Salaries**: Employee compensation records
- **Payslips**: Generated payroll documents
- **Users**: Authentication credentials with bcrypt-hashed passwords

**Database Seeding:**
- 3 departments pre-configured (Engineering, Human Resources, Sales)
- 10 employees with realistic data across all departments

## Authentication & Authorization

**Comprehensive JWT-based Authentication System:**

**Backend Security:**
- **Username-based Login**: Users authenticate with username (not email) and password
- **Password Hashing**: Bcrypt with DefaultCost (10) for all passwords
- **JWT Token Management**: 24-hour token expiration with JWT_SECRET validation
- **User-Employee Linking**: Users table with unique username field; employees link to users via UserID
- **Protected Routes**: JWT middleware validates tokens on all protected endpoints
- **Environment Security**: JWT_SECRET required in environment variables; server refuses to start without it

**Frontend Security:**
- **Login Page**: Professional Ant Design form with username/password fields and autocomplete attributes
- **Protected Routes**: ProtectedRoute component validates token via /api/me on every access
- **Token Verification**: Calls backend /api/me to verify token validity (catches expired/forged tokens)
- **Loading States**: Displays spinner while verifying authentication
- **Auto-Logout**: 401 responses trigger automatic token cleanup and redirect to login
- **Centralized API**: Axios interceptors automatically attach JWT token to all requests
- **Session Management**: Token and user data stored in localStorage, cleared on logout

**User Experience:**
- **Login Flow**: Username/password → JWT token → redirect to dashboard
- **Logout Flow**: Clear storage → redirect to login page
- **User Dropdown**: Shows username in header with logout option
- **Test Accounts**: 10 seeded accounts (alice, bob, carol, david, emma, frank, grace, henry, iris, jack) with password "password"

**Security Features:**
- JWT_SECRET validation prevents token generation without proper secret
- Token validation on every protected route prevents unauthorized access
- Expired/invalid tokens automatically cleaned up
- 401 responses trigger immediate logout
- Username-based authentication (more secure than email-based)
- Centralized API client prevents hardcoded URLs and CORS issues

## Environment Configuration

**Multi-Environment Support:**
- **Development**: Uses `http://localhost:8080/api`
- **Replit**: Dynamically constructs API URL by detecting 'replit.dev' hostname and replacing '-5000' with '-8080' to create proper backend URL (e.g., https://uuid-8080.spock.replit.dev/api)
- **Custom**: Supports `VITE_API_URL` environment variable override

**Port Configuration:**
- Frontend runs on port 5000 (accessible via https://uuid-5000.spock.replit.dev on Replit)
- Backend runs on port 8080 (accessible via https://uuid-8080.spock.replit.dev on Replit)

This approach allows seamless deployment across different hosting environments without code changes.

# External Dependencies

## Third-Party Libraries

**Frontend UI & UX:**
- `react` (19.1.1) - Core UI library
- `react-dom` (19.1.1) - React DOM rendering
- `react-router-dom` (7.9.3) - Client-side routing
- `antd` - Ant Design component library for professional enterprise UI
- `@ant-design/icons` - Ant Design icon library
- `dayjs` - Date utility library for Ant Design DatePicker
- `tailwindcss` (4.1.13) - Utility-first CSS framework
- `autoprefixer` (10.4.21) - CSS vendor prefixing
- `postcss` (8.5.6) - CSS transformation

**HTTP Communication:**
- `axios` (1.12.2) - Promise-based HTTP client for API requests

**Development Tools:**
- `vite` (7.1.7) - Build tool and dev server
- `@vitejs/plugin-react` (5.0.3) - React support for Vite
- `eslint` (9.36.0) - Code linting
- `@tailwindcss/postcss` (4.1.13) - Tailwind CSS processing

## Backend Dependencies

**Go Modules:**
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM for database operations
- `gorm.io/driver/postgres` - PostgreSQL driver for GORM
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/joho/godotenv` - Environment variable management
- `github.com/openai/openai-go/v2` - Official OpenAI SDK for Go

**Go Version:**
- Go 1.24 (upgraded from 1.19 for OpenAI SDK compatibility)

**Services:**
- Employee management with CRUD operations
- Attendance tracking with clock-in/out functionality
- Leave management with request/approval workflow
- Salary/payroll processing with payslip generation
- AI-powered chat assistant using OpenAI GPT-4o-mini for intelligent HR query responses

## Hosting & Deployment

**Supported Platforms:**
- Local development environment
- Replit cloud platform (with automatic URL resolution)
- Any environment supporting Node.js and static file serving

The system uses intelligent hostname detection to automatically configure API endpoints based on the deployment environment.