package handler

import (
	"os"
	"time"

	"regonline-backend/internal/model"
	"regonline-backend/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResetDBHandler struct {
	db       *gorm.DB
	dbPath   string
	backupDir string
}

func NewResetDBHandler(db *gorm.DB, dbPath, backupDir string) *ResetDBHandler {
	return &ResetDBHandler{
		db:        db,
		dbPath:    dbPath,
		backupDir: backupDir,
	}
}

func (h *ResetDBHandler) ResetDatabase(c *gin.Context) {
	if _, err := os.Stat(h.dbPath); os.IsNotExist(err) {
		response.BadRequest(c, "数据库文件不存在")
		return
	}

	backupFilename := "before_reset_" + time.Now().Format("20060102_150405") + ".db"
	backupPath := h.backupDir + "/" + backupFilename

	if err := os.MkdirAll(h.backupDir, 0755); err != nil {
		response.InternalError(c, "无法创建备份目录")
		return
	}

	if err := os.Rename(h.dbPath, backupPath); err != nil {
		response.InternalError(c, "重置前备份失败: "+err.Error())
		return
	}

	if err := h.db.AutoMigrate(&model.Admin{}, &model.Class{}, &model.Registration{}, &model.SiteSetting{}); err != nil {
		response.InternalError(c, "数据库结构重建失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":  "数据库已重置",
		"backup":   backupFilename,
		"location": backupPath,
	})
}