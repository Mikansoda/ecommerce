package middleware

import (
	"net/http"
	"strings"
	"time"

	"ecommerce/service"

	"github.com/gin-gonic/gin"
)

func Auth(requiredRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(requiredRoles))
	for _, r := range requiredRoles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing bearer token",
			})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		// Cek blacklist
		if exp, ok := service.AccessBlacklistLookup(token); ok {
			if time.Now().Before(exp) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token blacklisted",
				})
				return
			}
		}

		claims, err := service.ParseAccessForMiddleware(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}
		// inject to context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// role check
		if len(roleSet) > 0 {
			if _, ok := roleSet[claims.Role]; !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Forbidden",
				})
				return
			}
		}
		c.Next()
	}
}
