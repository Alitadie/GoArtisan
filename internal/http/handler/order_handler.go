package handler

import (
	"go-artisan/pkg/response"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type OrderHandler struct {
	logger *slog.Logger
	// 这里可以添加 service 依赖，例如: svc *service.OrderService
}

// NewOrderHandler 构造函数
func NewOrderHandler(logger *slog.Logger) *OrderHandler {
	return &OrderHandler{
		logger: logger,
	}
}

// Index 示例方法
func (h *OrderHandler) Index(c *gin.Context) {
	// 示例：使用统一响应
	h.logger.Info("Accessing Order Index")
	response.Success(c, gin.H{"module": "Order", "action": "index"})
}
