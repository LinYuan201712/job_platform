package dto

// ==================== 获取档案信息 ====================

// GetProfileResponse 获取学生档案的响应
type GetProfileResponse struct {
	AvatarURL        string         `json:"avatar_url"`
	BasicInfo        BasicInfo      `json:"basic_info"`
	PrimaryEducation []EducationDTO `json:"primary_education"`
	ExpectedJob      ExpectedJob    `json:"expected_job"`
	PersonalTags     []PersonalTag  `json:"personal_tags"`
}

// BasicInfo 基本信息
type BasicInfo struct {
	FullName         string `json:"full_name"`
	Gender           string `json:"gender"`             // "男" 或 "女"
	DateOfBirth      string `json:"date_of_birth"`      // YYYY-MM-DD
	JobSeekingStatus string `json:"job_seeking_status"` // 中文描述
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	StudentID        string `json:"student_id"`
}

// EducationDTO 教育经历
type EducationDTO struct {
	ID          *int64 `json:"id,omitempty"` // 用于更新时标识，新增时为nil
	SchoolName  string `json:"school_name"`
	DegreeLevel string `json:"degree_level"` // "本科"、"硕士"、"博士"
	Major       string `json:"major"`
	StartDate   string `json:"start_date"` // YYYY-MM-DD
	EndDate     string `json:"end_date"`   // YYYY-MM-DD
	MajorRank   string `json:"major_rank"`
}

// ExpectedJob 期望岗位信息
type ExpectedJob struct {
	ExpectedPosition  string `json:"expected_position"`
	ExpectedMinSalary int    `json:"expected_min_salary"`
	ExpectedMaxSalary int    `json:"expected_max_salary"`
}

// PersonalTag 个人标签
type PersonalTag struct {
	TagID int    `json:"tag_id"`
	Name  string `json:"name"`
}

// ==================== 更新档案信息 ====================

// UpdateProfileRequest 更新档案信息的请求
type UpdateProfileRequest struct {
	BasicInfo        BasicInfoUpdate `json:"basic_info"`
	PrimaryEducation []EducationDTO  `json:"primary_education"`
	ExpectedJob      ExpectedJob     `json:"expected_job"`
	PersonalTagIDs   []int           `json:"personal_tag_ids"`
}

// BasicInfoUpdate 更新基本信息（包含所有可编辑字段）
type BasicInfoUpdate struct {
	FullName         string `json:"full_name"`
	Gender           string `json:"gender"`
	DateOfBirth      string `json:"date_of_birth"`
	JobSeekingStatus string `json:"job_seeking_status"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	StudentID        string `json:"student_id"`
}

// ==================== 欢迎信息 ====================

// WelcomeInfoResponse 欢迎信息响应
type WelcomeInfoResponse struct {
	FullName     string        `json:"full_name"`
	SchoolName   string        `json:"school_name"`
	PhoneNumber  string        `json:"phone_number"`
	Email        string        `json:"email"`
	LastLoginAt  string        `json:"last_login_at"`
	StudentID    string        `json:"student_id"`
	PersonalTags []PersonalTag `json:"personal_tags"`
}

// ==================== 修改密码 ====================

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

// ==================== 简历预览 ====================

// ResumePreviewResponse 简历预览响应（含年龄计算）
type ResumePreviewResponse struct {
	AvatarURL        string           `json:"avatar_url"`
	BasicInfo        BasicInfoWithAge `json:"basic_info"`
	PrimaryEducation []EducationDTO   `json:"primary_education"`
	ExpectedJob      ExpectedJob      `json:"expected_job"`
	PersonalTags     []PersonalTag    `json:"personal_tags"`
}

// BasicInfoWithAge 基本信息（带年龄字段）
type BasicInfoWithAge struct {
	FullName         string `json:"full_name"`
	Gender           string `json:"gender"`
	Age              int    `json:"age"`          // 计算得出
	DegreeLevel      string `json:"degree_level"` // 最高学历
	JobSeekingStatus string `json:"job_seeking_status"`
}
