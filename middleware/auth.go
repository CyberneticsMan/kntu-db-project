package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.Fields(authorization)
		tokenString := ""
		switch {
		case len(parts) == 2 && strings.EqualFold(parts[0], "bearer"):
			tokenString = parts[1]
		case len(parts) == 1:
			tokenString = parts[0]
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		if sub, ok := claims["sub"]; ok {
			c.Set(ContextUserIDKey, sub)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set(ContextRoleKey, models.UserRole(role))
		}

		c.Next()
	}
}
