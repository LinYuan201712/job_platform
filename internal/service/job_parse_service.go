package service

import (
	"job-platform-go2/internal/dto"
)

type JobParseService struct {
	// 可以在这里添加大模型 API 客户端
	// 例如: AIClient *some_ai_sdk.Client
}

// ParseJobFromImage 从图片解析岗位信息 (暂时返回模拟数据)
func (s *JobParseService) ParseJobFromImage(imagePath string) (*dto.ParseJobResponse, error) {
	// TODO: 接入真实的多模态大模型 API
	// 目前返回模拟数据
	return s.getMockParseResponse(), nil
}

// ParseJobFromText 从文本解析岗位信息 (暂时返回模拟数据)
func (s *JobParseService) ParseJobFromText(text string) (*dto.ParseJobResponse, error) {
	// TODO: 接入真实的大模型 API
	// 目前返回模拟数据
	return s.getMockParseResponse(), nil
}

// getMockParseResponse 返回模拟的解析结果
func (s *JobParseService) getMockParseResponse() *dto.ParseJobResponse {
	// 使用指针以支持 null 值
	headcount := 5
	minSalary := 15
	maxSalary := 25
	requiredDegree := 0
	deadline := "2026-05-01"
	requiredStartDate := "2026-06-01"

	return &dto.ParseJobResponse{
		JobDetails: dto.ParsedJobDetails{
			Title:             "后端开发工程师",
			Description:       "负责后端系统开发与维护",
			Department:        "技术部",
			Headcount:         &headcount,
			Deadline:          &deadline,
			TechRequirements:  "熟悉 Go 语言,了解微服务架构",
			BonusPoints:       "有大厂实习经验者优先",
			MinSalary:         &minSalary,
			MaxSalary:         &maxSalary,
			ProvinceID:        19, // 广东
			ProvinceName:      "广东省",
			CityID:            3, // 广州
			CityName:          "广州市",
			AddressDetail:     "天河区",
			WorkNature:        2, // 2=全职
			Type:              1, // 假设1=后端开发
			RequiredDegree:    &requiredDegree,
			RequiredStartDate: &requiredStartDate,
		},
	}
}
