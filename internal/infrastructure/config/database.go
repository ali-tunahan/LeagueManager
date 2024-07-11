package config

import (
	"LeagueManager/internal/domain/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("league_manager.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, err
	}

	// Perform migrations
	if err := db.AutoMigrate(&models.Team{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
		return nil, err
	}

	return db, nil
}
