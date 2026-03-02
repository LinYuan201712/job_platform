package repository

import (
	"fmt"
	"job-platform-go2/internal/model"

	"gorm.io/gorm"
)

type EventRepository struct {
	DB *gorm.DB
}

// GetEventsByMonth 按年月查询活动
func (r *EventRepository) GetEventsByMonth(year int, month int) ([]model.Event, error) {
	var events []model.Event
	// 构建该月的开始和结束时间
	startDate := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
	// 计算下个月的第一天
	var endDate string
	if month == 12 {
		endDate = fmt.Sprintf("%d-01-01 00:00:00", year+1)
	} else {
		endDate = fmt.Sprintf("%d-%02d-01 00:00:00", year, month+1)
	}

	err := r.DB.Where("event_start_time >= ? AND event_start_time < ?", startDate, endDate).
		Order("event_start_time ASC").
		Find(&events).Error
	return events, err
}

// GetUpcomingEvents 获取近期活动（按开始时间降序）
func (r *EventRepository) GetUpcomingEvents(limit int) ([]model.Event, error) {
	var events []model.Event
	err := r.DB.Order("event_start_time DESC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

// GetEventList 分页查询活动列表（支持关键词搜索）
func (r *EventRepository) GetEventList(page, pageSize int, keyword string) ([]model.Event, int64, error) {
	var events []model.Event
	var total int64

	query := r.DB.Model(&model.Event{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("event_title LIKE ? OR event_summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("event_start_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&events).Error

	return events, total, err
}

// GetEventByID 根据ID获取活动详情
func (r *EventRepository) GetEventByID(eventID int64) (*model.Event, error) {
	var event model.Event
	err := r.DB.Where("id = ?", eventID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}
