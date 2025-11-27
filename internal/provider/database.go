package provider

import (
	"go-artisan/internal/config"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(NewDatabase),
	fx.Provide(NewRedis),          // 👈 注册 Redis
	fx.Provide(NewCasbinEnforcer), // 👈 注册 Casbin
)

// NewDatabase 负责初始化 DB 并设置连接池参数
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
		// 关闭默认事务以提升性能（业务层按需开启）
		SkipDefaultTransaction: true,
		// 准备语句缓存，类似 PreparedStatement，提升重复 SQL 执行效率
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 连接池核心配置
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	// 防止连接持有太久导致 MySQL 服务器端超时断开
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db, nil
}
