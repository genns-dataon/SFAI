# Overview

This Employee Management System is a comprehensive HR solution featuring a React + Vite frontend and a Go backend API. It supports essential HR functions including employee record management with hierarchical reporting, attendance tracking, leave management, salary processing, organization chart visualization, and an integrated AI chat assistant. The system aims to streamline HR operations, enhance employee data accessibility, and provide intelligent assistance for HR queries.

# User Preferences

Preferred communication style: Simple, everyday language.

# System Architecture

## UI/UX Decisions

The frontend is built with React 19.1.1 and Vite, leveraging Ant Design for a professional enterprise UI and Tailwind CSS 4.1.13 for custom styling. Key UI/UX features include responsive design elements like a collapsible sidebar, adaptive layouts for various screen sizes, and interactive components for data display and input. `react-organizational-chart` is used for visualizing the employee hierarchy.

## Technical Implementations

### Frontend

- **Technology Stack**: React 19.1.1, Vite 7.1.7, React Router DOM 7.9.3, Ant Design, Tailwind CSS 4.1.13, Axios 1.12.2, Day.js.
- **Design Decisions**: Prioritizes fast build times (Vite), consistent UI (Ant Design), and robust API communication (Axios).
- **Responsive Design**: Implements a collapsible sidebar, responsive tables, modals, and adaptive layouts for mobile optimization.

### Backend

- **Technology Stack**: Go with Gin web framework, GORM for ORM, PostgreSQL database, JWT for authentication, bcrypt for password hashing.
- **API Structure**: A RESTful API provides endpoints for:
    - **Authentication**: User registration, login, and current user retrieval.
    - **Employee Management**: CRUD operations for comprehensive employee records, including manager hierarchy validation and 26 detailed fields.
    - **Attendance Tracking**: Clock-in/out functionality with duration calculation.
    - **Leave Management**: Creation and retrieval of leave requests.
    - **Salary/Payroll**: Export salary data and generate payslips.
    - **AI Chat Assistant**: Integrates OpenAI GPT-4o-mini for natural language queries, leveraging OpenAI Function Calling for secure, local processing of employee data. It supports:
        - Employee queries (list all, by department, reporting structure, specific details like email)
        - Attendance management (self clock-in/out, manager recording for direct reports, view today's attendance)
        - Leave requests (creation with natural language date parsing, view requests by month)
        - Feedback system with thumbs up/down ratings and comment collection
        - Configurable settings (key-value pairs injected into system prompt)
    - **Feedback Management**: Stores user feedback on AI responses.
    - **Settings Management**: CRUD operations for chatbot configuration (note: unsecured by user request).

## System Design Choices

- **Client-Server Architecture**: Frontend communicates with the backend via REST API.
- **Data Model**: PostgreSQL database with 9 tables: employees, departments, attendances, leave_requests, salaries, payslips, users, chat_feedbacks, and chatbot_settings.
- **Authentication**: Robust JWT-based authentication with bcrypt hashing, token expiration, and protected routes. Frontend handles token storage (localStorage) and validation.
- **Environment Configuration**: Supports multi-environment deployment (development, Replit, custom) with dynamic API URL resolution. Frontend on port 5000, backend on port 8080.
- **AI Integration**: AI chat assistant uses OpenAI's Function Calling to interact with backend functions securely, preventing PII exposure to external AI services.

# External Dependencies

## Third-Party Libraries

- **Frontend**:
    - `react`, `react-dom`, `react-router-dom`: Core UI and routing.
    - `antd`, `@ant-design/icons`, `dayjs`: UI components, icons, and date utilities.
    - `tailwindcss`, `autoprefixer`, `postcss`: Styling.
    - `axios`: HTTP client.
    - `vite`, `@vitejs/plugin-react`, `eslint`: Development and build tools.
- **Backend (Go Modules)**:
    - `github.com/gin-gonic/gin`: Web framework.
    - `gorm.io/gorm`, `gorm.io/driver/postgres`: ORM and PostgreSQL driver.
    - `github.com/golang-jwt/jwt/v5`: JWT authentication.
    - `golang.org/x/crypto/bcrypt`: Password hashing.
    - `github.com/joho/godotenv`: Environment variables.
    - `github.com/openai/openai-go/v2`: OpenAI SDK.

## Services

- **Database**: PostgreSQL.
- **AI**: OpenAI GPT-4o-mini.

## Hosting & Deployment

- **Platforms**: Local development, Replit cloud platform.
- **Go Version**: 1.24.