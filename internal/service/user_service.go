package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-artisan/internal/config"
	"go-artisan/internal/domain"

	"go-artisan/pkg/auth"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   domain.UserRepository
	config *config.Config
	redis  *redis.Client // 👈 新增依赖
}

// LoginDTO 输入对象
type LoginDTO struct {
	Email    string
	Password string
}

// LoginResponse
type LoginResponse struct {
	Token     string       `json:"token"`
	User      *domain.User `json:"user"`
	ExpiresIn int          `json:"expires_in"`
}

func NewUserService(repo domain.UserRepository, cfg *config.Config) *UserService {
	return &UserService{repo: repo, config: cfg}
}

// RegisterDTO 输入对象
type RegisterDTO struct {
	Name     string
	Email    string
	Password string
}

func (s *UserService) Register(req RegisterDTO) (*domain.User, error) {
	// 1. 检查邮箱
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil && existing.ID > 0 {
		return nil, errors.New("email already taken")
	}

	// 2. 密码加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. 构造模型
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPwd),
	}

	// 4. 落库
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req LoginDTO) (*LoginResponse, error) {
	// 1. 查用户
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. 比对密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials") // 模糊报错为了安全
	}

	// 3. 签发 Token
	// 注意：实际项目中建议在 config.go 中解析成 AuthConfig 结构体
	// secret := s.config.App.Name + "Secret" // MVP简化，建议用 s.config.Auth.Secret
	// if val := s.config.Database.DSN; val != "" {
	// 	// 这里为了演示演示，实际上你应该在 Viper 加载好 Auth 配置
	// 	// secret = "TempSecretInDev"
	// }
	// 使用 .env 中加载的配置 (这里为了不大量改动 config.go，写死示例，请自行优化)
	jwtSecret := "KeepItSecretKeepItSafe!GoArtisanKey"

	token, err := auth.GenerateToken(user.ID, jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token:     token,
		User:      user,
		ExpiresIn: 86400,
	}, nil
}

// GetUserProfile 示例：带缓存查询
func (s *UserService) GetUserProfile(id uint) (*domain.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:profile:%d", id)

	// 1. 查缓存
	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// 命中！反序列化
		var user domain.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil // 直接返回
		}
	}

	// 2. 未命中或出错 -> 查数据库
	// 注意：Repo 接口可能需要添加 FindByID 方法，这里假设已有
	// user, err := s.repo.FindByID(id)
	// 为了演示代码编译通过，先 mock 一个空操作
	var user *domain.User = nil // 实际应调用 s.repo
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 3. 回填缓存 (设置 10 分钟过期，防止缓存雪崩可加随机值)
	data, _ := json.Marshal(user)
	s.redis.Set(ctx, cacheKey, data, 10*time.Minute)

	return user, nil
}
