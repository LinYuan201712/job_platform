package service

import (
	"errors"
	"fmt"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"
	"time"

	"gorm.io/gorm"
)

type JobService struct {
	JobRepo *repository.JobRepository
	TagRepo *repository.TagRepository
	DB      *gorm.DB
}

// CreateJob 创建岗位(草稿或提交申请)
func (s *JobService) CreateJob(req dto.CreateJobRequest, userID int, companyID int) (*dto.CreateJobResponse, error) {
	// 1. 业务校验: 最低薪资不能高于最高薪资
	if req.MinSalary > req.MaxSalary {
		return nil, errors.New("最低薪资不能高于最高薪资")
	}

	// 2. 解析日期
	requiredStartDate, err := time.Parse("2006-01-02", req.RequiredStartDate)
	if err != nil {
		return nil, errors.New("要求到岗时间格式错误")
	}
	deadline, err := time.Parse("2006-01-02", req.Deadline)
	if err != nil {
		return nil, errors.New("招聘截止时间格式错误")
	}

	// 3. 获取标签
	tags, err := s.TagRepo.GetTagsByIDs(req.Tags)
	if err != nil {
		return nil, err
	}

	// 4. 构建 Job 对象
	job := &model.Job{
		CompanyID:         companyID,
		PostedByUserID:    userID,
		Title:             req.Title,
		Description:       req.Description,
		TechRequirements:  req.TechRequirements,
		MinSalary:         req.MinSalary,
		MaxSalary:         req.MaxSalary,
		ProvinceID:        req.ProvinceID,
		CityID:            req.CityID,
		AddressDetail:     req.AddressDetail,
		WorkNature:        req.WorkNature,
		Deadline:          deadline,
		Status:            req.Status, // 1=draft, 10=pending
		Type:              req.Type,
		Department:        req.Department,
		Headcount:         req.Headcount,
		RequiredDegree:    req.RequiredDegree,
		RequiredStartDate: requiredStartDate,
		BonusPoints:       req.BonusPoints,
		Tags:              tags,
	}

	// 5. 保存到数据库
	if err := s.JobRepo.CreateJob(job); err != nil {
		return nil, fmt.Errorf("创建岗位失败: %v", err)
	}

	// 6. 构建响应
	statusStr := "draft"
	if req.Status == 10 {
		statusStr = "pending"
	}

	resp := &dto.CreateJobResponse{}
	resp.NewJob.JobID = job.ID
	resp.NewJob.Title = job.Title
	resp.NewJob.Status = statusStr

	return resp, nil
}

// GetJobDetail 获取岗位详情
func (s *JobService) GetJobDetail(jobID int) (*dto.JobDetailResponse, error) {
	// 1. 从数据库获取岗位(含标签)
	job, err := s.JobRepo.GetJobByID(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("岗位不存在或无权访问")
		}
		return nil, err
	}

	// 2. 转换为响应格式
	resp := &dto.JobDetailResponse{
		JobID:             job.ID,
		Title:             job.Title,
		Status:            job.Status,
		Description:       job.Description,
		TechRequirements:  &job.TechRequirements,
		BonusPoints:       job.BonusPoints,
		MinSalary:         job.MinSalary,
		MaxSalary:         job.MaxSalary,
		ProvinceID:        job.ProvinceID,
		CityID:            job.CityID,
		AddressDetail:     job.AddressDetail,
		WorkNature:        job.WorkNature,
		Department:        job.Department,
		Headcount:         job.Headcount,
		Type:              job.Type,
		RequiredDegree:    job.RequiredDegree,
		RequiredStartDate: job.RequiredStartDate.Format("2006-01-02"),
		Deadline:          job.Deadline.Format("2006-01-02"),
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
	}

	// 3. 转换标签
	for _, tag := range job.Tags {
		resp.Tags = append(resp.Tags, dto.TagInfo{
			TagID:   tag.ID,
			TagName: tag.Name,
		})
	}

	return resp, nil
}

