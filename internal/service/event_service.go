package service

import (
	"fmt"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/repository"
)

type EventService struct {
	EventRepo *repository.EventRepository
}

// GetEventList 获取活动列表（分页+搜索）
func (s *EventService) GetEventList(req dto.EventListRequest) (*dto.EventListResponse, error) {
	// 默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	events, total, err := s.EventRepo.GetEventList(req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, fmt.Errorf("查询活动列表失败: %w", err)
	}

	var eventItems []dto.EventListItem
	for _, event := range events {
		location := ""
		if event.EventLocation != nil {
			location = *event.EventLocation
		}
		summary := ""
		if event.EventSummary != nil {
			summary = *event.EventSummary
		}

		eventItems = append(eventItems, dto.EventListItem{
			EventID:        int(event.ID),
			EventTitle:     event.EventTitle,
			EventStartTime: event.EventStartTime.Format("2006-01-02 15:04:05"),
			EventLocation:  location,
			EventSummary:   summary,
		})
	}

	return &dto.EventListResponse{
		Total:  total,
		Events: eventItems,
	}, nil
}

// GetEventDetail 获取活动详情
func (s *EventService) GetEventDetail(eventID int64) (*dto.EventDetailResponse, error) {
	event, err := s.EventRepo.GetEventByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("查询活动详情失败: %w", err)
	}

	location := ""
	if event.EventLocation != nil {
		location = *event.EventLocation
	}
	targetAudience := ""
	if event.TargetAudience != nil {
		targetAudience = *event.TargetAudience
	}
	summary := ""
	if event.EventSummary != nil {
		summary = *event.EventSummary
	}

	return &dto.EventDetailResponse{
		EventID:             fmt.Sprintf("uuid-%03d", event.ID),
		EventTitle:          event.EventTitle,
		EventStartTime:      event.EventStartTime.Format("2006-01-02 15:04:05"),
		EventLocation:       location,
		EventTargetAudience: targetAudience,
		EventSummary:        summary,
	}, nil
}
