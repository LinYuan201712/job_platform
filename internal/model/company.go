package model

import "time"

// Company 对应数据库表: company
type Company struct {
	CompanyID          int       `gorm:"column:company_id;primaryKey;autoIncrement" json:"company_id"`
	UserID             int       `gorm:"column:user_id;not null" json:"user_id"` // 关联 User.ID
	CompanyName        string    `gorm:"column:company_name;not null" json:"company_name"`
	Description        string    `gorm:"column:description;type:text" json:"description"`
	LogoUrl            string    `gorm:"column:logo_url" json:"logo_url"`
	IndustryID         int       `gorm:"column:industry_id" json:"industry_id"`
	NatureID           int       `gorm:"column:nature_id" json:"nature_id"`
	CompanyScaleID     int       `gorm:"column:company_scale_id" json:"company_scale_id"`
	CompanyAddress     string    `gorm:"column:company_address" json:"company_address"`
	ContactPersonName  string    `gorm:"column:contact_person_name" json:"contact_person_name"`
	ContactPersonPhone string    `gorm:"column:contact_person_phone" json:"contact_person_phone"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Company) TableName() string {
	return "companies"
}
