package service

import (
	"errors"
	"fmt"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"
	"job-platform-go2/pkg/utils"
	"mime/multipart"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// ResumeService 简历服务
type ResumeService struct {
	DB               *gorm.DB
	StudentRepo      *repository.StudentRepository
	EducationRepo    *repository.EducationExperienceRepository
	WorkRepo         *repository.WorkExperienceRepository
	ProjectRepo      *repository.ProjectExperienceRepository
	OrganizationRepo *repository.OrganizationExperienceRepository
	CompetitionRepo  *repository.CompetitionExperienceRepository
	ResumeRepo       *repository.ResumeRepository
}

// GetResumeDraft 获取简历草稿
func (s *ResumeService) GetResumeDraft(studentUserID int) (*dto.ResumeDraftResponse, error) {
	// 1. 查询学生基本信息
	student, err := s.StudentRepo.FindByUserID(studentUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学生信息不存在")
		}
		return nil, err
	}

	// 2. 查询用户邮箱
	var user model.User
	if err := s.DB.Where("id = ?", studentUserID).First(&user).Error; err != nil {
		return nil, err
	}

	// 3. 构建个人资料
	profile := dto.ProfileInfo{
		FullName:         student.FullName,
		DateOfBirth:      utils.FormatYYYYMMPtr(student.DateOfBirth),
		Email:            user.Email,
		Gender:           convertGenderToString(student.Gender),
		JobSeekingStatus: convertJobSeekingStatus(student.JobSeekingStatus),
		PhoneNumber:      student.PhoneNumber,
		AvatarURL:        student.AvatarUrl,
	}

	// 4. 查询教育经历(按degree_level升序)
	educations, err := s.EducationRepo.FindByStudentUserIDOrderByDegree(studentUserID)
	if err != nil {
		return nil, err
	}

	educationInfos := make([]dto.EducationInfo, 0, len(educations))
	for _, edu := range educations {
		educationInfos = append(educationInfos, dto.EducationInfo{
			Degree:     convertDegreeLevel(edu.DegreeLevel),
			SchoolName: edu.SchoolName,
			Major:      edu.Major,
			MajorRank:  edu.MajorRank,
			StartDate:  utils.FormatYYYYMM(edu.StartDate),
			EndDate:    utils.FormatYYYYMM(edu.EndDate),
		})
	}

	// 5. 查询工作经历
	workExps, err := s.WorkRepo.FindByStudentUserID(studentUserID)
	if err != nil {
		return nil, err
	}

	workResponses := make([]dto.WorkExperienceResponse, 0, len(workExps))
	for _, work := range workExps {
		workResponses = append(workResponses, dto.WorkExperienceResponse{
			ID:            work.ID,
			CompanyName:   work.CompanyName,
			PositionTitle: work.PositionTitle,
			StartDate:     utils.FormatYYYYMM(work.StartDate),
			EndDate:       utils.FormatYYYYMM(work.EndDate),
			Description:   work.Description,
		})
	}

	// 6. 查询项目经历
	projectExps, err := s.ProjectRepo.FindByStudentUserID(studentUserID)
	if err != nil {
		return nil, err
	}

	projectResponses := make([]dto.ProjectResponse, 0, len(projectExps))
	for _, proj := range projectExps {
		projectResponses = append(projectResponses, dto.ProjectResponse{
			ID:          proj.ID,
			ProjectName: proj.ProjectName,
			Role:        proj.Role,
			ProjectLink: proj.ProjectLink,
			StartDate:   utils.FormatYYYYMM(proj.StartDate),
			EndDate:     utils.FormatYYYYMM(proj.EndDate),
			Description: proj.Description,
		})
	}

	// 7. 查询组织经历
	orgExps, err := s.OrganizationRepo.FindByStudentUserID(studentUserID)
	if err != nil {
		return nil, err
	}

	orgResponses := make([]dto.OrganizationResponse, 0, len(orgExps))
	for _, org := range orgExps {
		orgResponses = append(orgResponses, dto.OrganizationResponse{
			ID:               org.ID,
			OrganizationName: org.OrganizationName,
			Role:             org.Role,
			StartDate:        utils.FormatYYYYMM(org.StartDate),
			EndDate:          utils.FormatYYYYMM(org.EndDate),
			Description:      org.Description,
		})
	}

	// 8. 查询竞赛经历
	compExps, err := s.CompetitionRepo.FindByStudentUserID(studentUserID)
	if err != nil {
		return nil, err
	}

	compResponses := make([]dto.CompetitionResponse, 0, len(compExps))
	for _, comp := range compExps {
		compResponses = append(compResponses, dto.CompetitionResponse{
			ID:              comp.ID,
			CompetitionName: comp.CompetitionName,
			Role:            comp.Role,
			Award:           comp.Award,
			Date:            comp.Date,
		})
	}

	return &dto.ResumeDraftResponse{
		Profile:       profile,
		Education:     educationInfos,
		SkillsSummary: student.SkillsSummary,
		WorkExp:       workResponses,
		Projects:      projectResponses,
		Organizations: orgResponses,
		Competitions:  compResponses,
	}, nil
}

