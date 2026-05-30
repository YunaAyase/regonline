package handler

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"regonline-backend/internal/response"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

type ServerInfoHandler struct {
	dbPath string
}

func NewServerInfoHandler(dbPath string) *ServerInfoHandler {
	return &ServerInfoHandler{dbPath: dbPath}
}

func (h *ServerInfoHandler) GetServerInfo(c *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	uptime := int64(time.Since(startTime).Seconds())

	var dbSize int64
	if info, err := os.Stat(h.dbPath); err == nil {
		dbSize = info.Size()
	}

	response.Success(c, gin.H{
		"go_version":    runtime.Version(),
		"os":            runtime.GOOS,
		"arch":          runtime.GOARCH,
		"cpu_cores":     runtime.NumCPU(),
		"goroutines":    runtime.NumGoroutine(),
		"memory_alloc":  fmt.Sprintf("%.2f MB", float64(mem.Alloc)/1024/1024),
		"uptime_seconds": uptime,
		"db_size":       fmt.Sprintf("%.2f KB", float64(dbSize)/1024),
	})
}