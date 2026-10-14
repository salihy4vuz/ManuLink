package services

import (
	"fmt"
)

type MachiningQuoteRequest struct {
	StepFileBase64   string
	Material         string
	Quantity         int
	MachineType      string `json:"machine_type"` // milling_3axis, lathe, swiss_lathe
	WhiteLabelMargin float64
	TargetCurrency   string
}

type MachiningQuoteResponse struct {
	MachineType          string  `json:"machine_type"`
	BoundingBoxVolumeMm3 float64 `json:"bounding_box_volume_mm3"`
	RemovedWeightKg      float64 `json:"removed_weight_kg"`
	MaterialCostUSD      float64 `json:"material_cost_usd"`
	MachiningCostUSD     float64 `json:"machining_cost_usd"`
	SetupCostUSD         float64 `json:"setup_cost_usd"`
	TotalCostUSD         float64 `json:"total_cost_usd"`
	FinalSellingPrice    float64 `json:"final_selling_price"`
	Currency             string  `json:"currency"`
}

func CalculateMachiningQuote(req MachiningQuoteRequest) (MachiningQuoteResponse, error) {
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	matData, err := GetMaterialPricing(req.Material)
	if err != nil {
		return MachiningQuoteResponse{}, fmt.Errorf("material not found: %w", err)
	}

	if req.MachineType == "welding" || req.MachineType == "press" {
		return MachiningQuoteResponse{}, fmt.Errorf("non-cutting operations (welding/press) are moved to their own dedicated services")
	}

	rustData, err := CallRustMachiningEngine(req.StepFileBase64, matData.DensityKgM3, req.MachineType)
	if err != nil {
		return MachiningQuoteResponse{}, fmt.Errorf("geometry engine failed: %w", err)
	}

	// 1. Material Cost (Full Bounding Box block)
	rawVolumeM3 := rustData.RawBoundingBoxVolumeMm3 / 1_000_000_000.0
	rawWeightKg := rawVolumeM3 * matData.DensityKgM3
	unitMaterialCost := rawWeightKg * matData.PricePerKg

	// 2. Machining Cost (CNC Hourly Rate based on Machine Type)
	cncHourlyRateUSD := 65.0 // Generic 3-Axis CNC rate default
	setupCostUSD := 120.0    // Default setup cost

	switch req.MachineType {
	case "lathe":
		cncHourlyRateUSD = 45.0 // Torna is generally cheaper
		setupCostUSD = 80.0
	case "swiss_lathe":
		cncHourlyRateUSD = 80.0 // Kayar Otomat is expensive but fast
		setupCostUSD = 150.0
	case "welding":
		cncHourlyRateUSD = 55.0 // Kaynak işçiliği (TIG/Gazaltı)
		setupCostUSD = 100.0    // Fikstür hazırlığı
	case "drilling":
		cncHourlyRateUSD = 30.0 // Sütunlu Matkap/Kılavuz daha ucuz
		setupCostUSD = 50.0
	case "press":
		cncHourlyRateUSD = 90.0 // Eksantrik Pres (Saatlik kalıp vurma)
		setupCostUSD = 250.0    // Kalıp bağlama maliyeti yüksek
	case "milling_3axis":
		cncHourlyRateUSD = 65.0
		setupCostUSD = 120.0
	default:
		req.MachineType = "milling_3axis"
	}

	unitMachiningCost := (rustData.EstimatedTimeMin / 60.0) * cncHourlyRateUSD

	totalVariableCost := (unitMaterialCost + unitMachiningCost) * float64(req.Quantity)
	totalCostUSD := setupCostUSD + totalVariableCost

	finalPriceUSD := totalCostUSD * (1.0 + req.WhiteLabelMargin)

	targetCurrency := "USD"
	if req.TargetCurrency != "" {
		targetCurrency = req.TargetCurrency
	}
	exchangeRate := GetExchangeRate(targetCurrency)
	finalPriceTarget := finalPriceUSD * exchangeRate

	return MachiningQuoteResponse{
		MachineType:          req.MachineType,
		BoundingBoxVolumeMm3: rustData.RawBoundingBoxVolumeMm3,
		RemovedWeightKg:      rustData.RemovedWeightKg * float64(req.Quantity),
		MaterialCostUSD:      unitMaterialCost * float64(req.Quantity),
		MachiningCostUSD:     unitMachiningCost * float64(req.Quantity),
		SetupCostUSD:         setupCostUSD,
		TotalCostUSD:         totalCostUSD,
		FinalSellingPrice:    finalPriceTarget,
		Currency:             targetCurrency,
	}, nil
}


