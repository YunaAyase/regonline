package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type QRCodeHandler struct {
	apiBase string
}

func NewQRCodeHandler(apiBase string) *QRCodeHandler {
	return &QRCodeHandler{apiBase: apiBase}
}

type QRCodeRequest struct {
	URL string `json:"url"`
}

func (h *QRCodeHandler) Generate(c *gin.Context) {
	var req QRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	content := req.URL
	if content == "" {
		content = h.apiBase
	}

	pngData, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR code"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename=\"qrcode.png\"")
	c.Data(http.StatusOK, "image/png", pngData)
}
