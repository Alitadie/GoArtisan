package provider

import (
	"go-artisan/internal/config"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(NewDatabase),
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// 连接池配置应在此处设置
	return db, nil
}
