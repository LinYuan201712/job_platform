package service

import (
	"errors"
	"fmt"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyService struct {
	CompanyRepo *repository.CompanyRepository
	DB          *gorm.DB
}

// ... existing methods ...

// UploadCompanyLogo 上传企业Logo
func (s *CompanyService) UploadCompanyLogo(companyID int, fileHeader *multipart.FileHeader) (string, error) {
	// 1. 文件大小校验（最大5MB）
	maxSize := int64(5 * 1024 * 1024)
	if fileHeader.Size > maxSize {
		return "", errors.New("文件大小不能超过5MB")
	}

	// 2. 文件格式校验
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		return "", errors.New("仅支持 JPEG、PNG、WebP 图片格式")
	}

	// 3. 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadDir := "./uploads/company_logos"
	filePath := filepath.Join(uploadDir, filename)

	// 4. 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 5. 打开上传的文件
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	// 6. 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// 7. 复制文件内容
	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 8. 生成访问URL
	logoURL := fmt.Sprintf("/uploads/company_logos/%s", filename)

	// 9. 更新数据库
	err = s.DB.Model(&model.Company{}).
		Where("company_id = ?", companyID).
		Update("logo_url", logoURL).Error
	if err != nil {
		return "", err
	}

	return logoURL, nil
}

// GetCompanyDetail 获取企业详情
func (s *CompanyService) GetCompanyDetail(companyID int) (*dto.CompanyDetailResponse, error) {
	// 1. 获取企业基本信息(含字典表名称)
	companyWithDict, err := s.CompanyRepo.GetCompanyWithDictionaries(companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("企业不存在")
		}
		return nil, err
	}
	company := &companyWithDict.Company

	// 2. 获取企业链接
	companyLinks, err := s.CompanyRepo.GetCompanyLinks(companyID)
	if err != nil {
		return nil, err
	}

	// 3. 获取企业岗位列表
	jobs, err := s.CompanyRepo.GetCompanyJobsList(companyID)
	if err != nil {
		return nil, err
	}

	// 4. 计算统计信息
	activePositions, err := s.CompanyRepo.GetCompanyJobsCount(companyID, 20) // 20=approved
	if err != nil {
		return nil, err
	}

	// TODO: 计算简历处理率和最近活跃时间
	resumeProcessRate := "70%" // 示例值,需要实际计算
	recentActivity := calculateRecentActivity(company.CreatedAt)

	// 5. 构建响应
	links := make([]dto.CompanyLinkInfo, 0, len(companyLinks))
	for _, link := range companyLinks {
		links = append(links, dto.CompanyLinkInfo{
			LinkName: link.LinkName,
			LinkURL:  link.LinkURL,
		})
	}

	jobItems := make([]dto.CompanyJobItem, 0, len(jobs))
	for _, job := range jobs {
		jobItems = append(jobItems, dto.CompanyJobItem{
			JobID:    job.ID,
			JobTitle: job.Title,
			PostedAt: job.CreatedAt.Format("2006-01-02"),
		})
	}

	return &dto.CompanyDetailResponse{
		CompanyID:          company.CompanyID,
		CompanyName:        company.CompanyName,
		CompanyIndustry:    companyWithDict.IndustryName,
		CompanyNature:      companyWithDict.NatureName,
		CompanyScale:       companyWithDict.CompanyScaleName,
		CompanyAddress:     company.CompanyAddress,
		CompanyLogoURL:     company.LogoUrl,
		ContactPersonName:  company.ContactPersonName,
		ContactPersonPhone: company.ContactPersonPhone,
		Description:        company.Description,
		Statistics: dto.CompanyStatistics{
			ActivePositions:   int(activePositions),
			ResumeProcessRate: resumeProcessRate,
			RecentActivity:    recentActivity,
		},
		CompanyLinks: links,
		Jobs:         jobItems,
	}, nil
}

// calculateRecentActivity 计算最近活跃时间
func calculateRecentActivity(createdAt time.Time) string {
	now := time.Now()
	duration := now.Sub(createdAt)

	days := int(duration.Hours() / 24)
	if days == 0 {
		return "今天"
	} else if days == 1 {
		return "1天前"
	} else if days < 7 {
		return "本周"
	} else if days < 30 {
		return "本月"
	} else {
		return "较早"
	}
}

