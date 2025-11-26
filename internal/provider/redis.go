package provider

import (
	"context"
	"fmt"
	"time"

	"go-artisan/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	// 创建 Redis 客户端配置
	opt := &redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username, // 👈 就算是空字符串，go-redis 也会处理好
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,

		// 生产环境建议配置以下超时和连接池参数，不要用默认值
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10, // 连接池大小，根据并发量调整
		MinIdleConns: 5,  // 最小空闲连接
	}

	client := redis.NewClient(opt)

	// Ping 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		// 返回带有具体上下文的错误
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
