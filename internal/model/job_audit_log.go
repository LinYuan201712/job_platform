package model

import "time"

type JobAuditLog struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobID           int       `gorm:"column:job_id;not null" json:"job_id"`
	OperatorID      int64     `gorm:"column:operator_id" json:"operator_id"`
	OperatorContact string    `gorm:"column:operator_contact" json:"operator_contact"`
	AuditStatus     int       `gorm:"column:audit_status" json:"audit_status"`
	Remark          string    `gorm:"column:remark" json:"remark"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (JobAuditLog) TableName() string {
	return "job_audit_log"
}
