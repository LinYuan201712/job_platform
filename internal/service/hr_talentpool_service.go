package service

import (
	"errors"
	"fmt"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"
	"math"
	"time"

	"gorm.io/gorm"
)

type HRTalentpoolService struct {
	ApplicationRepo *repository.ApplicationRepository
	JobRepo         *repository.JobRepository
	StudentRepo     *repository.StudentRepository
	ResumeRepo      *repository.ResumeRepository
	EducationRepo   *repository.EducationExperienceRepository
	TagRepo         *repository.TagRepository
	DB              *gorm.DB
}

// MapStatusCodeToText 状态码 → 文本映射
func MapStatusCodeToText(code int) string {
	switch code {
	case model.AppStatusSubmitted:
		return "已投递"
	case model.AppStatusCandidate:
		return "候选人"
	case model.AppStatusInterview:
		return "面试邀请"
	case model.AppStatusPassed:
		return "通过"
	case model.AppStatusRejected:
		return "拒绝"
	default:
		return "未知状态"
	}
}

// MapStatusTextToCode 文本 → 状态码映射(防御性检查)
func MapStatusTextToCode(text string) (int, error) {
	switch text {
	case "已投递":
		return model.AppStatusSubmitted, nil
	case "候选人", "设置为候选人":
		return model.AppStatusCandidate, nil
	case "面试邀请", "邀请面试":
		return model.AppStatusInterview, nil
	case "通过":
		return model.AppStatusPassed, nil
	case "拒绝":
		return model.AppStatusRejected, nil
	default:
		return 0, fmt.Errorf("无效的状态文本: %s", text)
	}
}

// MapWorkNatureToText 岗位性质映射
func MapWorkNatureToText(workNature int) string {
	if workNature == 1 {
		return "实习"
	} else if workNature == 2 {
		return "全职"
	}
	return ""
}

// MapJobStatusToText 岗位状态映射
func MapJobStatusToText(status int) string {
	switch status {
	case 1:
		return "draft"
	case 10:
		return "pending"
	case 20:
		return "approved"
	case 30:
		return "rejected"
	case 40:
		return "closed"
	default:
		return "unknown"
	}
}

