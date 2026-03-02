package service

import (
	"encoding/json"
	"errors"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/repository"
	"job-platform-go2/pkg/utils"
	"strconv"

	"gorm.io/gorm"
)

type JobCenterService struct {
	JobRepo         *repository.JobRepository
	FavoriteRepo    *repository.JobFavoriteRepository
	ApplicationRepo *repository.ApplicationRepository
	CompanyRepo     *repository.CompanyRepository
	DB              *gorm.DB
}

// GetJobList 获取职位列表(含收藏状态和投递状态)
func (s *JobCenterService) GetJobList(req dto.JobListRequest, currentUserID int) (*dto.JobListResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建筛选条件
	filters := make(map[string]interface{})

	if req.Province != "" {
		filters["province"] = req.Province
	}
	if req.City != "" {
		filters["city"] = req.City
	}
	if req.Title != "" {
		filters["title"] = req.Title
	}
	if req.CompanyName != "" {
		filters["company_name"] = req.CompanyName
	}

	// 薪资筛选(字符串转int)
	if req.MinSalary != "" {
		if minSal, err := strconv.Atoi(req.MinSalary); err == nil {
			filters["min_salary"] = minSal
		}
	}
	if req.MaxSalary != "" {
		if maxSal, err := strconv.Atoi(req.MaxSalary); err == nil {
			filters["max_salary"] = maxSal
		}
	}

	// 工作性质转换 ("实习招聘" -> 1, "校园招聘" -> 2)
	if req.WorkNature != "" {
		workNatureCode := utils.ParseWorkNatureToCode(req.WorkNature)
		if workNatureCode > 0 {
			filters["work_nature"] = workNatureCode
		}
	}

	// 职能类别(需要先查询字典表获取ID,这里简化直接接受ID)
	if req.Type != "" {
		if typeID, err := strconv.Atoi(req.Type); err == nil {
			filters["type"] = typeID
		}
	}

	// 查询职位列表
	jobs, total, err := s.JobRepo.GetJobList(filters, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 批量获取当前用户的收藏job_id列表
	favoriteJobIDs, err := s.FavoriteRepo.GetFavoriteJobIDs(currentUserID)
	if err != nil {
		return nil, err
	}
	favoriteMap := make(map[int]bool)
	for _, jobID := range favoriteJobIDs {
		favoriteMap[jobID] = true
	}

	// 批量获取当前用户的投递记录映射
	applicationMap, err := s.ApplicationRepo.GetApplicationIDsByStudent(currentUserID)
	if err != nil {
		return nil, err
	}

	// 组装响应
	jobItems := make([]dto.JobItemResponse, 0, len(jobs))
	for _, job := range jobs {
		// 格式化薪资范围
		salaryRange := utils.FormatSalaryRange(job.MinSalary, job.MaxSalary)

		// 格式化工作性质
		workNature := utils.FormatWorkNature(job.WorkNature)

		// 检查是否收藏
		isFavorited := favoriteMap[job.ID]

		// 检查是否投递
		var applicationID *int
		if appID, exists := applicationMap[job.ID]; exists {
			applicationID = &appID
		}

		jobItem := dto.JobItemResponse{
			JobID:         job.ID,
			Title:         job.Title,
			CompanyName:   job.CompanyName,
			SalaryRange:   salaryRange,
			Address:       job.AddressDetail,
			WorkNature:    workNature,
			Department:    job.Department,
			Headcount:     job.Headcount,
			IsFavorited:   isFavorited,
			LogoURL:       job.LogoURL,
			ApplicationID: applicationID,
		}
		jobItems = append(jobItems, jobItem)
	}

	return &dto.JobListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Jobs:     jobItems,
	}, nil
}

// GetJobDetail 获取职位详情
func (s *JobCenterService) GetJobDetail(jobID, currentUserID int) (*dto.JobDetailFullResponse, error) {
	// 获取职位详情(含Tags和职位类型)
	jobWithType, err := s.JobRepo.GetJobDetailWithRelations(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("职位不存在")
		}
		return nil, err
	}
	job := &jobWithType.Job // 提取Job部分

	// 获取公司信息(含字典表名称)
	companyWithDict, err := s.CompanyRepo.GetCompanyWithDictionaries(job.CompanyID)
	if err != nil {
		return nil, errors.New("公司信息不存在")
	}
	company := &companyWithDict.Company // 提取Company部分

	// 获取公司链接
	companyLinks, err := s.CompanyRepo.GetCompanyLinks(job.CompanyID)
	if err != nil {
		return nil, err
	}

	// 检查收藏状态
	isFavorited, err := s.FavoriteRepo.CheckFavoriteStatus(currentUserID, jobID)
	if err != nil {
		return nil, err
	}

	// 解析required_skills (从Tags获取)
	requiredSkills := make([]string, 0, len(job.Tags))
	for _, tag := range job.Tags {
		requiredSkills = append(requiredSkills, tag.Name)
	}

	// 解析bonus_points（兼容JSON数组和纯文本两种格式）
	var bonusPoints []string
	if job.BonusPoints != "" && job.BonusPoints != "null" {
		// 尝试作为JSON数组解析
		err = json.Unmarshal([]byte(job.BonusPoints), &bonusPoints)
		if err != nil {
			// 如果不是JSON格式，当作纯文本处理，放入数组
			bonusPoints = []string{job.BonusPoints}
		}
	} else {
		bonusPoints = []string{}
	}

	// 构建公司链接
	links := make([]dto.CompanyLinkInfo, 0, len(companyLinks))
	for _, link := range companyLinks {
		links = append(links, dto.CompanyLinkInfo{
			LinkName: link.LinkName,
			LinkURL:  link.LinkURL,
		})
	}

	// 组装响应
	response := &dto.JobDetailFullResponse{
		JobID:               job.ID,
		Title:               job.Title,
		SalaryRange:         utils.FormatSalaryRange(job.MinSalary, job.MaxSalary),
		Address:             job.AddressDetail,
		WorkNature:          utils.FormatWorkNature(job.WorkNature),
		RequiredDegree:      utils.FormatDegreeRequirement(job.RequiredDegree),
		RequiredStartDate:   utils.FormatDate(job.RequiredStartDate),
		Times:               job.ViewCount,
		Type:                jobWithType.TypeName, // 使用查询到的职位类型名称
		RequiredSkills:      requiredSkills,
		Headcount:           job.Headcount,
		PostedAt:            utils.FormatDate(job.CreatedAt),
		PositionDescription: job.Description,
		PositionRequirement: job.TechRequirements,
		BonusPoints:         bonusPoints,
		AddressDetail:       job.AddressDetail,
		CompanyInfo: dto.CompanyInfoNested{
			CompanyID:          company.CompanyID,
			CompanyLogoURL:     company.LogoUrl,
			CompanyName:        company.CompanyName,
			CompanyIndustry:    companyWithDict.IndustryName,
			CompanyNature:      companyWithDict.NatureName,
			CompanyScale:       companyWithDict.CompanyScaleName,
			ContactPersonName:  company.ContactPersonName,
			ContactPersonPhone: company.ContactPersonPhone,
			CompanyLinks:       links,
		},
		IsFavorited: isFavorited,
	}

	return response, nil
}
