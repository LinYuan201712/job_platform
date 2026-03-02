package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type EducationExperienceRepository struct {
	DB *gorm.DB
}

// FindByStudentUserIDOrderByDegree 查询学生的教育经历,按学历等级升序排列
// 0=本科, 1=硕士, 2=博士
func (r *EducationExperienceRepository) FindByStudentUserIDOrderByDegree(studentUserID int) ([]model.EducationExperience, error) {
	var experiences []model.EducationExperience
	err := r.DB.Where("student_user_id = ?", studentUserID).
		Order("degree_level ASC").
		Find(&experiences).Error
	return experiences, err
}

// FindByStudentUserID 查询学生的所有教育经历（不排序）
func (r *EducationExperienceRepository) FindByStudentUserID(studentUserID int) ([]model.EducationExperience, error) {
	var experiences []model.EducationExperience
	err := r.DB.Where("student_user_id = ?", studentUserID).Find(&experiences).Error
	return experiences, err
}

// Create 创建教育经历
func (r *EducationExperienceRepository) Create(exp *model.EducationExperience) error {
	return r.DB.Create(exp).Error
}

// Update 更新教育经历
func (r *EducationExperienceRepository) Update(exp *model.EducationExperience) error {
	return r.DB.Save(exp).Error
}

// DeleteByID 根据ID删除教育经历
func (r *EducationExperienceRepository) DeleteByID(id int64) error {
	return r.DB.Delete(&model.EducationExperience{}, id).Error
}

// DeleteByIDAndUserID 根据ID和用户ID删除教育经历（含所有权检查）
func (r *EducationExperienceRepository) DeleteByIDAndUserID(id int64, userID int) error {
	return r.DB.Where("id = ? AND student_user_id = ?", id, userID).
		Delete(&model.EducationExperience{}).Error
}