// GetTalentPoolJobList 获取人才库岗位列表(带投递统计)
func (s *HRTalentpoolService) GetTalentPoolJobList(req dto.TalentPoolJobListRequest, companyID int) (*dto.TalentPoolJobListResponse, error) {
	// 默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 1. 查询企业岗位列表(带筛选)
	query := s.DB.Model(&model.Job{}).Where("company_id = ?", companyID)

	// 筛选条件
	if req.TitleKeyword != "" {
		query = query.Where("title LIKE ?", "%"+req.TitleKeyword+"%")
	}
	if req.WorkNature != "" {
		if req.WorkNature == "实习" || req.WorkNature == "internship" {
			query = query.Where("work_nature = ?", 1)
		} else if req.WorkNature == "全职" || req.WorkNature == "full-time" {
			query = query.Where("work_nature = ?", 2)
		}
	}
	if req.Status != "" {
		statusMap := map[string]int{
			"draft":    1,
			"pending":  10,
			"approved": 20,
			"rejected": 30,
			"closed":   40,
		}
		if statusCode, ok := statusMap[req.Status]; ok {
			query = query.Where("status = ?", statusCode)
		}
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询岗位
	var jobs []model.Job
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("updated_at DESC").Offset(offset).Limit(req.PageSize).Find(&jobs).Error; err != nil {
		return nil, err
	}

	// 2. 提取job_id列表
	jobIDs := make([]int, 0, len(jobs))
	for _, job := range jobs {
		jobIDs = append(jobIDs, job.ID)
	}

	// 3. **性能优化**: 批量查询投递数量
	receivedMap, err := s.ApplicationRepo.GetApplicationCountsByJobIDs(jobIDs)
	if err != nil {
		return nil, err
	}

	unreviewedMap, err := s.ApplicationRepo.GetUnreviewedCountsByJobIDs(jobIDs)
	if err != nil {
		return nil, err
	}

	// 4. 组装响应
	jobList := make([]dto.TalentPoolJobResponse, 0, len(jobs))
	for _, job := range jobs {
		jobList = append(jobList, dto.TalentPoolJobResponse{
			Title:       job.Title,
			Status:      MapJobStatusToText(job.Status),
			JobID:       job.ID,
			WorkNature:  MapWorkNatureToText(job.WorkNature),
			UpdatedAt:   job.UpdatedAt.Format("2006-01-02T15:04:05"),
			ReceivedNum: receivedMap[job.ID],
			NoReviewNum: unreviewedMap[job.ID],
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &dto.TalentPoolJobListResponse{
		JobList: jobList,
		Pagination: dto.PaginationResponse{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
		},
	}, nil
}

// GetCandidateListByJob 获取岗位候选人列表
func (s *HRTalentpoolService) GetCandidateListByJob(req dto.CandidateListRequest, companyID int) (*dto.CandidateListResponse, error) {
	// 默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 1. 验证job_id属于当前企业
	job, err := s.JobRepo.GetJobByID(req.JobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("岗位不存在")
		}
		return nil, err
	}
	if job.CompanyID != companyID {
		return nil, errors.New("无权访问该岗位")
	}

	// 2. 查询候选人列表
	candidates, total, err := s.ApplicationRepo.GetCandidatesByJobID(
		req.JobID, req.NameKeyword, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	candidateList := make([]dto.CandidateItem, 0, len(candidates))
	for _, c := range candidates {
		// 学历映射
		degreeText := ""
		switch c.Degree {
		case 0:
			degreeText = "bachelor"
		case 1:
			degreeText = "master"
		case 2:
			degreeText = "doctor"
		}

		candidateList = append(candidateList, dto.CandidateItem{
			Grade:         c.Grade,
			Degree:        degreeText,
			ApplicationID: c.ApplicationID,
			CandidateName: c.FullName,
			AvatarURL:     c.AvatarURL,
			UserID:        c.UserID,
			ResumeStatus:  MapStatusCodeToText(c.Status),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &dto.CandidateListResponse{
		CandidateList: candidateList,
		Pagination: dto.PaginationResponse{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
		},
	}, nil
}

// GetApplicationDetail 获取候选人简历详情(HR端)
func (s *HRTalentpoolService) GetApplicationDetail(appID int, companyID int) (map[string]interface{}, error) {
	// 1. 查询投递记录
	app, err := s.ApplicationRepo.GetApplicationDetailByID(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("投递记录不存在")
		}
		return nil, err
	}

	// 2. 验证权限: application的job属于当前企业
	job, err := s.JobRepo.GetJobByID(app.JobID)
	if err != nil {
		return nil, errors.New("岗位信息不存在")
	}
	if job.CompanyID != companyID {
		return nil, errors.New("无权访问该投递记录")
	}

	// 3. 查询简历文件
	resume, err := s.ResumeRepo.GetResumeByID(app.ResumeID)
	if err != nil {
		return nil, errors.New("简历文件不存在")
	}

	// 4. 返回数据
	return map[string]interface{}{
		"id":         app.ID,
		"status":     MapStatusCodeToText(app.Status),
		"resume_url": resume.FileUrl,
	}, nil
}

// UpdateApplicationStatus 更新人才状态
func (s *HRTalentpoolService) UpdateApplicationStatus(appID int, statusText string, companyID int) (*dto.UpdateStatusResponse, error) {
	// 1. 验证权限
	app, err := s.ApplicationRepo.GetApplicationDetailByID(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("投递记录不存在")
		}
		return nil, err
	}

	job, err := s.JobRepo.GetJobByID(app.JobID)
	if err != nil {
		return nil, errors.New("岗位信息不存在")
	}
	if job.CompanyID != companyID {
		return nil, errors.New("无权操作该投递记录")
	}

	// 2. 状态文本 → 状态码
	statusCode, err := MapStatusTextToCode(statusText)
	if err != nil {
		return nil, err
	}

	// 3. 更新状态
	if err := s.ApplicationRepo.UpdateApplicationStatus(appID, statusCode); err != nil {
		return nil, err
	}

	// 4. 返回
	return &dto.UpdateStatusResponse{
		Status:        statusText,
		ApplicationID: appID,
		StatusCode:    statusCode,
	}, nil
}

// GetStudentResumePreview HR端获取学生简历预览
func (s *HRTalentpoolService) GetStudentResumePreview(studentUserID int) (*dto.StudentResumePreviewResponse, error) {
	// 1. 查询学生基本信息
	student, err := s.StudentRepo.FindByUserID(studentUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学生不存在")
		}
		return nil, err
	}

	// 2. 基本信息
	genderText := "男"
	if student.Gender == 1 {
		genderText = "女"
	}

	age := 0
	if student.DateOfBirth != nil {
		age = time.Now().Year() - student.DateOfBirth.Year()
	}

	jobSeekingStatusText := ""
	switch student.JobSeekingStatus {
	case 0:
		jobSeekingStatusText = "在校-暂不考虑"
	case 1:
		jobSeekingStatusText = "在校-寻求实习"
	case 2:
		jobSeekingStatusText = "应届-寻求实习"
	case 3:
		jobSeekingStatusText = "应届-寻求校招"
	}

	basicInfo := dto.BasicInfoDTO{
		FullName:         student.FullName,
		Gender:           genderText,
		Age:              age,
		DegreeLevel:      "", // 从教育经历中获取最高学历
		JobSeekingStatus: jobSeekingStatusText,
	}

	// 3. 教育经历
	educations, err := s.EducationRepo.FindByStudentUserID(studentUserID)
	if err != nil {
		return nil, err
	}

	primaryEducation := make([]dto.EducationExperienceDTO, 0, len(educations))
	maxDegree := 0
	for _, edu := range educations {
		degreeText := ""
		switch edu.DegreeLevel {
		case 0:
			degreeText = "本科"
		case 1:
			degreeText = "硕士"
		case 2:
			degreeText = "博士"
		}

		if edu.DegreeLevel > maxDegree {
			maxDegree = edu.DegreeLevel
			basicInfo.DegreeLevel = degreeText
		}

		startDate := ""
		if !edu.StartDate.IsZero() {
			startDate = edu.StartDate.Format("2006-01-02")
		}
		endDate := ""
		if !edu.EndDate.IsZero() {
			endDate = edu.EndDate.Format("2006-01-02")
		}

		primaryEducation = append(primaryEducation, dto.EducationExperienceDTO{
			SchoolName:  edu.SchoolName,
			DegreeLevel: degreeText,
			Major:       edu.Major,
			StartDate:   startDate,
			EndDate:     endDate,
			MajorRank:   edu.MajorRank,
		})
	}

	// 4. 期望岗位
	expectedJob := dto.ExpectedJobDTO{
		ExpectedPosition:  student.ExpectedPosition,
		ExpectedMinSalary: student.ExpectedMinSalary,
		ExpectedMaxSalary: student.ExpectedMaxSalary,
	}

	// 5. 个人标签
	var tags []model.Tag
	err = s.DB.Table("student_tags").
		Select("tags.*").
		Joins("INNER JOIN tags ON student_tags.tag_id = tags.id").
		Where("student_tags.student_user_id = ?", studentUserID).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}

	personalTags := make([]dto.PersonalTagDTO, 0, len(tags))
	for _, tag := range tags {
		personalTags = append(personalTags, dto.PersonalTagDTO{
			TagID: tag.ID,
			Name:  tag.Name,
		})
	}

	return &dto.StudentResumePreviewResponse{
		AvatarURL:        student.AvatarUrl,
		BasicInfo:        basicInfo,
		PrimaryEducation: primaryEducation,
		ExpectedJob:      expectedJob,
		PersonalTags:     personalTags,
	}, nil
}
