package services

import (
	"testing"
)

// --- Machine Service Unit Testleri ---

// Test 1: İnce malzeme → ucuz laser (2kW) seçilmeli
func TestSelectLaser_ThinMaterial_CheapLaser(t *testing.T) {
	laser, err := SelectLaser(2.0) // 2mm → laser_2kw (40$/h) yeterlil
	if err != nil {
		t.Fatalf("SelectLaser(2.0) failed: %v", err)
	}
	if laser.Type != "laser" {
		t.Errorf("Expected type=laser, got %s", laser.Type)
	}
	if laser.MaxThicknessMm < 2.0 {
		t.Errorf("Selected laser max_thickness (%f) must be >= 2.0", laser.MaxThicknessMm)
	}
	// 2mm için en ucuz (40$/h) seçilmeli
	if laser.HourlyRateUSD > 50.0 {
		t.Errorf("For 2mm, cheap laser (<=50$/h) expected, got %f$/h", laser.HourlyRateUSD)
	}
}

// Test 2: Kalın malzeme → güçlü laser (6kW) seçilmeli
func TestSelectLaser_ThickMaterial_PowerfulLaser(t *testing.T) {
	laser, err := SelectLaser(15.0) // 15mm → sadece 6kW (20mm kapasiteli) yeterli
	if err != nil {
		t.Fatalf("SelectLaser(15.0) failed: %v", err)
	}
	if laser.MaxThicknessMm < 15.0 {
		t.Errorf("Selected laser max_thickness (%f) must be >= 15.0", laser.MaxThicknessMm)
	}
}

// Test 3: Aşırı kalınlık → hata beklenir
func TestSelectLaser_ExtremeThickness_Error(t *testing.T) {
	_, err := SelectLaser(999.0)
	if err == nil {
		t.Error("expected error for thickness=999mm, got nil")
	}
}

// Test 4: Sıfır kalınlık → hata veya düşük kapasiteli laser
func TestSelectLaser_ZeroThickness_Handled(t *testing.T) {
	_, err := SelectLaser(0.0)
	// 0mm için de bir makine bulunabilir (tüm makineler 0'ı kesiyor), hata değil
	// Ama eğer hata da dönerse kabul edilebilir (implementasyon spesifik)
	_ = err // Sonuç ne olursa olsun panic olmamalı
}

// Test 5: Press Brake seçimi — desteklenen kalınlık
func TestSelectPressBrake_SupportedThickness(t *testing.T) {
	brake, err := SelectPressBrake(5.0) // 5mm → brake_100t (10mm max) yeterli
	if err != nil {
		t.Fatalf("SelectPressBrake(5.0) failed: %v", err)
	}
	if brake.Type != "press_brake" {
		t.Errorf("Expected type=press_brake, got %s", brake.Type)
	}
	if brake.MaxThicknessMm < 5.0 {
		t.Errorf("Selected brake max_thickness (%f) must be >= 5.0", brake.MaxThicknessMm)
	}
	if brake.HourlyRateUSD <= 0 {
		t.Errorf("HourlyRateUSD must be > 0, got %f", brake.HourlyRateUSD)
	}
	if brake.SetupTimeMin <= 0 {
		t.Errorf("SetupTimeMin must be > 0, got %f", brake.SetupTimeMin)
	}
}

// Test 6: Press Brake — desteklenmeyen kalınlık
func TestSelectPressBrake_UnsupportedThickness_Error(t *testing.T) {
	_, err := SelectPressBrake(999.0)
	if err == nil {
		t.Error("expected error for thickness=999mm press brake, got nil")
	}
}

// Test 7: GetAllMachines — tam katalog dönmeli
func TestGetAllMachines_ReturnsCatalog(t *testing.T) {
	db := GetAllMachines()
	if len(db) == 0 {
		t.Error("GetAllMachines returned empty catalog")
	}

	hasLaser := false
	hasBrake := false
	for _, m := range db {
		if m.Type == "laser" {
			hasLaser = true
		}
		if m.Type == "press_brake" {
			hasBrake = true
		}
		// Her makinenin geçerli alanları olmalı
		if m.ID == "" {
			t.Errorf("Machine has empty ID: %+v", m)
		}
		if m.HourlyRateUSD <= 0 {
			t.Errorf("Machine %s has invalid HourlyRateUSD: %f", m.ID, m.HourlyRateUSD)
		}
	}
	if !hasLaser {
		t.Error("Catalog must contain at least one laser machine")
	}
	if !hasBrake {
		t.Error("Catalog must contain at least one press_brake machine")
	}
}

// Test 8: GetAllMachines — dönen map orijinalden bağımsız kopya olmalı
func TestGetAllMachines_ReturnsCopy_NoMutation(t *testing.T) {
	db := GetAllMachines()
	originalCount := len(db)

	// Dönen map'e ekleme yapalım — orijinali etkilememeli
	db["__test_machine__"] = MachineProfile{ID: "__test__", Type: "laser", HourlyRateUSD: 1.0}

	db2 := GetAllMachines()
	if len(db2) != originalCount {
		t.Errorf("GetAllMachines returned map is not a copy: original=%d, after mutation=%d",
			originalCount, len(db2))
	}
}

// Test 9: UpdateMachine — geçerli ID güncelleme
func TestUpdateMachine_ValidID_Success(t *testing.T) {
	original := GetAllMachines()
	var testID string
	var originalRate float64
	for id, m := range original {
		testID = id
		originalRate = m.HourlyRateUSD
		break
	}
	if testID == "" {
		t.Skip("no machines in catalog to update")
	}

	updated := original[testID]
	updated.HourlyRateUSD = originalRate + 99.99
	err := UpdateMachine(testID, updated)
	if err != nil {
		t.Errorf("UpdateMachine(%s) failed: %v", testID, err)
	}

	// Geri döndür (test izolasyonu)
	updated.HourlyRateUSD = originalRate
	_ = UpdateMachine(testID, updated)
}

// Test 10: UpdateMachine — geçersiz ID → hata
func TestUpdateMachine_InvalidID_Error(t *testing.T) {
	err := UpdateMachine("nonexistent_machine_xyz", MachineProfile{})
	if err == nil {
		t.Error("expected error for invalid machine ID, got nil")
	}
}
