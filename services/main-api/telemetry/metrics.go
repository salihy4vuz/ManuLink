package telemetry

import (
	"sync"
	"time"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// In a real production setup, we would use prometheus/client_golang/prometheus
// Here we simulate the prometheus /metrics endpoint natively for the architecture

var (
	TotalRequests      int64
	TotalQuotes        int64
	RabbitMQueueLength int64
	mutex              sync.Mutex
)

// Middleware to track total API requests
func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		mutex.Lock()
		TotalRequests++
		mutex.Unlock()
		return c.Next()
	}
}

// IncrementQuote tracks business logic metrics
func IncrementQuote() {
	mutex.Lock()
	TotalQuotes++
	mutex.Unlock()
}

// SimulateQueueMonitoring simulates watching RabbitMQ lengths
func SimulateQueueMonitoring() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			// Randomly fluctuate between 0 and 15 items in queue
			RabbitMQueueLength = (RabbitMQueueLength + 3) % 15 
			mutex.Unlock()
		}
	}()
}

// MetricsHandler exposes the /metrics endpoint for Prometheus scrapers
func MetricsHandler(c *fiber.Ctx) error {
	mutex.Lock()
	defer mutex.Unlock()

	metricsData := ""
	metricsData += "# HELP api_total_requests_count Total number of HTTP requests\n"
	metricsData += "# TYPE api_total_requests_count counter\n"
	metricsData += "api_total_requests_count " + formatInt(TotalRequests) + "\n\n"

	metricsData += "# HELP api_total_quotes_generated Total number of B2B quotes generated\n"
	metricsData += "# TYPE api_total_quotes_generated counter\n"
	metricsData += "api_total_quotes_generated " + formatInt(TotalQuotes) + "\n\n"

	metricsData += "# HELP rabbitmq_geometry_queue_length Current number of jobs waiting for Rust Engine\n"
	metricsData += "# TYPE rabbitmq_geometry_queue_length gauge\n"
	metricsData += "rabbitmq_geometry_queue_length " + formatInt(RabbitMQueueLength) + "\n"

	c.Set("Content-Type", "text/plain; version=0.0.4")
	return c.SendString(metricsData)
}

func formatInt(n int64) string {
	return fmt.Sprintf("%d", n)
}
