package service

import (
	"fmt"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/repository"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StudentDashboardService struct {
	StudentRepo *repository.StudentRepository
	EventRepo   *repository.EventRepository
	JobRepo     *repository.JobRepository
	DB          *gorm.DB
}

// GetUserName 获取当前用户姓名
func (s *StudentDashboardService) GetUserName(userID int) (*dto.UserNameResponse, error) {
	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取学生信息失败: %w", err)
	}

	return &dto.UserNameResponse{
		FullName: student.FullName,
	}, nil
}

// GetCalendar 获取求职日历数据（按日期聚合）
func (s *StudentDashboardService) GetCalendar(month string, userID int) (*dto.CalendarResponse, error) {
	// 解析月份参数，默认为当前月
	var year, monthNum int
	var err error

	if month == "" {
		now := time.Now()
		year = now.Year()
		monthNum = int(now.Month())
	} else {
		// 解析 YYYY-MM 格式
		parts := strings.Split(month, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("无效的月份格式，应为 YYYY-MM")
		}
		year, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("无效的年份")
		}
		monthNum, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("无效的月份")
		}
	}

	// 获取该月份的所有活动
	events, err := s.EventRepo.GetEventsByMonth(year, monthNum)
	if err != nil {
		return nil, fmt.Errorf("查询活动失败: %w", err)
	}
	fmt.Printf("[Debug] GetCalendar: year=%d, month=%d, found %d events\n", year, monthNum, len(events))

	// 按日期聚合活动
	dailyEventsMap := make(map[string][]dto.EventSummaryItem)
	for _, event := range events {
		dateStr := event.EventStartTime.Format("2006-01-02")
		timeStr := event.EventStartTime.Format("15:04")

		summary := fmt.Sprintf("%s - %s", timeStr, event.EventTitle)
		// 截取摘要（如果太长）
		if len(summary) > 30 {
			summary = summary[:27] + "..."
		}

		dailyEventsMap[dateStr] = append(dailyEventsMap[dateStr], dto.EventSummaryItem{
			EventID: fmt.Sprintf("evt-uuid-%03d", event.ID), // 格式调整为 evt-uuid-00X
			Summary: summary,
		})
	}

	// 转换为数组格式
	var dailyEvents []dto.DailyEvent
	for date, summaries := range dailyEventsMap {
		dailyEvents = append(dailyEvents, dto.DailyEvent{
			EventDate:      date,
			EventSummaries: summaries,
		})
	}

	// 构建响应
	displayMonth := fmt.Sprintf("%d年 %02d月", year, monthNum)
	currentDay := time.Now().Day()

	return &dto.CalendarResponse{
		DisplayMonth: displayMonth,
		CurrentDay:   currentDay,
		DailyEvents:  dailyEvents,
	}, nil
}

// GetUpcomingEvents 获取近期求职活动
func (s *StudentDashboardService) GetUpcomingEvents(limit int) (*dto.UpcomingEventsResponse, error) {
	if limit <= 0 {
		limit = 5
	}

	events, err := s.EventRepo.GetUpcomingEvents(limit)
	if err != nil {
		return nil, fmt.Errorf("查询近期活动失败: %w", err)
	}

	var eventItems []dto.UpcomingEventItem
	for _, event := range events {
		location := "暂无地点"
		if event.EventLocation != nil && *event.EventLocation != "" {
			location = *event.EventLocation
		}

		eventItems = append(eventItems, dto.UpcomingEventItem{
			EventID:        fmt.Sprintf("uuid-%03d", event.ID), // 格式调整为 uuid-00X
			EventTitle:     event.EventTitle,
			EventStartTime: event.EventStartTime.Format("2006-01-02 15:04"),
			EventLocation:  location,
		})
	}

	return &dto.UpcomingEventsResponse{
		Events: eventItems,
	}, nil
}

// GetRankedJobs 获取岗位热度排行榜
func (s *StudentDashboardService) GetRankedJobs(limit int) (*dto.RankedJobsResponse, error) {
	if limit <= 0 {
		limit = 5
	}

	jobs, err := s.JobRepo.GetRankedJobs(limit)
	if err != nil {
		return nil, fmt.Errorf("查询热度排行失败: %w", err)
	}

	var rankedJobs []dto.RankedJobItem
	for idx, job := range jobs {
		// 构建地点字符串
		location := job.Job.AddressDetail

		rankedJobs = append(rankedJobs, dto.RankedJobItem{
			Rank:        idx + 1,
			JobID:       job.Job.ID,
			JobTitle:    job.Job.Title,
			CompanyName: job.CompanyName,
			Location:    location,
		})
	}

	return &dto.RankedJobsResponse{
		RankedJobs: rankedJobs,
	}, nil
}

// GetRecentJobs 获取近期招聘信息列表
func (s *StudentDashboardService) GetRecentJobs(jobTypeFilter string, limit int) (*dto.RecentJobsResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	jobs, err := s.JobRepo.GetRecentJobs(jobTypeFilter, limit)
	if err != nil {
		return nil, fmt.Errorf("查询近期岗位失败: %w", err)
	}

	var jobPostings []dto.RecentJobItem
	for _, job := range jobs {
		jobPostings = append(jobPostings, dto.RecentJobItem{
			JobID:    job.ID,
			JobTitle: job.Title,
			PostedAt: job.CreatedAt.Format(time.RFC3339),
		})
	}

	return &dto.RecentJobsResponse{
		JobPostings: jobPostings,
	}, nil
}
