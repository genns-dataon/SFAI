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

        r.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{"status": "healthy"})
        })

        api := r.Group("/api")
        {
                api.POST("/auth/signup", handlers.Signup)
                api.POST("/auth/login", handlers.Login)

                protected := api.Group("/")
                protected.Use(middleware.AuthMiddleware())
                {
                        protected.GET("/me", handlers.GetMe)

                        protected.GET("/employees", handlers.GetEmployees)
                        protected.GET("/employees/:id", handlers.GetEmployee)
                        protected.POST("/employees", handlers.CreateEmployee)
                        protected.PUT("/employees/:id", handlers.UpdateEmployee)

                        protected.POST("/attendance/clockin", handlers.ClockIn)
                        protected.POST("/attendance/clockout", handlers.ClockOut)
                        protected.GET("/attendance", handlers.GetAttendance)

                        protected.POST("/leave", handlers.CreateLeaveRequest)
                        protected.GET("/leave", handlers.GetLeaveRequests)

                        protected.GET("/salary/export", handlers.ExportSalary)
                        protected.POST("/salary/payslip", handlers.GeneratePayslip)

                        protected.POST("/chat", handlers.Chat)
                        
                        protected.POST("/feedback", handlers.CreateFeedback)
                        protected.GET("/feedback", handlers.GetAllFeedback)
                        
                        protected.GET("/settings", handlers.GetSettings)
                        protected.GET("/settings/:key", handlers.GetSetting)
                        protected.POST("/settings", handlers.UpsertSetting)
                        protected.DELETE("/settings/:key", handlers.DeleteSetting)
                }
        }

        r.Static("/assets", "../client/dist/assets")
        r.NoRoute(func(c *gin.Context) {
                c.File("../client/dist/index.html")
        })

        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }

        log.Printf("Server starting on port %s...", port)
        if err := r.Run("0.0.0.0:" + port); err != nil {
                log.Fatal("Failed to start server:", err)
        }
}
