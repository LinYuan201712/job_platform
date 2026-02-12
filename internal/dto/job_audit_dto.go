package dto

// ==================== 岗位审核相关 DTO ====================

// JobAuditResponse 岗位审核详情响应
type JobAuditResponse struct {
	JobID           int    `json:"job_id"`
	AuditStatus     int    `json:"audit_status"`     // 20=通过, 30=拒绝
	Remark          string `json:"remark"`           // 审核备注/拒绝原因
	OperatorContact string `json:"operator_contact"` // 审核人联系方式
	AuditTime       string `json:"audit_time"`       // 审核时间 (格式: YYYY-MM-DD HH:mm:ss)
}
