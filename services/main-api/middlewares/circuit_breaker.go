package middlewares

import (
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	failureCount int32
	isCircuitOpen int32 // 0 = Closed (Healthy), 1 = Open (Failing)
	lastFailure time.Time
)

// CircuitBreakerMiddleware prevents cascading failures if downstream services (Rust/DB) crash
func CircuitBreakerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if circuit is open
		if atomic.LoadInt32(&isCircuitOpen) == 1 {
			// Wait 30 seconds before trying to close the circuit (Half-Open state)
			if time.Since(lastFailure) > 30*time.Second {
				atomic.StoreInt32(&isCircuitOpen, 0)
				atomic.StoreInt32(&failureCount, 0)
			} else {
				// Fast fail
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
					"error": "Service temporarily unavailable. Circuit Breaker is OPEN.",
				})
			}
		}

		err := c.Next()

		if err != nil {
			// If it's a 500 Internal Server Error, increment failure count
			if e, ok := err.(*fiber.Error); ok && e.Code >= 500 {
				fails := atomic.AddInt32(&failureCount, 1)
				if fails > 5 { // Trip the breaker after 5 consecutive failures
					atomic.StoreInt32(&isCircuitOpen, 1)
					lastFailure = time.Now()
				}
			}
		} else {
			// Reset on success
			atomic.StoreInt32(&failureCount, 0)
		}

		return err
	}
}
