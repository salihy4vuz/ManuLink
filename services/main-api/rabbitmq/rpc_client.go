package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"
)

// CallRustEngine sends a task to the Rust Geometry Engine via RabbitMQ RPC
func CallRustEngine(ctx context.Context, taskType string, payload interface{}) (string, error) {
	// 10 saniyelik bir timeout belirliyoruz (eğer dışarıdan ctx ile timeout gelmemişse)
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	queueName := "geometry_tasks"
	
	log.Printf("[RabbitMQ] Publishing RPC request to '%s' with task: %s\n", queueName, taskType)
	
	response, err := PublishRPC(ctx, queueName, taskType, payload)
	if err != nil {
		return "", fmt.Errorf("failed to call rust engine: %w", err)
	}
	
	log.Printf("[RabbitMQ] Received successful response from Rust Engine for task: %s\n", taskType)
	
	return response, nil
}
