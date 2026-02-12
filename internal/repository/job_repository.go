package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type JobRepository struct {
	DB *gorm.DB
}

// GetJobByID 根据ID获取岗位(包含标签)
func (r *JobRepository) GetJobByID(id int) (*model.Job, error) {
	var job model.Job
	// Preload 会自动加载 many2many 关联的 Tags
	err := r.DB.Preload("Tags").Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// CreateJob 创建岗位
func (r *JobRepository) CreateJob(job *model.Job) error {
	// Create 会自动处理 many2many 关联的 Tags
	return r.DB.Create(job).Error
}

// UpdateJob 更新岗位
func (r *JobRepository) UpdateJob(job *model.Job) error {
	// 使用 Session 更新关联
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(job).Error
}

// DeleteJob 硬删除岗位
func (r *JobRepository) DeleteJob(id int) error {
	return r.DB.Delete(&model.Job{}, id).Error
}

// UpdateJobStatus 更新岗位状态
func (r *JobRepository) UpdateJobStatus(id int, status int) error {
	return r.DB.Model(&model.Job{}).Where("id = ?", id).Update("status", status).Error
}

// CheckJobOwnership 验证岗位是否属于指定用户
func (r *JobRepository) CheckJobOwnership(jobID, userID int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Job{}).Where("id = ? AND posted_by_user_id = ?", jobID, userID).Count(&count).Error
	return count > 0, err
}

// GetJobWithCompany 获取岗位及其公司信息
func (r *JobRepository) GetJobWithCompany(jobID int) (*model.Job, error) {
	var job model.Job
	err := r.DB.Preload("Tags").Where("id = ?", jobID).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}
