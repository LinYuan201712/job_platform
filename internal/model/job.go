package model

import "time"

// Job 对应数据库表: job
type Job struct {
	ID                int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CompanyID         int       `gorm:"column:company_id;not null" json:"company_id"`
	PostedByUserID    int       `gorm:"column:posted_by_user_id;not null" json:"posted_by_user_id"`
	Title             string    `gorm:"column:title;not null" json:"title"`
	Description       string    `gorm:"column:description;type:text" json:"description"`
	TechRequirements  string    `gorm:"column:tech_requirements;type:text" json:"tech_requirements"`
	MinSalary         int       `gorm:"column:min_salary" json:"min_salary"`
	MaxSalary         int       `gorm:"column:max_salary" json:"max_salary"`
	ProvinceID        int       `gorm:"column:province_id" json:"province_id"`
	CityID            int       `gorm:"column:city_id" json:"city_id"`
	AddressDetail     string    `gorm:"column:address_detail" json:"address_detail"`
	WorkNature        int       `gorm:"column:work_nature" json:"work_nature"`
	Deadline          time.Time `gorm:"column:deadline;type:date" json:"deadline"` // type:date 表示只存日期
	Status            int       `gorm:"column:status;default:0" json:"status"`
	Type              int       `gorm:"column:type" json:"type"`
	Department        string    `gorm:"column:department" json:"department"`
	Headcount         int       `gorm:"column:headcount" json:"headcount"`
	ViewCount         int       `gorm:"column:view_count;default:0" json:"view_count"`
	RequiredDegree    int       `gorm:"column:required_degree" json:"required_degree"`
	RequiredStartDate time.Time `gorm:"column:required_start_date;type:date" json:"required_start_date"`
	BonusPoints       string    `gorm:"column:bonus_points;type:text" json:"bonus_points"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Tags      []Tag     `gorm:"many2many:job_tags;" json:"tags"` // 关联标签表
}

func (Job) TableName() string {
	return "jobs"
}
