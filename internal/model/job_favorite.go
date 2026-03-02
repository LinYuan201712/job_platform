package model

import "time"

// JobFavorite 对应数据库表: job_favorites
type JobFavorite struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentUserID int       `gorm:"column:student_user_id;not null" json:"student_user_id"`
	JobID         int       `gorm:"column:job_id;not null" json:"job_id"`
	SavedAt       time.Time `gorm:"column:saved_at;autoCreateTime" json:"saved_at"`
}

func (JobFavorite) TableName() string {
	return "job_favorites"
}
