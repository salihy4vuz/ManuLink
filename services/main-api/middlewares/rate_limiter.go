package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig holds configuration for the Redis-backed Sliding Window Rate Limiter
type RateLimiterConfig struct {
	RedisClient *redis.Client
	MaxRequests int
	Window      time.Duration
}

// NewRateLimiter creates a new Redis rate limiter middleware
func NewRateLimiter(config RateLimiterConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.RedisClient == nil {
			// If Redis is not configured, skip rate limiting (fallback)
			return c.Next()
		}

		ip := c.IP()
		if ip == "" {
			ip = "unknown_ip"
		}

		key := fmt.Sprintf("rate_limit:%s:%s", c.Path(), ip)

		ctx := context.Background()
		now := time.Now().UnixNano()
		windowStart := now - config.Window.Nanoseconds()

		// Sliding Window implementation using Redis ZSET
		pipe := config.RedisClient.TxPipeline()
		
		// 1. Remove timestamps outside the current window
		pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", windowStart))
		
		// 2. Add current timestamp
		pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
		
		// 3. Count requests in the window
		countCmd := pipe.ZCard(ctx, key)
		
		// 4. Update TTL
		pipe.Expire(ctx, key, config.Window)
		
		_, err := pipe.Exec(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limiter internal error",
			})
		}

		reqCount := countCmd.Val()
		if reqCount > int64(config.MaxRequests) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too Many Requests",
				"details": "You have exceeded the rate limit. Please try again later.",
			})
		}

		return c.Next()
	}
}
