package services

type DrillingQuoteRequest struct {
	Material         string
	Quantity         int
	HoleCount        int
	MaxDiameterMm    float64
	WhiteLabelMargin float64
	TargetCurrency   string
}

type DrillingQuoteResponse struct {
	MachineType       string  `json:"machine_type"`
	HoleCount         int     `json:"hole_count"`
	ToolWearCostUSD   float64 `json:"tool_wear_cost_usd"`
	OperationCostUSD  float64 `json:"operation_cost_usd"`
	SetupCostUSD      float64 `json:"setup_cost_usd"`
	TotalCostUSD      float64 `json:"total_cost_usd"`
	FinalSellingPrice float64 `json:"final_selling_price"`
	Currency          string  `json:"currency"`
}

func CalculateDrillingQuote(req DrillingQuoteRequest) DrillingQuoteResponse {
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	setupCostUSD := 50.0 // Sütunlu matkap kurulumu

	costPerHole := 0.20
	if req.MaxDiameterMm > 10.0 {
		costPerHole = 0.50
	}
	if req.Material == "stainless" || req.Material == "titanium" {
		costPerHole *= 2.0 // Sert malzeme
	}

	toolWearCost := (costPerHole * 0.3) * float64(req.HoleCount*req.Quantity)
	operationCost := (costPerHole * 0.7) * float64(req.HoleCount*req.Quantity)

	totalCostUSD := setupCostUSD + toolWearCost + operationCost
	finalPriceUSD := totalCostUSD * (1.0 + req.WhiteLabelMargin)

	targetCurrency := "USD"
	if req.TargetCurrency != "" {
		targetCurrency = req.TargetCurrency
	}
	exchangeRate := GetExchangeRate(targetCurrency)

	return DrillingQuoteResponse{
		MachineType:       "drilling",
		HoleCount:         req.HoleCount,
		ToolWearCostUSD:   toolWearCost,
		OperationCostUSD:  operationCost,
		SetupCostUSD:      setupCostUSD,
		TotalCostUSD:      totalCostUSD,
		FinalSellingPrice: finalPriceUSD * exchangeRate,
		Currency:          targetCurrency,
	}
}
