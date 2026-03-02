package dto

import "time"

// ==================== 求职中心相关 DTO ====================

// JobListRequest 职位列表查询请求
type JobListRequest struct {
	Province      string `form:"province"`       // 省份
	City          string `form:"city"`           // 城市
	Title         string `form:"title"`          // 职位名称
	CompanyName   string `form:"company_name"`   // 公司名称
	MinSalary     string `form:"min_salary"`     // 最低薪资
	MaxSalary     string `form:"max_salary"`     // 最高薪资
	WorkNature    string `form:"work_nature"`    // 工作性质: "实习招聘" / "校园招聘"
	Type          string `form:"type"`           // 职能类别
	CompanyNature string `form:"company_nature"` // 公司性质: "事业单位招聘" / "企业招聘"
	Page          int    `form:"page"`           // 页码 default 1
	PageSize      int    `form:"page_size"`      // 每页数量 default 10
}

// JobItemResponse 职位列表中的单个职位
type JobItemResponse struct {
	JobID         int    `json:"job_id"`
	Title         string `json:"title"`
	CompanyName   string `json:"company_name"`
	SalaryRange   string `json:"salary_range"`
	Address       string `json:"address"`
	WorkNature    string `json:"work_nature"`
	Department    string `json:"department"`
	Headcount     int    `json:"headcount"`
	IsFavorited   bool   `json:"is_favorited"`
	LogoURL       string `json:"logo_url"`
	ApplicationID *int   `json:"application_id"` // 投递记录ID,未投递为null
}

// JobListResponse 职位列表响应
type JobListResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Jobs     []JobItemResponse `json:"jobs"`
}

// CompanyLinkInfo 企业链接信息
type CompanyLinkInfo struct {
	LinkName string `json:"link_name"`
	LinkURL  string `json:"link_url"`
}

// CompanyInfoNested 职位详情中的公司信息(嵌套对象)
type CompanyInfoNested struct {
	CompanyID          int               `json:"company_id"`
	CompanyLogoURL     string            `json:"company_logo_url"`
	CompanyName        string            `json:"company_name"`
	CompanyIndustry    string            `json:"company_industry"`
	CompanyNature      string            `json:"company_nature"`
	CompanyScale       string            `json:"company_scale"`
	ContactPersonName  string            `json:"contact_person_name"`
	ContactPersonPhone string            `json:"contact_person_phone"`
	CompanyLinks       []CompanyLinkInfo `json:"company_links"`
}

// JobDetailFullResponse 职位详情完整响应
type JobDetailFullResponse struct {
	JobID               int               `json:"job_id"`
	Title               string            `json:"title"`
	SalaryRange         string            `json:"salary_range"`
	Address             string            `json:"address"`
	WorkNature          string            `json:"work_nature"`
	RequiredDegree      string            `json:"required_degree"`
	RequiredStartDate   string            `json:"required_start_date"` // YYYY-MM-DD
	Times               int               `json:"times"`               // 浏览次数(view_count)
	Type                string            `json:"type"`                // 职能类别名称
	RequiredSkills      []string          `json:"required_skills"`     // 技能标签数组
	Headcount           int               `json:"headcount"`
	PostedAt            string            `json:"posted_at"`             // YYYY-MM-DD
	PositionDescription string            `json:"position_description"`  // 职位描述
	PositionRequirement string            `json:"position_requirements"` // 职位要求
	BonusPoints         []string          `json:"bonus_points"`          // 加分项数组
	AddressDetail       string            `json:"address_detail"`        // 详细地址
	CompanyInfo         CompanyInfoNested `json:"company_info"`          // 公司信息嵌套对象
	IsFavorited         bool              `json:"is_favorited"`
}

// FavoriteListResponse 收藏列表响应
type FavoriteListResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Jobs     []JobItemResponse `json:"jobs"`
}

// FavoriteStatusResponse 收藏状态响应
type FavoriteStatusResponse struct {
	JobID       int  `json:"job_id"`
	IsFavorited bool `json:"is_favorited"`
}

// FavoriteSearchRequest 收藏列表搜索请求
type FavoriteSearchRequest struct {
	Province      string `form:"province"`
	City          string `form:"city"`
	Title         string `form:"title"`
	CompanyName   string `form:"company_name"`
	MinSalary     string `form:"min_salary"`
	MaxSalary     string `form:"max_salary"`
	WorkNature    string `form:"work_nature"`
	Type          string `form:"type"`
	CompanyNature string `form:"company_nature"`
	Page          int    `form:"page"`
	PageSize      int    `form:"page_size"`
}

// ApplyJobRequest 岗位投递请求
type ApplyJobRequest struct {
	JobID    int   `json:"job_id" binding:"required"`
	ResumeID int64 `json:"resume_id" binding:"required"`
}

// ApplicationResponse 投递响应
type ApplicationResponse struct {
	ApplicationID int       `json:"application_id"`
	JobID         int       `json:"job_id"`
	StudentUserID int       `json:"student_user_id"`
	ResumeID      int64     `json:"resume_id"`
	Status        int       `json:"status"`
	SubmittedAt   time.Time `json:"submitted_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CompanyStatistics 企业统计信息
type CompanyStatistics struct {
	ActivePositions   int    `json:"active_positions"`    // 在招岗位数
	ResumeProcessRate string `json:"resume_process_rate"` // 简历处理率
	RecentActivity    string `json:"recent_activity"`     // 最近活跃时间
}

// CompanyJobItem 企业岗位列表项
type CompanyJobItem struct {
	JobID    int    `json:"job_id"`
	JobTitle string `json:"job_title"`
	PostedAt string `json:"posted_at"` // YYYY-MM-DD
}

// CompanyDetailResponse 企业详情响应
type CompanyDetailResponse struct {
	CompanyID          int               `json:"company_id"`
	CompanyName        string            `json:"company_name"`
	CompanyIndustry    string            `json:"company_industry"`
	CompanyNature      string            `json:"company_nature"`
	CompanyScale       string            `json:"company_scale"`
	CompanyAddress     string            `json:"company_address"`
	CompanyLogoURL     string            `json:"company_logo_url"`
	ContactPersonName  string            `json:"contact_person_name"`
	ContactPersonPhone string            `json:"contact_person_phone"`
	Description        string            `json:"description"`
	Statistics         CompanyStatistics `json:"statistics"`
	CompanyLinks       []CompanyLinkInfo `json:"company_links"`
	Jobs               []CompanyJobItem  `json:"jobs"`
}
