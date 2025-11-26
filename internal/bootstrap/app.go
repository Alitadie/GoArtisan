package bootstrap

import (
	"context"
	"fmt"

	"go-artisan/internal/config"
	"go-artisan/internal/http/handler"
	"go-artisan/internal/http/router"
	"go-artisan/internal/provider"
	"go-artisan/internal/repository"
	"go-artisan/internal/service"

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
)

func NewConfig() (*config.Config, error) {
	return config.Load("configs/config.yaml")
}

func NewLogger(cfg *config.Config) *slog.Logger {
	// 简单实现，企业级可换 Zap
	return slog.Default()
}

// Start 启动 HTTP Server
func Start(lifecycle fx.Lifecycle, cfg *config.Config, r *gin.Engine) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.App.Port)
			fmt.Printf("🚀 Artisan Server running on %s\n", addr)
			go func() {
				if err := r.Run(addr); err != nil {
					fmt.Printf("Error starting server: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("🛑 Stopping Server...")
			return nil
		},
	})
}
