package repository

import (
	"errors"
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type OrganizationExperienceRepository struct {
	DB *gorm.DB
}

// Create 创建组织经历
func (r *OrganizationExperienceRepository) Create(exp *model.OrganizationExperience) error {
	return r.DB.Create(exp).Error
}

// Update 更新组织经历
func (r *OrganizationExperienceRepository) Update(exp *model.OrganizationExperience) error {
	return r.DB.Save(exp).Error
}

// Delete 删除组织经历(需校验student_user_id防止越权)
func (r *OrganizationExperienceRepository) Delete(id int64, studentUserID int) error {
	result := r.DB.Where("id = ? AND student_user_id = ?", id, studentUserID).
		Delete(&model.OrganizationExperience{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("组织经历不存在或无权删除")
	}

	return nil
}

// FindByID 根据ID查询组织经历
func (r *OrganizationExperienceRepository) FindByID(id int64) (*model.OrganizationExperience, error) {
	var exp model.OrganizationExperience
	err := r.DB.Where("id = ?", id).First(&exp).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

// FindByStudentUserID 查询学生的所有组织经历
func (r *OrganizationExperienceRepository) FindByStudentUserID(studentUserID int) ([]model.OrganizationExperience, error) {
	var experiences []model.OrganizationExperience
	err := r.DB.Where("student_user_id = ?", studentUserID).
		Order("start_date DESC").
		Find(&experiences).Error
	return experiences, err
}
