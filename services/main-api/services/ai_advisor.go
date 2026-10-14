package services

import (
	"fmt"
	"math"
)

type DFMRequest struct {
	Material    string  `json:"material"`
	ThicknessMm float64 `json:"thickness_mm"`
	BendCount   int     `json:"bend_count"`
	FlangeMm    float64 `json:"flange_length_mm"` // Shortest flange
	HoleDistMm  float64 `json:"min_hole_distance_to_bend_mm"` // Shortest distance from hole to bend line
}

type DFMResponse struct {
	ManufacturabilityScore int      `json:"score"` // 0-100
	AdviceList             []string `json:"advices"`
	Warning                bool     `json:"warning"`
}

// AnalyzeDFM applies advanced sheet metal statistical and deterministic rules 
// to evaluate if a part can be manufactured without defects (cracking, deformation, V-die dropping).
func AnalyzeDFM(req DFMRequest) DFMResponse {
	score := 100
	var advices []string
	warning := false

	// Rule 1: Minimum Flange Length (V-Die rule)
	// Usually, minimum flange length should be ~4x the thickness for standard V-dies.
	minRequiredFlange := req.ThicknessMm * 4.0
	if req.FlangeMm > 0 && req.FlangeMm < minRequiredFlange {
		score -= 30
		warning = true
		advices = append(advices, fmt.Sprintf("Kritik: En kısa flanş (%.1f mm), sac kalınlığına (%.1f mm) göre çok kısa. V-Kalıp içine düşme riski var. Minimum %.1f mm önerilir.", req.FlangeMm, req.ThicknessMm, minRequiredFlange))
	} else if req.FlangeMm > 0 && req.FlangeMm < (req.ThicknessMm * 5.0) {
		score -= 10
		advices = append(advices, fmt.Sprintf("Uyarı: Flanş uzunluğu (%.1f mm) sınır değerlere yakın. Operatör zorlanabilir.", req.FlangeMm))
	}

	// Rule 2: Minimum Hole Distance to Bend Line
	// Holes too close to the bend line will deform (become oval). 
	// Safe distance is usually 1.5x to 2x thickness + bend radius.
	minSafeHoleDist := (req.ThicknessMm * 1.5) + req.ThicknessMm // assuming R=T
	if req.HoleDistMm > 0 && req.HoleDistMm < minSafeHoleDist {
		score -= 25
		warning = true
		advices = append(advices, fmt.Sprintf("Kritik: Delik, büküm çizgisine çok yakın (%.1f mm). Büküm sırasında delikte ovalleşme (deformasyon) yaşanacaktır. Minimum güvenli mesafe: %.1f mm.", req.HoleDistMm, minSafeHoleDist))
	}

	// Rule 3: Material Specific Bend Constraints
	materialLower := req.Material
	if materialLower == "aluminum" && req.ThicknessMm > 4.0 {
		score -= 15
		advices = append(advices, "Uyarı: 4mm üzeri Alüminyum bükümlerinde çatlama (cracking) riski yüksektir. Büküm radyüsünü (İç R) artırmanızı tavsiye ederiz.")
	}
	if materialLower == "stainless" && req.ThicknessMm > 3.0 {
		score -= 10
		advices = append(advices, "Bilgi: Kalın paslanmaz çelik bükümü, abkant preste yüksek tonaj (Springback) gerektirecektir. Tonaj limitlerinizi kontrol edin.")
	}

	// Rule 4: Complexity / Bend Count Penalty
	if req.BendCount > 10 {
		penalty := math.Min(float64(req.BendCount-10)*2.0, 20.0)
		score -= int(penalty)
		advices = append(advices, fmt.Sprintf("Uyarı: Büküm sayısı çok fazla (%d). Operatör hata payı ve birikimli tolerans sapması riski yüksek.", req.BendCount))
	}

	if score < 0 {
		score = 0
	}

	if score == 100 {
		advices = append(advices, "Tasarım üretilebilirlik (DFM) açısından kusursuz. Makine parkuru için herhangi bir risk tespit edilmedi.")
	} else if score > 70 && !warning {
		warning = false
	} else {
		warning = true
	}

	return DFMResponse{
		ManufacturabilityScore: score,
		AdviceList:             advices,
		Warning:                warning,
	}
}
