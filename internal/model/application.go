package model

import "time"

// Application 对应数据库表: application
type Application struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobID         int       `gorm:"column:job_id;not null" json:"job_id"`
	StudentUserID int       `gorm:"column:student_user_id;not null" json:"student_user_id"`
	ResumeID      int64     `gorm:"column:resume_id" json:"resume_id"`
	Status        int       `gorm:"column:status;not null" json:"status"` // 对应 constants.go 中的 ApplicationStatus
	SubmittedAt   time.Time `gorm:"column:submitted_at;autoCreateTime" json:"submitted_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Application) TableName() string {
	return "application"
}
