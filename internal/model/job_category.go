package model

type JobCategory struct {
	ID   int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (JobCategory) TableName() string {
	return "job_category"
}
