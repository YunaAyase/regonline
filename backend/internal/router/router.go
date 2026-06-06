package router

import (
	"regonline-backend/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	classHandler *handler.ClassHandler,
	regHandler *handler.RegistrationHandler,
	ocrHandler *handler.OCRHandler,
	statsHandler *handler.StatsHandler,
	authHandler *handler.AuthHandler,
	settingsHandler *handler.SiteSettingHandler,
	serverInfoHandler *handler.ServerInfoHandler,
	serverIPHandler *handler.ServerIPHandler,
	photoHandler *handler.PhotoHandler,
	backupHandler *handler.BackupHandler,
	resetDBHandler *handler.ResetDBHandler,
	qrcodeHandler *handler.QRCodeHandler,
) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/logout", authHandler.Logout)
		api.GET("/auth/me", authHandler.Me)
		api.GET("/auth/init-status", authHandler.InitStatus)
		api.PUT("/admin/account", authHandler.UpdateAccount)

		api.GET("/classes", classHandler.ListClasses)
		api.GET("/class/info", classHandler.GetClassInfo)
		api.POST("/classes", classHandler.CreateClass)
		api.PUT("/classes/:id", classHandler.UpdateClass)
		api.PUT("/classes/:id/toggle", classHandler.ToggleClass)
		api.DELETE("/classes/:id", classHandler.DeleteClass)

		api.POST("/registrations", regHandler.Create)
		api.GET("/registrations", regHandler.ListRegistrations)
		api.GET("/registrations/:id", regHandler.GetRegistration)
		api.DELETE("/registrations/:id", regHandler.DeleteRegistration)

		api.POST("/ocr/recognize", ocrHandler.RecognizeID)

		api.GET("/stats", statsHandler.GetStats)
		api.GET("/date-range", statsHandler.GetDateRange)

		api.GET("/settings", settingsHandler.GetSettings)
		api.PUT("/settings", settingsHandler.UpdateSettings)

		api.GET("/server-info", serverInfoHandler.GetServerInfo)
		api.GET("/server-ip", serverIPHandler.GetServerIP)

		api.GET("/photos/:filename", photoHandler.ServePhoto)

		api.POST("/backup", backupHandler.CreateBackup)
		api.POST("/reset-db", resetDBHandler.ResetDatabase)
		api.POST("/qrcode/generate", qrcodeHandler.Generate)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}