package repository

import (
	"errors"
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type CompetitionExperienceRepository struct {
	DB *gorm.DB
}

// Create 创建竞赛经历
func (r *CompetitionExperienceRepository) Create(exp *model.CompetitionExperience) error {
	return r.DB.Create(exp).Error
}

// Update 更新竞赛经历
func (r *CompetitionExperienceRepository) Update(exp *model.CompetitionExperience) error {
	return r.DB.Save(exp).Error
}

// Delete 删除竞赛经历(需校验student_user_id防止越权)
func (r *CompetitionExperienceRepository) Delete(id int64, studentUserID int) error {
	result := r.DB.Where("id = ? AND student_user_id = ?", id, studentUserID).
		Delete(&model.CompetitionExperience{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("竞赛经历不存在或无权删除")
	}

	return nil
}

// FindByID 根据ID查询竞赛经历
func (r *CompetitionExperienceRepository) FindByID(id int64) (*model.CompetitionExperience, error) {
	var exp model.CompetitionExperience
	err := r.DB.Where("id = ?", id).First(&exp).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

// FindByStudentUserID 查询学生的所有竞赛经历
func (r *CompetitionExperienceRepository) FindByStudentUserID(studentUserID int) ([]model.CompetitionExperience, error) {
	var experiences []model.CompetitionExperience
	err := r.DB.Where("student_user_id = ?", studentUserID).
		Order("date DESC").
		Find(&experiences).Error
	return experiences, err
}
