package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv" // 1. 引入库
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN             string        `mapstructure:"dsn"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func Load(configPath string) (*Config, error) {
	// 2. 核心修复：手动加载 .env 文件
	// 尝试加载根目录下的 .env，如果在生产环境没这个文件可以忽略错误
	if err := godotenv.Load(); err != nil {
		// 如果只是文件不存在，可以允许（也许完全依赖系统 env）
		if !os.IsNotExist(err) {
			log.Println("⚠️ .env file found but failed to load:", err)
		}
	}

	v := viper.New()

	// 3. 配置搜索路径优化
	// 允许传入路径，同时兜底 configs/ 和 根目录
	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath("configs") // 自动找 configs/config.yaml
	v.AddConfigPath(".")

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 4. 读取 YAML 文件 (如果文件不存在，也不应该恐慌，可能全靠 ENV 配置)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// 如果没找到 config.yaml，不 return error，继续尝试读环境变量
		log.Println("⚠️ Config file not found, using environment variables only")
	}

	// 5. 设置环境变量替换规则
	// 让 yaml 里的 database.dsn 能够被 DATABASE_DSN 覆盖
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 6. 手动绑定，确保精准匹配 (双保险)
	// 将 Config 结构体的 database.dsn 字段 绑定到 DB_DSN 环境变量
	_ = v.BindEnv("database.dsn", "DB_DSN")
	_ = v.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	_ = v.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")

	// 7. 解析
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
