package dto

// ==================== 岗位解析相关 DTO (大模型) ====================

// ParseJobRequest 岗位解析请求
type ParseJobRequest struct {
	InputType string `form:"input_type" binding:"required,oneof=image text"` // "image" 或 "text"
	Text      string `form:"text"`                                           // 当 input_type=text 时必填
	// Image 字段由 gin 的 multipart/form-data 处理,不在此定义
}

// ParseJobResponse 岗位解析响应
type ParseJobResponse struct {
	JobDetails ParsedJobDetails `json:"job_details"`
}

// ParsedJobDetails 解析出的岗位详情
type ParsedJobDetails struct {
	Title             string  `json:"title"`               // 岗位名称
	Description       string  `json:"description"`         // 岗位描述
	Department        string  `json:"department"`          // 所属部门
	Headcount         *int    `json:"headcount"`           // 招聘人数(可为null)
	Deadline          *string `json:"deadline"`            // 招聘截止日期(可为null)
	TechRequirements  string  `json:"tech_requirements"`   // 岗位要求
	BonusPoints       string  `json:"bonus_points"`        // 岗位加分项
	MinSalary         *int    `json:"min_salary"`          // 最低薪资(单位k,可为null)
	MaxSalary         *int    `json:"max_salary"`          // 最高薪资(单位k,可为null)
	ProvinceID        int     `json:"province_id"`         // 省份ID
	ProvinceName      string  `json:"province_name"`       // 省份名称
	CityID            int     `json:"city_id"`             // 城市ID
	CityName          string  `json:"city_name"`           // 城市名称
	AddressDetail     string  `json:"address_detail"`      // 详细地址
	WorkNature        int     `json:"work_nature"`         // 岗位性质 1=internship, 2=full-time
	Type              int     `json:"type"`                // 岗位类别id
	RequiredDegree    *int    `json:"required_degree"`     // 学历要求(可为null)
	RequiredStartDate *string `json:"required_start_date"` // 岗位要求入职时间(可为null)
}
