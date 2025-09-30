package main

import (
        "log"
        "os"

        "hcm-backend/database"
        "hcm-backend/handlers"
        "hcm-backend/middleware"

        "github.com/gin-contrib/cors"
        "github.com/gin-gonic/gin"
        "github.com/joho/godotenv"
)

func main() {
        if err := godotenv.Load(); err != nil {
                log.Println("No .env file found, using environment variables")
        }

        database.Connect()
        database.Migrate()
        database.SeedData()

        r := gin.Default()

        r.Use(cors.New(cors.Config{
                AllowOriginFunc: func(origin string) bool {
                        return true
                },
                AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
                AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
                ExposeHeaders:    []string{"Content-Length"},
                AllowCredentials: true,
        }))

        api := r.Group("/api")
        {
                api.POST("/auth/signup", handlers.Signup)
                api.POST("/auth/login", handlers.Login)

                api.GET("/employees", handlers.GetEmployees)
                api.GET("/employees/:id", handlers.GetEmployee)
                api.POST("/employees", handlers.CreateEmployee)
                api.PUT("/employees/:id", handlers.UpdateEmployee)

                api.POST("/attendance/clockin", handlers.ClockIn)
                api.GET("/attendance", handlers.GetAttendance)

                api.POST("/leave", handlers.CreateLeaveRequest)
                api.GET("/leave", handlers.GetLeaveRequests)

                api.GET("/salary/export", handlers.ExportSalary)
                api.POST("/salary/payslip", handlers.GeneratePayslip)

                api.POST("/chat", handlers.Chat)

                protected := api.Group("/")
                protected.Use(middleware.AuthMiddleware())
                {
                        protected.GET("/me", handlers.GetMe)
                }
        }

        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }

        log.Printf("Server starting on port %s...", port)
        if err := r.Run("0.0.0.0:" + port); err != nil {
                log.Fatal("Failed to start server:", err)
        }
}
