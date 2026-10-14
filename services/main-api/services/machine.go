package services

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// MachineProfile represents the configuration and cost of a specific machine.
type MachineProfile struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Type               string  `json:"type"` // "laser", "press_brake"
	MaxThicknessMm     float64 `json:"max_thickness_mm"`
	HourlyRateUSD      float64 `json:"hourly_rate_usd"`
	SetupTimeMin       float64 `json:"setup_time_min"`
	BaseSpeedMmPerMin  float64 `json:"base_speed_mm_per_min"` // Used for lasers
}

var (
	localMachineDB = map[string]MachineProfile{
		"laser_2kw": {
			ID:                "laser_2kw",
			Name:              "2kW Fiber Laser",
			Type:              "laser",
			MaxThicknessMm:    6.0,
			HourlyRateUSD:     40.0,
			SetupTimeMin:      5.0,
			BaseSpeedMmPerMin: 2500.0,
		},
		"laser_6kw": {
			ID:                "laser_6kw",
			Name:              "6kW Fiber Laser",
			Type:              "laser",
			MaxThicknessMm:    20.0,
			HourlyRateUSD:     85.0,
			SetupTimeMin:      10.0,
			BaseSpeedMmPerMin: 6000.0,
		},
		"brake_100t": {
			ID:                "brake_100t",
			Name:              "100-Ton Press Brake",
			Type:              "press_brake",
			MaxThicknessMm:    10.0,
			HourlyRateUSD:     50.0,
			SetupTimeMin:      15.0,
			BaseSpeedMmPerMin: 0.0,
		},
	}
	machineMutex sync.RWMutex
)

// getMachineDB gets the machine catalog from Redis or LocalDB
func getMachineDB() map[string]MachineProfile {
	// Try Redis
	if redisClient != nil {
		cached, err := redisClient.Get(context.Background(), "machines:catalog").Bytes()
		if err == nil {
			var catalog map[string]MachineProfile
			if json.Unmarshal(cached, &catalog) == nil {
				return catalog
			}
		}
	}

	machineMutex.RLock()
	defer machineMutex.RUnlock()
	
	// Set Redis
	if redisClient != nil {
		bytes, _ := json.Marshal(localMachineDB)
		redisClient.Set(context.Background(), "machines:catalog", bytes, 1*time.Hour)
	}

	return localMachineDB
}

// SelectLaser selects the most cost-effective laser machine capable of cutting the given thickness.
func SelectLaser(thickness float64) (MachineProfile, error) {
	db := getMachineDB()
	var selected MachineProfile
	found := false

	for _, m := range db {
		if m.Type == "laser" && m.MaxThicknessMm >= thickness {
			// Select the cheapest capable machine
			if !found || m.HourlyRateUSD < selected.HourlyRateUSD {
				selected = m
				found = true
			}
		}
	}

	if !found {
		return MachineProfile{}, errors.New("no laser machine capable of cutting this thickness")
	}
	return selected, nil
}

// SelectPressBrake selects the appropriate press brake for bending.
func SelectPressBrake(thickness float64) (MachineProfile, error) {
	db := getMachineDB()
	var selected MachineProfile
	found := false

	for _, m := range db {
		if m.Type == "press_brake" && m.MaxThicknessMm >= thickness {
			if !found || m.HourlyRateUSD < selected.HourlyRateUSD {
				selected = m
				found = true
			}
		}
	}

	if !found {
		return MachineProfile{}, errors.New("no press brake capable of bending this thickness")
	}
	return selected, nil
}

// GetAllMachines returns the full catalog of machines.
func GetAllMachines() map[string]MachineProfile {
	db := getMachineDB()
	copyDB := make(map[string]MachineProfile)
	for k, v := range db {
		copyDB[k] = v
	}
	return copyDB
}

// UpdateMachine allows the Admin to update machine properties (e.g. Hourly Rate).
func UpdateMachine(id string, profile MachineProfile) error {
	machineMutex.Lock()
	defer machineMutex.Unlock()

	if _, exists := localMachineDB[id]; exists {
		localMachineDB[id] = profile
		// Clear cache
		if redisClient != nil {
			redisClient.Del(context.Background(), "machines:catalog")
		}
		return nil
	}
	return errors.New("machine not found")
}
