package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type TagRepository struct {
	DB *gorm.DB
}

// GetAllTagsGrouped 获取所有标签并按分类分组
func (r *TagRepository) GetAllTagsGrouped() (map[int][]model.Tag, []model.TagCategory, error) {
	// 获取所有分类
	var categories []model.TagCategory
	if err := r.DB.Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	// 获取所有标签
	var tags []model.Tag
	if err := r.DB.Find(&tags).Error; err != nil {
		return nil, nil, err
	}

	// 按 category_id 分组
	tagsByCategory := make(map[int][]model.Tag)
	for _, tag := range tags {
		tagsByCategory[tag.CategoryID] = append(tagsByCategory[tag.CategoryID], tag)
	}

	return tagsByCategory, categories, nil
}

// CreateTag 创建新标签
func (r *TagRepository) CreateTag(tag *model.Tag) error {
	return r.DB.Create(tag).Error
}

// CheckTagExists 检查标签名称是否已存在
func (r *TagRepository) CheckTagExists(name string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Tag{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// GetTagsByIDs 根据ID数组获取标签列表
func (r *TagRepository) GetTagsByIDs(tagIDs []int) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.DB.Where("id IN ?", tagIDs).Find(&tags).Error
	return tags, err
}

// FindTagsByStudentUserID 根据学生用户ID获取其所有标签
func (r *TagRepository) FindTagsByStudentUserID(userID int) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.DB.Table("tags").
		Joins("INNER JOIN students_tags ON tags.id = students_tags.tag_id").
		Where("students_tags.user_id = ?", userID).
		Find(&tags).Error
	return tags, err
}

// DeleteStudentTagsByUserID 删除学生的所有标签关联
func (r *TagRepository) DeleteStudentTagsByUserID(userID int) error {
	return r.DB.Exec("DELETE FROM students_tags WHERE user_id = ?", userID).Error
}

// BatchInsertStudentTags 批量插入学生标签关联
func (r *TagRepository) BatchInsertStudentTags(userID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// 构建批量插入的values
	type StudentTag struct {
		UserID int `gorm:"column:user_id"`
		TagID  int `gorm:"column:tag_id"`
	}

	var records []StudentTag
	for _, tagID := range tagIDs {
		records = append(records, StudentTag{
			UserID: userID,
			TagID:  tagID,
		})
	}

	return r.DB.Table("students_tags").Create(&records).Error
}
