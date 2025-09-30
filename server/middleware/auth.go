package middleware

import (
        "net/http"
        "os"
        "strings"

        "github.com/gin-gonic/gin"
        "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
                authHeader := c.GetHeader("Authorization")
                if authHeader == "" {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
                        c.Abort()
                        return
                }

                bearerToken := strings.Split(authHeader, " ")
                if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
                        c.Abort()
                        return
                }

                tokenString := bearerToken[1]
                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                        return []byte(os.Getenv("JWT_SECRET")), nil
                })

                if err != nil || !token.Valid {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
                        c.Abort()
                        return
                }

                if claims, ok := token.Claims.(jwt.MapClaims); ok {
                        if userID, exists := claims["user_id"]; exists {
                                c.Set("user_id", uint(userID.(float64)))
                        }
                }

                c.Next()
        }
}
