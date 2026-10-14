package services

import (
	"testing"
)

// buildBaseRequest returns a standard valid PricingRequest for steel 2mm.
func buildBaseRequest() PricingRequest {
	return PricingRequest{
		TenantID:         1,
		Material:         "steel",
		ThicknessMm:      2.0,
		FlatAreaMm2:      200_000.0,
		CuttingLengthMm:  1500.0,
		BendCount:        2,
		Quantity:         10,
		BoundingBoxX:     500,
		BoundingBoxY:     400,
		WhiteLabelMargin: 0.20,
		TargetCurrency:   "USD",
	}
}

// --- Test 1: Temel Fiyatlandırma Bütünlüğü ---
func TestCalculateQuote_Basic(t *testing.T) {
	req := buildBaseRequest()
	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("CalculateQuote failed: %v", err)
	}
	if resp.FinalSellingPrice <= 0 {
		t.Errorf("FinalSellingPrice must be > 0, got %f", resp.FinalSellingPrice)
	}
	if resp.BendingCostUSD <= 0 {
		t.Errorf("BendingCostUSD must be > 0 (BendCount=2), got %f", resp.BendingCostUSD)
	}
	if resp.LaserCostUSD <= 0 {
		t.Errorf("LaserCostUSD must be > 0, got %f", resp.LaserCostUSD)
	}
	if resp.MaterialCostUSD <= 0 {
		t.Errorf("MaterialCostUSD must be > 0, got %f", resp.MaterialCostUSD)
	}
	if resp.TotalCostUSD <= 0 {
		t.Errorf("TotalCostUSD must be > 0, got %f", resp.TotalCostUSD)
	}
	if resp.UnitCostUSD <= 0 {
		t.Errorf("UnitCostUSD must be > 0, got %f", resp.UnitCostUSD)
	}
	// Margin uygulandıktan sonra FinalSellingPrice > TotalCostUSD olmali
	if resp.FinalSellingPrice <= resp.TotalCostUSD {
		t.Errorf("FinalSellingPrice (%f) must exceed TotalCostUSD (%f) with 20%% margin",
			resp.FinalSellingPrice, resp.TotalCostUSD)
	}
}

// --- Test 2: Sifir Miktar → 1'e normalize edilmeli ---
func TestCalculateQuote_ZeroQuantityNormalized(t *testing.T) {
	req := buildBaseRequest()
	req.Quantity = 0
	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.FinalSellingPrice <= 0 {
		t.Errorf("expected positive price for zero-quantity fallback, got %f", resp.FinalSellingPrice)
	}
}

// --- Test 3: Geçersiz Malzeme → hata beklenir ---
func TestCalculateQuote_InvalidMaterial(t *testing.T) {
	req := buildBaseRequest()
	req.Material = "unobtanium_9000"
	_, err := CalculateQuote(req)
	if err == nil {
		t.Error("expected error for unknown material, got nil")
	}
}

// --- Test 4: Asiri Kalinlik → lazer seçim hatasi ---
func TestCalculateQuote_ThicknessTooHigh(t *testing.T) {
	req := buildBaseRequest()
	req.ThicknessMm = 999.0
	_, err := CalculateQuote(req)
	if err == nil {
		t.Error("expected laser selection error for extreme thickness, got nil")
	}
}

// --- Test 5: Bükme Yok → BendingCost sifir olmali ---
func TestCalculateQuote_NoBending(t *testing.T) {
	req := buildBaseRequest()
	req.BendCount = 0
	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.BendingCostUSD != 0 {
		t.Errorf("BendingCostUSD must be 0 when BendCount=0, got %f", resp.BendingCostUSD)
	}
}

// --- Test 6: Hacim Indirimi Kademeleri (99 / 100 / 500 / 1000 adet) ---
func TestCalculateQuote_VolumeDiscountTiers(t *testing.T) {
	type tierCase struct {
		qty          int
		expectDiscount bool
	}
	tiers := []tierCase{
		{99, false},
		{100, true},
		{500, true},
		{1000, true},
	}
	for _, tc := range tiers {
		req := buildBaseRequest()
		req.Quantity = tc.qty
		resp, err := CalculateQuote(req)
		if err != nil {
			t.Fatalf("qty=%d: unexpected error: %v", tc.qty, err)
		}
		hasDiscount := resp.VolumeDiscountUSD > 0
		if hasDiscount != tc.expectDiscount {
			t.Errorf("qty=%d: expectDiscount=%v, got VolumeDiscountUSD=%f",
				tc.qty, tc.expectDiscount, resp.VolumeDiscountUSD)
		}
	}
}