// UpdateSkills 更新技能信息
func (s *ResumeService) UpdateSkills(studentUserID int, req *dto.UpdateSkillsRequest) (*dto.UpdateSkillsResponse, error) {
	if err := s.StudentRepo.UpdateSkills(studentUserID, req.SkillsSummary); err != nil {
		return nil, err
	}

	return &dto.UpdateSkillsResponse{
		SkillsSummary: req.SkillsSummary,
	}, nil
}

// SetTemplate 设置简历模板
func (s *ResumeService) SetTemplate(studentUserID int, req *dto.SetTemplateRequest) (*dto.SetTemplateResponse, error) {
	if err := s.StudentRepo.UpdateTemplate(studentUserID, req.TemplateID); err != nil {
		return nil, err
	}

	// TODO: 后续如果有模板表，可以查询模板名称和预览图
	return &dto.SetTemplateResponse{
		CurrentTemplateID: req.TemplateID,
		TemplateName:      fmt.Sprintf("模板 %d", req.TemplateID),
		PreviewURL:        "", // 暂时为空
	}, nil
}

// UploadResumeFile 上传简历文件
// 注意: 实际生产环境应该上传到云存储(如阿里云OSS)，这里简化处理
func (s *ResumeService) UploadResumeFile(studentUserID int, file *multipart.FileHeader) (*dto.ResumeFileResponse, error) {
	// 1. 验证文件类型
	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" {
		return nil, errors.New("仅支持 PDF 格式的简历文件")
	}

	// 2. TODO: 实际应上传到云存储，这里简化为生成一个模拟的URL
	fileURL := fmt.Sprintf("/uploads/resumes/%d_%s", time.Now().Unix(), file.Filename)

	// 3. 创建数据库记录
	resume := &model.Resume{
		StudentUserID: studentUserID,
		FileName:      file.Filename,
		FileUrl:       fileURL,
		FileSize:      file.Size,
		UsageType:     "resume_pdf",
		UploadedAt:    time.Now(),
	}

	if err := s.ResumeRepo.Create(resume); err != nil {
		return nil, err
	}

	return &dto.ResumeFileResponse{
		ID:         resume.ID,
		FileName:   resume.FileName,
		FileURL:    resume.FileUrl,
		FileSize:   resume.FileSize,
		TemplateID: nil,
		Usage:      resume.UsageType,
		UploadedAt: resume.UploadedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteResumeFile 删除简历文件(需检查是否被引用)
func (s *ResumeService) DeleteResumeFile(resumeID int64, studentUserID int) (*dto.DeleteFileResponse, error) {
	// 1. 检查简历是否被投递使用
	isUsed, err := s.ResumeRepo.CheckResumeUsedInApplications(resumeID)
	if err != nil {
		return nil, err
	}

	if isUsed {
		return nil, errors.New("该简历已被投递使用，无法删除")
	}

	// 2. 删除数据库记录(会校验student_user_id)
	if err := s.ResumeRepo.Delete(resumeID, studentUserID); err != nil {
		return nil, err
	}

	// TODO: 删除云端文件

	return &dto.DeleteFileResponse{
		DeletedID: resumeID,
		Status:    "deleted",
	}, nil
}

// GetResumeFiles 获取简历文件列表
func (s *ResumeService) GetResumeFiles(studentUserID int) ([]dto.ResumeFileResponse, error) {
	resumes, err := s.ResumeRepo.FindByStudentUserID(studentUserID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ResumeFileResponse, 0, len(resumes))
	for _, resume := range resumes {
		var templateID *int64
		if resume.TemplateID != 0 {
			tempID := int64(resume.TemplateID)
			templateID = &tempID
		}

		responses = append(responses, dto.ResumeFileResponse{
			ID:         resume.ID,
			FileName:   resume.FileName,
			FileURL:    resume.FileUrl,
			FileSize:   resume.FileSize,
			TemplateID: templateID,
			Usage:      resume.UsageType,
			UploadedAt: resume.UploadedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return responses, nil
}

// ==================== 辅助函数 ====================

// convertGenderToString 将性别代码转换为字符串
func convertGenderToString(gender int) string {
	switch gender {
	case 0:
		return "male"
	case 1:
		return "female"
	default:
		return "unknown"
	}
}

// convertJobSeekingStatus 将求职状态代码转换为字符串
func convertJobSeekingStatus(status int) string {
	switch status {
	case 0:
		return "on_campus_not_seeking"
	case 1:
		return "on_campus_seeking_internship"
	case 2:
		return "graduating_seeking_internship"
	case 3:
		return "graduating_seeking_fulltime"
	default:
		return "unknown"
	}
}

// convertDegreeLevel 将学历等级转换为字符串
func convertDegreeLevel(level int) string {
	switch level {
	case 0:
		return "本科"
	case 1:
		return "硕士"
	case 2:
		return "博士"
	default:
		return "其他"
	}
}
