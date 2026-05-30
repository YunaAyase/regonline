package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	photoDir string
}

func NewPhotoHandler(photoDir string) *PhotoHandler {
	return &PhotoHandler{photoDir: photoDir}
}

func (h *PhotoHandler) ServePhoto(c *gin.Context) {
	filename := filepath.Base(c.Param("filename"))
	filePath := filepath.Join(h.photoDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(filePath)
}