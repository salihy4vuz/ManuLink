package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

type MaterialPricing struct {
	DensityKgM3 float64 `json:"density_kg_m3"`
	PricePerKg  float64 `json:"price_per_kg"`
}

var (
	localMaterialDB = map[string]MaterialPricing{
		"steel": {
			DensityKgM3: 7850.0,
			PricePerKg:  1.20, // Base price in USD
		},
		"stainless": {
			DensityKgM3: 8000.0,
			PricePerKg:  3.50,
		},
		"aluminum": {
			DensityKgM3: 2700.0,
			PricePerKg:  2.80,
		},
	}
	lmeMutex sync.RWMutex
)

// GetMaterialPricing returns the density and price for a given material.
func GetMaterialPricing(material string) (MaterialPricing, error) {
	mat := strings.ToLower(material)

	// 1. Try Redis Cache
	if redisClient != nil {
		cached, err := redisClient.Get(context.Background(), "material_price:"+mat).Bytes()
		if err == nil {
			var pricing MaterialPricing
			if json.Unmarshal(cached, &pricing) == nil {
				return pricing, nil
			}
		}
	}

	// 2. Cache Miss - Acquire Distributed Lock to prevent Cache Stampede
	if redisClient != nil {
		pool := goredis.NewPool(redisClient)
		rs := redsync.New(pool)
		mutex := rs.NewMutex("lock:material_price:"+mat)
		
		if err := mutex.Lock(); err == nil {
			defer mutex.Unlock()

			// Double-check cache inside lock
			cached, err := redisClient.Get(context.Background(), "material_price:"+mat).Bytes()
			if err == nil {
				var pricing MaterialPricing
				if json.Unmarshal(cached, &pricing) == nil {
					return pricing, nil
				}
			}
		}
	}

	// 3. Fallback to Local/Postgres DB
	lmeMutex.RLock()
	defer lmeMutex.RUnlock()

	if data, exists := localMaterialDB[mat]; exists {
		// Dinamik LME fiyatını al
		dynamicPrice := GetMaterialPrice(mat)
		data.PricePerKg = dynamicPrice
		
		// Set to Redis
		if redisClient != nil {
			bytes, _ := json.Marshal(data)
			redisClient.Set(context.Background(), "material_price:"+mat, bytes, 15*time.Minute)
		}
		return data, nil
	}
	return MaterialPricing{}, errors.New("material not found in local pricing database")
}


// UpdateMaterialPrice allows the Admin to update the base price per Kg
func UpdateMaterialPrice(material string, newPrice float64) error {
	if newPrice <= 0 {
		return errors.New("price must be greater than zero")
	}

	mat := strings.ToLower(material)
	lmeMutex.Lock()
	defer lmeMutex.Unlock()

	if entry, exists := localMaterialDB[mat]; exists {
		entry.PricePerKg = newPrice
		localMaterialDB[mat] = entry
		return nil
	}
	return errors.New("material not found")
}

// GetAllMaterials returns the full catalog of materials and their current prices
func GetAllMaterials() map[string]MaterialPricing {
	lmeMutex.RLock()
	defer lmeMutex.RUnlock()

	// Return a copy to prevent external mutation
	copyDB := make(map[string]MaterialPricing)
	for k, v := range localMaterialDB {
		copyDB[k] = v
	}
	return copyDB
}
