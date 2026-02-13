package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type StudentRepository struct {
	DB *gorm.DB
}

// FindByUserID 根据用户ID查询学生信息
func (r *StudentRepository) FindByUserID(userID int) (*model.Student, error) {
	var student model.Student
	err := r.DB.Where("user_id = ?", userID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// UpdateSkills 更新学生的技能掌握信息
func (r *StudentRepository) UpdateSkills(userID int, skillsSummary string) error {
	return r.DB.Model(&model.Student{}).
		Where("user_id = ?", userID).
		Update("skills_summary", skillsSummary).Error
}

// UpdateTemplate 更新学生当前使用的简历模板
func (r *StudentRepository) UpdateTemplate(userID int, templateID int64) error {
	return r.DB.Model(&model.Student{}).
		Where("user_id = ?", userID).
		Update("current_template_id", templateID).Error
}
