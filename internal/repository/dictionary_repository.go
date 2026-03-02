package repository

import (
	"gorm.io/gorm"
)

type DictionaryRepository struct {
	DB *gorm.DB
}

// IndustryResult 行业查询结果
type IndustryResult struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

// NatureResult 企业性质查询结果
type NatureResult struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

// ScaleResult 企业规模查询结果
type ScaleResult struct {
	ID    int    `gorm:"column:id"`
	Scale string `gorm:"column:scale"`
}

// GetAllIndustries 获取所有行业
func (r *DictionaryRepository) GetAllIndustries() ([]string, error) {
	var results []IndustryResult
	err := r.DB.Table("t_industries").
		Select("id, name").
		Order("id ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(results))
	for _, result := range results {
		names = append(names, result.Name)
	}

	return names, nil
}

// GetAllNatures 获取所有企业性质
func (r *DictionaryRepository) GetAllNatures() ([]string, error) {
	var results []NatureResult
	err := r.DB.Table("t_company_natures").
		Select("id, name").
		Order("id ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(results))
	for _, result := range results {
		names = append(names, result.Name)
	}

	return names, nil
}

// GetAllScales 获取所有企业规模
func (r *DictionaryRepository) GetAllScales() ([]string, error) {
	var results []ScaleResult
	err := r.DB.Table("t_company_scales").
		Select("id, scale").
		Order("id ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	scales := make([]string, 0, len(results))
	for _, result := range results {
		scales = append(scales, result.Scale)
	}

	return scales, nil
}

// GetIndustryIDByName 根据名称获取行业ID
func (r *DictionaryRepository) GetIndustryIDByName(name string) (int, error) {
	var result IndustryResult
	err := r.DB.Table("t_industries").
		Select("id").
		Where("name = ?", name).
		First(&result).Error

	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

// GetNatureIDByName 根据名称获取性质ID
func (r *DictionaryRepository) GetNatureIDByName(name string) (int, error) {
	var result NatureResult
	err := r.DB.Table("t_company_natures").
		Select("id").
		Where("name = ?", name).
		First(&result).Error

	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

// GetScaleIDByName 根据名称获取规模ID
func (r *DictionaryRepository) GetScaleIDByName(scale string) (int, error) {
	var result ScaleResult
	err := r.DB.Table("t_company_scales").
		Select("id").
		Where("scale = ?", scale).
		First(&result).Error

	if err != nil {
		return 0, err
	}

	return result.ID, nil
}
