package middleware

import (
	"net/http"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/gin-gonic/gin"
)

func RequireRoles(allowed ...models.UserRole) gin.HandlerFunc {
	allowedSet := make(map[models.UserRole]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		roleValue, exists := c.Get(ContextRoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user role"})
			return
		}

		role, ok := roleValue.(models.UserRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user role"})
			return
		}

		if _, ok := allowedSet[role]; !ok {
			// Admins inherit staff permissions.
			if role == models.RoleAdmin {
				if _, staffAllowed := allowedSet[models.RoleStaff]; staffAllowed {
					c.Next()
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
