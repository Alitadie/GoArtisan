package handler

import (
	"net/http"

	"go-artisan/internal/service"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc    *service.UserService
	logger *slog.Logger
}

func NewUserHandler(svc *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	// 1. 验证参数
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Register bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 调用业务服务
	user, err := h.svc.Register(service.RegisterDTO{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		h.logger.Error("Register failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 返回响应 (隐式屏蔽了 Password 字段)
	c.JSON(http.StatusCreated, user)
}
