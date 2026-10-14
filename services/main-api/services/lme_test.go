package services

import (
	"testing"
)

// TestLMEFetch: GetExchangeRate fonksiyonunun bilinen kurlar icin dogru calistigi dogrulanir.
func TestLMEFetch(t *testing.T) {
	// USD kuru her zaman 1.0 olmali
	rate := GetExchangeRate("USD")
	if rate != 1.0 {
		t.Errorf("GetExchangeRate(USD) expected 1.0, got %f", rate)
	}

	// TRY kuru 1'den buyuk olmali (USD/TRY > 1)
	rateTRY := GetExchangeRate("TRY")
	if rateTRY <= 1.0 {
		t.Errorf("GetExchangeRate(TRY) expected > 1.0, got %f (TRY/USD exchange check)", rateTRY)
	}

	// EUR kuru 0 olmamali
	rateEUR := GetExchangeRate("EUR")
	if rateEUR <= 0 {
		t.Errorf("GetExchangeRate(EUR) expected > 0, got %f", rateEUR)
	}
}

// TestLMECache: Bilinmeyen para birimi icin GetExchangeRate 1.0 (fallback) dondurmeli.
func TestLMECache(t *testing.T) {
	rateFallback := GetExchangeRate("XYZ_UNKNOWN")
	if rateFallback != 1.0 {
		t.Errorf("GetExchangeRate(unknown) should fallback to 1.0, got %f", rateFallback)
	}
}

