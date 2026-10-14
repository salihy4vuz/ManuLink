package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Dynamic Feature Flag Manager (Redis Backed)
// This completely replaces static config-based toggles for isolating ERP/HR modules.

// FeatureFlagConfig
type FeatureFlagConfig struct {
	RedisClient *redis.Client
}

// RequireFeature checks Redis if a specific feature is enabled globally or for a tenant
func RequireFeature(config FeatureFlagConfig, featureName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.RedisClient == nil {
			// Fallback if Redis is down: In a "Blue Ocean" CAD-to-MES scope,
			// we assume legacy features (HR, Accounting) are DISABLED by default.
			if featureName == "hr_module" || featureName == "accounting" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Feature Disabled (Fallback Mode)",
				})
			}
			return c.Next()
		}

		ctx := context.Background()
		// Check Redis key (e.g. "feature:hr_module")
		// "1" means enabled, "0" or missing means disabled
		val, err := config.RedisClient.Get(ctx, "feature:"+featureName).Result()

		// If key doesn't exist or is 0, block the request
		if err == redis.Nil || strings.TrimSpace(val) != "1" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This module is disabled as part of the CAD-to-MES strategic scope.",
				"feature": featureName,
			})
		}

		return c.Next()
	}
}
