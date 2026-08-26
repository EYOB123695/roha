package initializers

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // Cache prepared statements to avoid re-parsing SQL
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure underlying SQL connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}

	sqlDB.SetMaxOpenConns(50)                  // Cap open connections to avoid saturating DB pooler
	sqlDB.SetMaxIdleConns(25)                  // Keep warm connections ready for reuse
	sqlDB.SetConnMaxLifetime(5 * time.Minute)  // Recycle connections before server-side timeouts
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)  // Close idle connections after inactivity
}

