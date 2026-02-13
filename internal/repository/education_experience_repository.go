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
