package model

import "time"

type EducationExperience struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentUserID int       `gorm:"column:student_user_id;not null" json:"student_user_id"`
	SchoolName    string    `gorm:"column:school_name" json:"school_name"`
	DegreeLevel   int       `gorm:"column:degree_level" json:"degree_level"` // 0=本科, 1=硕士, 2=博士
	Major         string    `gorm:"column:major" json:"major"`
	StartDate     time.Time `gorm:"column:start_date;type:date" json:"start_date"`
	EndDate       time.Time `gorm:"column:end_date;type:date" json:"end_date"`
	MajorRank     string    `gorm:"column:major_rank" json:"major_rank"`
}

func (EducationExperience) TableName() string {
	return "education_experience"
}
