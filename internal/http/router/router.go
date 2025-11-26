package router

import (
	"go-artisan/internal/config"
	"go-artisan/internal/http/handler"
	"go-artisan/internal/http/middleware"

	"log/slog"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module 将 Handler 和 Router 导出给 FX 容器
var Module = fx.Options(
	// 注册 Handler (类似于 Laravel 的 Service Provider)
	fx.Provide(handler.NewWelcomeHandler),

	// 注册 Router 构造函数
	fx.Provide(NewRouter),
)

// NewRouter 生成并配置 Gin Engine
// Fx 会自动注入：配置对象、Logger、以及我们编写的 Handler
func NewRouter(
	cfg *config.Config,
	logger *slog.Logger,
	welcomeHandler *handler.WelcomeHandler,
	userHandler *handler.UserHandler, // <-- 新增注入参数
) *gin.Engine {

	// 设置运行模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 1. 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger)) // 自定义结构化日志中间件

	// 2. 路由定义 (像 Laravel 的 routes/web.php)
	api := r.Group("/api")
	{
		api.GET("/hello", welcomeHandler.Index)
		// 未来可以通过 go generate 自动往这里追加代码
		api.POST("/register", userHandler.Register) // <-- 注册路由
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
