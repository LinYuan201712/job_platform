package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type JobRepository struct {
	DB *gorm.DB
}

// GetJobByID 根据ID获取岗位(包含标签)
func (r *JobRepository) GetJobByID(id int) (*model.Job, error) {
	var job model.Job
	// Preload 会自动加载 many2many 关联的 Tags
	err := r.DB.Preload("Tags").Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// CreateJob 创建岗位
func (r *JobRepository) CreateJob(job *model.Job) error {
	// Create 会自动处理 many2many 关联的 Tags
	return r.DB.Create(job).Error
}

// UpdateJob 更新岗位
func (r *JobRepository) UpdateJob(job *model.Job) error {
	// 使用 Session 更新关联
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(job).Error
}

// DeleteJob 硬删除岗位
func (r *JobRepository) DeleteJob(id int) error {
	return r.DB.Delete(&model.Job{}, id).Error
}

// UpdateJobStatus 更新岗位状态
func (r *JobRepository) UpdateJobStatus(id int, status int) error {
	return r.DB.Model(&model.Job{}).Where("id = ?", id).Update("status", status).Error
}

// CheckJobOwnership 验证岗位是否属于指定用户
func (r *JobRepository) CheckJobOwnership(jobID, userID int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Job{}).Where("id = ? AND posted_by_user_id = ?", jobID, userID).Count(&count).Error
	return count > 0, err
}

// GetJobWithCompany 获取岗位及其公司信息
func (r *JobRepository) GetJobWithCompany(jobID int) (*model.Job, error) {
	var job model.Job
	err := r.DB.Preload("Tags").Where("id = ?", jobID).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// ======================== 求职中心相关查询 ========================

// JobListItem 职位列表查询结果项(带公司logo)
type JobListItem struct {
	model.Job
	CompanyName string `gorm:"column:company_name"`
	LogoURL     string `gorm:"column:logo_url"`
}

// GetJobList 多条件动态查询职位列表
func (r *JobRepository) GetJobList(filters map[string]interface{}, page, pageSize int) ([]JobListItem, int64, error) {
	var jobs []JobListItem
	var total int64

	// 构建基础查询(JOIN companies表获取logo_url)
	query := r.DB.Table("jobs").
		Select("jobs.*, companies.company_name, companies.logo_url").
		Joins("LEFT JOIN companies ON jobs.company_id = companies.company_id").
		Where("jobs.status = ?", 20) // 只查询已审核通过的职位

	// 动态添加筛选条件
	if province, ok := filters["province"].(string); ok && province != "" {
		// 需要JOIN t_provinces表获取province名称
		query = query.Joins("LEFT JOIN t_provinces ON jobs.province_id = t_provinces.id").
			Where("t_provinces.name = ?", province)
	}

	if city, ok := filters["city"].(string); ok && city != "" {
		// 需要JOIN t_cities表获取city名称
		query = query.Joins("LEFT JOIN t_cities ON jobs.city_id = t_cities.id").
			Where("t_cities.name = ?", city)
	}

	if title, ok := filters["title"].(string); ok && title != "" {
		query = query.Where("jobs.title LIKE ?", "%"+title+"%")
	}

	if companyName, ok := filters["company_name"].(string); ok && companyName != "" {
		query = query.Where("companies.company_name LIKE ?", "%"+companyName+"%")
	}

	if minSalary, ok := filters["min_salary"].(int); ok && minSalary > 0 {
		query = query.Where("jobs.max_salary >= ?", minSalary)
	}

	if maxSalary, ok := filters["max_salary"].(int); ok && maxSalary > 0 {
		query = query.Where("jobs.min_salary <= ?", maxSalary)
	}

	if workNature, ok := filters["work_nature"].(int); ok && workNature > 0 {
		query = query.Where("jobs.work_nature = ?", workNature)
	}

	if jobType, ok := filters["type"].(int); ok && jobType > 0 {
		query = query.Where("jobs.type = ?", jobType)
	}

	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("jobs.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&jobs).Error

	return jobs, total, err
}

// JobWithType 包含职位详情和职位类型名称
type JobWithType struct {
	model.Job
	TypeName string `gorm:"column:type_name"`
}

// GetJobDetailWithRelations 获取职位详情(含Tags和职位类型)
func (r *JobRepository) GetJobDetailWithRelations(jobID int) (*JobWithType, error) {
	var result JobWithType

	err := r.DB.Table("jobs").
		Select("jobs.*, COALESCE(t_job_categories.name, '') as type_name").
		Joins("LEFT JOIN t_job_categories ON jobs.type = t_job_categories.id").
		Where("jobs.id = ?", jobID).
		First(&result).Error

	if err != nil {
		return nil, err
	}

	// 手动加载Tags
	r.DB.Model(&result.Job).Association("Tags").Find(&result.Tags)

	return &result, nil
}

// GetRankedJobs 获取岗位热度排行榜（按view_count降序）
func (r *JobRepository) GetRankedJobs(limit int) ([]JobListItem, error) {
	var jobs []JobListItem
	err := r.DB.Table("jobs").
		Select("jobs.*, companies.company_name, companies.logo_url").
		Joins("LEFT JOIN companies ON jobs.company_id = companies.company_id").
		Where("jobs.status = ?", 20). // 只查询已审核通过的职位
		Order("jobs.view_count DESC").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}

// GetRecentJobs 获取近期招聘信息（按created_at降序，支持类型筛选）
func (r *JobRepository) GetRecentJobs(jobTypeFilter string, limit int) ([]model.Job, error) {
	var jobs []model.Job
	query := r.DB.Where("jobs.status = ?", 20) // 只查询已审核通过的职位

	// 根据岗位类型筛选 - 前端传来的 filter 是场景化的，需要映射到不同字段
	switch jobTypeFilter {
	case "实习招聘":
		// 对应 work_nature = 1
		query = query.Where("jobs.work_nature = ?", 1)
	case "校园招聘":
		// 对应 work_nature = 1(实习) 或 2(全职) 且 包含"校招"相关关键词
		query = query.Where("(jobs.work_nature = ? OR jobs.work_nature = ?) AND (jobs.title LIKE ? OR jobs.description LIKE ? OR jobs.title LIKE ?)", 1, 2, "%校招%", "%校招%", "%应届%")
	case "企业招聘":
		// 关联 companies 表，查询 nature_id 对应的 t_company_natures
		// 假设包含"企业"或"民营"等关键词，或者排除"事业单位"
		query = query.Joins("LEFT JOIN companies ON jobs.company_id = companies.company_id").
			Joins("LEFT JOIN t_company_natures ON companies.nature_id = t_company_natures.id").
			Where("t_company_natures.name LIKE ? OR t_company_natures.name LIKE ?", "%企业%", "%民营%")
	case "事业单位招聘":
		// 关联 companies 表，查询 nature_id 对应的 t_company_natures
		query = query.Joins("LEFT JOIN companies ON jobs.company_id = companies.company_id").
			Joins("LEFT JOIN t_company_natures ON companies.nature_id = t_company_natures.id").
			Where("t_company_natures.name LIKE ? OR t_company_natures.name LIKE ?", "%事业单位%", "%机关%")
	default:
		// 默认按 t_job_categories 名称匹配
		if jobTypeFilter != "" {
			query = query.Joins("LEFT JOIN t_job_categories ON jobs.type = t_job_categories.id").
				Where("t_job_categories.name = ?", jobTypeFilter)
		}
	}

	err := query.Order("jobs.created_at DESC").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}
