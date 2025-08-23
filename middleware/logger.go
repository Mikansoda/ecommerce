package middleware

import (
	"log"
	"strings"

	"ecommerce/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ActionLogger(logService service.ActionLogService) gin.HandlerFunc {
	// helper function to detect entity type from route path
	detectEntityType := func(path string) string {
		switch {
		case strings.Contains(path, "/products"):
			return "products"
		case strings.Contains(path, "/orders"):
			return "orders"
		case strings.Contains(path, "/users"):
			return "users"
		case strings.Contains(path, "/categories"):
			return "categories"
		case strings.Contains(path, "/payments"):
			return "payments"
		case strings.Contains(path, "/product_images"):
			return "product_images"
		case strings.Contains(path, "/carts"):
			return "carts"
		case strings.Contains(path, "/addresses"):
			return "addresses"
		case strings.Contains(path, "/auth"):
			return "auth"
		default:
			return "unknown"
		}
	}

	return func(c *gin.Context) {
		c.Next()

		if c.Request.Method != "POST" && c.Request.Method != "PUT" &&
			c.Request.Method != "PATCH" && c.Request.Method != "DELETE" {
			return
		}

		var parsedActorID *uuid.UUID
		if actorIDVal, exists := c.Get("userID"); exists {
			if actorIDStr, ok := actorIDVal.(string); ok {
				if id, err := uuid.Parse(actorIDStr); err == nil {
					parsedActorID = &id
				}
			}
		}

		roleStr := "unknown"
		if roleVal, exists := c.Get("role"); exists {
			if r, ok := roleVal.(string); ok {
				roleStr = r
			}
		}

		path := c.FullPath()
		entityType := detectEntityType(path)

		entityID := ""
		if idParam := c.Param("productId"); idParam != "" {
			entityID = idParam
		} else if idParam := c.Param("id"); idParam != "" {
			entityID = idParam
		}

		methodToAction := map[string]string{
			"POST":   "create",
			"PUT":    "update",
			"PATCH":  "update",
			"DELETE": "delete",
		}
		action := methodToAction[c.Request.Method]

		if err := logService.Log(
			c.Request.Context(),
			roleStr,
			parsedActorID,
			action,
			entityType,
			entityID,
		); err != nil {
			log.Printf("Failed to save action to log: %v", err)
		}
	}
}
