package models

import "gorm.io/gorm"

// PricingConfig tüm fiyatlandırma parametrelerini tenant bazında saklar.
// Bu model sayesinde admin panel üzerinden tüm sabit değerler DB'den yönetilebilir.
type PricingConfig struct {
	gorm.Model
	TenantID uint `gorm:"uniqueIndex;not null" json:"tenant_id"`

	// --- Faz 3.6: CAM Maliyetleri ---
	PunchCostPerHit    float64 `gorm:"default:0.05"  json:"punch_cost_per_hit"`    // $/hit
	HardwareCostUnit   float64 `gorm:"default:0.15"  json:"hardware_cost_unit"`    // $/fastener
	FilmCostPerM2      float64 `gorm:"default:2.50"  json:"film_cost_per_m2"`      // $/m²
	PaintCostPerM2     float64 `gorm:"default:15.00" json:"paint_cost_per_m2"`     // $/m²
	CMMSetupCost       float64 `gorm:"default:50.00" json:"cmm_setup_cost"`        // $ (setup)
	CMMUnitCost        float64 `gorm:"default:2.00"  json:"cmm_unit_cost"`         // $/part
	PackagingPalletCost float64 `gorm:"default:40.00" json:"packaging_pallet_cost"` // $ (setup)
	PackagingWoodCost  float64 `gorm:"default:80.00" json:"packaging_wood_cost"`   // $ (setup)
	PackagingCartonCost float64 `gorm:"default:2.00"  json:"packaging_carton_cost"` // $/unit
	LaserGasO2PerMin   float64 `gorm:"default:0.45"  json:"laser_gas_o2_per_min"`  // $/min
	LaserGasN2PerMin   float64 `gorm:"default:1.80"  json:"laser_gas_n2_per_min"`  // $/min
	LaserGasAirPerMin  float64 `gorm:"default:0.15"  json:"laser_gas_air_per_min"` // $/min
	WeldingCostPerMeter float64 `gorm:"default:10.00" json:"welding_cost_per_meter"` // $/m
	LogisticsCostPerKm float64 `gorm:"default:1.50"  json:"logistics_cost_per_km"` // $/km

	// --- Faz 3.7: Sekonder Operasyonlar ---
	EtchingCostPerMm   float64 `gorm:"default:0.01"  json:"etching_cost_per_mm"`   // $/mm
	BevelingCostPerMm  float64 `gorm:"default:0.05"  json:"beveling_cost_per_mm"`  // $/mm
	DeburringCostUnit  float64 `gorm:"default:1.00"  json:"deburring_cost_unit"`   // $/part
	HemmingCostPerOp   float64 `gorm:"default:1.00"  json:"hemming_cost_per_op"`   // $/op
	MicroJointCostUnit float64 `gorm:"default:0.10"  json:"micro_joint_cost_unit"` // $/joint
	AssemblyCostPerPart float64 `gorm:"default:5.00"  json:"assembly_cost_per_part"` // $/part
	HAZReportCost      float64 `gorm:"default:10.00" json:"haz_report_cost"`       // $ (setup)
	GrainDirectionPenaltyPct float64 `gorm:"default:10.0" json:"grain_direction_penalty_pct"` // %

	// --- Faz 3.8: Ticari Değişkenler ---
	TappingCostPerHole    float64 `gorm:"default:0.25"  json:"tapping_cost_per_hole"`     // $/hole
	CountersinkCostPerHole float64 `gorm:"default:0.15"  json:"countersink_cost_per_hole"` // $/hole
	FormingCostPerHit     float64 `gorm:"default:0.50"  json:"forming_cost_per_hit"`      // $/form
	ComplexBendingSetup   float64 `gorm:"default:30.00" json:"complex_bending_setup"`    // $ (setup)
	CornerReliefCostUnit  float64 `gorm:"default:0.02"  json:"corner_relief_cost_unit"`  // $/relief
	UnfoldingServiceCost  float64 `gorm:"default:15.00" json:"unfolding_service_cost"`   // $ (setup)
	MaterialCertCost      float64 `gorm:"default:25.00" json:"material_cert_cost"`       // $ (setup)
	RushOrderMarkupPct    float64 `gorm:"default:50.0"  json:"rush_order_markup_pct"`   // %

	// --- Hacim İndirimi Kademeleri ---
	VolumeDiscount100Pct  float64 `gorm:"default:5.0"  json:"volume_discount_100_pct"`  // % (>=100 qty)
	VolumeDiscount500Pct  float64 `gorm:"default:10.0" json:"volume_discount_500_pct"`  // % (>=500 qty)
	VolumeDiscount1000Pct float64 `gorm:"default:15.0" json:"volume_discount_1000_pct"` // % (>=1000 qty)

	// --- Varsayılan Margin ---
	DefaultWhiteLabelMarginPct float64 `gorm:"default:20.0" json:"default_white_label_margin_pct"` // %
}
