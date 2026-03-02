package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type JobFavoriteRepository struct {
	DB *gorm.DB
}

// AddFavorite 收藏职位 (幂等性由unique key保证)
func (r *JobFavoriteRepository) AddFavorite(studentUserID, jobID int) error {
	favorite := &model.JobFavorite{
		StudentUserID: studentUserID,
		JobID:         jobID,
	}
	// 使用FirstOrCreate实现幂等性,如果已存在则不插入
	result := r.DB.Where("student_user_id = ? AND job_id = ?", studentUserID, jobID).FirstOrCreate(favorite)
	return result.Error
}

// RemoveFavorite 取消收藏
func (r *JobFavoriteRepository) RemoveFavorite(studentUserID, jobID int) error {
	return r.DB.Where("student_user_id = ? AND job_id = ?", studentUserID, jobID).
		Delete(&model.JobFavorite{}).Error
}

// CheckFavoriteStatus 检查单个职位收藏状态
func (r *JobFavoriteRepository) CheckFavoriteStatus(studentUserID, jobID int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.JobFavorite{}).
		Where("student_user_id = ? AND job_id = ?", studentUserID, jobID).
		Count(&count).Error
	return count > 0, err
}

// GetFavoriteJobIDs 批量获取用户收藏的所有job_id (用于列表is_favorited判断)
func (r *JobFavoriteRepository) GetFavoriteJobIDs(studentUserID int) ([]int, error) {
	var jobIDs []int
	err := r.DB.Model(&model.JobFavorite{}).
		Where("student_user_id = ?", studentUserID).
		Pluck("job_id", &jobIDs).Error
	return jobIDs, err
}

// GetFavoriteList 获取收藏列表(分页)
func (r *JobFavoriteRepository) GetFavoriteList(studentUserID int, page, pageSize int) ([]model.JobFavorite, int64, error) {
	var favorites []model.JobFavorite
	var total int64

	// 计算总数
	err := r.DB.Model(&model.JobFavorite{}).
		Where("student_user_id = ?", studentUserID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.DB.Where("student_user_id = ?", studentUserID).
		Order("saved_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&favorites).Error

	return favorites, total, err
}
