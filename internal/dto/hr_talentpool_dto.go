package dto

// TalentPoolJobListRequest 人才库岗位列表请求
type TalentPoolJobListRequest struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"pageSize"`
	TitleKeyword string `form:"titleKeyword"`
	WorkNature   string `form:"workNature"` // "实习" or "校招"
	Status       string `form:"status"`     // "draft", "pending", "approved", "closed"
}

// TalentPoolJobResponse 人才库岗位响应
type TalentPoolJobResponse struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	JobID       int    `json:"job_id"`
	WorkNature  string `json:"work_nature"`
	UpdatedAt   string `json:"updated_at"`
	ReceivedNum int64  `json:"received_num"`
	NoReviewNum int64  `json:"no_review_num"`
}

// TalentPoolJobListResponse 人才库岗位列表响应(带分页)
type TalentPoolJobListResponse struct {
	JobList    []TalentPoolJobResponse `json:"job_list"`
	Pagination PaginationResponse      `json:"pagination"`
}

// PaginationResponse 分页信息
type PaginationResponse struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
}

// CandidateListRequest 候选人列表请求
type CandidateListRequest struct {
	JobID       int    `uri:"job_id" binding:"required"`
	NameKeyword string `form:"name_keyword"`
	Status      *int   `form:"status"` // 使用指针以区分0和未传值
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

// CandidateItem 候选人信息
type CandidateItem struct {
	Grade         string `json:"grade"`
	Degree        string `json:"degree"`
	ApplicationID int    `json:"application_id"`
	CandidateName string `json:"candidate_name"`
	AvatarURL     string `json:"avatar_url"`
	UserID        int    `json:"user_id"`
	ResumeStatus  string `json:"resume_status"`
}

// CandidateListResponse 候选人列表响应(带分页)
type CandidateListResponse struct {
	CandidateList []CandidateItem    `json:"candidate_list"`
	Pagination    PaginationResponse `json:"pagination"`
}

// UpdateApplicationStatusRequest 更新人才状态请求
type UpdateApplicationStatusRequest struct {
	Status string `json:"status" binding:"required"` // 文本: "已投递", "候选人", "面试邀请", "通过", "拒绝"
}

// UpdateStatusResponse 状态更新响应
type UpdateStatusResponse struct {
	Status        string `json:"status"`
	ApplicationID int    `json:"application_id"`
	StatusCode    int    `json:"status_code"`
}

// StudentResumePreviewResponse 学生简历预览响应
type StudentResumePreviewResponse struct {
	AvatarURL        string                   `json:"avatar_url"`
	BasicInfo        BasicInfoDTO             `json:"basic_info"`
	PrimaryEducation []EducationExperienceDTO `json:"primary_education"`
	ExpectedJob      ExpectedJobDTO           `json:"expected_job"`
	PersonalTags     []PersonalTagDTO         `json:"personal_tags"`
}

// BasicInfoDTO 基本信息
type BasicInfoDTO struct {
	FullName         string `json:"full_name"`
	Gender           string `json:"gender"`
	Age              int    `json:"age"`
	DegreeLevel      string `json:"degree_level"`
	JobSeekingStatus string `json:"job_seeking_status"`
}

// EducationExperienceDTO 教育经历
type EducationExperienceDTO struct {
	SchoolName  string `json:"school_name"`
	DegreeLevel string `json:"degree_level"`
	Major       string `json:"major"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	MajorRank   string `json:"major_rank"`
}

// ExpectedJobDTO 期望岗位
type ExpectedJobDTO struct {
	ExpectedPosition  string `json:"expected_position"`
	ExpectedMinSalary int    `json:"expected_min_salary"`
	ExpectedMaxSalary int    `json:"expected_max_salary"`
}

// PersonalTagDTO 个人标签
type PersonalTagDTO struct {
	TagID int    `json:"tag_id"`
	Name  string `json:"name"`
}
