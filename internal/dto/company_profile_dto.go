package dto

// UpdateCompanyProfileRequest 企业信息更新请求
type UpdateCompanyProfileRequest struct {
	Description        string            `json:"description" binding:"required"`
	CompanyAddress     string            `json:"company_address" binding:"required"`
	Nature             string            `json:"nature" binding:"required"`
	Industry           string            `json:"industry" binding:"required"`
	CompanyScale       string            `json:"company_scale" binding:"required"`
	ContactPersonName  string            `json:"contact_person_name" binding:"required"`
	ContactPersonPhone string            `json:"contact_person_phone" binding:"required"`
	ExternalLinks      []ExternalLinkDTO `json:"external_links"`
}

// ExternalLinkDTO 外部链接
type ExternalLinkDTO struct {
	ID       int    `json:"id,omitempty"` // 用于GET响应
	LinkName string `json:"link_name" binding:"required"`
	LinkURL  string `json:"link_url" binding:"required"`
}

// CompanyProfileResponse 企业信息响应
type CompanyProfileResponse struct {
	OpenJobsCount      int               `json:"open_jobs_count"`
	ResumeProcessRate  float64           `json:"resume_process_rate"`
	LastLoginAt        string            `json:"last_login_at"`
	CompanyName        string            `json:"company_name"`
	Description        string            `json:"description"`
	LogoURL            string            `json:"logo_url"`
	Nature             string            `json:"nature"`
	Industry           string            `json:"industry"`
	CompanyScale       string            `json:"company_scale"`
	ContactPersonName  string            `json:"contact_person_name"`
	ContactPersonPhone string            `json:"contact_person_phone"`
	CompanyAddress     string            `json:"company_address"`
	ExternalLinks      []ExternalLinkDTO `json:"external_links"`
}

// CompanyOptionsResponse 企业选项响应
type CompanyOptionsResponse struct {
	Industries []string `json:"industries"`
	Natures    []string `json:"natures"`
	Scales     []string `json:"scales"`
}
