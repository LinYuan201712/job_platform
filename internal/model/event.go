package model

import "time"

// Event 对应数据库表: events
type Event struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AdminUserID    int        `gorm:"column:admin_user_id;not null" json:"admin_user_id"`
	EventTitle     string     `gorm:"column:event_title;not null" json:"event_title"`
	EventSummary   *string    `gorm:"column:event_summary" json:"event_summary"`
	EventStartTime time.Time  `gorm:"column:event_start_time;not null" json:"event_start_time"`
	EventEndTime   *time.Time `gorm:"column:event_end_time" json:"event_end_time"`
	EventLocation  *string    `gorm:"column:event_location" json:"event_location"`
	TargetAudience *string    `gorm:"column:target_audience" json:"target_audience"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Event) TableName() string {
	return "events"
}
