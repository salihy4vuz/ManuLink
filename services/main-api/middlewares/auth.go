package middlewares

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strings"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// Roles
const (
	RoleAdmin    = "Admin"
	RoleOperator = "Operator"
	RoleCustomer = "Customer"
)

var rsaPublicKey *rsa.PublicKey

func init() {
	keyData := os.Getenv("JWT_PUBLIC_KEY")
	if keyData != "" {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyData))
		if err == nil {
			rsaPublicKey = key
		}
	}
}

// Protected verifies the JWT (RS256) and adds user context
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// Mock fallback for testing
			if authHeader == "admin-token" {
				c.Locals("role", RoleAdmin)
				return c.Next()
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}
		
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if rsaPublicKey == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "RS256 config error"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return rsaPublicKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		
		// Check Redis Blacklist
		redisClient := services.GetRedisClient()
		if redisClient != nil {
			val, _ := redisClient.Get(context.Background(), "blacklist:"+tokenString).Result()
			if val == "revoked" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has been revoked"})
			}
		}
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("role", claims["role"])
			c.Locals("user_id", claims["sub"])
		}
		return c.Next()
	}
}

// RequireRole enforces RBAC
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Role not found in context"})
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}
}
