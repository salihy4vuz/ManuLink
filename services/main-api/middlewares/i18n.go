package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

var translations = map[string]map[string]string{
	"en": {
		"error_too_many_requests": "Too many requests. Please try again later.",
		"error_service_unavailable": "Service temporarily unavailable.",
	},
	"tr": {
		"error_too_many_requests": "Çok fazla istek gönderdiniz. Lütfen daha sonra tekrar deneyin.",
		"error_service_unavailable": "Servis geçici olarak kullanılamıyor.",
	},
	"de": {
		"error_too_many_requests": "Zu viele Anfragen. Bitte versuchen Sie es später noch einmal.",
		"error_service_unavailable": "Dienst vorübergehend nicht verfügbar.",
	},
}

// I18nMiddleware sets the client's language preference
func I18nMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.Get("Accept-Language", "en")
		// Basic parsing (e.g., tr-TR, tr;q=0.9 -> tr)
		langPrefix := "en"
		if len(lang) >= 2 {
			parsedLang := strings.ToLower(lang[:2])
			if _, exists := translations[parsedLang]; exists {
				langPrefix = parsedLang
			}
		}
		
		c.Locals("lang", langPrefix)
		return c.Next()
	}
}

// T translates a key based on the context's language
func T(c *fiber.Ctx, key string) string {
	lang := c.Locals("lang")
	if langStr, ok := lang.(string); ok {
		if msg, exists := translations[langStr][key]; exists {
			return msg
		}
	}
	// Fallback to English
	return translations["en"][key]
}
