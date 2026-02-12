package model

import "time"

type ProjectExperience struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentUserID int       `gorm:"column:student_user_id;not null" json:"student_user_id"`
	ProjectName   string    `gorm:"column:project_name" json:"project_name"`
	Role          string    `gorm:"column:role" json:"role"`
	ProjectLink   string    `gorm:"column:project_link" json:"project_link"`
	StartDate     time.Time `gorm:"column:start_date;type:date" json:"start_date"`
	EndDate       time.Time `gorm:"column:end_date;type:date" json:"end_date"`
	Description   string    `gorm:"column:description;type:text" json:"description"`
}

func (ProjectExperience) TableName() string {
	return "project_experience"
}
