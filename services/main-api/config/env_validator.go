package config

import (
	"log"
	"os"
)

// RequiredEnvVars, uygulamanın çalışması için zorunlu ortam değişkenleri
var RequiredEnvVars = []string{
	"DB_HOST",
	"DB_NAME",
	"DB_USER",
	"DB_PASSWORD",
	"REDIS_HOST",
	"RABBITMQ_HOST",
	"RABBITMQ_USER",
	"RABBITMQ_PASSWORD",
	"JWT_SECRET",
	"MINIO_ENDPOINT",
	"MINIO_ACCESS_KEY",
	"MINIO_SECRET_KEY",
}

// OptionalEnvVars, eksikse uyarı basılan değişkenler
var OptionalEnvVars = []string{
	"LOGO_SQL_HOST",
	"MIKRO_SQL_HOST",
	"NETSIS_DB_DSN",
	"STRIPE_SECRET_KEY",
	"IYZICO_API_KEY",
	"GEMINI_API_KEY",
	"SMTP_HOST",
}

// ValidateEnv, ortam değişkenlerini doğrular; zorunlu olanlar eksikse panic atar.
func ValidateEnv() {
	var missing []string
	for _, key := range RequiredEnvVars {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		log.Fatalf("[Config] KRITIK: Zorunlu ortam degiskenleri eksik: %v\nLutfen .env dosyanizi veya K8s Secret'larinizi kontrol edin.", missing)
	}
	log.Println("[Config] ✅ Tum zorunlu ortam degiskenleri dogrulandi.")

	for _, key := range OptionalEnvVars {
		if os.Getenv(key) == "" {
			log.Printf("[Config] ⚠️ Istege bagli ortam degiskeni eksik (devre disi olacak): %s", key)
		}
	}
}
