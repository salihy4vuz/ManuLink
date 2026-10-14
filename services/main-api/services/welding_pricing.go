package services

type WeldingQuoteRequest struct {
	Material         string
	Quantity         int
	WeldLengthMm     float64
	WeldType         string // TIG, MIG, Spot
	WhiteLabelMargin float64
	TargetCurrency   string
}

type WeldingQuoteResponse struct {
	MachineType       string  `json:"machine_type"`
	WeldLengthMm      float64 `json:"weld_length_mm"`
	LaborCostUSD      float64 `json:"labor_cost_usd"`
	GasCostUSD        float64 `json:"gas_cost_usd"`
	SetupCostUSD      float64 `json:"setup_cost_usd"`
	TotalCostUSD      float64 `json:"total_cost_usd"`
	FinalSellingPrice float64 `json:"final_selling_price"`
	Currency          string  `json:"currency"`
}

func CalculateWeldingQuote(req WeldingQuoteRequest) WeldingQuoteResponse {
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	setupCostUSD := 100.0 // Fikstür hazırlığı

	laborRatePerMm := 0.05
	gasRatePerMm := 0.01

	if req.WeldType == "TIG" {
		laborRatePerMm = 0.08 // TIG işçiliği daha pahalıdır
		gasRatePerMm = 0.02   // Argon gazı
	}

	laborCost := req.WeldLengthMm * laborRatePerMm * float64(req.Quantity)
	gasCost := req.WeldLengthMm * gasRatePerMm * float64(req.Quantity)

	totalCostUSD := setupCostUSD + laborCost + gasCost
	finalPriceUSD := totalCostUSD * (1.0 + req.WhiteLabelMargin)

	targetCurrency := "USD"
	if req.TargetCurrency != "" {
		targetCurrency = req.TargetCurrency
	}
	exchangeRate := GetExchangeRate(targetCurrency)

	return WeldingQuoteResponse{
		MachineType:       "welding_" + req.WeldType,
		WeldLengthMm:      req.WeldLengthMm,
		LaborCostUSD:      laborCost,
		GasCostUSD:        gasCost,
		SetupCostUSD:      setupCostUSD,
		TotalCostUSD:      totalCostUSD,
		FinalSellingPrice: finalPriceUSD * exchangeRate,
		Currency:          targetCurrency,
	}
}
