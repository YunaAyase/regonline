package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"regonline-backend/internal/config"
	"regonline-backend/internal/database"
	"regonline-backend/internal/model"
	"regonline-backend/internal/repository"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestServiceForClass(t *testing.T) (*ClassService, *gorm.DB) {
	tmpFile := filepath.Join(t.TempDir(), "test.db")

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver:      "sqlite",
			Path:        tmpFile,
			AutoMigrate: true,
		},
	}

	db, err := database.NewDB(cfg)
	require.NoError(t, err)

	err = database.AutoMigrate(db, cfg)
	require.NoError(t, err)

	classRepo := repository.NewClassRepository(db)
	regRepo := repository.NewRegistrationRepository(db)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		time.Sleep(100 * time.Millisecond)
		os.Remove(tmpFile)
	})

	return NewClassService(classRepo, regRepo), db
}

func seedClass(t *testing.T, db *gorm.DB, name string, maxStudents, minAge, maxAge int) uint {
	class := model.Class{
		Name:        name,
		MaxStudents: maxStudents,
		MinAge:      minAge,
		MaxAge:      maxAge,
	}
	db.Create(&class)
	return class.ID
}
