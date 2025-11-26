package handler

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// WelcomeHandler 结构体
type WelcomeHandler struct {
	logger *slog.Logger
}

// NewWelcomeHandler 构造函数
// 你可以在这里声明需要注入 Service 或 Repository
func NewWelcomeHandler(logger *slog.Logger) *WelcomeHandler {
	return &WelcomeHandler{
		logger: logger,
	}
}

// Index 方法 (对应 GET /api/hello)
func (h *WelcomeHandler) Index(c *gin.Context) {
	// 使用结构化日志记录业务行为
	h.logger.Info("Hello API called", "client_ip", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to GoArtisan Framework!",
		"style":   "Laravel-like Developer Experience",
	})
}