// GetCompanyProfile 获取企业信息(含统计数据)
func (s *CompanyService) GetCompanyProfile(companyID int) (*dto.CompanyProfileResponse, error) {
	// 1. 查询企业信息(含字典名称)
	companyData, err := s.CompanyRepo.GetCompanyWithDictionaries(companyID)
	if err != nil {
		return nil, err
	}

	// 2. 在招岗位数 (status=20 approved)
	openJobsCount, err := s.CompanyRepo.GetCompanyJobsCount(companyData.CompanyID, 20)
	if err != nil {
		return nil, err
	}

	// 3. 简历处理率
	processRate, err := s.CompanyRepo.GetCompanyResumeProcessRate(companyData.CompanyID)
	if err != nil {
		return nil, err
	}

	// 4. 最近登录时间
	lastLoginAt, err := s.CompanyRepo.GetUserLastLoginTime(companyData.UserID)
	if err != nil {
		return nil, err
	}

	// 5. 外部链接
	links, err := s.CompanyRepo.GetCompanyLinks(companyData.CompanyID)
	if err != nil {
		return nil, err
	}

	externalLinks := make([]dto.ExternalLinkDTO, 0, len(links))
	for _, link := range links {
		externalLinks = append(externalLinks, dto.ExternalLinkDTO{
			ID:       link.ID,
			LinkName: link.LinkName,
			LinkURL:  link.LinkURL,
		})
	}

	return &dto.CompanyProfileResponse{
		OpenJobsCount:      int(openJobsCount),
		ResumeProcessRate:  processRate,
		LastLoginAt:        lastLoginAt,
		CompanyName:        companyData.CompanyName,
		Description:        companyData.Description,
		LogoURL:            companyData.LogoUrl,
		Nature:             companyData.NatureName,
		Industry:           companyData.IndustryName,
		CompanyScale:       companyData.CompanyScaleName,
		ContactPersonName:  companyData.ContactPersonName,
		ContactPersonPhone: companyData.ContactPersonPhone,
		CompanyAddress:     companyData.CompanyAddress,
		ExternalLinks:      externalLinks,
	}, nil
}

// UpdateCompanyProfile 更新企业信息(事务化)
func (s *CompanyService) UpdateCompanyProfile(companyID int, req dto.UpdateCompanyProfileRequest, dictRepo *repository.DictionaryRepository) error {
	// 1. 字典转换: 文本 → ID
	industryID, err := dictRepo.GetIndustryIDByName(req.Industry)
	if err != nil {
		return errors.New("行业信息不存在: " + req.Industry)
	}

	natureID, err := dictRepo.GetNatureIDByName(req.Nature)
	if err != nil {
		return errors.New("企业性质不存在: " + req.Nature)
	}

	scaleID, err := dictRepo.GetScaleIDByName(req.CompanyScale)
	if err != nil {
		return errors.New("企业规模不存在: " + req.CompanyScale)
	}

	// 2. 事务化更新
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 2.1 更新企业基本信息
		updates := map[string]interface{}{
			"description":          req.Description,
			"company_address":      req.CompanyAddress,
			"nature_id":            natureID,
			"industry_id":          industryID,
			"company_scale_id":     scaleID,
			"contact_person_name":  req.ContactPersonName,
			"contact_person_phone": req.ContactPersonPhone,
		}

		companyRepo := &repository.CompanyRepository{DB: tx}
		if err := companyRepo.UpdateCompanyProfile(companyID, updates); err != nil {
			return err
		}

		// 2.2 全量替换外部链接
		newLinks := make([]model.CompanyLink, 0, len(req.ExternalLinks))
		for _, linkDTO := range req.ExternalLinks {
			newLinks = append(newLinks, model.CompanyLink{
				LinkName: linkDTO.LinkName,
				LinkURL:  linkDTO.LinkURL,
			})
		}

		if err := companyRepo.ReplaceCompanyLinks(tx, companyID, newLinks); err != nil {
			return err
		}

		return nil
	})
}

// GetCompanyOptions 获取企业信息选项(行业/性质/规模)
func (s *CompanyService) GetCompanyOptions(dictRepo *repository.DictionaryRepository) (*dto.CompanyOptionsResponse, error) {
	// 1. 查询所有行业
	industries, err := dictRepo.GetAllIndustries()
	if err != nil {
		return nil, err
	}

	// 2. 查询所有性质
	natures, err := dictRepo.GetAllNatures()
	if err != nil {
		return nil, err
	}

	// 3. 查询所有规模
	scales, err := dictRepo.GetAllScales()
	if err != nil {
		return nil, err
	}

	return &dto.CompanyOptionsResponse{
		Industries: industries,
		Natures:    natures,
		Scales:     scales,
	}, nil
}
