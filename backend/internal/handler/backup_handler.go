package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"regonline-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	dbPath     string
	backupDir  string
}

func NewBackupHandler(dbPath, backupDir string) *BackupHandler {
	return &BackupHandler{
		dbPath:    dbPath,
		backupDir: backupDir,
	}
}

func (h *BackupHandler) CreateBackup(c *gin.Context) {
	if err := os.MkdirAll(h.backupDir, 0755); err != nil {
		response.InternalError(c, "无法创建备份目录")
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("regonline_%s.db", timestamp)
	backupPath := filepath.Join(h.backupDir, backupFilename)

	src, err := os.Open(h.dbPath)
	if err != nil {
		response.InternalError(c, "无法打开数据库文件")
		return
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		response.InternalError(c, "无法创建备份文件")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		response.InternalError(c, "备份文件写入失败")
		return
	}

	if err := dst.Sync(); err != nil {
		response.InternalError(c, "备份文件同步失败")
		return
	}

	response.Success(c, gin.H{
		"message":   "数据库备份成功",
		"file":      backupFilename,
		"location":  backupPath,
	})
}