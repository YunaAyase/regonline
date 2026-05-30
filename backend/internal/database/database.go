package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"regonline-backend/internal/config"
	"regonline-backend/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dbPath := cfg.Database.Path

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	gormLogger := logger.Default.LogMode(logger.Silent)
	if cfg.Server.Mode == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)

	return db, nil
}

func AutoMigrate(db *gorm.DB, cfg *config.Config) error {
	if !cfg.Database.AutoMigrate {
		return nil
	}

	log.Println("Running database auto migration...")
	if err := db.AutoMigrate(&model.Class{}, &model.Registration{}, &model.Admin{}, &model.SiteSetting{}); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	log.Println("Database migration completed")
	return nil
}

func SeedClasses(db *gorm.DB, cfg *config.Config) error {
	if !cfg.Classes.SeedEnabled {
		return nil
	}

	var count int64
	db.Model(&model.Class{}).Count(&count)
	if count > 0 {
		log.Printf("Database already has %d classes, skipping seed", count)
		return nil
	}

	log.Println("Seeding class presets...")

	classes := make([]model.Class, 0, len(cfg.Classes.Presets))
	for _, preset := range cfg.Classes.Presets {
		classes = append(classes, model.Class{
			Name:        preset.Name,
			MaxStudents: preset.MaxStudents,
			MinAge:      preset.MinAge,
			MaxAge:      preset.MaxAge,
		})
	}

	if err := db.Create(&classes).Error; err != nil {
		return fmt.Errorf("failed to seed classes: %w", err)
	}

	log.Printf("Seeded %d class presets", len(classes))
	return nil
}