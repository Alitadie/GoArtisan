package config

import (
	"strings"
	"time"

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

func Load(path string) (*Config, error) {
	v := viper.New()

	// 1. 设置配置文件路径和格式
	v.AddConfigPath(path)     // e.g. "configs"
	v.AddConfigPath(".")      // 支持从根目录找
	v.SetConfigName("config") // 文件名 config.yaml
	v.SetConfigType("yaml")

	// 2. 读取 YAML 默认值
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err // 只有文件存在但格式错误才报错
		}
	}

	// 3. 环境变量覆盖 (.env 或 系统变量)
	// Laravel: DB_DSN -> Go: database.dsn
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("") // 如果需要前缀如 APP_，这里设置
	v.AutomaticEnv()

	// 4. 手动绑定环境变量映射 (解决 Viper 嵌套结构映射问题)
	// 这里将 yaml 里的 "database.dsn" 绑定到 环境变量 "DB_DSN"
	_ = v.BindEnv("database.dsn", "DB_DSN")
	_ = v.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	_ = v.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
