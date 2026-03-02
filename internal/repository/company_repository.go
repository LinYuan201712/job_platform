package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type CompanyRepository struct {
	DB *gorm.DB
}

// GetCompanyByID 获取企业基本信息
func (r *CompanyRepository) GetCompanyByID(companyID int) (*model.Company, error) {
	var company model.Company
	err := r.DB.Where("company_id = ?", companyID).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// GetCompanyByUserID 根据UserID获取企业信息
func (r *CompanyRepository) GetCompanyByUserID(userID int) (*model.Company, error) {
	var company model.Company
	err := r.DB.Where("user_id = ?", userID).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// GetCompanyLinks 获取企业链接列表
func (r *CompanyRepository) GetCompanyLinks(companyID int) ([]model.CompanyLink, error) {
	var links []model.CompanyLink
	err := r.DB.Where("company_id = ?", companyID).Find(&links).Error
	return links, err
}

// GetCompanyJobsCount 获取企业在招岗位数(status=20 approved)
func (r *CompanyRepository) GetCompanyJobsCount(companyID int, status int) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Job{}).
		Where("company_id = ? AND status = ?", companyID, status).
		Count(&count).Error
	return count, err
}

// GetCompanyJobsList 获取企业所有发布的岗位
func (r *CompanyRepository) GetCompanyJobsList(companyID int) ([]model.Job, error) {
	var jobs []model.Job
	err := r.DB.Where("company_id = ? AND status = ?", companyID, 20). // 只返回已审核通过的岗位
										Order("created_at DESC").
										Find(&jobs).Error
	return jobs, err
}

// CompanyWithDictionaries 包含字典表名称的公司信息
type CompanyWithDictionaries struct {
	model.Company
	IndustryName     string `gorm:"column:industry_name"`
	NatureName       string `gorm:"column:nature_name"`
	CompanyScaleName string `gorm:"column:scale_name"`
}

// GetCompanyWithDictionaries 获取企业信息(含行业、性质、规模名称)
// 使用两步查询确保稳定性
func (r *CompanyRepository) GetCompanyWithDictionaries(companyID int) (*CompanyWithDictionaries, error) {
	var result CompanyWithDictionaries

	// 第一步：查询公司基本信息
	err := r.DB.Where("company_id = ?", companyID).First(&result.Company).Error
	if err != nil {
		return nil, err
	}

	// 第二步：查询字典名称（使用Scan,即使字典表不存在也不报错）
	type DictNames struct {
		IndustryName     string `gorm:"column:industry_name"`
		NatureName       string `gorm:"column:nature_name"`
		CompanyScaleName string `gorm:"column:scale_name"`
	}

	var names DictNames
	r.DB.Table("companies").
		Select(`COALESCE(t_industries.name, '') as industry_name,
			COALESCE(t_company_natures.name, '') as nature_name,
			COALESCE(t_company_scales.scale, '') as scale_name`).
		Joins("LEFT JOIN t_industries ON companies.industry_id = t_industries.id").
		Joins("LEFT JOIN t_company_natures ON companies.nature_id = t_company_natures.id").
		Joins("LEFT JOIN t_company_scales ON companies.company_scale_id = t_company_scales.id").
		Where("companies.company_id = ?", companyID).
		Scan(&names) // 使用Scan而不是First

	result.IndustryName = names.IndustryName
	result.NatureName = names.NatureName
	result.CompanyScaleName = names.CompanyScaleName

	return &result, nil
}

// UpdateCompanyProfile 更新企业基本信息
func (r *CompanyRepository) UpdateCompanyProfile(companyID int, updates map[string]interface{}) error {
	return r.DB.Model(&model.Company{}).
		Where("company_id = ?", companyID).
		Updates(updates).Error
}

// ReplaceCompanyLinks 全量替换企业链接(事务化,先删后增)
func (r *CompanyRepository) ReplaceCompanyLinks(tx *gorm.DB, companyID int, newLinks []model.CompanyLink) error {
	// 1. 删除旧链接
	if err := tx.Where("company_id = ?", companyID).Delete(&model.CompanyLink{}).Error; err != nil {
		return err
	}

	// 2. 插入新链接
	if len(newLinks) > 0 {
		for i := range newLinks {
			newLinks[i].CompanyID = companyID
		}
		if err := tx.Create(&newLinks).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetCompanyResumeProcessRate 计算简历处理率 (已审核数 / 总投递数)
func (r *CompanyRepository) GetCompanyResumeProcessRate(companyID int) (float64, error) {
	// 查询该企业所有岗位的投递记录
	type Stats struct {
		Total     int64 `gorm:"column:total"`
		Processed int64 `gorm:"column:processed"`
	}

	var stats Stats
	err := r.DB.Table("applications").
		Select(`COUNT(*) as total,
			SUM(CASE WHEN applications.status != 10 THEN 1 ELSE 0 END) as processed`).
		Joins("INNER JOIN jobs ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", companyID).
		Scan(&stats).Error

	if err != nil {
		return 0, err
	}

	if stats.Total == 0 {
		return 0, nil
	}

	return float64(stats.Processed) / float64(stats.Total), nil
}

// GetUserLastLoginTime 获取用户最近登录时间
func (r *CompanyRepository) GetUserLastLoginTime(userID int) (string, error) {
	type UserInfo struct {
		LastLoginAt *string `gorm:"column:last_login_at"`
	}

	var userInfo UserInfo
	err := r.DB.Table("users").
		Select("last_login_at").
		Where("id = ?", userID).
		Scan(&userInfo).Error

	if err != nil {
		return "", err
	}

	if userInfo.LastLoginAt == nil {
		return "", nil
	}

	return *userInfo.LastLoginAt, nil
}
