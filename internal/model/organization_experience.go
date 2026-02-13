package model

import "time"

type OrganizationExperience struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentUserID    int       `gorm:"column:student_user_id;not null" json:"student_user_id"`
	OrganizationName string    `gorm:"column:organization_name" json:"organization_name"`
	Role             string    `gorm:"column:role" json:"role"`
	StartDate        time.Time `gorm:"column:start_date;type:date" json:"start_date"`
	EndDate          time.Time `gorm:"column:end_date;type:date" json:"end_date"`
	Description      string    `gorm:"column:description;type:text" json:"description"`
}

func (OrganizationExperience) TableName() string {
	return "organization_experiences"
}
