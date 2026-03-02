package model

import "time"

// CompanyLink 对应数据库表: company_links
type CompanyLink struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CompanyID int       `gorm:"column:company_id;not null" json:"company_id"`
	LinkName  string    `gorm:"column:link_name;not null" json:"link_name"`
	LinkURL   string    `gorm:"column:link_url;not null" json:"link_url"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (CompanyLink) TableName() string {
	return "company_links"
}
