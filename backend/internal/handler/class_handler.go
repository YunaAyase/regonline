package handler

import (
	"errors"
	"fmt"

	"regonline-backend/internal/service"
	"regonline-backend/internal/response"
	errs "regonline-backend/internal/error"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	classService *service.ClassService
}

func NewClassHandler(classService *service.ClassService) *ClassHandler {
	return &ClassHandler{
		classService: classService,
	}
}

type ClassListResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	MaxStudents  int    `json:"max_students"`
	MinAge       int    `json:"min_age"`
	MaxAge       int    `json:"max_age"`
	Enabled      bool   `json:"enabled"`
	CurrentCount int64  `json:"current_count"`
}

func (h *ClassHandler) ListClasses(c *gin.Context) {
	classes, err := h.classService.ListClasses()
	if err != nil {
		response.InternalError(c, "获取班级列表失败")
		return
	}

	result := make([]ClassListResponse, len(classes))
	for i, class := range classes {
		result[i] = ClassListResponse{
			ID:           class.ID,
			Name:         class.Name,
			Description:  class.Description,
			MaxStudents:  class.MaxStudents,
			MinAge:       class.MinAge,
			MaxAge:       class.MaxAge,
			Enabled:      class.Enabled,
			CurrentCount: class.CurrentCount,
		}
	}

	response.Success(c, result)
}

type ClassInfoResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	MaxStudents  int    `json:"max_students"`
	MinAge       int    `json:"min_age"`
	MaxAge       int    `json:"max_age"`
	CurrentCount int64  `json:"current_count"`
}

func (h *ClassHandler) GetClassInfo(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "请提供班级名称参数: ?name=班级名")
		return
	}

	classInfo, err := h.classService.GetClassByName(name)
	if err != nil {
		var notFoundErr *errs.NotFoundError
		if errors.As(err, &notFoundErr) {
			response.NotFound(c, "班级不存在")
			return
		}
		response.InternalError(c, "获取班级信息失败")
		return
	}

	response.Success(c, ClassInfoResponse{
		ID:           classInfo.ID,
		Name:         classInfo.Name,
		MaxStudents:  classInfo.MaxStudents,
		MinAge:       classInfo.MinAge,
		MaxAge:       classInfo.MaxAge,
		CurrentCount: classInfo.CurrentCount,
	})
}

type CreateClassRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description"`
	MaxStudents int    `json:"max_students" binding:"required,min=1"`
	MinAge      int    `json:"min_age" binding:"required,min=0"`
	MaxAge      int    `json:"max_age" binding:"required,min=0"`
}

func (h *ClassHandler) CreateClass(c *gin.Context) {
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	classInfo, err := h.classService.CreateClass(req.Name, req.Description, req.MaxStudents, req.MinAge, req.MaxAge)
	if err != nil {
		var dupErr *errs.DuplicateError
		if errors.As(err, &dupErr) {
			response.BadRequest(c, "班级名称已存在")
			return
		}
		response.InternalError(c, "创建班级失败")
		return
	}

	response.Success(c, ClassListResponse{
		ID:           classInfo.ID,
		Name:         classInfo.Name,
		Description:  classInfo.Description,
		MaxStudents:  classInfo.MaxStudents,
		MinAge:       classInfo.MinAge,
		MaxAge:       classInfo.MaxAge,
		Enabled:      classInfo.Enabled,
		CurrentCount: classInfo.CurrentCount,
	})
}

type UpdateClassRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description"`
	MaxStudents int    `json:"max_students" binding:"required,min=1"`
	MinAge      int    `json:"min_age" binding:"required,min=0"`
	MaxAge      int    `json:"max_age" binding:"required,min=0"`
}

func (h *ClassHandler) UpdateClass(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.BadRequest(c, "无效的班级ID")
		return
	}

	var req UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	classInfo, err := h.classService.UpdateClassFull(id, req.Name, req.Description, req.MaxStudents, req.MinAge, req.MaxAge)
	if err != nil {
		var notFoundErr *errs.NotFoundError
		if errors.As(err, &notFoundErr) {
			response.NotFound(c, "班级不存在")
			return
		}
		var dupErr *errs.DuplicateError
		if errors.As(err, &dupErr) {
			response.BadRequest(c, "班级名称已存在")
			return
		}
		response.InternalError(c, "更新班级失败")
		return
	}

	response.Success(c, ClassListResponse{
		ID:           classInfo.ID,
		Name:         classInfo.Name,
		Description:  classInfo.Description,
		MaxStudents:  classInfo.MaxStudents,
		MinAge:       classInfo.MinAge,
		MaxAge:       classInfo.MaxAge,
		Enabled:      classInfo.Enabled,
		CurrentCount: classInfo.CurrentCount,
	})
}

type ToggleClassRequest struct {
	Enabled bool `json:"enabled"`
}

func (h *ClassHandler) ToggleClass(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.BadRequest(c, "无效的班级ID")
		return
	}

	var req ToggleClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	classInfo, err := h.classService.ToggleClass(id, req.Enabled)
	if err != nil {
		var notFoundErr *errs.NotFoundError
		if errors.As(err, &notFoundErr) {
			response.NotFound(c, "班级不存在")
			return
		}
		response.InternalError(c, "更新班级状态失败")
		return
	}

	response.Success(c, ClassListResponse{
		ID:           classInfo.ID,
		Name:         classInfo.Name,
		Description:  classInfo.Description,
		MaxStudents:  classInfo.MaxStudents,
		MinAge:       classInfo.MinAge,
		MaxAge:       classInfo.MaxAge,
		Enabled:      classInfo.Enabled,
		CurrentCount: classInfo.CurrentCount,
	})
}

func (h *ClassHandler) DeleteClass(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.BadRequest(c, "无效的班级ID")
		return
	}

	if err := h.classService.DeleteClass(id); err != nil {
		var valErr *errs.ValidationError
		if errors.As(err, &valErr) {
			response.BadRequest(c, err.Error())
			return
		}
		var notFoundErr *errs.NotFoundError
		if errors.As(err, &notFoundErr) {
			response.NotFound(c, "班级不存在")
			return
		}
		response.InternalError(c, "删除班级失败")
		return
	}

	response.Success(c, gin.H{"message": "班级已删除"})
}

func parseIDParam(c *gin.Context) (uint, error) {
	var id uint
	_, err := fmt.Sscan(c.Param("id"), &id)
	return id, err
}
