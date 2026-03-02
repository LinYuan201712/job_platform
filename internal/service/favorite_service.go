package service

import (
	"errors"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/repository"

	"gorm.io/gorm"
)

type FavoriteService struct {
	FavoriteRepo *repository.JobFavoriteRepository
	JobRepo      *repository.JobRepository
	DB           *gorm.DB
}

// AddFavorite 收藏职位
func (s *FavoriteService) AddFavorite(studentUserID, jobID int) error {
	// 验证job_id是否存在
	_, err := s.JobRepo.GetJobByID(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("职位不存在")
		}
		return err
	}

	// 添加收藏(幂等性由repository处理)
	return s.FavoriteRepo.AddFavorite(studentUserID, jobID)
}

// RemoveFavorite 取消收藏
func (s *FavoriteService) RemoveFavorite(studentUserID, jobID int) error {
	return s.FavoriteRepo.RemoveFavorite(studentUserID, jobID)
}

// GetFavoriteStatus 获取收藏状态
func (s *FavoriteService) GetFavoriteStatus(studentUserID, jobID int) (*dto.FavoriteStatusResponse, error) {
	isFavorited, err := s.FavoriteRepo.CheckFavoriteStatus(studentUserID, jobID)
	if err != nil {
		return nil, err
	}

	return &dto.FavoriteStatusResponse{
		JobID:       jobID,
		IsFavorited: isFavorited,
	}, nil
}

// GetFavoriteList 获取收藏列表
func (s *FavoriteService) GetFavoriteList(studentUserID int, page, pageSize int) (*dto.FavoriteListResponse, error) {
	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 获取收藏列表
	favorites, total, err := s.FavoriteRepo.GetFavoriteList(studentUserID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 获取每个收藏的职位详细信息
	jobItems := make([]dto.JobItemResponse, 0, len(favorites))
	for _, fav := range favorites {
		// 获取职位信息
		job, err := s.JobRepo.GetJobByID(fav.JobID)
		if err != nil {
			continue // 跳过不存在的职位
		}

		// 获取公司信息(需要查询)
		var companyName, logoURL string
		// TODO: 查询company表获取company_name和logo_url

		jobItem := dto.JobItemResponse{
			JobID:         job.ID,
			Title:         job.Title,
			CompanyName:   companyName,
			SalaryRange:   formatSalaryRange(job.MinSalary, job.MaxSalary),
			Address:       job.AddressDetail,
			WorkNature:    formatWorkNature(job.WorkNature),
			Department:    job.Department,
			Headcount:     job.Headcount,
			IsFavorited:   true, // 收藏列表中全部为true
			LogoURL:       logoURL,
			ApplicationID: nil, // TODO: 查询是否已投递
		}
		jobItems = append(jobItems, jobItem)
	}

	return &dto.FavoriteListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Jobs:     jobItems,
	}, nil
}

// SearchFavorites 在收藏列表中搜索
func (s *FavoriteService) SearchFavorites(req dto.FavoriteSearchRequest, studentUserID int) (*dto.FavoriteListResponse, error) {
	// TODO: 实现收藏列表的二次搜索
	// 这里可以复用GetFavoriteList,然后在应用层过滤
	return s.GetFavoriteList(studentUserID, req.Page, req.PageSize)
}

// ==================== 工具函数导入 ====================

// formatSalaryRange 格式化薪资范围
func formatSalaryRange(min, max int) string {
	if max == 0 || max >= 15000 {
		return formatInt(min) + "及以上"
	}
	return formatInt(min) + "-" + formatInt(max)
}

// formatWorkNature 工作性质转换
func formatWorkNature(code int) string {
	switch code {
	case 1:
		return "实习招聘"
	case 2:
		return "校园招聘"
	default:
		return "未知"
	}
}

// formatInt 整数转字符串
func formatInt(i int) string {
	return string(rune(i))
}