// --- Test 7: Acil Siparis %50 Markup ---
func TestCalculateQuote_RushOrderMarkup(t *testing.T) {
	base := buildBaseRequest()
	base.RushOrder = false
	respNormal, err := CalculateQuote(base)
	if err != nil {
		t.Fatalf("normal order failed: %v", err)
	}
	rush := buildBaseRequest()
	rush.RushOrder = true
	respRush, err := CalculateQuote(rush)
	if err != nil {
		t.Fatalf("rush order failed: %v", err)
	}
	if respRush.RushOrderMarkupUSD <= 0 {
		t.Errorf("RushOrderMarkupUSD must be > 0, got %f", respRush.RushOrderMarkupUSD)
	}
	if respRush.TotalCostUSD <= respNormal.TotalCostUSD {
		t.Errorf("Rush total (%f) must exceed normal total (%f)",
			respRush.TotalCostUSD, respNormal.TotalCostUSD)
	}
}

// --- Test 8: Faz 3.6 — 9 Gelismis CAM Maliyeti ---
func TestCalculateQuote_AdvancedCAMCosts(t *testing.T) {
	req := buildBaseRequest()
	req.PunchCount = 10
	req.HardwareCount = 5
	req.RequiresFilm = true
	req.PaintingM2 = 0.2
	req.RequiresCMM = true
	req.PackagingType = "pallet"
	req.LaserGasType = "N2"
	req.ToolWearPct = 5.0
	req.WeldingLengthMm = 100
	req.ShippingDistanceKm = 200

	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checks := []struct {
		name string
		val  float64
	}{
		{"PunchingCostUSD", resp.PunchingCostUSD},
		{"HardwareCostUSD", resp.HardwareCostUSD},
		{"SurfaceFilmCostUSD", resp.SurfaceFilmCostUSD},
		{"PaintingCostUSD", resp.PaintingCostUSD},
		{"QualityControlCostUSD", resp.QualityControlCostUSD},
		{"LaserGasCostUSD", resp.LaserGasCostUSD},
		{"WeldingCostUSD", resp.WeldingCostUSD},
		{"LogisticsCostUSD", resp.LogisticsCostUSD},
	}
	for _, c := range checks {
		if c.val <= 0 {
			t.Errorf("%s must be > 0, got %f", c.name, c.val)
		}
	}
}

// --- Test 9: Faz 3.7 — Sekonder Operasyonlar ---
func TestCalculateQuote_SecondaryOperations(t *testing.T) {
	req := buildBaseRequest()
	req.EtchingLengthMm = 500
	req.BevelingLengthMm = 200
	req.RequiresDeburring = true
	req.HemmingCount = 3
	req.MicroJointCount = 8
	req.AssemblyParts = 2
	req.RequiresHAZReport = true

	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checks := []struct {
		name string
		val  float64
	}{
		{"EtchingCostUSD", resp.EtchingCostUSD},
		{"BevelingCostUSD", resp.BevelingCostUSD},
		{"DeburringCostUSD", resp.DeburringCostUSD},
		{"HemmingCostUSD", resp.HemmingCostUSD},
		{"MicroJointCostUSD", resp.MicroJointCostUSD},
		{"AssemblyCostUSD", resp.AssemblyCostUSD},
		{"HAZReportCostUSD", resp.HAZReportCostUSD},
	}
	for _, c := range checks {
		if c.val <= 0 {
			t.Errorf("%s must be > 0, got %f", c.name, c.val)
		}
	}
}

// --- Test 10: Faz 3.8 — Ticari Degiskenler ---
func TestCalculateQuote_CommercialVariances(t *testing.T) {
	req := buildBaseRequest()
	req.TappingCount = 6
	req.CountersinkCount = 4
	req.FormingCount = 3
	req.ComplexBending = true
	req.CornerReliefCount = 10
	req.Unfolding3DService = true
	req.MaterialCertificate = true

	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checks := []struct {
		name string
		val  float64
	}{
		{"TappingCostUSD", resp.TappingCostUSD},
		{"CountersinkCostUSD", resp.CountersinkCostUSD},
		{"FormingCostUSD", resp.FormingCostUSD},
		{"ComplexBendingCostUSD", resp.ComplexBendingCostUSD},
		{"CornerReliefCostUSD", resp.CornerReliefCostUSD},
		{"UnfoldingServiceCostUSD", resp.UnfoldingServiceCostUSD},
		{"MaterialCertificateCostUSD", resp.MaterialCertificateCostUSD},
	}
	for _, c := range checks {
		if c.val <= 0 {
			t.Errorf("%s must be > 0, got %f", c.name, c.val)
		}
	}
}

