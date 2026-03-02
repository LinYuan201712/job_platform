package repository

import (
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type ApplicationRepository struct {
	DB *gorm.DB
}

// CreateApplication 创建投递记录
func (r *ApplicationRepository) CreateApplication(app *model.Application) error {
	return r.DB.Create(app).Error
}

// CheckApplicationExists 检查是否已投递
func (r *ApplicationRepository) CheckApplicationExists(studentUserID, jobID int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Application{}).
		Where("student_user_id = ? AND job_id = ?", studentUserID, jobID).
		Count(&count).Error
	return count > 0, err
}

// GetApplicationByJobAndStudent 获取投递记录
func (r *ApplicationRepository) GetApplicationByJobAndStudent(jobID, studentUserID int) (*model.Application, error) {
	var app model.Application
	err := r.DB.Where("job_id = ? AND student_user_id = ?", jobID, studentUserID).
		First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetApplicationIDsByStudent 获取用户所有投递记录的映射 job_id->application_id
func (r *ApplicationRepository) GetApplicationIDsByStudent(studentUserID int) (map[int]int, error) {
	var applications []model.Application
	err := r.DB.Where("student_user_id = ?", studentUserID).
		Find(&applications).Error
	if err != nil {
		return nil, err
	}

	// 构建 job_id -> application_id 的映射
	mapping := make(map[int]int)
	for _, app := range applications {
		mapping[app.JobID] = app.ID
	}
	return mapping, nil
}

// GetApplicationByID 根据ID获取投递记录（验证用户权限）
func (r *ApplicationRepository) GetApplicationByID(appID int, studentUserID int) (*model.Application, error) {
	var app model.Application
	err := r.DB.Where("id = ? AND student_user_id = ?", appID, studentUserID).
		First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetApplicationListByStudent 查询学生的投递列表（支持筛选）
func (r *ApplicationRepository) GetApplicationListByStudent(studentUserID int, jobTitle, companyName string, page, pageSize int) ([]model.Application, int64, error) {
	var applications []model.Application
	var total int64

	query := r.DB.Model(&model.Application{}).
		Joins("LEFT JOIN jobs ON applications.job_id = jobs.id").
		Joins("LEFT JOIN companies ON jobs.company_id = companies.company_id").
		Where("applications.student_user_id = ?", studentUserID)

	// 筛选条件
	if jobTitle != "" {
		query = query.Where("jobs.title LIKE ?", "%"+jobTitle+"%")
	}
	if companyName != "" {
		query = query.Where("companies.company_name LIKE ?", "%"+companyName+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Select("applications.*").
		Order("applications.submitted_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&applications).Error

	return applications, total, err
}

// GetApplicationCountsByJobIDs 批量查询岗位投递数量(性能优化,避免N+1)
func (r *ApplicationRepository) GetApplicationCountsByJobIDs(jobIDs []int) (map[int]int64, error) {
	if len(jobIDs) == 0 {
		return make(map[int]int64), nil
	}

	type CountResult struct {
		JobID int   `gorm:"column:job_id"`
		Count int64 `gorm:"column:count"`
	}

	var results []CountResult
	err := r.DB.Model(&model.Application{}).
		Select("job_id, COUNT(*) as count").
		Where("job_id IN ?", jobIDs).
		Group("job_id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[int]int64)
	for _, result := range results {
		countMap[result.JobID] = result.Count
	}

	return countMap, nil
}

// GetUnreviewedCountsByJobIDs 批量查询未审核数量(status=10 已投递)
func (r *ApplicationRepository) GetUnreviewedCountsByJobIDs(jobIDs []int) (map[int]int64, error) {
	if len(jobIDs) == 0 {
		return make(map[int]int64), nil
	}

	type CountResult struct {
		JobID int   `gorm:"column:job_id"`
		Count int64 `gorm:"column:count"`
	}

	var results []CountResult
	err := r.DB.Model(&model.Application{}).
		Select("job_id, COUNT(*) as count").
		Where("job_id IN ? AND status = ?", jobIDs, model.AppStatusSubmitted). // 10 = 已投递
		Group("job_id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[int]int64)
	for _, result := range results {
		countMap[result.JobID] = result.Count
	}

	return countMap, nil
}

// CandidateQueryResult 候选人查询结果
type CandidateQueryResult struct {
	ApplicationID int    `gorm:"column:application_id"`
	UserID        int    `gorm:"column:user_id"`
	FullName      string `gorm:"column:full_name"`
	AvatarURL     string `gorm:"column:avatar_url"`
	Degree        int    `gorm:"column:degree"`
	Grade         string `gorm:"column:grade"`
	Status        int    `gorm:"column:status"`
}

// GetCandidatesByJobID 获取岗位候选人列表(HR端,带筛选和分页)
func (r *ApplicationRepository) GetCandidatesByJobID(jobID int, nameKeyword string, status *int, page, pageSize int) ([]CandidateQueryResult, int64, error) {
	var candidates []CandidateQueryResult
	var total int64

	query := r.DB.Table("applications").
		Select(`applications.id as application_id,
			applications.student_user_id as user_id,
			students.full_name,
			students.avatar_url,
			COALESCE((SELECT degree_level FROM education_experiences WHERE student_user_id = applications.student_user_id ORDER BY degree_level DESC LIMIT 1), 0) as degree,
			COALESCE(CONCAT(YEAR(students.date_of_birth), '级'), '') as grade,
			applications.status`).
		Joins("LEFT JOIN students ON applications.student_user_id = students.user_id").
		Where("applications.job_id = ?", jobID)

	// 筛选条件
	if nameKeyword != "" {
		query = query.Where("students.full_name LIKE ?", "%"+nameKeyword+"%")
	}
	if status != nil {
		query = query.Where("applications.status = ?", *status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("applications.submitted_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&candidates).Error

	return candidates, total, err
}

// GetApplicationDetailByID 获取投递详情(HR端,不验证studentUserID)
func (r *ApplicationRepository) GetApplicationDetailByID(appID int) (*model.Application, error) {
	var app model.Application
	err := r.DB.Where("id = ?", appID).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// UpdateApplicationStatus 更新投递状态
func (r *ApplicationRepository) UpdateApplicationStatus(appID int, status int) error {
	return r.DB.Model(&model.Application{}).
		Where("id = ?", appID).
		Update("status", status).Error
}
