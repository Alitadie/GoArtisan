package router

import (
	"go-artisan/internal/config"
	"go-artisan/internal/http/handler"
	"go-artisan/internal/http/middleware"
	"go-artisan/pkg/response"

	"log/slog"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module 将 Handler 和 Router 导出给 FX 容器
var Module = fx.Options(

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
	r.Use(middleware.VersionMiddleware())      // 👈 新增

	// 公开路由
	public := r.Group("/api")
	{
		public.GET("/hello", func(ctx *gin.Context) {
			response.Success(ctx, gin.H{"status": "public"})
		})
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login) // 👈 新增
	}

	// 保护路由 (类似 Laravel Route::middleware('auth:api'))
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/user/profile", func(c *gin.Context) {
			// 获取中间件塞入的 userID
			uid, _ := c.Get("userID")
			response.Success(c, gin.H{
				"message": "You are accessing protected data",
				"your_id": uid,
			})
		})
	}

	return r
}
