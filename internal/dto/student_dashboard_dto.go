package dto

// ==================== 用户名 ====================

type UserNameResponse struct {
	FullName string `json:"full_name"`
}

// ==================== 求职日历 ====================

type CalendarRequest struct {
	Month string `form:"month"` // YYYY-MM 格式
}

type EventSummaryItem struct {
	EventID string `json:"event_id"`
	Summary string `json:"summary"`
}

type DailyEvent struct {
	EventDate      string             `json:"event_date"` // YYYY-MM-DD
	EventSummaries []EventSummaryItem `json:"event_summaries"`
}

type CalendarResponse struct {
	DisplayMonth string       `json:"display_month"` // 如: "2025年 10月"
	CurrentDay   int          `json:"current_day"`
	DailyEvents  []DailyEvent `json:"daily_events"`
}

// ==================== 近期求职活动 ====================

type UpcomingEventsRequest struct {
	Limit int `form:"limit"`
}

type UpcomingEventItem struct {
	EventID        string `json:"event_id"`
	EventTitle     string `json:"event_title"`
	EventStartTime string `json:"event_start_time"` // 格式: "2025-10-29 14:00"
	EventLocation  string `json:"event_location"`
}

type UpcomingEventsResponse struct {
	Events []UpcomingEventItem `json:"events"`
}

// ==================== 岗位热度排行榜 ====================

type RankedJobsRequest struct {
	Limit int `form:"limit"`
}

type RankedJobItem struct {
	Rank        int    `json:"rank"`
	JobID       int    `json:"job_id"`
	JobTitle    string `json:"job_title"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
}

type RankedJobsResponse struct {
	RankedJobs []RankedJobItem `json:"ranked_jobs"`
}

// ==================== 近期招聘信息列表 ====================

type RecentJobsRequest struct {
	JobTypeFilter string `form:"job_type_filter"` // 企业招聘, 事业单位招聘, 实习招聘, 校园招聘
	Limit         int    `form:"limit"`
}

type RecentJobItem struct {
	JobID    int    `json:"job_id"`
	JobTitle string `json:"job_title"`
	PostedAt string `json:"posted_at"` // ISO 8601 格式
}

type RecentJobsResponse struct {
	JobPostings []RecentJobItem `json:"job_postings"`
}

// ==================== 招聘活动列表 ====================

type EventListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
}

type EventListItem struct {
	EventID        int    `json:"event_id"`
	EventTitle     string `json:"event_title"`
	EventStartTime string `json:"event_start_time"`
	EventLocation  string `json:"event_location"`
	EventSummary   string `json:"event_summary"`
}

type EventListResponse struct {
	Total  int64           `json:"total"`
	Events []EventListItem `json:"events"`
}

// ==================== 招聘活动详情 ====================

type EventDetailResponse struct {
	EventID             string `json:"event_id"`
	EventTitle          string `json:"event_title"`
	EventStartTime      string `json:"event_start_time"`
	EventLocation       string `json:"event_location"`
	EventTargetAudience string `json:"event_target_audience"`
	EventSummary        string `json:"event_summary"`
}
