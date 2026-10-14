package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB initializes the PostgreSQL database using GORM.
func ConnectDB() {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" { host = "localhost" }
	user := os.Getenv("POSTGRES_USER")
	if user == "" { user = "postgres" }
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" { pass = "postgres" }
	dbn := os.Getenv("POSTGRES_DB")
	if dbn == "" { dbn = "sheetmetal" }
	port := os.Getenv("POSTGRES_PORT")
	if port == "" { port = "5432" }
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul", host, user, pass, dbn, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Helpful for debugging slow queries
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL via GORM")
	DB = db
}

// AutoMigrate automates the creation of schemas based on Go models.
// Warning: In highly critical systems, we'd use golang-migrate, but GORM auto-migrate
// is acceptable for MVP and Phase 3.4.
func AutoMigrate(models ...interface{}) {
	if DB == nil {
		log.Fatal("Database not initialized")
	}
	
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}
	
	log.Println("🚀 Database schemas migrated successfully")
}
