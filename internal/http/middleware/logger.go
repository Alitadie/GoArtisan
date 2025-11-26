package middleware

import (
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 1. 生成并传递 Trace ID
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)

		// 2. 处理请求
		c.Next()

		// 3. 记录响应日志
		duration := time.Since(start)
		status := c.Writer.Status()

		logAttr := []any{
			slog.String("trace_id", traceID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Duration("latency", duration),
			slog.String("ip", c.ClientIP()),
		}

		if status >= 500 {
			logger.Error("Request Failed", logAttr...)
		} else {
			logger.Info("Request Success", logAttr...)
		}
	}
}
