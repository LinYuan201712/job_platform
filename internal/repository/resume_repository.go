package repository

import (
	"errors"
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type ResumeRepository struct {
	DB *gorm.DB
}

// Create 创建简历文件记录
func (r *ResumeRepository) Create(resume *model.Resume) error {
	return r.DB.Create(resume).Error
}

// FindByID 根据ID查询简历文件
func (r *ResumeRepository) FindByID(id int64) (*model.Resume, error) {
	var resume model.Resume
	err := r.DB.Where("id = ?", id).First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

// FindByStudentUserID 查询学生的所有简历文件
func (r *ResumeRepository) FindByStudentUserID(studentUserID int) ([]model.Resume, error) {
	var resumes []model.Resume
	err := r.DB.Where("student_user_id = ?", studentUserID).
		Order("uploaded_at DESC").
		Find(&resumes).Error
	return resumes, err
}

// Delete 删除简历文件(需校验student_user_id防止越权)
func (r *ResumeRepository) Delete(id int64, studentUserID int) error {
	result := r.DB.Where("id = ? AND student_user_id = ?", id, studentUserID).
		Delete(&model.Resume{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("简历文件不存在或无权删除")
	}

	return nil
}

// CheckResumeUsedInApplications 检查简历是否被投递使用
// 如果存在引用返回 true,否则返回 false
func (r *ResumeRepository) CheckResumeUsedInApplications(resumeID int64) (bool, error) {
	var count int64
	err := r.DB.Table("applications").
		Where("resume_id = ?", resumeID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
