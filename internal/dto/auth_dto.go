package dto

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required,min=3"`
	VerificationCode string `json:"verification_code" binding:"required"`
	Name             string `json:"name" binding:"required"`
	// 修改为 string，且只允许 student 或 hr
	Role string `json:"role" binding:"required,oneof=student hr"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserInfo 响应中的用户信息部分
type UserInfo struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Role   string `json:"role"`   // "student" 或 "hr"
	Status string `json:"status"` // "active" 或 "pending"
}

// LoginResponseData 登录响应的 data 字段
type LoginResponseData struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"user_info"`
}
