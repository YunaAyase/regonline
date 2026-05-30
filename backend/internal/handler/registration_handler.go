package handler

import (
	"fmt"
	"time"

	"regonline-backend/internal/service"
	"regonline-backend/internal/response"
	errs "regonline-backend/internal/error"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	regService *service.RegistrationService
}

func NewRegistrationHandler(regService *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		regService: regService,
	}
}

type CreateRegistrationRequest struct {
	Name        string `form:"student_name" binding:"required"`
	Gender      string `form:"gender" binding:"required,oneof=男 女"`
	BirthDate   string `form:"birth_date" binding:"required"`
	Grade       string `form:"grade" binding:"required"`
	ClassID     uint   `form:"class_id" binding:"required"`
	ParentName  string `form:"parent_name" binding:"required"`
	ParentPhone string `form:"parent_phone" binding:"required"`
	Address     string `form:"address" binding:"required"`
	IDNumber    string `form:"id_number" binding:"required"`
}

func (h *RegistrationHandler) Create(c *gin.Context) {
	var req CreateRegistrationRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		response.BadRequest(c, "出生日期格式错误，请使用 YYYY-MM-DD 格式")
		return
	}

	photo, photoName := getPhotoFromForm(c)

	createReq := &service.CreateRegistrationRequest{
		Name:        req.Name,
		Gender:      req.Gender,
		BirthDate:   birthDate,
		Grade:       req.Grade,
		ClassID:     req.ClassID,
		ParentName:  req.ParentName,
		ParentPhone: req.ParentPhone,
		Address:     req.Address,
		IDNumber:    req.IDNumber,
		Photo:       photo,
		PhotoName:   photoName,
	}

	registration, err := h.regService.Create(createReq)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Created(c, gin.H{
		"id":       registration.ID,
		"name":     registration.Name,
		"class_id": registration.ClassID,
		"message":  "报名成功！",
	})
}

func (h *RegistrationHandler) ListRegistrations(c *gin.Context) {
	name := c.Query("name")
	classID := c.Query("class_id")

	var cid uint
	if classID != "" {
		fmt.Sscanf(classID, "%d", &cid)
	}

	registrations, err := h.regService.SearchRegistrations(name, cid)
	if err != nil {
		response.InternalError(c, "获取报名列表失败")
		return
	}

	response.Success(c, registrations)
}

func (h *RegistrationHandler) GetRegistration(c *gin.Context) {
	id := c.Param("id")
	var regID uint
	fmt.Sscanf(id, "%d", &regID)

	registration, err := h.regService.GetRegistrationByID(regID)
	if err != nil {
		response.NotFound(c, "报名记录不存在")
		return
	}

	response.Success(c, registration)
}

func (h *RegistrationHandler) DeleteRegistration(c *gin.Context) {
	id := c.Param("id")
	var regID uint
	fmt.Sscanf(id, "%d", &regID)

	if err := h.regService.DeleteRegistration(regID); err != nil {
		response.InternalError(c, "删除报名记录失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

func getPhotoFromForm(c *gin.Context) ([]byte, string) {
	file, err := c.FormFile("photo")
	if err != nil {
		return nil, ""
	}

	f, err := file.Open()
	if err != nil {
		return nil, ""
	}
	defer f.Close()

	buf := make([]byte, file.Size)
	if _, err := f.Read(buf); err != nil {
		return nil, ""
	}

	return buf, file.Filename
}

func handleServiceError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errs.ValidationError:
		response.BadRequest(c, e.Error())
	case *errs.DuplicateError:
		response.Conflict(c, e.Error())
	case *errs.CapacityError:
		response.Conflict(c, e.Error())
	case *errs.NotFoundError:
		response.NotFound(c, e.Error())
	default:
		response.InternalError(c, "服务器内部错误: "+err.Error())
	}
}
