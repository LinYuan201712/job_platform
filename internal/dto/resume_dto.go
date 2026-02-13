package dto

// ==================== 简历草稿相关 DTO ====================

// UpdateSkillsRequest 更新技能信息请求
type UpdateSkillsRequest struct {
	SkillsSummary string `json:"skills_summary" binding:"required"` // 技能掌握描述
}

// UpdateSkillsResponse 更新技能信息响应
type UpdateSkillsResponse struct {
	SkillsSummary string `json:"skills_summary"` // 更新后的技能描述
}

// ==================== 工作经历相关 DTO ====================

// CreateWorkExperienceRequest 创建工作经历请求
type CreateWorkExperienceRequest struct {
	CompanyName   string `json:"company_name" binding:"required"`   // 公司名称
	PositionTitle string `json:"position_title" binding:"required"` // 职位名称
	StartDate     string `json:"start_date" binding:"required"`     // 开始时间 YYYY-MM
	EndDate       string `json:"end_date"`                          // 结束时间 YYYY-MM（可选）
	Description   string `json:"description"`                       // 工作内容描述
}

// UpdateWorkExperienceRequest 更新工作经历请求
type UpdateWorkExperienceRequest struct {
	CompanyName   string `json:"company_name" binding:"required"`   // 公司名称
	PositionTitle string `json:"position_title" binding:"required"` // 职位名称
	StartDate     string `json:"start_date" binding:"required"`     // 开始时间 YYYY-MM
	EndDate       string `json:"end_date"`                          // 结束时间 YYYY-MM（可选）
	Description   string `json:"description"`                       // 工作内容描述
}

// WorkExperienceResponse 工作经历响应
type WorkExperienceResponse struct {
	ID            int64  `json:"id"`             // 工作经历ID
	CompanyName   string `json:"company_name"`   // 公司名称
	PositionTitle string `json:"position_title"` // 职位名称
	StartDate     string `json:"start_date"`     // 开始时间 YYYY-MM
	EndDate       string `json:"end_date"`       // 结束时间 YYYY-MM
	Description   string `json:"description"`    // 工作内容描述
}

// ==================== 项目经历相关 DTO ====================

// CreateProjectRequest 创建项目经历请求
type CreateProjectRequest struct {
	ProjectName string `json:"project_name" binding:"required"` // 项目名称
	Role        string `json:"role"`                            // 担任角色
	ProjectLink string `json:"project_link"`                    // 项目链接
	StartDate   string `json:"start_date"`                      // 开始时间 YYYY-MM
	EndDate     string `json:"end_date"`                        // 结束时间 YYYY-MM
	Description string `json:"description"`                     // 项目描述
}

// UpdateProjectRequest 更新项目经历请求
type UpdateProjectRequest struct {
	ProjectName string `json:"project_name" binding:"required"` // 项目名称
	Role        string `json:"role"`                            // 担任角色
	ProjectLink string `json:"project_link"`                    // 项目链接
	StartDate   string `json:"start_date"`                      // 开始时间 YYYY-MM
	EndDate     string `json:"end_date"`                        // 结束时间 YYYY-MM
	Description string `json:"description"`                     // 项目描述
}

// ProjectResponse 项目经历响应
type ProjectResponse struct {
	ID          int64  `json:"id"`           // 项目经历ID
	ProjectName string `json:"project_name"` // 项目名称
	Role        string `json:"role"`         // 担任角色
	ProjectLink string `json:"project_link"` // 项目链接
	StartDate   string `json:"start_date"`   // 开始时间 YYYY-MM
	EndDate     string `json:"end_date"`     // 结束时间 YYYY-MM
	Description string `json:"description"`  // 项目描述
}

// ==================== 组织经历相关 DTO ====================

// CreateOrganizationRequest 创建组织经历请求
type CreateOrganizationRequest struct {
	OrganizationName string `json:"organization_name" binding:"required"` // 组织名称
	Role             string `json:"role"`                                 // 担任角色
	StartDate        string `json:"start_date"`                           // 开始时间 YYYY-MM
	EndDate          string `json:"end_date"`                             // 结束时间 YYYY-MM
	Description      string `json:"description"`                          // 主要工作内容
}

// UpdateOrganizationRequest 更新组织经历请求
type UpdateOrganizationRequest struct {
	OrganizationName string `json:"organization_name" binding:"required"` // 组织名称
	Role             string `json:"role"`                                 // 担任角色
	StartDate        string `json:"start_date"`                           // 开始时间 YYYY-MM
	EndDate          string `json:"end_date"`                             // 结束时间 YYYY-MM
	Description      string `json:"description"`                          // 主要工作内容
}

