package handler

import (
	// 引入新包 (起别名避免冲突)

	"go-artisan/internal/service"
	"go-artisan/pkg/response"

	"log/slog"

	myvalidator "go-artisan/pkg/validator" // 引入新包 (起别名避免冲突)

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
	// 1. 参数绑定与验证
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Register validation error", "err", err)

		// 使用我们的翻译工具和响应包
		errMsgs := myvalidator.Translate(err)
		response.ValidationError(c, errMsgs)
		return
	}

	// 2. 调用服务
	user, err := h.svc.Register(service.RegisterDTO{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// 这里的 err 是业务错误（如：邮箱已存在）
		// 在 V2 中我们会定义统一的 ErrEmailExists 业务错误码
		h.logger.Error("Register service error", "err", err)
		response.Error(c, 500, err.Error())
		return
	}

	// 3. 成功返回
	// 注意：user 里面可能包含一些你不想要额外字段，V2 里我们会做 DTO->VO 转换
	response.Success(c, user)
}
