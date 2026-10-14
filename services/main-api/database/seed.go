package database

import (
	"log"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

// SeedDatabase inserts default values for the MVP Multi-Tenant architecture
func SeedDatabase() {
	if DB == nil {
		log.Fatal("Database not initialized")
	}

	var tenantCount int64
	DB.Model(&models.Tenant{}).Count(&tenantCount)

	if tenantCount == 0 {
		log.Println("🌱 Seeding default Tenant and Admin User...")
		
		tenant := models.Tenant{
			Name:   "Apex Sheet Metal",
			Domain: "localhost", // For local testing
			Theme: models.TenantTheme{
				PrimaryColor:   "#10b981",
				SecondaryColor: "#0f172a",
				LogoURL:        "/vite.svg",
			},
		}
		DB.Create(&tenant)

		admin := models.User{
			TenantID: tenant.ID,
			Email:    "admin@apex.com",
			Password: "hashed_password", // In production use bcrypt
			Name:     "System Administrator",
			Role:     "admin",
		}
		DB.Create(&admin)

		log.Println("✅ Seeding complete!")
	} else {
		log.Println("ℹ️ Database already seeded, skipping.")
	}
}
