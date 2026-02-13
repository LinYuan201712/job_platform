package repository

import (
	"errors"
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type ProjectExperienceRepository struct {
	DB *gorm.DB
}

// Create 创建项目经历
func (r *ProjectExperienceRepository) Create(exp *model.ProjectExperience) error {
	return r.DB.Create(exp).Error
}

// Update 更新项目经历
func (r *ProjectExperienceRepository) Update(exp *model.ProjectExperience) error {
	return r.DB.Save(exp).Error
}

// Delete 删除项目经历(需校验student_user_id防止越权)
func (r *ProjectExperienceRepository) Delete(id int64, studentUserID int) error {
	result := r.DB.Where("id = ? AND student_user_id = ?", id, studentUserID).
		Delete(&model.ProjectExperience{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("项目经历不存在或无权删除")
	}

	return nil
}

// FindByID 根据ID查询项目经历
func (r *ProjectExperienceRepository) FindByID(id int64) (*model.ProjectExperience, error) {
	var exp model.ProjectExperience
	err := r.DB.Where("id = ?", id).First(&exp).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

// FindByStudentUserID 查询学生的所有项目经历
func (r *ProjectExperienceRepository) FindByStudentUserID(studentUserID int) ([]model.ProjectExperience, error) {
	var experiences []model.ProjectExperience
	err := r.DB.Where("student_user_id = ?", studentUserID).
		Order("start_date DESC").
		Find(&experiences).Error
	return experiences, err
}
