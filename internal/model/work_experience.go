package model

import "time"

type WorkExperience struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentUserID int       `gorm:"column:student_user_id;not null" json:"student_user_id"`
	CompanyName   string    `gorm:"column:company_name" json:"company_name"`
	PositionTitle string    `gorm:"column:position_title" json:"position_title"`
	StartDate     time.Time `gorm:"column:start_date;type:date" json:"start_date"`
	EndDate       time.Time `gorm:"column:end_date;type:date" json:"end_date"`
	Description   string    `gorm:"column:description;type:text" json:"description"`
}

func (WorkExperience) TableName() string {
	return "work_experience"
}
