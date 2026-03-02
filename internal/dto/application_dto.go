package dto

// ==================== 投递详情 ====================

type JobInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CompanyInfo struct {
	Name string `json:"name"`
}

type ApplicationDetailResponse struct {
	ID           int         `json:"id"`
	Status       string      `json:"status"`        // 状态名称: 已投递, 候选人, 面试邀请, 通过, 拒绝
	StatusDetail string      `json:"status_detail"` // 状态详细文案
	SubmittedAt  string      `json:"submitted_at"`  // 格式: "2025-10-29 09:15:00"
	UpdatedAt    string      `json:"updated_at"`    // 格式: "2025-10-31 14:30:00"
	Job          JobInfo     `json:"job"`
	Company      CompanyInfo `json:"company"`
}

// ==================== 已投递岗位列表 ====================

type DeliveryListRequest struct {
	JobTitle    string `form:"job_title"`
	CompanyName string `form:"company_name"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

type DeliveryJobItem struct {
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
	ApplicationID *int   `json:"application_id"` // 投递记录ID
}

type DeliveryListResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Jobs     []DeliveryJobItem `json:"jobs"`
}