// OrganizationResponse 组织经历响应
type OrganizationResponse struct {
	ID               int64  `json:"id"`                // 组织经历ID
	OrganizationName string `json:"organization_name"` // 组织名称
	Role             string `json:"role"`              // 担任角色
	StartDate        string `json:"start_date"`        // 开始时间 YYYY-MM
	EndDate          string `json:"end_date"`          // 结束时间 YYYY-MM
	Description      string `json:"description"`       // 主要工作内容
}

// ==================== 竞赛经历相关 DTO ====================

// CreateCompetitionRequest 创建竞赛经历请求
type CreateCompetitionRequest struct {
	CompetitionName string `json:"competition_name" binding:"required"` // 竞赛名称
	Role            string `json:"role" binding:"required"`             // 担任角色
	Award           string `json:"award" binding:"required"`            // 奖项
	Date            string `json:"date" binding:"required"`             // 竞赛时间 YYYY-MM
}

// UpdateCompetitionRequest 更新竞赛经历请求
type UpdateCompetitionRequest struct {
	CompetitionName string `json:"competition_name" binding:"required"` // 竞赛名称
	Role            string `json:"role" binding:"required"`             // 担任角色
	Award           string `json:"award" binding:"required"`            // 奖项
	Date            string `json:"date" binding:"required"`             // 竞赛时间 YYYY-MM
}

// CompetitionResponse 竞赛经历响应
type CompetitionResponse struct {
	ID              int64  `json:"id"`               // 竞赛经历ID
	CompetitionName string `json:"competition_name"` // 竞赛名称
	Role            string `json:"role"`             // 担任角色
	Award           string `json:"award"`            // 奖项
	Date            string `json:"date"`             // 竞赛时间 YYYY-MM
}

// ==================== 模板设置相关 DTO ====================

// SetTemplateRequest 设置简历模板请求
type SetTemplateRequest struct {
	TemplateID int64 `json:"template_id" binding:"required"` // 所选模板ID
}

// SetTemplateResponse 设置简历模板响应
type SetTemplateResponse struct {
	CurrentTemplateID int64  `json:"current_template_id"` // 当前已设置的模板ID
	TemplateName      string `json:"template_name"`       // 模板名称
	PreviewURL        string `json:"preview_url"`         // 模板预览图URL
}

// ==================== 简历文件相关 DTO ====================

// ResumeFileResponse 简历文件响应
type ResumeFileResponse struct {
	ID         int64  `json:"id"`          // 简历文件ID
	FileName   string `json:"file_name"`   // 文件名
	FileURL    string `json:"file_url"`    // 文件URL
	FileSize   int64  `json:"file_size"`   // 文件大小（字节）
	TemplateID *int64 `json:"template_id"` // 使用的模板ID（可为空）
	Usage      string `json:"usage"`       // 用途类型
	UploadedAt string `json:"uploaded_at"` // 上传时间
}

// DeleteFileResponse 删除文件响应
type DeleteFileResponse struct {
	DeletedID int64  `json:"deleted_id"` // 已删除的文件ID
	Status    string `json:"status"`     // 删除状态
}

// ==================== 简历草稿响应 DTO ====================

// ProfileInfo 个人资料信息（不可编辑）
type ProfileInfo struct {
	FullName         string  `json:"full_name"`          // 学生姓名
	DateOfBirth      *string `json:"date_of_birth"`      // 出生日期 YYYY-MM
	Email            string  `json:"email"`              // 登录邮箱
	Gender           string  `json:"gender"`             // 性别 male/female
	JobSeekingStatus string  `json:"job_seeking_status"` // 求职状态
	PhoneNumber      string  `json:"phone_number"`       // 手机号
	AvatarURL        string  `json:"avatar_url"`         // 头像URL
}

// EducationInfo 教育经历信息（不可编辑）
type EducationInfo struct {
	Degree     string `json:"degree"`      // 学位层次
	SchoolName string `json:"school_name"` // 学校名称
	Major      string `json:"major"`       // 专业
	MajorRank  string `json:"major_rank"`  // 专业排名
	StartDate  string `json:"start_date"`  // 开始时间 YYYY-MM
	EndDate    string `json:"end_date"`    // 结束时间 YYYY-MM
}

// ResumeDraftResponse 简历草稿完整响应
type ResumeDraftResponse struct {
	Profile       ProfileInfo              `json:"profile"`        // 个人资料（不可编辑）
	Education     []EducationInfo          `json:"education"`      // 教育经历列表（不可编辑）
	SkillsSummary string                   `json:"skills_summary"` // 技能掌握（可编辑）
	WorkExp       []WorkExperienceResponse `json:"work_experiences"`
	Projects      []ProjectResponse        `json:"projects"`
	Organizations []OrganizationResponse   `json:"organizations"`
	Competitions  []CompetitionResponse    `json:"competitions"`
}
