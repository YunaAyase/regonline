package handler

import (
	"regonline-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type OCRHandler struct {
	ocrService *service.OCRService
}

func NewOCRHandler(ocrService *service.OCRService) *OCRHandler {
	return &OCRHandler{
		ocrService: ocrService,
	}
}

type OCRResponse struct {
	Code    int                 `json:"code"`
	Data    *OCRRecognizeResult `json:"data,omitempty"`
	Message string              `json:"message,omitempty"`
}

type OCRRecognizeResult struct {
	IDNumber   string                       `json:"id_number"`
	Candidates []service.OCRNumberCandidate `json:"candidates,omitempty"`
}

func (h *OCRHandler) RecognizeID(c *gin.Context) {
	if !h.ocrService.IsAvailable() {
		c.JSON(200, OCRResponse{
			Code:    1,
			Message: "未配置 OCR 云识别服务，请在管理后台 → 网站设置 → OCR 云识别配置中配置",
		})
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(400, OCRResponse{
			Code:    1,
			Message: "未上传照片",
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(500, OCRResponse{
			Code:    1,
			Message: "读取照片失败",
		})
		return
	}
	defer f.Close()

	buf := make([]byte, file.Size)
	f.Read(buf)

	candidates, err := h.ocrService.RecognizeAllNumbers(buf)
	if err != nil {
		c.JSON(200, OCRResponse{
			Code:    1,
			Message: err.Error(),
		})
		return
	}

	var bestID string
	for _, c := range candidates {
		if c.IsIDNumber {
			bestID = c.Value
			break
		}
	}

	c.JSON(200, OCRResponse{
		Code: 0,
		Data: &OCRRecognizeResult{
			IDNumber:   bestID,
			Candidates: candidates,
		},
	})
}