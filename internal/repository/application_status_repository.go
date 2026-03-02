package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type ApplicationStatusRepository struct {
	DB *gorm.DB
}

// GetStatusByCode 根据状态码查询状态详情
func (r *ApplicationStatusRepository) GetStatusByCode(code int) (*model.ApplicationStatus, error) {
	var status model.ApplicationStatus
	err := r.DB.Where("code = ?", code).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}
