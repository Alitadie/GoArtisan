package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-artisan/internal/config"
	"go-artisan/internal/domain"
	"go-artisan/internal/domain/mocks"
	"go-artisan/internal/service"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"go.uber.org/mock/gomock"
)

// setupRedisContainer 负责启动一个真实的 Redis Docker 容器
func setupRedisContainer(t *testing.T) (*redis.Client, func()) {
	ctx := context.Background()

	// 1. 启动容器 (使用 redis:7-alpine 镜像)
	redisContainer, err := rediscontainer.Run(ctx, "docker.io/redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start redis container: %s", err)
	}

	// 2. 获取容器的连接地址 (TestContainers 会把端口映射到本机随机端口)
	endpoint, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("failed to get redis endpoint: %s", err)
	}

	// 3. 创建客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	// 4. 返回 Cleanup 函数 (关闭容器)
	return rdb, func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate redis container: %s", err)
		}
	}
}

func TestUserService_GetUserProfile_Integration(t *testing.T) {
	// ⚠️ 这一步会通过 Docker 启动 Redis，第一次运行可能需要下载镜像，稍微慢一点
	realRedis, cleanup := setupRedisContainer(t)
	defer cleanup() // 测试结束后，容器会被杀死，完全不污染环境

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 依赖准备
	mockRepo := mocks.NewMockUserRepository(ctrl)
	// Integration Test 中我们通常不 mock config，直接造个结构体
	mockConfig := &config.Config{
		App: config.AppConfig{Name: "IntegrationTestApp"},
	}

	svc := service.NewUserService(mockRepo, mockConfig, realRedis)

	// 测试数据
	userID := uint(101)
	expectedUser := &domain.User{
		ID:    userID,
		Name:  "Redis Master",
		Email: "redis@example.com",
	}

	t.Run("Scenario 1: Cache Miss (Redis无数据 -> 查DB -> 回填Redis)", func(t *testing.T) {
		// Step A: 确保 Redis 里是空的
		// (这是一个真 Redis，所以我们可以用 FlushDB)
		realRedis.FlushDB(context.Background())

		// Step B: 设置 Mock Repo 期望被调用 (因为缓存没有)
		mockRepo.EXPECT().
			FindByID(userID). // 假设你在接口里加了这个方法
			Return(expectedUser, nil).
			Times(1) // 预期只会调用一次 DB

		// C: ⚠️ 必须执行调用，否则 Mock 会报错 Missing Call
		user, err := svc.GetUserProfile(userID)

		// D: 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Name, user.Name)

		// E: 额外验证缓存是否真的回填了 (这是集成测试的精髓)
		ctx := context.Background()
		cachedVal, err := realRedis.Get(ctx, "user:profile:101").Result()
		assert.NoError(t, err, "Should have set cache key")
		assert.Contains(t, cachedVal, "Redis Master", "Cache content mismatch")

	})

	// --- ⚠️ 重要提示：为了让这段测试能跑通，你需要完成下面的 TODO ⚠️ ---

	/*
	 * 现在你的 UserRepo 接口里可能没有 FindByID。
	 * 你需要在测试这个之前：
	 * 1. 修改 internal/domain/user.go -> 接口增加 FindByID(id uint) (*User, error)
	 * 2. 重新运行 mockgen
	 * 3. 运行这个测试
	 */

	// 为了演示“缓存生效”的逻辑，我们手动往 Redis 塞数据，模拟 Scenario 2
	t.Run("Scenario 2: Cache Hit (Redis有数据 -> 直接返回 -> 不查DB)", func(t *testing.T) {
		ctx := context.Background()
		cacheKey := "user:profile:101"

		// 1. 直接往真实 Redis 预热数据
		data, _ := json.Marshal(expectedUser)
		err := realRedis.Set(ctx, cacheKey, data, time.Minute).Err()
		assert.NoError(t, err)

		// 2. 这里的 Mock 不设任何 Expect，意味着：如果不小心查了数据库，Mock 就会报错失败！
		// mockRepo.EXPECT().FindByID... (不需要写！)

		// 3. 调用 Service
		user, err := svc.GetUserProfile(userID)

		// 4. 验证
		assert.NoError(t, err)
		assert.Equal(t, "Redis Master", user.Name)
		assert.Equal(t, "redis@example.com", user.Email)
	})
}
