package model

import "time"

// TagCategory 对应数据库表: tag_categories
// 用于标签分组,如"编程语言"、"技术框架"等
type TagCategory struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"column:code;not null;uniqueIndex" json:"code"`       // 分类短码, 如 'lang', 'skill'
	Name        string    `gorm:"column:name;not null" json:"name"`                   // 分类名称, 如 '编程语言'
	Description string    `gorm:"column:description" json:"description"`              // 分类描述
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"` // 更新时间
}

func (TagCategory) TableName() string {
	return "tag_categories"
}
