package service

import (
	"errors"
	"fmt"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"

	"gorm.io/gorm"
)

type ApplicationService struct {
	ApplicationRepo       *repository.ApplicationRepository
	ApplicationStatusRepo *repository.ApplicationStatusRepository
	JobRepo               *repository.JobRepository
	ResumeRepo            *repository.ResumeRepository
	CompanyRepo           *repository.CompanyRepository
	FavoriteRepo          *repository.JobFavoriteRepository
	DB                    *gorm.DB
}

// ApplyJob 投递职位
func (s *ApplicationService) ApplyJob(req dto.ApplyJobRequest, studentUserID int) (*dto.ApplicationResponse, error) {
	// 1. 验证job_id是否存在且状态为approved(20)
	job, err := s.JobRepo.GetJobByID(req.JobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("职位不存在")
		}
		return nil, err
	}

	if job.Status != 20 { // 20 = approved
		return nil, errors.New("该职位暂不可投递")
	}

	// 2. 验证resume_id是否属于当前学生
	resume, err := s.ResumeRepo.GetResumeByID(req.ResumeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("简历不存在")
		}
		return nil, err
	}

	if resume.StudentUserID != studentUserID {
		return nil, errors.New("简历不存在或无权使用")
	}

	// 3. 检查是否已投递(防重复)
	exists, err := s.ApplicationRepo.CheckApplicationExists(studentUserID, req.JobID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("您已投递过该职位")
	}

	// 4. 使用事务创建投递记录
	var application *model.Application
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		application = &model.Application{
			JobID:         req.JobID,
			StudentUserID: studentUserID,
			ResumeID:      req.ResumeID,
			Status:        model.AppStatusSubmitted, // 初始状态: 10(已投递)
		}

		// 创建投递记录
		appRepo := &repository.ApplicationRepository{DB: tx}
		if err := appRepo.CreateApplication(application); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 5. 返回投递结果
	return &dto.ApplicationResponse{
		ApplicationID: application.ID,
		JobID:         application.JobID,
		StudentUserID: application.StudentUserID,
		ResumeID:      application.ResumeID,
		Status:        application.Status,
		SubmittedAt:   application.SubmittedAt,
		UpdatedAt:     application.UpdatedAt,
	}, nil
}

// GetApplicationDetail 获取投递详情（含状态文案映射）
func (s *ApplicationService) GetApplicationDetail(appID int, studentUserID int) (*dto.ApplicationDetailResponse, error) {
	// 1. 查询投递记录（验证用户权限）
	app, err := s.ApplicationRepo.GetApplicationByID(appID, studentUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到该投递记录")
		}
		return nil, err
	}

	// 2. 查询状态文案映射
	statusName, statusDetail, err := s.getAppStatusDetail(app.Status)
	if err != nil {
		return nil, err
	}

	// 3. 关联查询岗位信息
	job, err := s.JobRepo.GetJobByID(app.JobID)
	if err != nil {
		return nil, errors.New("岗位信息不存在")
	}

	// 4. 关联查询公司信息
	company, err := s.CompanyRepo.GetCompanyByID(job.CompanyID)
	if err != nil {
		return nil, errors.New("公司信息不存在")
	}

	// 5. 构建响应
	return &dto.ApplicationDetailResponse{
		ID:           app.ID,
		Status:       statusName,
		StatusDetail: statusDetail,
		SubmittedAt:  app.SubmittedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    app.UpdatedAt.Format("2006-01-02 15:04:05"),
		Job: dto.JobInfo{
			ID:    job.ID,
			Title: job.Title,
		},
		Company: dto.CompanyInfo{
			Name: company.CompanyName,
		},
	}, nil
}

// GetDeliveryList 获取已投递岗位列表
func (s *ApplicationService) GetDeliveryList(req dto.DeliveryListRequest, studentUserID int) (*dto.DeliveryListResponse, error) {
	// 默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 1. 查询投递记录列表
	applications, total, err := s.ApplicationRepo.GetApplicationListByStudent(
		studentUserID, req.JobTitle, req.CompanyName, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 2. 获取岗位ID列表
	jobIDs := make([]int, 0, len(applications))
	appIDMap := make(map[int]int) // job_id -> application_id
	for _, app := range applications {
		jobIDs = append(jobIDs, app.JobID)
		appIDMap[app.JobID] = app.ID
	}

	// 3. 批量查询岗位信息
	var jobs []repository.JobListItem
	err = s.DB.Table("jobs").
		Select("jobs.*, companies.company_name, companies.logo_url").
		Joins("LEFT JOIN companies ON jobs.company_id = companies.company_id").
		Where("jobs.id IN ?", jobIDs).
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	// 4. 查询收藏状态
	favoriteMap := make(map[int]bool)
	var favorites []model.JobFavorite
	err = s.DB.Where("student_user_id = ? AND job_id IN ?", studentUserID, jobIDs).
		Find(&favorites).Error
	if err == nil {
		for _, fav := range favorites {
			favoriteMap[fav.JobID] = true
		}
	}

	// 5. 构建响应
	var jobItems []dto.DeliveryJobItem
	for _, job := range jobs {
		salaryRange := ""
		if job.Job.MinSalary > 0 && job.Job.MaxSalary > 0 {
			salaryRange = fmt.Sprintf("%d-%dK", job.Job.MinSalary, job.Job.MaxSalary)
		}

		address := job.Job.AddressDetail

		workNature := ""
		if job.Job.WorkNature == 1 {
			workNature = "实习"
		} else if job.Job.WorkNature == 2 {
			workNature = "全职"
		}

		department := job.Job.Department
		headcount := job.Job.Headcount

		appID := appIDMap[job.Job.ID]

		jobItems = append(jobItems, dto.DeliveryJobItem{
			JobID:         job.Job.ID,
			Title:         job.Job.Title,
			CompanyName:   job.CompanyName,
			SalaryRange:   salaryRange,
			Address:       address,
			WorkNature:    workNature,
			Department:    department,
			Headcount:     headcount,
			IsFavorited:   favoriteMap[job.Job.ID],
			LogoURL:       job.LogoURL,
			ApplicationID: &appID,
		})
	}

	return &dto.DeliveryListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Jobs:     jobItems,
	}, nil
}

// getAppStatusDetail 私有方法：根据状态码查询状态文案映射
func (s *ApplicationService) getAppStatusDetail(code int) (string, string, error) {
	status, err := s.ApplicationStatusRepo.GetStatusByCode(code)
	if err != nil {
		// 如果数据库中没有对应记录，返回默认值
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "未知状态", "暂无详情", nil
		}
		return "", "", err
	}
	return status.Name, status.Detail, nil
}
