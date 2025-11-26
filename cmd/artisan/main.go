package main

import (
	"fmt"
	"os"

	"go-artisan/cmd/artisan/commands" // 引入新的包
	"go-artisan/internal/config"

	"github.com/spf13/cobra"
)

func main() {
	// 1. 尝试加载配置
	// CLI 环境允许部分配置加载失败（比如只要运行 help），但在涉及 DB 时通过 validate 检查
	cfg, err := config.Load(".")
	if err != nil {
		// 这里我们不 panic，而是打印警告，因为用户可能正在执行不依赖配置的命令 (如 make:controller)
		fmt.Println("⚠️  Config load warning (ignore if running non-db commands):", err)
	}

	// 2. 根命令
	var rootCmd = &cobra.Command{
		Use:   "artisan",
		Short: "GoArtisan CLI Tool",
		Long:  "Command line utility for GoArtisan Framework",
	}

	// 3. 注册子命令
	// 将 Config 注入到需要的命令中
	rootCmd.AddCommand(
		commands.NewMakeControllerCommand(),
		commands.NewMakeMigrationCommand(),
		commands.NewMigrateCommand(cfg),         // 注入 Config
		commands.NewMigrateRollbackCommand(cfg), // 注入 Config
	)

	// 4. 执行
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