// UpdateJob 更新岗位
func (s *JobService) UpdateJob(jobID int, req dto.UpdateJobRequest, userID int) (*dto.UpdateJobResponse, error) {
	// 1. 校验薪资
	if req.MinSalary > req.MaxSalary {
		return nil, errors.New("最低薪资不能高于最高薪资")
	}

	// 2. 验证岗位归属
	isOwner, err := s.JobRepo.CheckJobOwnership(jobID, userID)
	if err != nil || !isOwner {
		return nil, errors.New("岗位不存在或无权访问")
	}

	// 3. 获取现有岗位
	job, err := s.JobRepo.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}

	// 4. 解析日期
	requiredStartDate, err := time.Parse("2006-01-02", req.RequiredStartDate)
	if err != nil {
		return nil, errors.New("要求到岗时间格式错误")
	}
	deadline, err := time.Parse("2006-01-02", req.Deadline)
	if err != nil {
		return nil, errors.New("招聘截止时间格式错误")
	}

	// 5. 获取新的标签
	tags, err := s.TagRepo.GetTagsByIDs(req.Tags)
	if err != nil {
		return nil, err
	}

	// 6. 更新字段
	job.Title = req.Title
	job.Description = req.Description
	job.TechRequirements = req.TechRequirements
	job.MinSalary = req.MinSalary
	job.MaxSalary = req.MaxSalary
	job.ProvinceID = req.ProvinceID
	job.CityID = req.CityID
	job.AddressDetail = req.AddressDetail
	job.WorkNature = req.WorkNature
	job.Deadline = deadline
	job.Status = req.Status
	job.Type = req.Type
	job.Department = req.Department
	job.Headcount = req.Headcount
	job.RequiredDegree = req.RequiredDegree
	job.RequiredStartDate = requiredStartDate
	job.BonusPoints = req.BonusPoints
	job.Tags = tags

	// 7. 保存
	if err := s.JobRepo.UpdateJob(job); err != nil {
		return nil, errors.New("更新失败")
	}

	// 8. 构建响应
	statusStr := s.getStatusString(job.Status)
	resp := &dto.UpdateJobResponse{}
	resp.UpdatedJob.JobID = job.ID
	resp.UpdatedJob.Title = job.Title
	resp.UpdatedJob.Status = statusStr

	return resp, nil
}

// DeleteJobDraft 删除岗位草稿(硬删除)
func (s *JobService) DeleteJobDraft(jobID int, userID int) error {
	// 1. 验证岗位归属
	isOwner, err := s.JobRepo.CheckJobOwnership(jobID, userID)
	if err != nil || !isOwner {
		return errors.New("岗位不存在或无权访问")
	}

	// 2. 获取岗位检查状态
	job, err := s.JobRepo.GetJobByID(jobID)
	if err != nil {
		return errors.New("岗位不存在")
	}

	// 3. 只能删除草稿(status=1)
	if job.Status != 1 {
		return errors.New("操作失败:只能删除草稿状态的岗位")
	}

	// 4. 执行删除
	return s.JobRepo.DeleteJob(jobID)
}

// OfflineJob 下线岗位(只能下线已发布的岗位)
func (s *JobService) OfflineJob(jobID int, userID int) (*dto.OfflineJobResponse, error) {
	// 1. 验证岗位归属
	isOwner, err := s.JobRepo.CheckJobOwnership(jobID, userID)
	if err != nil || !isOwner {
		return nil, errors.New("岗位不存在或无权访问")
	}

	// 2. 获取岗位检查状态
	job, err := s.JobRepo.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}

	// 3. 只能下线已发布的岗位(status=20 表示approved/published)
	// 注意: 数据库注释说 20=approved, 但API文档里提到"已发布",这里按数据库为准
	if job.Status != 20 {
		return nil, errors.New("仅已发布岗位可下线")
	}

	// 4. 更新状态为 40=closed
	if err := s.JobRepo.UpdateJobStatus(jobID, 40); err != nil {
		return nil, errors.New("下线失败")
	}

	// 5. 构建响应
	return &dto.OfflineJobResponse{
		JobID:  jobID,
		Status: "closed",
	}, nil
}

// getStatusString 将状态码转为字符串
func (s *JobService) getStatusString(status int) string {
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
		return fmt.Sprintf("unknown_%d", status)
	}
}
