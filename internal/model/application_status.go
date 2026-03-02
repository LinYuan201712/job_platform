package model

import "time"

// ApplicationStatus 对应数据库表: application_statuses
type ApplicationStatus struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code      int       `gorm:"column:code;not null;uniqueIndex" json:"code"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Detail    string    `gorm:"column:detail;not null" json:"detail"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ApplicationStatus) TableName() string {
	return "application_statuses"
}
