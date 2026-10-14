package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v5" // Note: Requires go get github.com/golang-jwt/jwt/v5
)

// JWTAuthMiddleware verifies Keycloak RSA-signed JWT tokens
// In a real production scenario, this will fetch the Keycloak JWKS (JSON Web Key Set)
// and verify the RSA signature of the token to authenticate the user and get their TOTP claims.
func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or Invalid Authorization Header (Expected Bearer Token)",
			})
		}

		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		/* 
		// Real Keycloak Validation Logic:
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Fetch RSA Public Key from Keycloak: 
			// http://localhost:8080/realms/master/protocol/openid-connect/certs
			return rsaPublicKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT Token. Please sign in via Keycloak with TOTP."})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["sub"])
		c.Locals("username", claims["preferred_username"])
		*/

		// For the FOSS Scaffold, we pass the user so the rest of the app doesn't break
		// Replace this with the real parsed claims when go mod is downloaded.
		c.Locals("username", "admin")

		return c.Next()
	}
}
