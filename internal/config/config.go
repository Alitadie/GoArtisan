package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 聚合所有配置项
// mapstructure 标签用于 viper 将 yaml 字段映射到结构体
type Config struct {
	App      App      `mapstructure:"app"`
	Database Database `mapstructure:"database"`
	Log      Log      `mapstructure:"log"`
}

type App struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"` // local, production
	Port int    `mapstructure:"port"`
	URL  string `mapstructure:"url"`
}

type Database struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type Log struct {
	Level  string `mapstructure:"level"`  // debug, info, error
	Format string `mapstructure:"format"` // json, text
}

// Load 读取配置文件的逻辑
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 1. 设置配置文件的路径和名称
	// 如果传入 configPath 为空，默认查找 configs/config.yaml
	v.SetConfigFile(configPath)

	// 2. 允许环境变量覆盖 (像 Laravel 的 .env)
	// 例如: export APP_PORT=9090 会覆盖 config.yaml 里的 port
	v.AutomaticEnv()
	// 将点号分隔符转换为下划线 (App.Port -> APP_PORT)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 3. 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 4. 将读取到的值映射到 Config 结构体
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}
