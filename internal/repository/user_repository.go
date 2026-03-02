package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户仓库
type UserRepository struct {
	DB *gorm.DB
}

// FindByID 根据ID查询用户
func (r *UserRepository) FindByID(id int) (*model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

// UpdatePassword 更新用户密码hash
func (r *UserRepository) UpdatePassword(userID int, newPasswordHash string) error {
	return r.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("password_hash", newPasswordHash).Error
}

// VerifyPassword 验证密码（通过查询password_hash并在service层比对）
func (r *UserRepository) VerifyPassword(userID int, passwordHash string) (bool, error) {
	var user model.User
	err := r.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return false, err
	}
	// 返回hash值让service层进行bcrypt比对
	return user.PasswordHash == passwordHash, nil
}

// GetPasswordHash 获取用户的密码hash（用于service层进行bcrypt验证）
func (r *UserRepository) GetPasswordHash(userID int) (string, error) {
	var user model.User
	err := r.DB.Select("password_hash").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.PasswordHash, nil
}
