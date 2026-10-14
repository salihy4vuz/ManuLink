package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
	"os"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

// GetRedisClient returns the global redis client instance
func GetRedisClient() *redis.Client {
	return redisClient
}

// CacheLMEPrices, LME metal fiyatlarını 15 dakika cache'ler.
func CacheLMEPrices(ctx context.Context, prices interface{}) error {
	ttl, _ := time.ParseDuration(os.Getenv("REDIS_LME_TTL_SECONDS") + "s")
	if ttl == 0 { ttl = 15 * time.Minute }
	data, _ := json.Marshal(prices)
	return redisClient.Set(ctx, "lme:prices:current", data, ttl).Err()
}

func GetCachedLMEPrices(ctx context.Context) ([]byte, error) {
	return redisClient.Get(ctx, "lme:prices:current").Bytes()
}

// CachePriceCalculation, fiyat hesaplama sonuçlarını 5 dakika cache'ler.
func CachePriceCalculation(ctx context.Context, key string, result interface{}) error {
	ttl := 5 * time.Minute
	data, _ := json.Marshal(result)
	return redisClient.Set(ctx, "price:"+key, data, ttl).Err()
}

func GetCachedPrice(ctx context.Context, key string) ([]byte, error) {
	return redisClient.Get(ctx, "price:"+key).Bytes()
}
