package service

import (
	"errors"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"
	"job-platform-go2/pkg/utils"

	"gorm.io/gorm"
)

// ExperienceService 处理各类经历的业务逻辑
type ExperienceService struct {
	WorkRepo         *repository.WorkExperienceRepository
	ProjectRepo      *repository.ProjectExperienceRepository
	OrganizationRepo *repository.OrganizationExperienceRepository
	CompetitionRepo  *repository.CompetitionExperienceRepository
}

// ==================== 工作经历 ====================

// CreateWorkExperience 创建工作经历
func (s *ExperienceService) CreateWorkExperience(studentUserID int, req *dto.CreateWorkExperienceRequest) (*dto.WorkExperienceResponse, error) {
	// 日期转换
	startDate, err := utils.ParseYYYYMM(req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，应为 YYYY-MM")
	}

	endDate, err := utils.ParseYYYYMMOptional(req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，应为 YYYY-MM")
	}

	exp := &model.WorkExperience{
		StudentUserID: studentUserID,
		CompanyName:   req.CompanyName,
		PositionTitle: req.PositionTitle,
		StartDate:     startDate,
		EndDate:       endDate,
		Description:   req.Description,
	}

	if err := s.WorkRepo.Create(exp); err != nil {
		return nil, err
	}

	return &dto.WorkExperienceResponse{
		ID:            exp.ID,
		CompanyName:   exp.CompanyName,
		PositionTitle: exp.PositionTitle,
		StartDate:     utils.FormatYYYYMM(exp.StartDate),
		EndDate:       utils.FormatYYYYMM(exp.EndDate),
		Description:   exp.Description,
	}, nil
}

// UpdateWorkExperience 更新工作经历
func (s *ExperienceService) UpdateWorkExperience(id int64, studentUserID int, req *dto.UpdateWorkExperienceRequest) (*dto.WorkExperienceResponse, error) {
	// 先查询记录，校验权限
	existing, err := s.WorkRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工作经历不存在")
		}
		return nil, err
	}

	if existing.StudentUserID != studentUserID {
		return nil, errors.New("无权操作此工作经历")
	}

	// 日期转换
	startDate, err := utils.ParseYYYYMM(req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，应为 YYYY-MM")
	}

	endDate, err := utils.ParseYYYYMMOptional(req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，应为 YYYY-MM")
	}

	// 更新字段
	existing.CompanyName = req.CompanyName
	existing.PositionTitle = req.PositionTitle
	existing.StartDate = startDate
	existing.EndDate = endDate
	existing.Description = req.Description

	if err := s.WorkRepo.Update(existing); err != nil {
		return nil, err
	}

	return &dto.WorkExperienceResponse{
		ID:            existing.ID,
		CompanyName:   existing.CompanyName,
		PositionTitle: existing.PositionTitle,
		StartDate:     utils.FormatYYYYMM(existing.StartDate),
		EndDate:       utils.FormatYYYYMM(existing.EndDate),
		Description:   existing.Description,
	}, nil
}

// DeleteWorkExperience 删除工作经历
func (s *ExperienceService) DeleteWorkExperience(id int64, studentUserID int) error {
	return s.WorkRepo.Delete(id, studentUserID)
}

// ==================== 项目经历 ====================

// CreateProjectExperience 创建项目经历
func (s *ExperienceService) CreateProjectExperience(studentUserID int, req *dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	startDate, err := utils.ParseYYYYMMOptional(req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，应为 YYYY-MM")
	}

	endDate, err := utils.ParseYYYYMMOptional(req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，应为 YYYY-MM")
	}

	exp := &model.ProjectExperience{
		StudentUserID: studentUserID,
		ProjectName:   req.ProjectName,
		Role:          req.Role,
		ProjectLink:   req.ProjectLink,
		StartDate:     startDate,
		EndDate:       endDate,
		Description:   req.Description,
	}

	if err := s.ProjectRepo.Create(exp); err != nil {
		return nil, err
	}

	return &dto.ProjectResponse{
		ID:          exp.ID,
		ProjectName: exp.ProjectName,
		Role:        exp.Role,
		ProjectLink: exp.ProjectLink,
		StartDate:   utils.FormatYYYYMM(exp.StartDate),
		EndDate:     utils.FormatYYYYMM(exp.EndDate),
		Description: exp.Description,
	}, nil
}

// UpdateProjectExperience 更新项目经历
func (s *ExperienceService) UpdateProjectExperience(id int64, studentUserID int, req *dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	existing, err := s.ProjectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目经历不存在")
		}
		return nil, err
	}

	if existing.StudentUserID != studentUserID {
		return nil, errors.New("无权操作此项目经历")
	}

	startDate, err := utils.ParseYYYYMMOptional(req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，应为 YYYY-MM")
	}

	endDate, err := utils.ParseYYYYMMOptional(req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，应为 YYYY-MM")
	}

	existing.ProjectName = req.ProjectName
	existing.Role = req.Role
	existing.ProjectLink = req.ProjectLink
	existing.StartDate = startDate
	existing.EndDate = endDate
	existing.Description = req.Description

	if err := s.ProjectRepo.Update(existing); err != nil {
		return nil, err
	}

	return &dto.ProjectResponse{
		ID:          existing.ID,
		ProjectName: existing.ProjectName,
		Role:        existing.Role,
		ProjectLink: existing.ProjectLink,
		StartDate:   utils.FormatYYYYMM(existing.StartDate),
		EndDate:     utils.FormatYYYYMM(existing.EndDate),
		Description: existing.Description,
	}, nil
}

