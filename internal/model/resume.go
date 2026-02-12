package model

import "time"

// Resume 对应数据库表: resume
type Resume struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"` // 注意：Java 中是 Long，Go 中用 int64
	StudentUserID int       `gorm:"column:student_user_id;not null" json:"student_user_id"`
	FileName      string    `gorm:"column:file_name" json:"file_name"`
	FileUrl       string    `gorm:"column:file_url" json:"file_url"`
	FileSize      int64     `gorm:"column:file_size" json:"file_size"`
	UsageType     string    `gorm:"column:usage_type" json:"usage_type"`
	TemplateID    int       `gorm:"column:template_id" json:"template_id"`
	UploadedAt    time.Time `gorm:"column:uploaded_at;autoCreateTime" json:"uploaded_at"`
}

func (Resume) TableName() string {
	return "resume"
}
