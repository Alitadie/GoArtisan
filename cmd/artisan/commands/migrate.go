package commands

import (
	"fmt"
	"os"

	"go-artisan/internal/config"

	_ "github.com/go-sql-driver/mysql" // 必须在这里导入驱动
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

// NewMigrateCommand 运行迁移
func NewMigrateCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			ensureDB(cfg) // 确保有配置

			db, err := goose.OpenDBWithDriver("mysql", cfg.Database.DSN)
			if err != nil {
				fmt.Printf("❌ Connection failed: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			if err := goose.Up(db, "migrations"); err != nil {
				fmt.Printf("❌ Migration failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Database migrated successfully")
		},
	}
}

// NewMigrateRollbackCommand 回滚
func NewMigrateRollbackCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate:rollback",
		Short: "Rollback the last migration",
		Run: func(cmd *cobra.Command, args []string) {
			ensureDB(cfg)

			db, err := goose.OpenDBWithDriver("mysql", cfg.Database.DSN)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			if err := goose.Down(db, "migrations"); err != nil {
				panic(err)
			}
			fmt.Println("✅ Rollback successful")
		},
	}
}

// NewMakeMigrationCommand 创建 SQL 文件 (无需数据库连接，无需 cfg，但为了统一风格可保留参数接口，这里省略 cfg 即可)
func NewMakeMigrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:migration [name]",
		Short: "Create a new migration file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			// goose Create 使用的是本地文件系统
			if err := goose.Create(nil, "migrations", name, "sql"); err != nil {
				fmt.Printf("❌ Failed to create migration: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// 辅助函数：检查配置是否为空
func ensureDB(cfg *config.Config) {
	if cfg == nil {
		fmt.Println("❌ Error: Database config is missing. Check your .env file.")
		os.Exit(1)
	}
}
