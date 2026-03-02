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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// StudentProfileService 学生档案服务
type StudentProfileService struct {
	StudentRepo   *repository.StudentRepository
	EducationRepo *repository.EducationExperienceRepository
	TagRepo       *repository.TagRepository
	UserRepo      *repository.UserRepository
	DB            *gorm.DB
}

// ==================== 获取档案信息 ====================

// GetMyProfile 获取学生档案信息（聚合）
func (s *StudentProfileService) GetMyProfile(userID int) (*dto.GetProfileResponse, error) {
	// 1. 获取学生和用户信息
	student, user, err := s.StudentRepo.FindByUserIDWithUser(userID)
	if err != nil {
		return nil, fmt.Errorf("学生信息不存在: %w", err)
	}

	// 2. 获取教育经历（按学历等级排序）
	educations, err := s.EducationRepo.FindByStudentUserIDOrderByDegree(userID)
	if err != nil {
		return nil, fmt.Errorf("获取教育经历失败: %w", err)
	}

	// 3. 获取个人标签
	tags, err := s.TagRepo.FindTagsByStudentUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取标签失败: %w", err)
	}

	// 4. 组装响应VO
	response := &dto.GetProfileResponse{
		AvatarURL: student.AvatarUrl,
		BasicInfo: dto.BasicInfo{
			FullName:         student.FullName,
			Gender:           mapGenderToString(student.Gender),
			DateOfBirth:      formatDate(student.DateOfBirth),
			JobSeekingStatus: mapJobStatusToString(student.JobSeekingStatus),
			Email:            user.Email,
			PhoneNumber:      student.PhoneNumber,
			StudentID:        student.StudentID,
		},
		ExpectedJob: dto.ExpectedJob{
			ExpectedPosition:  student.ExpectedPosition,
			ExpectedMinSalary: student.ExpectedMinSalary,
			ExpectedMaxSalary: student.ExpectedMaxSalary,
		},
	}

	// 转换教育经历
	for _, edu := range educations {
		response.PrimaryEducation = append(response.PrimaryEducation, dto.EducationDTO{
			ID:          &edu.ID,
			SchoolName:  edu.SchoolName,
			DegreeLevel: mapDegreeLevelToString(edu.DegreeLevel),
			Major:       edu.Major,
			StartDate:   formatDate(&edu.StartDate),
			EndDate:     formatDate(&edu.EndDate),
			MajorRank:   edu.MajorRank,
		})
	}

	// 转换标签
	for _, tag := range tags {
		response.PersonalTags = append(response.PersonalTags, dto.PersonalTag{
			TagID: tag.ID,
			Name:  tag.Name,
		})
	}

	return response, nil
}

// ==================== 更新档案信息 ====================

