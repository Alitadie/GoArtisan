package service_test

import (
	"errors"
	"testing"

	"go-artisan/internal/config"
	"go-artisan/internal/domain"
	"go-artisan/internal/domain/mocks" // 引入刚才生成的 mock 包
	"go-artisan/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

// 辅助函数：快速生成 bcrypt 密码（因为 service 里会做校验）
func hashPassword(t *testing.T, raw string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

func TestUserService_Login(t *testing.T) {
	// 1. 初始化 Mock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 确保所有 EXPECT 都在函数结束前被调用过

	// 2. 创建 Mock 对象
	mockRepo := mocks.NewMockUserRepository(ctrl)

	// 3. 准备配置 (Login 需要读取配置里的 JWT 密钥)
	mockConfig := &config.Config{
		App: config.AppConfig{Name: "TestApp"},
	}

	// 4. 初始化被测 Service
	svc := service.NewUserService(mockRepo, mockConfig, nil)

	// 5. 准备测试数据
	validEmail := "test@example.com"
	validPass := "secret123"
	hashedPass := hashPassword(t, validPass)

	// --- 表格驱动测试 ---
	tests := []struct {
		name        string
		req         service.LoginDTO
		setupMock   func() // 这里定义每种场景下 Mock 应该怎么表现
		expectError bool
		expectToken bool
	}{
		{
			name: "Happy Path - 登录成功",
			req: service.LoginDTO{
				Email:    validEmail,
				Password: validPass,
			},
			setupMock: func() {
				// 期望：Repo.FindByEmail 会被调用一次，参数是 validEmail
				// 动作：返回一个正常的 User 对象，密码是哈希过的
				mockRepo.EXPECT().
					FindByEmail(validEmail).
					Return(&domain.User{
						ID:       1,
						Email:    validEmail,
						Password: hashedPass, // 注意这里必须是真哈希
					}, nil)
			},
			expectError: false,
			expectToken: true,
		},
		{
			name: "Fail - 用户不存在",
			req: service.LoginDTO{
				Email:    "missing@example.com",
				Password: "any",
			},
			setupMock: func() {
				// 模拟数据库找不到用户，返回错误
				mockRepo.EXPECT().
					FindByEmail("missing@example.com").
					Return(nil, errors.New("record not found"))
			},
			expectError: true, // 应该报错 "invalid credentials"
			expectToken: false,
		},
		{
			name: "Fail - 密码错误",
			req: service.LoginDTO{
				Email:    validEmail,
				Password: "wrongpassword",
			},
			setupMock: func() {
				// 用户找得到，但是密码校验会在 Service 层失败
				mockRepo.EXPECT().
					FindByEmail(validEmail).
					Return(&domain.User{
						ID:       1,
						Email:    validEmail,
						Password: hashedPass,
					}, nil)
			},
			expectError: true,
			expectToken: false,
		},
	}

	// 执行循环测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock() // 设置 Mock 行为

			resp, err := svc.Login(tt.req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.expectToken {
					assert.NotEmpty(t, resp.Token)
					assert.Equal(t, uint(1), resp.User.ID)
				}
			}
		})
	}
}
