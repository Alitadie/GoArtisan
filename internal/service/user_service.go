package service

import (
	"errors"
	"fmt"

	"go-artisan/internal/domain"
	"go-artisan/internal/repository" // 注意不要让 Service 依赖具体 DB，而应依赖 Interface，这里为了 MVP 简化直接依赖 struct

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
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
