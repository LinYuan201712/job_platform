package model

// CompetitionExperience 竞赛经历模型
// 注意: date 字段在数据库中是 varchar(20),直接存储 YYYY-MM 字符串格式
type CompetitionExperience struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentUserID   int    `gorm:"column:student_user_id;not null" json:"student_user_id"`
	CompetitionName string `gorm:"column:competition_name" json:"competition_name"`
	Role            string `gorm:"column:role" json:"role"`
	Award           string `gorm:"column:award" json:"award"`
	Date            string `gorm:"column:date" json:"date"` // 存储 YYYY-MM 格式字符串
}

func (CompetitionExperience) TableName() string {
	return "competition_experiences"
}
