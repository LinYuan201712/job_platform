package model

import "time"

// Student 对应数据库表: student
type Student struct {
	UserID            int        `gorm:"column:user_id;primaryKey" json:"user_id"` // 注意：主键不是自增，而是关联 User 表的 ID
	StudentID         string     `gorm:"column:student_id" json:"student_id"`      // 学号
	AvatarUrl         string     `gorm:"column:avatar_url" json:"avatar_url"`
	FullName          string     `gorm:"column:full_name" json:"full_name"`
	PhoneNumber       string     `gorm:"column:phone_number" json:"phone_number"`
	Gender            int        `gorm:"column:gender" json:"gender"`                         // 0=男, 1=女
	DateOfBirth       *time.Time `gorm:"column:date_of_birth;type:date" json:"date_of_birth"` // 使用指针支持 NULL
	JobSeekingStatus  int        `gorm:"column:job_seeking_status" json:"job_seeking_status"`
	ExpectedPosition  string     `gorm:"column:expected_position" json:"expected_position"`
	ExpectedMinSalary int        `gorm:"column:expected_min_salary" json:"expected_min_salary"`
	ExpectedMaxSalary int        `gorm:"column:expected_max_salary" json:"expected_max_salary"`
	SkillsSummary     string     `gorm:"column:skills_summary;type:text" json:"skills_summary"`
	CurrentTemplateID int64      `gorm:"column:current_template_id" json:"current_template_id"`
}

func (Student) TableName() string {
	return "students"
}