// UpdateMyBaseProfile 更新学生档案基本信息（含教育经历和标签同步）
func (s *StudentProfileService) UpdateMyBaseProfile(userID int, req *dto.UpdateProfileRequest) error {
	// 使用事务保证数据一致性
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新学生基本信息
		student, err := s.StudentRepo.FindByUserID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果学生记录不存在，创建新记录
				student = &model.Student{
					UserID: userID,
				}
			} else {
				return fmt.Errorf("查询学生信息失败: %w", err)
			}
		}

		// 更新基本信息字段
		student.FullName = req.BasicInfo.FullName
		student.PhoneNumber = req.BasicInfo.PhoneNumber
		student.StudentID = req.BasicInfo.StudentID
		student.Gender = mapStringToGender(req.BasicInfo.Gender)
		student.JobSeekingStatus = mapStringToJobStatus(req.BasicInfo.JobSeekingStatus)

		// 解析出生日期
		if req.BasicInfo.DateOfBirth != "" {
			dob, err := parseDate(req.BasicInfo.DateOfBirth)
			if err == nil {
				student.DateOfBirth = &dob
			}
		}

		// 更新期望岗位信息
		student.ExpectedPosition = req.ExpectedJob.ExpectedPosition
		student.ExpectedMinSalary = req.ExpectedJob.ExpectedMinSalary
		student.ExpectedMaxSalary = req.ExpectedJob.ExpectedMaxSalary

		// 保存或创建学生记录
		if err == nil {
			// 学生记录已存在，更新
			if err := s.StudentRepo.Update(student); err != nil {
				return fmt.Errorf("更新学生信息失败: %w", err)
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// 学生记录不存在，创建
			if err := s.StudentRepo.Create(student); err != nil {
				return fmt.Errorf("创建学生信息失败: %w", err)
			}
		} else {
			return fmt.Errorf("查询学生信息失败: %w", err)
		}

		// 2. 更新用户邮箱
		if req.BasicInfo.Email != "" {
			user, err := s.UserRepo.FindByID(userID)
			if err != nil {
				return fmt.Errorf("查询用户信息失败: %w", err)
			}
			user.Email = req.BasicInfo.Email
			if err := s.UserRepo.Update(user); err != nil {
				return fmt.Errorf("更新用户邮箱失败: %w", err)
			}
		}

		// 3. 同步教育经历
		if err := s.syncEducationExperiences(tx, userID, req.PrimaryEducation); err != nil {
			return fmt.Errorf("同步教育经历失败: %w", err)
		}

		// 4. 同步标签
		if err := s.syncPersonalTags(tx, userID, req.PersonalTagIDs); err != nil {
			return fmt.Errorf("同步标签失败: %w", err)
		}

		return nil
	})
}

