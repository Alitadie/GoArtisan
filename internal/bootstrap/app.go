package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-artisan/internal/config"
	"go-artisan/internal/http/handler"
	"go-artisan/internal/http/router"
	"go-artisan/internal/provider"
	"go-artisan/internal/repository"
	"go-artisan/internal/service"
	"go-artisan/pkg/validator"
	"go-artisan/pkg/version"

	"log/slog"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// RepositoryModule 定义仓储层的所有注入
var RepositoryModule = fx.Options(
	fx.Provide(repository.NewUserRepo),
)

// ServiceModule 定义服务层的所有注入
var ServiceModule = fx.Options(
	fx.Provide(service.NewUserService),
)

// HandlerModule 定义控制器层
var HandlerModule = fx.Options(
	fx.Provide(handler.NewWelcomeHandler), // 原来的
	fx.Provide(handler.NewUserHandler),    // 新增的
)

var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewLogger),
	provider.Module, // DB

	RepositoryModule, // 注入 Repo
	ServiceModule,    // 注入 Service
	HandlerModule,    // 注入 Handler

	router.Module, // 注入 Router (它现在依赖上面的 Handlers)

	// fx.Invoke(Start), // 调用启动逻辑
)

func NewConfig() (*config.Config, error) {
	return config.Load("configs/config.yaml")
}

func NewLogger(cfg *config.Config) *slog.Logger {
	// 简单实现，企业级可换 Zap
	return slog.Default()
}

// Start 启动 HTTP Server 现在变得更强壮
func Start(lifecycle fx.Lifecycle, cfg *config.Config, r *gin.Engine) {

	// 核心修复点：在这里调用独立的初始化
	validator.Init()

	// 打印版本信息 (炫酷一点)
	fmt.Println("---------------------------------------------------------")
	fmt.Printf("🚀 App: %s  Env: %s\n", cfg.App.Name, cfg.App.Env)
	fmt.Println(version.FullVersion())
	fmt.Println("---------------------------------------------------------")

	// 构造 HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: r,
		// 生产环境必须设置超时，防止 Slowloris 攻击
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	// 注册生命周期钩子
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 在 Goroutine 中启动服务器，因为 srv.ListenAndServe 是阻塞的
			// 如果在 OnStart 里直接调，会卡死整个 Fx 容器
			go func() {
				fmt.Printf("🌐 Serving on port %d\n", cfg.App.Port)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					fmt.Printf("❌ Server failed: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("🛑 Interrupt signal received...")
			fmt.Println("⏳ Waiting for active connections to finish...")

			// 这里创建一个带有超时的上下文
			// 意思：给你 5 秒钟处理正在进行的请求，处理完就停；如果 5 秒还在忙，强制杀掉。
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("server shutdown failed: %w", err)
			}

			fmt.Println("✅ Server exited gracefully")
			return nil
		},
	})
}
