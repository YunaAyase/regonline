package handler

import (
	"strconv"

	"regonline-backend/internal/service"
	"regonline-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	regService *service.RegistrationService
}

func NewStatsHandler(regService *service.RegistrationService) *StatsHandler {
	return &StatsHandler{
		regService: regService,
	}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.regService.GetStats()
	if err != nil {
		response.InternalError(c, "获取统计信息失败")
		return
	}

	type ClassStat struct {
		Name         string `json:"name"`
		CurrentCount int64  `json:"current_count"`
	}

	result := make([]ClassStat, 0, len(stats))
	for name, count := range stats {
		result = append(result, ClassStat{
			Name:         name,
			CurrentCount: count,
		})
	}

	response.Success(c, result)
}

func (h *StatsHandler) GetDateRange(c *gin.Context) {
	classIDStr := c.Query("class_id")
	if classIDStr == "" {
		minDate, maxDate := h.regService.GetDefaultDateRange()
		response.Success(c, gin.H{
			"min_date": minDate.Format("2006-01-02"),
			"max_date": maxDate.Format("2006-01-02"),
			"note":     "未指定班级，使用默认年龄范围 5-18 岁",
		})
		return
	}

	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "班级ID格式错误")
		return
	}

	minDate, maxDate, err := h.regService.GetDateRangeForClass(uint(classID))
	if err != nil {
		response.NotFound(c, "班级不存在")
		return
	}

	response.Success(c, gin.H{
		"min_date": minDate.Format("2006-01-02"),
		"max_date": maxDate.Format("2006-01-02"),
	})
}
