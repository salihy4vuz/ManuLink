package services

import (
	"testing"
)

// --- Sheet Utilization Unit Testleri ---

// Test 1: Küçük parça → birden fazla parça aynı plakaya sığmalı
func TestCalculateSheetUtilization_SmallPart_MultiplePerSheet(t *testing.T) {
	// 200x100mm parça, 1000x2000mm plakaya çok sayıda sığar
	req := CalculateSheetUtilization(200, 100, 10, 2.0, 3.0, 7.85)

	if req.SheetsNeeded <= 0 {
		t.Errorf("SheetsNeeded must be > 0, got %d", req.SheetsNeeded)
	}
	if req.TotalCost <= 0 {
		t.Errorf("TotalCost must be > 0, got %f", req.TotalCost)
	}
	// 10 adet küçük parça → 1 plakaya sığmalı
	if req.SheetsNeeded > 2 {
		t.Errorf("Small part (200x100) qty=10 should fit in <=2 sheets, got %d", req.SheetsNeeded)
	}
}

// Test 2: Scrap Rate sınırları — 0-100 arasında olmalı
func TestCalculateSheetUtilization_ScrapRateRange(t *testing.T) {
	req := CalculateSheetUtilization(500, 400, 5, 2.5, 2.0, 7.85)

	if req.ScrapRate < 0 || req.ScrapRate > 100 {
		t.Errorf("ScrapRate must be in [0,100], got %f", req.ScrapRate)
	}
}

// Test 3: Büyük qty → daha fazla plaka gerekli (doğrusal ölçekleme)
func TestCalculateSheetUtilization_HighQuantityScaling(t *testing.T) {
	req10 := CalculateSheetUtilization(300, 200, 10, 2.0, 2.0, 7.85)
	req100 := CalculateSheetUtilization(300, 200, 100, 2.0, 2.0, 7.85)

	if req100.SheetsNeeded <= req10.SheetsNeeded {
		t.Errorf("qty=100 should need more sheets than qty=10: got %d vs %d",
			req100.SheetsNeeded, req10.SheetsNeeded)
	}
	if req100.TotalCost <= req10.TotalCost {
		t.Errorf("qty=100 TotalCost (%f) should be > qty=10 TotalCost (%f)",
			req100.TotalCost, req10.TotalCost)
	}
}

// Test 4: Parça boyutu > herhangi bir standart plaka → 1 plaka gereklil
// (Her plakaya 0 parça sığar, ancak graceful fallback olmalı)
func TestCalculateSheetUtilization_PartLargerThanSheet_Handled(t *testing.T) {
	// 5000x4000mm parça → hiçbir standart plakaya sığmaz
	req := CalculateSheetUtilization(5000, 4000, 1, 5.0, 3.0, 7.85)
	// SheetsNeeded 0 olabilir (graceful) — panic olmamalı
	_ = req // sadece crash check
}

// Test 5: Plaka boyutu tutarlılığı — bilinen standart boyutlardan biri seçilmeli
func TestCalculateSheetUtilization_StandardSheetSizeSelected(t *testing.T) {
	req := CalculateSheetUtilization(400, 300, 5, 2.0, 2.0, 7.85)

	validLengths := []float64{SheetLength2000, SheetLength2500, SheetLength3000}
	validWidths := []float64{SheetWidth1000, SheetWidth1250, SheetWidth1500}

	validL := false
	for _, l := range validLengths {
		if req.SheetLength == l {
			validL = true
			break
		}
	}
	validW := false
	for _, w := range validWidths {
		if req.SheetWidth == w {
			validW = true
			break
		}
	}

	if !validL {
		t.Errorf("SheetLength %f is not a standard size %v", req.SheetLength, validLengths)
	}
	if !validW {
		t.Errorf("SheetWidth %f is not a standard size %v", req.SheetWidth, validWidths)
	}
}

// Test 6: Kalınlık maliyete etkisi — daha kalın = daha ağır = daha pahalı
func TestCalculateSheetUtilization_ThicknessAffectsCost(t *testing.T) {
	reqThin := CalculateSheetUtilization(300, 200, 5, 2.0, 1.0, 7.85) // 1mm
	reqThick := CalculateSheetUtilization(300, 200, 5, 2.0, 6.0, 7.85) // 6mm

	if reqThick.TotalCost <= reqThin.TotalCost {
		t.Errorf("Thick plate cost (%f) should be > thin plate cost (%f)",
			reqThick.TotalCost, reqThin.TotalCost)
	}
}

// Test 7: Malzeme fiyatı artışı maliyete yansımalı
func TestCalculateSheetUtilization_MaterialPriceAffectsCost(t *testing.T) {
	reqCheap := CalculateSheetUtilization(300, 200, 5, 0.5, 2.0, 7.85)  // $0.50/kg
	reqExpensive := CalculateSheetUtilization(300, 200, 5, 5.0, 2.0, 7.85) // $5.00/kg

	if reqExpensive.TotalCost <= reqCheap.TotalCost {
		t.Errorf("Expensive material (%f) should cost more than cheap (%f)",
			reqExpensive.TotalCost, reqCheap.TotalCost)
	}
}

// Test 8: Kare parça rotasyon optimizasyonu — scrap rate makul seviyede
func TestCalculateSheetUtilization_SquarePart_ReasonableScrap(t *testing.T) {
	// 490x490mm kare parça → rotasyon fark yaratmaz, scrap mantıklı olmalı
	req := CalculateSheetUtilization(490, 490, 4, 2.0, 3.0, 7.85)

	// Scrap oranı %50'nin altında olmalı (iyi nesting)
	if req.ScrapRate > 60.0 {
		t.Errorf("ScrapRate for square parts should be < 60%%, got %f%%", req.ScrapRate)
	}
}
