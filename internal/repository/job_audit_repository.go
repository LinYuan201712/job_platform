package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type JobAuditRepository struct {
	DB *gorm.DB
}

// GetLatestAuditByJobID 获取岗位的最新审核记录
func (r *JobAuditRepository) GetLatestAuditByJobID(jobID int) (*model.JobAuditLog, error) {
	var audit model.JobAuditLog
	err := r.DB.Where("job_id = ?", jobID).
		Order("created_at DESC").
		First(&audit).Error
	if err != nil {
		return nil, err
	}
	return &audit, nil
}
