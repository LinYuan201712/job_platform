package model

import (
	"time"
)

// User 对应数据库表: user
type User struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email        string     `gorm:"column:email;unique;not null" json:"email"`
	PasswordHash string     `gorm:"column:password_hash;not null" json:"-"` // json:"-" 表示这个字段永远不会返回给前端
	Role         int        `gorm:"column:role;not null" json:"role"`       // 对应 constants.go 中的常量
	Status       int        `gorm:"column:status;not null;default:1" json:"status"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at" json:"last_login_at"` // 使用指针 *time.Time 允许数据库存 NULL
}

// TableName 强制指定表名为 "user"，防止 GORM 自动复数化为 "users"
func (User) TableName() string {
	return "users"
}
