package model

// UserRole
const (
	UserRoleStudent = 1
	UserRoleHr      = 2
	UserRoleAdmin   = 0 // 预留管理员角色
)

// UserStatus
const (
	UserStatusActive   = 1
	UserStatusInactive = 0
)

// JobStatus
const (
	JobStatusDraft    = 0
	JobStatusPending  = 1
	JobStatusApproved = 2
	JobStatusRejected = 3
	JobStatusClosed   = 4
)

// ApplicationStatus (对应数据库 applications.status)
const (
	AppStatusSubmitted = 10 // 已投递
	AppStatusCandidate = 20 // 候选人/初筛通过
	AppStatusInterview = 30 // 面试邀请
	AppStatusPassed    = 40 // 通过/录用
	AppStatusRejected  = 50 // 拒绝
)

// WorkNature
const (
	WorkNatureFullTime = 0
	WorkNatureIntern   = 1
)
