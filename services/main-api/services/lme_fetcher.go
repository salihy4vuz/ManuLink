package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// LMEPrice represents London Metal Exchange price for a material
type LMEPrice struct {
	Material string  `json:"material"`
	PriceEUR float64 `json:"price_eur_per_kg"`
	Date     string  `json:"date"`
}

var lmeCache map[string]LMEPrice

func init() {
	lmeCache = make(map[string]LMEPrice)
}

// FetchLMELivePrices simulates fetching live prices from an external API, now actually uses MetalPriceAPI
func FetchLMELivePrices() error {
	log.Println("[LME Fetcher] Fetching live metal prices from MetalPriceAPI...")
	
	apiKey := os.Getenv("LME_API_KEY")
	if apiKey == "" {
		log.Println("[LME Fetcher] WARNING: LME_API_KEY is not set. Using fallback prices.")
		return fetchFallbackPrices()
	}

	url := fmt.Sprintf("https://api.metalpriceapi.com/v1/latest?api_key=%s&base=EUR&currencies=XAL,XCU,NI,ZNC", apiKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("[LME Fetcher] HTTP GET Error: %v", err)
		return fetchFallbackPrices()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[LME Fetcher] API returned status: %d", resp.StatusCode)
		return fetchFallbackPrices()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[LME Fetcher] Error reading response: %v", err)
		return fetchFallbackPrices()
	}

	var apiResponse struct {
		Success   bool               `json:"success"`
		Rates     map[string]float64 `json:"rates"`
		Timestamp int64              `json:"timestamp"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("[LME Fetcher] JSON Unmarshal Error: %v", err)
		return fetchFallbackPrices()
	}

	// XAL = Aluminum Ounce. Convert to kg (1 Ounce = 0.0283495 kg => 1 kg = 35.274 Ounces)
	// Base is EUR. Rate is XAL per EUR. So EUR per kg = (1 / Rate) * 35.274
	updateCacheIfPresent(apiResponse.Rates, "XAL", "Aluminum", 35.274)
	// XCU = Copper
	updateCacheIfPresent(apiResponse.Rates, "XCU", "Copper", 35.274)
	
	// Default/fallback for Mild Steel and Stainless if not provided directly
	lmeCache["Mild Steel"] = LMEPrice{Material: "Mild Steel", PriceEUR: 1.20, Date: time.Now().Format(time.RFC3339)}
	lmeCache["Stainless Steel"] = LMEPrice{Material: "Stainless Steel", PriceEUR: 3.50, Date: time.Now().Format(time.RFC3339)}
	
	return nil
}

func updateCacheIfPresent(rates map[string]float64, symbol string, matName string, multiplier float64) {
	if rate, ok := rates[symbol]; ok && rate > 0 {
		pricePerKg := (1.0 / rate) * multiplier
		lmeCache[matName] = LMEPrice{
			Material: matName,
			PriceEUR: pricePerKg,
			Date:     time.Now().Format(time.RFC3339),
		}
		log.Printf("[LME Fetcher] Updated %s: %.2f EUR/kg\n", matName, pricePerKg)
	}
}

func fetchFallbackPrices() error {
	log.Println("[LME Fetcher] Using Fallback Metal Prices")
	lmeCache["Mild Steel"] = LMEPrice{Material: "Mild Steel", PriceEUR: 1.25, Date: time.Now().Format(time.RFC3339)}
	lmeCache["Stainless Steel"] = LMEPrice{Material: "Stainless Steel", PriceEUR: 3.55, Date: time.Now().Format(time.RFC3339)}
	lmeCache["Aluminum"] = LMEPrice{Material: "Aluminum", PriceEUR: 2.85, Date: time.Now().Format(time.RFC3339)}
	return nil
}

// GetMaterialPrice gets the current price from cache
func GetMaterialPrice(material string) float64 {
	mat := material
	if material == "steel" { mat = "Mild Steel" }
	if material == "stainless" { mat = "Stainless Steel" }
	if material == "aluminum" { mat = "Aluminum" }

	if price, ok := lmeCache[mat]; ok {
		return price.PriceEUR
	}
	// Fallback prices if not in cache
	if mat == "Stainless Steel" { return 3.5 }
	if mat == "Aluminum" { return 2.8 }
	return 1.20 // Default for Mild Steel
}

// StartLMECronJob starts a background ticker that fetches LME prices daily.
func StartLMECronJob() {
	FetchLMELivePrices() // Initial fetch

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			FetchLMELivePrices()
		}
	}()
}

type LMEPrices struct {
	Prices []LMEPrice `json:"prices"`
}

func fetchFromLMEAPI() (LMEPrices, error) {
	err := FetchLMELivePrices()
	if err != nil {
		return LMEPrices{}, err
	}
	
	var prices []LMEPrice
	for _, p := range lmeCache {
		prices = append(prices, p)
	}
	return LMEPrices{Prices: prices}, nil
}

// FetchCurrentPrices önce Redis'e bak, yoksa API'ye git, sonucu Redis'e yaz
func FetchCurrentPrices(ctx context.Context) (LMEPrices, error) {
	// 1. Cache kontrol
	cached, err := GetCachedLMEPrices(ctx)
	if err == nil && len(cached) > 0 {
		var prices LMEPrices
		json.Unmarshal(cached, &prices)
		return prices, nil
	}
	// 2. API'den çek
	prices, err := fetchFromLMEAPI()
	if err != nil {
		return LMEPrices{}, err
	}
	// 3. Cache'e yaz
	CacheLMEPrices(ctx, prices)
	return prices, nil
}
