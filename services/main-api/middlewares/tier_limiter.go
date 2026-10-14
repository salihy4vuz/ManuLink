package middlewares

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// UsageTracker holds how many quotes a user has made today
type UsageTracker struct {
	Count     int
	ResetTime time.Time
}

var (
	usageMap   = make(map[string]*UsageTracker)
	usageMutex sync.Mutex
)

// TierRateLimiter checks the user's subscription tier and limits quotes.
func TierRateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mock JWT Auth Extraction (In reality, extract from Authorization Header)
		// We'll pass username in a custom header for this FOSS mock
		username := c.Get("X-User", "admin")
		
		tier := services.GetUserTier(username)

		// Get max limits based on tier
		maxLimit := 5 // Free
		if tier == services.TierPro {
			maxLimit = 500
		} else if tier == services.TierEnterprise {
			return c.Next() // Unlimited
		}

		usageMutex.Lock()
		tracker, exists := usageMap[username]
		
		// Reset logic (daily)
		if !exists || time.Now().After(tracker.ResetTime) {
			tracker = &UsageTracker{
				Count:     0,
				ResetTime: time.Now().Add(24 * time.Hour), // Reset in 24h
			}
			usageMap[username] = tracker
		}

		// Check limit
		if tracker.Count >= maxLimit {
			usageMutex.Unlock()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Quota Exceeded",
				"message": "Please upgrade your subscription to get more quotes.",
				"tier":    tier,
				"limit":   maxLimit,
			})
		}

		// Increment usage
		tracker.Count++
		usageMutex.Unlock()

		return c.Next()
	}
}
