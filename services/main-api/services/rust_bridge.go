package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
)

// RustMachiningResponse represents the data coming back from the geometry-engine
type RustMachiningResponse struct {
	RawBoundingBoxVolumeMm3 float64 `json:"raw_bounding_box_volume_mm3"`
	FinalPartVolumeMm3      float64 `json:"final_part_volume_mm3"`
	RemovedVolumeMm3        float64 `json:"removed_volume_mm3"`
	RemovedWeightKg         float64 `json:"removed_weight_kg"`
	EstimatedTimeMin        float64 `json:"estimated_machining_time_min"`
	SuggestedRPM            float64 `json:"suggested_rpm"`
	SuggestedFeedRate       float64 `json:"suggested_feed_rate_mm_min"`
}

// CallRustMachiningEngine simulates a gRPC/HTTP call to the Rust Geometry Engine
func CallRustMachiningEngine(stepBase64 string, materialDensity float64, machineType string) (RustMachiningResponse, error) {
	log.Println("[Rust Bridge] Calling geometry-engine machining.rs via RabbitMQ...")

	payload := map[string]interface{}{
		"step_file_base64": stepBase64,
		"material_density": materialDensity,
		"machining_type":   machineType,
	}

	// Call Rust engine via RabbitMQ
	resStr, err := rabbitmq.CallRustEngine(context.Background(), "analyze_machining", payload)
	if err != nil {
		log.Printf("[Rust Bridge] Error calling analyze_machining: %v\n", err)
		return RustMachiningResponse{}, err
	}

	var response RustMachiningResponse
	if err := json.Unmarshal([]byte(resStr), &response); err != nil {
		return RustMachiningResponse{}, err
	}

	return response, nil
}
