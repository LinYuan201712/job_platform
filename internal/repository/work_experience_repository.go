package repository

import (
	"errors"
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type WorkExperienceRepository struct {
	DB *gorm.DB
}

// Create 创建工作经历
func (r *WorkExperienceRepository) Create(exp *model.WorkExperience) error {
	return r.DB.Create(exp).Error
}

// Update 更新工作经历
func (r *WorkExperienceRepository) Update(exp *model.WorkExperience) error {
	return r.DB.Save(exp).Error
}

// Delete 删除工作经历(需校验student_user_id防止越权)
func (r *WorkExperienceRepository) Delete(id int64, studentUserID int) error {
	result := r.DB.Where("id = ? AND student_user_id = ?", id, studentUserID).
		Delete(&model.WorkExperience{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("工作经历不存在或无权删除")
	}

	return nil
}

// FindByID 根据ID查询工作经历
func (r *WorkExperienceRepository) FindByID(id int64) (*model.WorkExperience, error) {
	var exp model.WorkExperience
	err := r.DB.Where("id = ?", id).First(&exp).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

// FindByStudentUserID 查询学生的所有工作经历
func (r *WorkExperienceRepository) FindByStudentUserID(studentUserID int) ([]model.WorkExperience, error) {
	var experiences []model.WorkExperience
	err := r.DB.Where("student_user_id = ?", studentUserID).
		Order("start_date DESC").
		Find(&experiences).Error
	return experiences, err
}
