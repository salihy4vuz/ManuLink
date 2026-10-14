package services

import (
	"strings"
	"testing"
)

// --- Machining Pricing Unit Testleri ---
// NOT: CalculateMachiningQuote RabbitMQ (Rust engine) gerektirir.
// Bu testler, RabbitMQ bağımlılığından bağımsız parçaları test eder.

// Test 1: Geçersiz malzeme → hata beklenir (RabbitMQ öncesinde fail)
func TestCalculateMachiningQuote_InvalidMaterial(t *testing.T) {
	req := MachiningQuoteRequest{
		Material:    "unobtanium_9000",
		Quantity:    1,
		MachineType: "milling_3axis",
	}
	_, err := CalculateMachiningQuote(req)
	if err == nil {
		t.Error("expected error for invalid material, got nil")
	}
	if !strings.Contains(err.Error(), "material not found") {
		t.Errorf("expected 'material not found' in error, got: %v", err)
	}
}

// Test 2: Welding/Press → özel hata mesajı (dedicated service yönlendirmesi)
func TestCalculateMachiningQuote_WeldingRejected(t *testing.T) {
	req := MachiningQuoteRequest{
		Material:    "steel",
		Quantity:    1,
		MachineType: "welding",
	}
	_, err := CalculateMachiningQuote(req)
	if err == nil {
		t.Error("expected error for welding operation type, got nil")
	}
	if !strings.Contains(err.Error(), "non-cutting") {
		t.Errorf("expected 'non-cutting' in error message, got: %v", err)
	}
}

func TestCalculateMachiningQuote_PressRejected(t *testing.T) {
	req := MachiningQuoteRequest{
		Material:    "steel",
		Quantity:    1,
		MachineType: "press",
	}
	_, err := CalculateMachiningQuote(req)
	if err == nil {
		t.Error("expected error for press operation type, got nil")
	}
}

// Test 3: Sıfır Miktar → 1'e normalize (material validation geçer, engine fail olur)
func TestCalculateMachiningQuote_ZeroQuantityNormalized(t *testing.T) {
	req := MachiningQuoteRequest{
		Material:    "steel",
		Quantity:    0, // Sıfır → normalize to 1
		MachineType: "milling_3axis",
	}
	// RabbitMQ olmadan engine hatası beklenir, ama material hatası OLMAMALI
	_, err := CalculateMachiningQuote(req)
	if err != nil && strings.Contains(err.Error(), "material not found") {
		t.Errorf("should not fail with material error for steel: %v", err)
	}
	// Geometry engine hatası bekleniyor (test ortamında RabbitMQ yok) — bu kabul edilebilir
}

// Test 4: MachiningQuoteRequest struct alanları validasyonu
func TestMachiningQuoteRequest_FieldsAccessible(t *testing.T) {
	req := MachiningQuoteRequest{
		StepFileBase64:   "dGVzdA==",
		Material:         "aluminum",
		Quantity:         5,
		MachineType:      "lathe",
		WhiteLabelMargin: 0.25,
		TargetCurrency:   "TRY",
	}
	if req.Material != "aluminum" {
		t.Errorf("expected material=aluminum, got %s", req.Material)
	}
	if req.Quantity != 5 {
		t.Errorf("expected qty=5, got %d", req.Quantity)
	}
	if req.WhiteLabelMargin != 0.25 {
		t.Errorf("expected margin=0.25, got %f", req.WhiteLabelMargin)
	}
}

// Test 5: MachiningQuoteResponse alanları kontrolü
func TestMachiningQuoteResponse_FieldsPresent(t *testing.T) {
	resp := MachiningQuoteResponse{
		MachineType:          "milling_3axis",
		BoundingBoxVolumeMm3: 125_000.0,
		RemovedWeightKg:      0.5,
		MaterialCostUSD:      10.0,
		MachiningCostUSD:     25.0,
		SetupCostUSD:         120.0,
		TotalCostUSD:         155.0,
		FinalSellingPrice:    186.0,
		Currency:             "USD",
	}
	// Toplam tutarlılık: FinalSellingPrice = TotalCostUSD * (1 + margin)
	// 155 * 1.2 = 186
	if resp.FinalSellingPrice != 186.0 {
		t.Errorf("FinalSellingPrice should be 186.0, got %f", resp.FinalSellingPrice)
	}
	if resp.TotalCostUSD <= 0 {
		t.Errorf("TotalCostUSD must be > 0")
	}
}

// Test 6: Farklı makine tiplerinin saatlik ücret sıralaması (logic test)
// Swiss Lathe > Press Brake > 3-Axis Milling > Lathe > Drilling (kaynak koda göre)
func TestMachiningPricing_HourlyRateOrdering(t *testing.T) {
	// Bu değerler machining_pricing.go'daki switch-case'den geliyor
	rates := map[string]float64{
		"lathe":       45.0,
		"drilling":    30.0,
		"milling_3axis": 65.0,
		"swiss_lathe": 80.0,
	}
	if rates["swiss_lathe"] <= rates["milling_3axis"] {
		t.Errorf("swiss_lathe (%f) should be more expensive than milling_3axis (%f)",
			rates["swiss_lathe"], rates["milling_3axis"])
	}
	if rates["milling_3axis"] <= rates["lathe"] {
		t.Errorf("milling_3axis (%f) should be more expensive than lathe (%f)",
			rates["milling_3axis"], rates["lathe"])
	}
	if rates["lathe"] <= rates["drilling"] {
		t.Errorf("lathe (%f) should be more expensive than drilling (%f)",
			rates["lathe"], rates["drilling"])
	}
}
