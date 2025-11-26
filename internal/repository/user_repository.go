package repository

import (
	"go-artisan/internal/domain"

	"gorm.io/gorm"
)

// UserRepo 实现
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 构造函数，自动注入 gorm.DB
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// 确保实现了接口
var _ domain.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
