package handler

import (
	"regexp"

	"regonline-backend/internal/service"
	"regonline-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateAccountRequest struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	admin, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	c.SetCookie("admin_token", admin.Username, 3600*24, "/", "", false, true)

	response.Success(c, gin.H{
		"username": admin.Username,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("admin_token", "", -1, "/", "", false, true)
	response.Success(c, gin.H{"message": "已退出登录"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	token, err := c.Cookie("admin_token")
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}

	if !isValidUsername(token) {
		response.Unauthorized(c, "无效会话")
		return
	}

	response.Success(c, gin.H{"username": token})
}

func (h *AuthHandler) InitStatus(c *gin.Context) {
	response.Success(c, gin.H{"admin_exists": true})
}

func (h *AuthHandler) UpdateAccount(c *gin.Context) {
	token, err := c.Cookie("admin_token")
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	admin, err := h.authService.UpdateAccount(token, req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.SetCookie("admin_token", admin.Username, 3600*24, "/", "", false, true)
	response.Success(c, gin.H{"username": admin.Username})
}

func isValidUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched && len(username) > 0
}