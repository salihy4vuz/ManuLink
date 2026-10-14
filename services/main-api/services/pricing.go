package services

import (
	"fmt"
	"strings"
)

type PricingRequest struct {
	TenantID          uint
	Material          string
	ThicknessMm       float64
	FlatAreaMm2       float64
	CuttingLengthMm   float64
	BendCount         int
	Quantity          int     // Number of parts ordered
	BoundingBoxX      float64
	BoundingBoxY      float64
	ScrapMarginPct    float64 // Actual waste percentage returned by Rust Nesting
	WhiteLabelMargin  float64 // Profit Margin (e.g., 0.20 for 20%)
	TargetCurrency    string  // Target output currency (e.g., "TRY", "EUR")
	LMEPriceEUR       float64
	CADSource         string  // FreeCAD, SolidWorks, Inventor etc.
	
	// Faz 3.6.1-9: Advanced CAM Commercial Parameters
	PunchCount         int
	HardwareCount      int
	RequiresFilm       bool
	PaintingM2         float64
	RequiresCMM        bool
	PackagingType      string // "carton", "wood", "pallet"
	LaserGasType       string // "O2", "N2", "Air"
	ToolWearPct        float64
	WeldingLengthMm    float64
	ShippingDistanceKm float64
	
	// Faz 3.7.1-9: Secondary Operations Parameters
	EtchingLengthMm    float64
	BevelingLengthMm   float64
	RequiresDeburring  bool
	HemmingCount       int
	MicroJointCount    int
	AssemblyParts      int
	RequiresHAZReport  bool
	GrainDirection     bool
	NestingComputeTax  float64
	
	// Faz 3.8.1-9: Advanced Bending & Commercial Variances
	TappingCount       int
	CountersinkCount   int
	FormingCount       int
	ComplexBending     bool
	CornerReliefCount  int
	Unfolding3DService bool
	RushOrder          bool
	MaterialCertificate bool
}

type PricingResponse struct {
	MaterialWeightKg  float64 `json:"material_weight_kg"`
	MaterialCostUSD   float64 `json:"material_cost_usd"`
	LaserCostUSD      float64 `json:"laser_cutting_cost_usd"`
	BendingCostUSD    float64 `json:"bending_cost_usd"`
	
	// Faz 3.6.1-9: Advanced Costs
	PunchingCostUSD       float64 `json:"punching_cost_usd"`
	HardwareCostUSD       float64 `json:"hardware_cost_usd"`
	SurfaceFilmCostUSD    float64 `json:"surface_film_cost_usd"`
	PaintingCostUSD       float64 `json:"painting_cost_usd"`
	QualityControlCostUSD float64 `json:"quality_control_cost_usd"`
	PackagingCostUSD      float64 `json:"packaging_cost_usd"`
	LaserGasCostUSD       float64 `json:"laser_gas_cost_usd"`
	ToolWearCostUSD       float64 `json:"tool_wear_cost_usd"`
	WeldingCostUSD        float64 `json:"welding_cost_usd"`
	LogisticsCostUSD      float64 `json:"logistics_cost_usd"`

	// Faz 3.7.1-9: Secondary Costs
	EtchingCostUSD    float64 `json:"etching_cost_usd"`
	BevelingCostUSD   float64 `json:"beveling_cost_usd"`
	DeburringCostUSD  float64 `json:"deburring_cost_usd"`
	HemmingCostUSD    float64 `json:"hemming_cost_usd"`
	MicroJointCostUSD float64 `json:"micro_joint_cost_usd"`
	AssemblyCostUSD   float64 `json:"assembly_cost_usd"`
	HAZReportCostUSD  float64 `json:"haz_report_cost_usd"`
	
	// Faz 3.8.1-9: Commercial Variances Costs
	TappingCostUSD            float64 `json:"tapping_cost_usd"`
	CountersinkCostUSD        float64 `json:"countersink_cost_usd"`
	FormingCostUSD            float64 `json:"forming_cost_usd"`
	ComplexBendingCostUSD     float64 `json:"complex_bending_cost_usd"`
	CornerReliefCostUSD       float64 `json:"corner_relief_cost_usd"`
	UnfoldingServiceCostUSD   float64 `json:"unfolding_service_cost_usd"`
	MaterialCertificateCostUSD float64 `json:"material_certificate_cost_usd"`
	RushOrderMarkupUSD        float64 `json:"rush_order_markup_usd"`
	VolumeDiscountUSD         float64 `json:"volume_discount_usd"`
	CADSourceDiscountUSD      float64 `json:"cad_source_discount_usd"`

	SetupCostUSD      float64 `json:"setup_cost_usd"` // Fixed cost amortized
	TotalCostUSD      float64 `json:"total_cost_usd"`
	UnitCostUSD       float64 `json:"unit_cost_usd"`
	FinalSellingPrice float64 `json:"final_selling_price"`
	Currency          string  `json:"currency"`
}

