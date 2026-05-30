package handler

import (
	"regonline-backend/internal/response"
	"regonline-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SiteSettingHandler struct {
	service *service.SiteSettingService
}

func NewSiteSettingHandler(service *service.SiteSettingService) *SiteSettingHandler {
	return &SiteSettingHandler{service: service}
}

func (h *SiteSettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.GetSettings()
	if err != nil {
		response.InternalError(c, "获取设置失败")
		return
	}
	response.Success(c, settings)
}

func (h *SiteSettingHandler) UpdateSettings(c *gin.Context) {
	var updates map[string]string
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	allowedKeys := map[string]bool{
		"site_name":        true,
		"site_description": true,
		"icp_record":       true,
		"copyright":        true,
	}

	for key := range updates {
		if !allowedKeys[key] {
			response.BadRequest(c, "无效的设置项: "+key)
			return
		}
	}

	if err := h.service.UpdateSettings(updates); err != nil {
		response.InternalError(c, "保存设置失败")
		return
	}

	settings, _ := h.service.GetSettings()
	response.Success(c, settings)
}