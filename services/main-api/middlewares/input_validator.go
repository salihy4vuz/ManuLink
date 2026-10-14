package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// InputValidator sanitizes simple inputs (preventing basic XSS/SQL Injection)
func InputValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Basic sanitizer (just an example, in real life use robust library)
		body := string(c.Body())
		if strings.Contains(body, "<script>") || strings.Contains(strings.ToLower(body), "drop table") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz karakterler algılandı",
			})
		}
		return c.Next()
	}
}
