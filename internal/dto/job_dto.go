package dto

import "time"

// ==================== 岗位相关 DTO ====================

// CreateJobRequest 创建岗位请求
type CreateJobRequest struct {
	Title             string `json:"title" binding:"required"`
	Department        string `json:"department" binding:"required"`
	ProvinceID        int    `json:"province_id" binding:"required"`
	CityID            int    `json:"city_id" binding:"required"`
	AddressDetail     string `json:"address_detail" binding:"required"`
	WorkNature        int    `json:"work_nature" binding:"required,min=1"`   // 1=实习, 2=全职
	Type              int    `json:"type" binding:"required,min=1"`          // 职能类别ID
	Headcount         int    `json:"headcount" binding:"required,min=1"`     // 招聘人数
	MinSalary         int    `json:"min_salary" binding:"required,min=0"`    // 最低薪资(k)
	MaxSalary         int    `json:"max_salary" binding:"required,min=0"`    // 最高薪资(k)
	RequiredDegree    int    `json:"required_degree" binding:"min=0,max=2"`  // 0=本科+, 1=硕士+, 2=博士+ (允许0)
	RequiredStartDate string `json:"required_start_date" binding:"required"` // YYYY-MM-DD
	Deadline          string `json:"deadline" binding:"required"`            // YYYY-MM-DD
	Description       string `json:"description" binding:"required"`         // 职位描述
	TechRequirements  string `json:"tech_requirements" binding:"required"`   // 职位要求
	BonusPoints       string `json:"bonus_points"`                           // 加分项(可选)
	Status            int    `json:"status" binding:"required,min=1"`        // 1=草稿, 10=提交申请
	Tags              []int  `json:"tags" binding:"required"`                // 标签ID数组
}

// UpdateJobRequest 更新岗位请求 (字段与CreateJobRequest相同)
type UpdateJobRequest struct {
	Title             string `json:"title" binding:"required"`
	Department        string `json:"department" binding:"required"`
	ProvinceID        int    `json:"province_id" binding:"required"`
	CityID            int    `json:"city_id" binding:"required"`
	AddressDetail     string `json:"address_detail" binding:"required"`
	WorkNature        int    `json:"work_nature" binding:"required,min=1"`
	Type              int    `json:"type" binding:"required,min=1"`
	Headcount         int    `json:"headcount" binding:"required,min=1"`
	MinSalary         int    `json:"min_salary" binding:"required,min=0"`
	MaxSalary         int    `json:"max_salary" binding:"required,min=0"`
	RequiredDegree    int    `json:"required_degree" binding:"min=0,max=2"`
	RequiredStartDate string `json:"required_start_date" binding:"required"`
	Deadline          string `json:"deadline" binding:"required"`
	Description       string `json:"description" binding:"required"`
	TechRequirements  string `json:"tech_requirements" binding:"required"`
	BonusPoints       string `json:"bonus_points"`
	Status            int    `json:"status" binding:"required,min=1"`
	Tags              []int  `json:"tags" binding:"required"`
}

// JobDetailResponse 岗位详情响应
type JobDetailResponse struct {
	JobID             int       `json:"job_id"`
	Title             string    `json:"title"`
	Status            int       `json:"status"`
	Description       string    `json:"description"`
	TechRequirements  *string   `json:"tech_requirements"` // 可为null
	BonusPoints       string    `json:"bonus_points"`
	MinSalary         int       `json:"min_salary"`
	MaxSalary         int       `json:"max_salary"`
	ProvinceID        int       `json:"province_id"`
	CityID            int       `json:"city_id"`
	AddressDetail     string    `json:"address_detail"`
	WorkNature        int       `json:"work_nature"`
	Department        string    `json:"department"`
	Headcount         int       `json:"headcount"`
	Type              int       `json:"type"`
	RequiredDegree    int       `json:"required_degree"`
	RequiredStartDate string    `json:"required_start_date"` // 格式: YYYY-MM-DD
	Deadline          string    `json:"deadline"`            // 格式: YYYY-MM-DD
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Tags              []TagInfo `json:"tags"` // 岗位标签列表
}

// TagInfo 标签信息
type TagInfo struct {
	TagID   int    `json:"tag_id"`
	TagName string `json:"tag_name"`
}

// CreateJobResponse 创建岗位响应
type CreateJobResponse struct {
	NewJob struct {
		JobID  int    `json:"job_id"`
		Title  string `json:"title"`
		Status string `json:"status"` // "draft" 或 "pending"
	} `json:"new_job"`
}

// UpdateJobResponse 更新岗位响应
type UpdateJobResponse struct {
	UpdatedJob struct {
		JobID  int    `json:"job_id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	} `json:"updated_job"`
}

// OfflineJobResponse 下线岗位响应
type OfflineJobResponse struct {
	JobID  int    `json:"job_id"`
	Status string `json:"status"` // "closed"
}
