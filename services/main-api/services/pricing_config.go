package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/config"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

const pricingConfigCachePrefix = "pricing_config:tenant:"
const pricingConfigCacheTTL = 1 * time.Hour

// GetDefaultPricingConfig sistem genelinde kullanılan varsayılan değerleri döndürür.
// Bu değerler DB bağlantısı olmadan da fiyatlandırma motorunun çalışmasını sağlar.
func GetDefaultPricingConfig() models.PricingConfig {
	return models.PricingConfig{
		// Faz 3.6
		PunchCostPerHit:     0.05,
		HardwareCostUnit:    0.15,
		FilmCostPerM2:       2.50,
		PaintCostPerM2:      15.00,
		CMMSetupCost:        50.00,
		CMMUnitCost:         2.00,
		PackagingPalletCost: 40.00,
		PackagingWoodCost:   80.00,
		PackagingCartonCost: 2.00,
		LaserGasO2PerMin:    0.45,
		LaserGasN2PerMin:    1.80,
		LaserGasAirPerMin:   0.15,
		WeldingCostPerMeter: 10.00,
		LogisticsCostPerKm:  1.50,
		// Faz 3.7
		EtchingCostPerMm:         0.01,
		BevelingCostPerMm:        0.05,
		DeburringCostUnit:        1.00,
		HemmingCostPerOp:         1.00,
		MicroJointCostUnit:       0.10,
		AssemblyCostPerPart:      5.00,
		HAZReportCost:            10.00,
		GrainDirectionPenaltyPct: 10.0,
		// Faz 3.8
		TappingCostPerHole:     0.25,
		CountersinkCostPerHole: 0.15,
		FormingCostPerHit:      0.50,
		ComplexBendingSetup:    30.00,
		CornerReliefCostUnit:   0.02,
		UnfoldingServiceCost:   15.00,
		MaterialCertCost:       25.00,
		RushOrderMarkupPct:     50.0,
		// Hacim indirimi
		VolumeDiscount100Pct:  5.0,
		VolumeDiscount500Pct:  10.0,
		VolumeDiscount1000Pct: 15.0,
		// Margin
		DefaultWhiteLabelMarginPct: 20.0,
	}
}

// GetPricingConfig tenant ID'ye göre fiyatlandırma konfigürasyonunu döndürür.
// Önce Redis önbelleğine, sonra DB'ye bakar; bulunamazsa hardcoded default döner.
func GetPricingConfig(tenantID uint) models.PricingConfig {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("%s%d", pricingConfigCachePrefix, tenantID)

	// 1. Redis'ten dene
	if rc := GetRedisClient(); rc != nil {
		if data, err := rc.Get(ctx, cacheKey).Bytes(); err == nil {
			var cfg models.PricingConfig
			if json.Unmarshal(data, &cfg) == nil {
				return cfg
			}
		}
	}

	// 2. DB'den dene
	if config.DB != nil {
		var cfg models.PricingConfig
		if result := config.DB.Where("tenant_id = ?", tenantID).First(&cfg); result.Error == nil {
			// Redis'e kaydet
			if rc := GetRedisClient(); rc != nil {
				if data, err := json.Marshal(cfg); err == nil {
					rc.Set(ctx, cacheKey, data, pricingConfigCacheTTL)
				}
			}
			return cfg
		}
	}

	// 3. Fallback: sistem geneli default
	return GetDefaultPricingConfig()
}

// UpdatePricingConfig bir tenant'ın fiyatlandırma konfigürasyonunu günceller.
func UpdatePricingConfig(tenantID uint, updated models.PricingConfig) error {
	if config.DB == nil {
		return fmt.Errorf("database not connected")
	}

	updated.TenantID = tenantID
	result := config.DB.Where("tenant_id = ?", tenantID).
		Assign(updated).
		FirstOrCreate(&updated)

	if result.Error != nil {
		return fmt.Errorf("failed to update pricing config: %w", result.Error)
	}

	// Redis önbelleğini temizle
	if rc := GetRedisClient(); rc != nil {
		cacheKey := fmt.Sprintf("%s%d", pricingConfigCachePrefix, tenantID)
		rc.Del(context.Background(), cacheKey)
	}

	return nil
}

// CreateDefaultPricingConfig yeni bir tenant için varsayılan config oluşturur.
// Genellikle tenant onboarding akışında çağrılır.
func CreateDefaultPricingConfig(tenantID uint) error {
	if config.DB == nil {
		return nil // DB yok, sorun değil
	}
	cfg := GetDefaultPricingConfig()
	cfg.TenantID = tenantID
	return config.DB.FirstOrCreate(&cfg, "tenant_id = ?", tenantID).Error
}
