package handlers

import (
        "fmt"
        "net/http"
        "os"
        "time"

        "hcm-backend/database"
        "hcm-backend/models"

        "github.com/gin-gonic/gin"
        "github.com/golang-jwt/jwt/v5"
        "golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
}

func hashPassword(password string) (string, error) {
        bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        return string(bytes), err
}

func checkPassword(password, hash string) bool {
        err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
        return err == nil
}

func generateToken(userID uint) (string, error) {
        jwtSecret := os.Getenv("JWT_SECRET")
        if jwtSecret == "" {
                return "", fmt.Errorf("JWT_SECRET environment variable is required")
        }

        claims := jwt.MapClaims{
                "user_id": userID,
                "exp":     time.Now().Add(time.Hour * 24).Unix(),
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return token.SignedString([]byte(jwtSecret))
}

func Signup(c *gin.Context) {
        var req SignupRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        // Check for existing username
        var existingUser models.User
        if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
                c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
                return
        }

        // Check for existing email
        if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
                c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
                return
        }

        hashedPassword, err := hashPassword(req.Password)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
                return
        }

        user := models.User{
                Username: req.Username,
                Email:    req.Email,
                Password: hashedPassword,
        }

        if err := database.DB.Create(&user).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
                return
        }

        token, err := generateToken(user.ID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
                return
        }

        c.JSON(http.StatusCreated, gin.H{
                "token": token,
                "user": gin.H{
                        "id":       user.ID,
                        "username": user.Username,
                        "email":    user.Email,
                },
        })
}

func Login(c *gin.Context) {
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        var user models.User
        if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
                return
        }

        if !checkPassword(req.Password, user.Password) {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
                return
        }

        token, err := generateToken(user.ID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "token": token,
                "user": gin.H{
                        "id":       user.ID,
                        "username": user.Username,
                        "email":    user.Email,
                },
        })
}

func GetMe(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
                return
        }

        var user models.User
        if err := database.DB.First(&user, userID).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "id":       user.ID,
                "username": user.Username,
                "email":    user.Email,
        })
}
