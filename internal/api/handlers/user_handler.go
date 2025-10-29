package handlers

import (
	"net/http"

	"github.com/Anning01/user-management/internal/domain"
	"github.com/Anning01/user-management/internal/service"
	"github.com/Anning01/user-management/internal/util"
	"github.com/Anning01/user-management/pkg/security"

	"github.com/Anning01/user-management/internal/config"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
	jwtConfig   *config.JWTConfig
}

func NewUserHandler(userService service.UserService, jwtConfig *config.JWTConfig) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtConfig:   jwtConfig,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 数据验证
	if err := util.ValidateStruct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 数据验证
	if err := util.ValidateStruct(loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := security.GenerateToken(user.ID, h.jwtConfig.SecretKey, h.jwtConfig.ExpirationHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// 其他控制器方法...
