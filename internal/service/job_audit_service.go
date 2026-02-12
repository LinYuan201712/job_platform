package service

import (
	"errors"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/repository"

	"gorm.io/gorm"
)

type JobAuditService struct {
	AuditRepo *repository.JobAuditRepository
}

// GetJobAudit 获取岗位最新审核详情
func (s *JobAuditService) GetJobAudit(jobID int) (*dto.JobAuditResponse, error) {
	// 1. 获取最新审核记录
	audit, err := s.AuditRepo.GetLatestAuditByJobID(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("岗位不存在或无权访问")
		}
		return nil, err
	}

	// 2. 构建响应
	return &dto.JobAuditResponse{
		JobID:           audit.JobID,
		AuditStatus:     audit.AuditStatus,
		Remark:          audit.Remark,
		OperatorContact: audit.OperatorContact,
		AuditTime:       audit.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
