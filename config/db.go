package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *AppConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", cfg.DBHost, cfg.DBUserName, cfg.DBName, cfg.DBPort)
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", cfg.DBHost, cfg.DBUserName, cfg.DBUserPassword, cfg.DBName, cfg.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", cfg.DBUserName, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err := RunMigrations(dbURL); err != nil {
		log.Fatal("Migration Failed: \n", err.Error())
		os.Exit(1)
	}

	log.Println("🚀 Connected Successfully to the Database")
}
