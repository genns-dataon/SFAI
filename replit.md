# Overview

This is an Employee Management System built with a React + Vite frontend and a backend API. The application handles core HR functions including employee records management, attendance tracking, leave management, salary processing, and includes an AI chat assistant feature. The system uses a client-server architecture with REST API communication.

# User Preferences

Preferred communication style: Simple, everyday language.

# System Architecture

## Frontend Architecture

**Technology Stack:**
- **Framework**: React 19.1.1 with Vite 7.1.7 as the build tool
- **Routing**: React Router DOM 7.9.3 for client-side navigation
- **Styling**: Tailwind CSS 4.1.13 with PostCSS for utility-first styling
- **HTTP Client**: Axios 1.12.2 for API communication
- **Icons**: Lucide React for UI iconography

**Design Decisions:**
- Chose Vite over Create React App for faster build times and better development experience with HMR (Hot Module Replacement)
- Implemented utility-first CSS with Tailwind to maintain consistent styling and reduce CSS bloat
- Used Axios over fetch API for better error handling, request/response interceptors, and cleaner syntax

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
   - GET all employees (protected)
   - GET employee by ID (protected)
   - POST create employee (protected)
   - PUT update employee (protected)

3. **Attendance Tracking** (`/api/attendance`)
   - GET all attendance records (protected)
   - POST clock-in records (protected)

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
   - **Supported Employee Queries**: List all employees, department-specific queries (Engineering, Sales, HR), employee contact information

**API Client Design:**
- Centralized axios instance with base configuration
- Environment-aware API URL resolution (local development vs Replit deployment)
- Replit hostname detection replaces '-5000' with '-8080' for proper backend URL
- Modular API service layer with domain-specific exports (employeeAPI, attendanceAPI, etc.)
- Automatic JSON content-type headers
- JWT token stored in localStorage and sent in Authorization header

## Data Storage

The system uses PostgreSQL database with GORM to persist:
- **Employees**: Full employee records with department relationships
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

JWT-based authentication system with the following features:
- User registration and login endpoints
- Bcrypt password hashing for secure credential storage
- JWT tokens for session management
- Protected API routes requiring valid authentication tokens
- Token-based middleware protecting all employee, attendance, leave, and salary endpoints

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
- `lucide-react` (0.544.0) - Icon library
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