// --- Test 11: Lazer Gaz Maliyet Sirasi Air < O2 < N2 ---
func TestCalculateQuote_LaserGasCostRanking(t *testing.T) {
	makeReq := func(gas string) PricingRequest {
		r := buildBaseRequest()
		r.LaserGasType = gas
		r.BendCount = 0
		return r
	}
	respAir, _ := CalculateQuote(makeReq("Air"))
	respO2, _ := CalculateQuote(makeReq("O2"))
	respN2, _ := CalculateQuote(makeReq("N2"))

	if respN2.LaserGasCostUSD <= respO2.LaserGasCostUSD {
		t.Errorf("N2 (%f) should cost more than O2 (%f)", respN2.LaserGasCostUSD, respO2.LaserGasCostUSD)
	}
	if respO2.LaserGasCostUSD <= respAir.LaserGasCostUSD {
		t.Errorf("O2 (%f) should cost more than Air (%f)", respO2.LaserGasCostUSD, respAir.LaserGasCostUSD)
	}
}

// --- Test 12: Tane Yonu Penaltisi ---
func TestCalculateQuote_GrainDirectionPenalty(t *testing.T) {
	base := buildBaseRequest()
	base.GrainDirection = false
	respNoGrain, err := CalculateQuote(base)
	if err != nil {
		t.Fatalf("no-grain failed: %v", err)
	}
	withGrain := buildBaseRequest()
	withGrain.GrainDirection = true
	respWithGrain, err := CalculateQuote(withGrain)
	if err != nil {
		t.Fatalf("with-grain failed: %v", err)
	}
	if respWithGrain.MaterialCostUSD <= respNoGrain.MaterialCostUSD {
		t.Errorf("GrainDirection=true cost (%f) should exceed base (%f)",
			respWithGrain.MaterialCostUSD, respNoGrain.MaterialCostUSD)
	}
}

// --- Test 13: Para Birimi Donusumu USD → TRY ---
func TestCalculateQuote_CurrencyConversion(t *testing.T) {
	usdReq := buildBaseRequest()
	usdReq.TargetCurrency = "USD"
	respUSD, _ := CalculateQuote(usdReq)

	tryReq := buildBaseRequest()
	tryReq.TargetCurrency = "TRY"
	respTRY, _ := CalculateQuote(tryReq)

	if respTRY.FinalSellingPrice <= respUSD.FinalSellingPrice {
		t.Errorf("TRY price (%f) should be greater than USD price (%f)",
			respTRY.FinalSellingPrice, respUSD.FinalSellingPrice)
	}
	if respTRY.Currency != "TRY" {
		t.Errorf("expected currency=TRY, got %s", respTRY.Currency)
	}
}

// --- Test 14: Toplam Maliyet Tutarliligi ---
func TestCalculateQuote_TotalCostConsistency(t *testing.T) {
	req := buildBaseRequest()
	req.PunchCount = 5
	req.HardwareCount = 3
	req.RequiresFilm = true
	req.RequiresDeburring = true
	req.TappingCount = 4
	req.ComplexBending = true
	req.MaterialCertificate = true

	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// UnitCostUSD * Qty ≈ TotalCostUSD
	reconstructed := resp.UnitCostUSD * float64(req.Quantity)
	delta := reconstructed - resp.TotalCostUSD
	if delta < -0.01 || delta > 0.01 {
		t.Errorf("UnitCostUSD*Qty (%.4f) != TotalCostUSD (%.4f), delta=%.4f",
			reconstructed, resp.TotalCostUSD, delta)
	}
	if resp.FinalSellingPrice <= resp.TotalCostUSD {
		t.Errorf("FinalSellingPrice (%f) must exceed TotalCostUSD (%f)",
			resp.FinalSellingPrice, resp.TotalCostUSD)
	}
}

// --- Test 15: RushOrder + VolumeDiscount (1500 adet aluminium) ---
func TestCalculateQuote_RushOrder_VolumeDiscount(t *testing.T) {
	req := PricingRequest{
		TenantID:         1,
		Material:         "aluminum",
		ThicknessMm:      1.5,
		FlatAreaMm2:      100000.0,
		CuttingLengthMm:  800.0,
		BendCount:        0,
		Quantity:         1500,
		BoundingBoxX:     300,
		BoundingBoxY:     200,
		RushOrder:        true,
		WhiteLabelMargin: 0.10,
		TargetCurrency:   "USD",
	}
	resp, err := CalculateQuote(req)
	if err != nil {
		t.Fatalf("CalculateQuote failed: %v", err)
	}
	if resp.VolumeDiscountUSD <= 0 {
		t.Errorf("Expected VolumeDiscountUSD > 0 for 1500 parts, got %f", resp.VolumeDiscountUSD)
	}
	if resp.RushOrderMarkupUSD <= 0 {
		t.Errorf("Expected RushOrderMarkupUSD > 0 for rush order, got %f", resp.RushOrderMarkupUSD)
	}
}