// syncEducationExperiences 同步教育经历（增删改）
func (s *StudentProfileService) syncEducationExperiences(tx *gorm.DB, userID int, educations []dto.EducationDTO) error {
	// 获取现有教育经历
	existingEdus, err := s.EducationRepo.FindByStudentUserID(userID)
	if err != nil {
		return err
	}

	// 构建现有ID集合
	existingIDs := make(map[int64]bool)
	for _, edu := range existingEdus {
		existingIDs[edu.ID] = true
	}

	// 构建传入ID集合
	incomingIDs := make(map[int64]bool)
	for _, eduDTO := range educations {
		if eduDTO.ID != nil && *eduDTO.ID > 0 {
			incomingIDs[*eduDTO.ID] = true
		}
	}

	// 删除不在传入列表中的记录
	for id := range existingIDs {
		if !incomingIDs[id] {
			if err := tx.Delete(&model.EducationExperience{}, id).Error; err != nil {
				return err
			}
		}
	}

	// 更新或创建教育经历
	for _, eduDTO := range educations {
		edu := model.EducationExperience{
			StudentUserID: userID,
			SchoolName:    eduDTO.SchoolName,
			DegreeLevel:   mapStringToDegreeLevel(eduDTO.DegreeLevel),
			Major:         eduDTO.Major,
			MajorRank:     eduDTO.MajorRank,
		}

		// 解析日期
		if eduDTO.StartDate != "" {
			startDate, err := parseDate(eduDTO.StartDate)
			if err == nil {
				edu.StartDate = startDate
			}
		}
		if eduDTO.EndDate != "" {
			endDate, err := parseDate(eduDTO.EndDate)
			if err == nil {
				edu.EndDate = endDate
			}
		}

		if eduDTO.ID != nil && *eduDTO.ID > 0 && existingIDs[*eduDTO.ID] {
			// 更新现有记录
			edu.ID = *eduDTO.ID
			if err := tx.Save(&edu).Error; err != nil {
				return err
			}
		} else {
			// 创建新记录
			if err := tx.Create(&edu).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// syncPersonalTags 同步个人标签
func (s *StudentProfileService) syncPersonalTags(tx *gorm.DB, userID int, tagIDs []int) error {
	// 删除所有现有标签关联
	if err := tx.Exec("DELETE FROM students_tags WHERE user_id = ?", userID).Error; err != nil {
		return err
	}

	// 批量插入新标签关联
	if len(tagIDs) > 0 {
		type StudentTag struct {
			UserID int
			TagID  int
		}
		var records []StudentTag
		for _, tagID := range tagIDs {
			records = append(records, StudentTag{UserID: userID, TagID: tagID})
		}
		if err := tx.Table("students_tags").Create(&records).Error; err != nil {
			return err
		}
	}

	return nil
}

// ==================== 头像上传 ====================

// UploadAvatar 上传学生头像
func (s *StudentProfileService) UploadAvatar(userID int, fileHeader *multipart.FileHeader) (string, error) {
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
	uploadDir := "./uploads/avatars"
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
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	// 9. 更新数据库
	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新学生记录
			student = &model.Student{
				UserID:    userID,
				AvatarUrl: avatarURL,
			}
			if err := s.StudentRepo.Create(student); err != nil {
				return "", fmt.Errorf("创建学生记录失败: %w", err)
			}
		} else {
			return "", fmt.Errorf("查询学生信息失败: %w", err)
		}
	} else {
		student.AvatarUrl = avatarURL
		if err := s.StudentRepo.Update(student); err != nil {
			return "", fmt.Errorf("更新头像失败: %w", err)
		}
	}

	return avatarURL, nil
}

// ==================== 欢迎信息 ====================

// GetWelcomeInfo 获取欢迎信息
func (s *StudentProfileService) GetWelcomeInfo(userID int) (*dto.WelcomeInfoResponse, error) {
	// 1. 获取学生和用户信息
	student, user, err := s.StudentRepo.FindByUserIDWithUser(userID)
	if err != nil {
		return nil, fmt.Errorf("学生信息不存在: %w", err)
	}

	// 2. 获取第一条教育经历（作为学校名）
	educations, err := s.EducationRepo.FindByStudentUserIDOrderByDegree(userID)
	schoolName := ""
	if err == nil && len(educations) > 0 {
		schoolName = educations[0].SchoolName
	}

	// 3. 获取个人标签
	tags, err := s.TagRepo.FindTagsByStudentUserID(userID)
	if err != nil {
		tags = []model.Tag{} // 标签为空也不影响
	}

	// 4. 组装响应
	response := &dto.WelcomeInfoResponse{
		FullName:    student.FullName,
		SchoolName:  schoolName,
		PhoneNumber: student.PhoneNumber,
		Email:       user.Email,
		StudentID:   student.StudentID,
		LastLoginAt: formatDateTime(user.LastLoginAt),
	}

	// 转换标签
	for _, tag := range tags {
		response.PersonalTags = append(response.PersonalTags, dto.PersonalTag{
			TagID: tag.ID,
			Name:  tag.Name,
		})
	}

	return response, nil
}

// ==================== 修改密码 ====================

// ChangePassword 修改密码
func (s *StudentProfileService) ChangePassword(userID int, req *dto.ChangePasswordRequest) error {
	// 1. 获取用户当前密码hash
	currentHash, err := s.UserRepo.GetPasswordHash(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 2. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 3. 生成新密码hash
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 4. 更新密码
	if err := s.UserRepo.UpdatePassword(userID, string(newHash)); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// ==================== 简历预览 ====================

// GetResumePreview 获取简历预览（含年龄计算）
func (s *StudentProfileService) GetResumePreview(userID int) (*dto.ResumePreviewResponse, error) {
	// 1. 获取学生和用户信息
	student, _, err := s.StudentRepo.FindByUserIDWithUser(userID)
	if err != nil {
		return nil, fmt.Errorf("学生信息不存在: %w", err)
	}

	// 2. 获取教育经历
	educations, err := s.EducationRepo.FindByStudentUserIDOrderByDegree(userID)
	if err != nil {
		return nil, fmt.Errorf("获取教育经历失败: %w", err)
	}

	// 3. 获取标签
	tags, err := s.TagRepo.FindTagsByStudentUserID(userID)
	if err != nil {
		tags = []model.Tag{}
	}

	// 4. 计算年龄
	age := 0
	if student.DateOfBirth != nil {
		age = time.Now().Year() - student.DateOfBirth.Year()
	}

	// 5. 获取最高学历
	highestDegree := "本科"
	if len(educations) > 0 {
		highestDegree = mapDegreeLevelToString(educations[len(educations)-1].DegreeLevel)
	}

	// 6. 组装响应
	response := &dto.ResumePreviewResponse{
		AvatarURL: student.AvatarUrl,
		BasicInfo: dto.BasicInfoWithAge{
			FullName:         student.FullName,
			Gender:           mapGenderToString(student.Gender),
			Age:              age,
			DegreeLevel:      highestDegree,
			JobSeekingStatus: mapJobStatusToString(student.JobSeekingStatus),
		},
		ExpectedJob: dto.ExpectedJob{
			ExpectedPosition:  student.ExpectedPosition,
			ExpectedMinSalary: student.ExpectedMinSalary,
			ExpectedMaxSalary: student.ExpectedMaxSalary,
		},
	}

	// 转换教育经历
	for _, edu := range educations {
		response.PrimaryEducation = append(response.PrimaryEducation, dto.EducationDTO{
			ID:          &edu.ID,
			SchoolName:  edu.SchoolName,
			DegreeLevel: mapDegreeLevelToString(edu.DegreeLevel),
			Major:       edu.Major,
			StartDate:   formatDate(&edu.StartDate),
			EndDate:     formatDate(&edu.EndDate),
			MajorRank:   edu.MajorRank,
		})
	}

	// 转换标签
	for _, tag := range tags {
		response.PersonalTags = append(response.PersonalTags, dto.PersonalTag{
			TagID: tag.ID,
			Name:  tag.Name,
		})
	}

	return response, nil
}

// ==================== 辅助函数 ====================

// mapGenderToString 性别整数转字符串
func mapGenderToString(gender int) string {
	if gender == 0 {
		return "男"
	}
	return "女"
}

// mapStringToGender 性别字符串转整数
func mapStringToGender(gender string) int {
	if strings.Contains(gender, "男") {
		return 0
	}
	return 1
}

// mapJobStatusToString 求职状态整数转字符串
func mapJobStatusToString(status int) string {
	switch status {
	case 0:
		return "在校-暂不考虑"
	case 1:
		return "在校-寻求实习"
	case 2:
		return "应届-寻求实习"
	case 3:
		return "应届-寻求校招"
	default:
		return ""
	}
}

// mapStringToJobStatus 求职状态字符串转整数
func mapStringToJobStatus(status string) int {
	if strings.Contains(status, "暂不考虑") {
		return 0
	}
	if strings.Contains(status, "在校") && strings.Contains(status, "实习") {
		return 1
	}
	if strings.Contains(status, "应届") && strings.Contains(status, "实习") {
		return 2
	}
	if strings.Contains(status, "校招") || (strings.Contains(status, "应届") && !strings.Contains(status, "实习")) {
		return 3
	}
	return 0
}

// mapDegreeLevelToString 学历等级整数转字符串
func mapDegreeLevelToString(level int) string {
	switch level {
	case 0:
		return "本科"
	case 1:
		return "硕士"
	case 2:
		return "博士"
	default:
		return "本科"
	}
}

// mapStringToDegreeLevel 学历等级字符串转整数
func mapStringToDegreeLevel(level string) int {
	if strings.Contains(level, "硕") {
		return 1
	}
	if strings.Contains(level, "博") {
		return 2
	}
	return 0
}

// parseDate 解析日期字符串（支持 YYYY-MM-DD 和 YYYY-MM 格式）
func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("日期字符串为空")
	}

	// 尝试 YYYY-MM-DD 格式
	t, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}

	// 尝试 YYYY-MM 格式（补充为月初）
	t, err = time.Parse("2006-01", dateStr)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("无效的日期格式: %s", dateStr)
}

// formatDate 格式化日期为 YYYY-MM-DD
func formatDate(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// formatDateTime 格式化日期时间
func formatDateTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z")
}
