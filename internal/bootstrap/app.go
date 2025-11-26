package bootstrap

import (
	"context"
	"fmt"

	"go-artisan/internal/config"
	"go-artisan/internal/http/router"
	"go-artisan/internal/provider"

	"log/slog"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module 将核心模块打包
var Module = fx.Options(
	fx.Provide(NewConfig), // 注入配置
	fx.Provide(NewLogger), // 注入日志
	provider.Module,       // 注入 DB/Redis
	router.Module,         // 注入路由与 Handler
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