// CalculateQuote determines the B2B SaaS price using machine profiles and volume amortization.
func CalculateQuote(req PricingRequest) (PricingResponse, error) {
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Faz 21: Tenant'a özel dinamik fiyatlandırma kofigürasyonunu getir (DB/Redis)
	cfg := GetPricingConfig(req.TenantID)

	// 1. Get Material Rates
	matData, err := GetMaterialPricing(req.Material)
	if err != nil {
		return PricingResponse{}, fmt.Errorf("pricing engine error: %w", err)
	}

	// 2. Select Appropriate Machines
	laser, err := SelectLaser(req.ThicknessMm)
	if err != nil {
		return PricingResponse{}, fmt.Errorf("laser selection error: %w", err)
	}
	
	brake, errB := SelectPressBrake(req.ThicknessMm)
	if req.BendCount > 0 && errB != nil {
		return PricingResponse{}, fmt.Errorf("press brake selection error: %w", errB)
	}

	// 3. Material Cost Calculation (per unit)
	areaM2 := req.FlatAreaMm2 / 1_000_000.0
	thicknessM := req.ThicknessMm / 1000.0
	volumeM3 := areaM2 * thicknessM
	weightKg := volumeM3 * matData.DensityKgM3
	
	// Phase 3.2: Standard Sheet Plate Utilization
	// Fabrikalar sacı metrekareyle değil, plakayla alır.
	sheetReq := CalculateSheetUtilization(req.BoundingBoxX, req.BoundingBoxY, req.Quantity, matData.PricePerKg, req.ThicknessMm, matData.DensityKgM3)
	
	unitMaterialCost := sheetReq.TotalCost / float64(req.Quantity) // Cost per part based on full plates
	
	// Faz 3.7.1-9: Grain Direction Penalty
	if req.GrainDirection {
		// Material usage increases based on strict alignment (e.g. 10%)
		penaltyMultiplier := 1.0 + (cfg.GrainDirectionPenaltyPct / 100.0)
		unitMaterialCost *= penaltyMultiplier
		sheetReq.ScrapRate += (cfg.GrainDirectionPenaltyPct / 2.0)
	}
	
	scrapPercentage := sheetReq.ScrapRate
	req.ScrapMarginPct = scrapPercentage

	// 4. Laser Cutting Cost Calculation (per unit)
	cuttingTimeMin := req.CuttingLengthMm / laser.BaseSpeedMmPerMin
	machineCostPerMin := laser.HourlyRateUSD / 60.0
	unitLaserCost := cuttingTimeMin * machineCostPerMin

	// 5. Bending Cost Calculation (per unit)
	unitBendingCost := 0.0
	if req.BendCount > 0 {
		timePerBendMin := 0.5 // Standard handling time per bend
		unitBendingCost = (float64(req.BendCount) * timePerBendMin) * (brake.HourlyRateUSD / 60.0)
	}

	// 6. Setup Costs (Fixed cost for the entire batch)
	totalSetupCost := (laser.SetupTimeMin * (laser.HourlyRateUSD / 60.0))
	if req.BendCount > 0 {
		totalSetupCost += (brake.SetupTimeMin * (brake.HourlyRateUSD / 60.0))
	}

	// Faz 3.6.1-9: 9 Advanced CAM Costs (Per Unit unless setup)
	
	// 1. Punching & Hardware
	unitPunchCost := float64(req.PunchCount) * cfg.PunchCostPerHit
	unitHardwareCost := float64(req.HardwareCount) * cfg.HardwareCostUnit

	// 2. Surface Protection (Film)
	unitFilmCost := 0.0
	if req.RequiresFilm {
		unitFilmCost = areaM2 * cfg.FilmCostPerM2
	}

	// 3. Painting / Coating
	unitPaintingCost := req.PaintingM2 * cfg.PaintCostPerM2

	// 4. Quality Control
	unitQCCost := 0.0
	if req.RequiresCMM {
		totalSetupCost += cfg.CMMSetupCost
		unitQCCost = cfg.CMMUnitCost
	}

	// 5. Packaging
	packagingCost := 0.0
	if req.PackagingType == "pallet" {
		totalSetupCost += cfg.PackagingPalletCost
	} else if req.PackagingType == "wood" {
		totalSetupCost += cfg.PackagingWoodCost
	} else if req.PackagingType == "carton" {
		packagingCost = cfg.PackagingCartonCost
	}

	// 6. Laser Gas
	var gasCostPerMin float64
	switch strings.ToUpper(req.LaserGasType) {
	case "O2":
		gasCostPerMin = cfg.LaserGasO2PerMin
	case "N2":
		gasCostPerMin = cfg.LaserGasN2PerMin
	case "AIR", "":
		gasCostPerMin = cfg.LaserGasAirPerMin
	case "FIBER_AIR":
		gasCostPerMin = cfg.LaserGasAirPerMin * 1.33 // Extra consumption factor
	default:
		gasCostPerMin = cfg.LaserGasO2PerMin
	}
	unitLaserGasCost := cuttingTimeMin * gasCostPerMin

	// 7. Tool Wear (Amortization)
	unitToolWearCost := (req.ToolWearPct / 100.0) * (unitLaserCost + unitBendingCost)

	// 8. Welding
	unitWeldingCost := (req.WeldingLengthMm / 1000.0) * cfg.WeldingCostPerMeter

	// 9. Logistics (Flat fee added to setup/total)
	logisticsCostUSD := req.ShippingDistanceKm * cfg.LogisticsCostPerKm
	totalSetupCost += logisticsCostUSD

	// Faz 3.7.1-9: Secondary Operations Costs (Per Unit unless setup)
	
	// 1. Etching & Beveling
	unitEtchingCost := req.EtchingLengthMm * cfg.EtchingCostPerMm
	unitBevelingCost := req.BevelingLengthMm * cfg.BevelingCostPerMm

	// 2. Deburring
	unitDeburringCost := 0.0
	if req.RequiresDeburring {
		unitDeburringCost = cfg.DeburringCostUnit
	}

	// 3. Hemming & Micro-Joints
	unitHemmingCost := float64(req.HemmingCount) * cfg.HemmingCostPerOp
	unitMicroJointCost := float64(req.MicroJointCount) * cfg.MicroJointCostUnit

	// 4. Assembly
	unitAssemblyCost := float64(req.AssemblyParts) * cfg.AssemblyCostPerPart

	// 5. HAZ Report & Compute Tax (Setup Costs)
	hazReportCost := 0.0
	if req.RequiresHAZReport {
		hazReportCost = cfg.HAZReportCost
		totalSetupCost += hazReportCost
	}
	totalSetupCost += req.NestingComputeTax // E.g., cloud CPU tax for very complex nesting

	// Faz 3.8.1-9: Advanced Bending & Holes Costs (Per Unit unless setup)
	
	// 1. Tapping & Countersink (Per hole costs)
	unitTappingCost := float64(req.TappingCount) * cfg.TappingCostPerHole
	unitCountersinkCost := float64(req.CountersinkCount) * cfg.CountersinkCostPerHole

	// 2. Forming (Dimple, Louver)
	unitFormingCost := float64(req.FormingCount) * cfg.FormingCostPerHit

	// 3. Complex Bending
	complexBendingCost := 0.0
	if req.ComplexBending {
		complexBendingCost = cfg.ComplexBendingSetup
		totalSetupCost += complexBendingCost
	}

	// 4. Corner Relief
	unitCornerReliefCost := float64(req.CornerReliefCount) * cfg.CornerReliefCostUnit

	// 5. 3D Unfolding Service
	unfoldingCost := 0.0
	if req.Unfolding3DService {
		unfoldingCost = cfg.UnfoldingServiceCost
		totalSetupCost += unfoldingCost
	}

	// 6. Material Certificate
	certCost := 0.0
	if req.MaterialCertificate {
		certCost = cfg.MaterialCertCost
		totalSetupCost += certCost
	}

	// 7. Total Cost Calculation
	totalVariableCost := (unitMaterialCost + unitLaserCost + unitBendingCost + unitPunchCost + unitHardwareCost + unitFilmCost + unitPaintingCost + unitQCCost + packagingCost + unitLaserGasCost + unitToolWearCost + unitWeldingCost + unitEtchingCost + unitBevelingCost + unitDeburringCost + unitHemmingCost + unitMicroJointCost + unitAssemblyCost + unitTappingCost + unitCountersinkCost + unitFormingCost + unitCornerReliefCost) * float64(req.Quantity)
	totalCostUSD := totalSetupCost + totalVariableCost
	unitCostUSD := totalCostUSD / float64(req.Quantity) // Amortization of setup cost

	// Faz 3.8.1-9: Volume Discount
	volumeDiscountPct := 0.0
	if req.Quantity >= 1000 {
		volumeDiscountPct = cfg.VolumeDiscount1000Pct / 100.0
	} else if req.Quantity >= 500 {
		volumeDiscountPct = cfg.VolumeDiscount500Pct / 100.0
	} else if req.Quantity >= 100 {
		volumeDiscountPct = cfg.VolumeDiscount100Pct / 100.0
	}
	
	discountAmount := totalCostUSD * volumeDiscountPct
	totalCostUSD -= discountAmount

	// Faz 3.8.1-10: FreeCAD / Open Source Loyalty Discount
	cadSourceDiscountAmount := 0.0
	if strings.EqualFold(req.CADSource, "FreeCAD") {
		// %5 additional discount for using open source FreeCAD
		cadSourceDiscountAmount = totalCostUSD * 0.05
		totalCostUSD -= cadSourceDiscountAmount
	}

	// Faz 3.8.1-9: Rush Order Markup
	rushMarkup := 0.0
	if req.RushOrder {
		rushMarkup = totalCostUSD * (cfg.RushOrderMarkupPct / 100.0)
		totalCostUSD += rushMarkup
	}

	// WhiteLabel Margin (Kendi default değerini kullanır veya tenant'ınkini)
	margin := req.WhiteLabelMargin
	if margin == 0 {
		margin = cfg.DefaultWhiteLabelMarginPct / 100.0
	}
	finalSellingPriceUSD := totalCostUSD * (1.0 + margin)

	// 8. Convert to Target Currency
	targetCurrency := "USD"
	if req.TargetCurrency != "" {
		targetCurrency = req.TargetCurrency
	}
	exchangeRate := GetExchangeRate(targetCurrency)
	finalSellingPriceTarget := finalSellingPriceUSD * exchangeRate

	return PricingResponse{
		MaterialWeightKg:  weightKg * float64(req.Quantity),
		MaterialCostUSD:   unitMaterialCost * float64(req.Quantity),
		LaserCostUSD:      unitLaserCost * float64(req.Quantity),
		BendingCostUSD:    unitBendingCost * float64(req.Quantity),
		
		PunchingCostUSD:       unitPunchCost * float64(req.Quantity),
		HardwareCostUSD:       unitHardwareCost * float64(req.Quantity),
		SurfaceFilmCostUSD:    unitFilmCost * float64(req.Quantity),
		PaintingCostUSD:       unitPaintingCost * float64(req.Quantity),
		QualityControlCostUSD: unitQCCost * float64(req.Quantity),
		PackagingCostUSD:      packagingCost * float64(req.Quantity),
		LaserGasCostUSD:       unitLaserGasCost * float64(req.Quantity),
		ToolWearCostUSD:       unitToolWearCost * float64(req.Quantity),
		WeldingCostUSD:        unitWeldingCost * float64(req.Quantity),
		LogisticsCostUSD:      logisticsCostUSD,

		EtchingCostUSD:    unitEtchingCost * float64(req.Quantity),
		BevelingCostUSD:   unitBevelingCost * float64(req.Quantity),
		DeburringCostUSD:  unitDeburringCost * float64(req.Quantity),
		HemmingCostUSD:    unitHemmingCost * float64(req.Quantity),
		MicroJointCostUSD: unitMicroJointCost * float64(req.Quantity),
		AssemblyCostUSD:   unitAssemblyCost * float64(req.Quantity),
		HAZReportCostUSD:  hazReportCost,

		TappingCostUSD:             unitTappingCost * float64(req.Quantity),
		CountersinkCostUSD:         unitCountersinkCost * float64(req.Quantity),
		FormingCostUSD:             unitFormingCost * float64(req.Quantity),
		ComplexBendingCostUSD:      complexBendingCost,
		CornerReliefCostUSD:        unitCornerReliefCost * float64(req.Quantity),
		UnfoldingServiceCostUSD:    unfoldingCost,
		MaterialCertificateCostUSD: certCost,
		RushOrderMarkupUSD:         rushMarkup,
		VolumeDiscountUSD:          discountAmount,
		CADSourceDiscountUSD:       cadSourceDiscountAmount,

		SetupCostUSD:      totalSetupCost,
		TotalCostUSD:      totalCostUSD,
		UnitCostUSD:       unitCostUSD,
		FinalSellingPrice: finalSellingPriceTarget,
		Currency:          targetCurrency,
	}, nil
}
