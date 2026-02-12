package model

// UserRole
const (
	UserRoleStudent = 0
	UserRoleHr      = 1
	UserRoleAdmin   = 2
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

// ApplicationStatus
const (
	AppStatusSubmitted = 0 // 已投递
	AppStatusCandidate = 1 // 候选人/初筛通过
	AppStatusInterview = 2 // 面试中
	AppStatusPassed    = 3 // 通过/录用
	AppStatusRejected  = 4 // 拒绝
)

// WorkNature
const (
	WorkNatureFullTime = 0
	WorkNatureIntern   = 1
)