// DeleteProjectExperience 删除项目经历
func (s *ExperienceService) DeleteProjectExperience(id int64, studentUserID int) error {
	return s.ProjectRepo.Delete(id, studentUserID)
}

// ==================== 组织经历 ====================

// CreateOrganizationExperience 创建组织经历
func (s *ExperienceService) CreateOrganizationExperience(studentUserID int, req *dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error) {
	startDate, err := utils.ParseYYYYMMOptional(req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，应为 YYYY-MM")
	}

	endDate, err := utils.ParseYYYYMMOptional(req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，应为 YYYY-MM")
	}

	exp := &model.OrganizationExperience{
		StudentUserID:    studentUserID,
		OrganizationName: req.OrganizationName,
		Role:             req.Role,
		StartDate:        startDate,
		EndDate:          endDate,
		Description:      req.Description,
	}

	if err := s.OrganizationRepo.Create(exp); err != nil {
		return nil, err
	}

	return &dto.OrganizationResponse{
		ID:               exp.ID,
		OrganizationName: exp.OrganizationName,
		Role:             exp.Role,
		StartDate:        utils.FormatYYYYMM(exp.StartDate),
		EndDate:          utils.FormatYYYYMM(exp.EndDate),
		Description:      exp.Description,
	}, nil
}

// UpdateOrganizationExperience 更新组织经历
func (s *ExperienceService) UpdateOrganizationExperience(id int64, studentUserID int, req *dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error) {
	existing, err := s.OrganizationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("组织经历不存在")
		}
		return nil, err
	}

	if existing.StudentUserID != studentUserID {
		return nil, errors.New("无权操作此组织经历")
	}

	startDate, err := utils.ParseYYYYMMOptional(req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，应为 YYYY-MM")
	}

	endDate, err := utils.ParseYYYYMMOptional(req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，应为 YYYY-MM")
	}

	existing.OrganizationName = req.OrganizationName
	existing.Role = req.Role
	existing.StartDate = startDate
	existing.EndDate = endDate
	existing.Description = req.Description

	if err := s.OrganizationRepo.Update(existing); err != nil {
		return nil, err
	}

	return &dto.OrganizationResponse{
		ID:               existing.ID,
		OrganizationName: existing.OrganizationName,
		Role:             existing.Role,
		StartDate:        utils.FormatYYYYMM(existing.StartDate),
		EndDate:          utils.FormatYYYYMM(existing.EndDate),
		Description:      existing.Description,
	}, nil
}

// DeleteOrganizationExperience 删除组织经历
func (s *ExperienceService) DeleteOrganizationExperience(id int64, studentUserID int) error {
	return s.OrganizationRepo.Delete(id, studentUserID)
}

// ==================== 竞赛经历 ====================

// CreateCompetitionExperience 创建竞赛经历
// 注意: date 字段直接存储 YYYY-MM 字符串
func (s *ExperienceService) CreateCompetitionExperience(studentUserID int, req *dto.CreateCompetitionRequest) (*dto.CompetitionResponse, error) {
	exp := &model.CompetitionExperience{
		StudentUserID:   studentUserID,
		CompetitionName: req.CompetitionName,
		Role:            req.Role,
		Award:           req.Award,
		Date:            req.Date,
	}

	if err := s.CompetitionRepo.Create(exp); err != nil {
		return nil, err
	}

	return &dto.CompetitionResponse{
		ID:              exp.ID,
		CompetitionName: exp.CompetitionName,
		Role:            exp.Role,
		Award:           exp.Award,
		Date:            exp.Date,
	}, nil
}

// UpdateCompetitionExperience 更新竞赛经历
func (s *ExperienceService) UpdateCompetitionExperience(id int64, studentUserID int, req *dto.UpdateCompetitionRequest) (*dto.CompetitionResponse, error) {
	existing, err := s.CompetitionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("竞赛经历不存在")
		}
		return nil, err
	}

	if existing.StudentUserID != studentUserID {
		return nil, errors.New("无权操作此竞赛经历")
	}

	existing.CompetitionName = req.CompetitionName
	existing.Role = req.Role
	existing.Award = req.Award
	existing.Date = req.Date

	if err := s.CompetitionRepo.Update(existing); err != nil {
		return nil, err
	}

	return &dto.CompetitionResponse{
		ID:              existing.ID,
		CompetitionName: existing.CompetitionName,
		Role:            existing.Role,
		Award:           existing.Award,
		Date:            existing.Date,
	}, nil
}

// DeleteCompetitionExperience 删除竞赛经历
func (s *ExperienceService) DeleteCompetitionExperience(id int64, studentUserID int) error {
	return s.CompetitionRepo.Delete(id, studentUserID)
}
