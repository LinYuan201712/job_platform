package model

import "time"

type Tag struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;not null" json:"name"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CategoryID int       `gorm:"column:category_id" json:"category_id"`
	CreatedBy  int       `gorm:"column:created_by" json:"created_by"`
}

func (Tag) TableName() string {
	return "tag"
}
