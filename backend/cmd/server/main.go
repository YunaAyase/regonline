package main

import (
	"fmt"
	"log"
	"path/filepath"

	"regonline-backend/internal/config"
	"regonline-backend/internal/database"
	"regonline-backend/internal/handler"
	"regonline-backend/internal/repository"
	"regonline-backend/internal/router"
	"regonline-backend/internal/service"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := config.EnsureDirs(cfg); err != nil {
		log.Fatalf("Failed to ensure directories: %v", err)
	}

	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.AutoMigrate(db, cfg); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	adminRepo := repository.NewAdminRepository(db)
	authService := service.NewAuthService(adminRepo)

	if adminRepo.AdminCount() == 0 {
		if err := authService.SeedAdmin("admin", "admin"); err != nil {
			log.Fatalf("Failed to seed admin user: %v", err)
		}

		log.Printf("============================================================")
		log.Printf("  Default admin account: admin")
		log.Printf("  Default admin password: admin")
		log.Printf("  Please change the password after first login!")
		log.Printf("============================================================")
	} else {
		log.Println("Admin user already exists, skipping seed")
		if err := authService.MigratePasswords(); err != nil {
			log.Fatalf("Failed to migrate admin passwords: %v", err)
		}
	}

	settingRepo := repository.NewSiteSettingRepository(db)
	settingService := service.NewSiteSettingService(settingRepo)

	if err := settingService.SeedDefaults(); err != nil {
		log.Fatalf("Failed to seed site settings: %v", err)
	}
	log.Println("Site settings initialized")

	classRepo := repository.NewClassRepository(db)
	regRepo := repository.NewRegistrationRepository(db)

	classService := service.NewClassService(classRepo, regRepo)
	regService := service.NewRegistrationService(regRepo, classRepo, cfg.Storage.PhotoDir)
	ocrService := service.NewOCRService(settingRepo)

	classHandler := handler.NewClassHandler(classService)
	regHandler := handler.NewRegistrationHandler(regService)
	ocrHandler := handler.NewOCRHandler(ocrService)
	statsHandler := handler.NewStatsHandler(regService)
	authHandler := handler.NewAuthHandler(authService)
	settingsHandler := handler.NewSiteSettingHandler(settingService)
	serverInfoHandler := handler.NewServerInfoHandler(cfg.Database.Path)
	photoHandler := handler.NewPhotoHandler(cfg.Storage.PhotoDir)
	backupHandler := handler.NewBackupHandler(cfg.Database.Path, filepath.Dir(cfg.Database.Path)+"/backups")
	resetDBHandler := handler.NewResetDBHandler(db, cfg.Database.Path, filepath.Dir(cfg.Database.Path)+"/backups")
	qrcodeHandler := handler.NewQRCodeHandler(fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port))
	serverIPHandler := handler.NewServerIPHandler()

	r := router.SetupRouter(classHandler, regHandler, ocrHandler, statsHandler, authHandler, settingsHandler, serverInfoHandler, serverIPHandler, photoHandler, backupHandler, resetDBHandler, qrcodeHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting at http://%s", addr)
	log.Printf("Health check: http://%s/health", addr)
	log.Printf("Classes API: http://%s/api/classes", addr)
	log.Printf("Registrations API: http://%s/api/registrations", addr)
	log.Printf("Auth API: http://%s/api/auth/login", addr)
	log.Printf("Settings API: http://%s/api/settings", addr)
	log.Printf("Server Info API: http://%s/api/server-info", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}