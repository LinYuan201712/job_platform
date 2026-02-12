package service

import (
	"errors"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/pkg/utils"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func (s *AuthService) Register(req dto.RegisterRequest) (int, string, error) {
	// 1. 验证码校验 (示例逻辑)
	if req.VerificationCode != "123456" {
		return 400, "验证码错误", errors.New("vcode error")
	}

	processedEmail := strings.TrimSpace(req.Email)
	var roleCode int
	var successCode int
	var successMsg string
	var status int

	// 2. 映射字符串角色到内部 code，并处理邮箱逻辑
	if req.Role == "student" {
		roleCode = model.UserRoleStudent
		status = 1
		successCode, successMsg = 201, "学生账户注册成功"
		if !strings.Contains(processedEmail, "@") {
			processedEmail = processedEmail + "@mail2.sysu.edu.cn"
		}
	} else {
		roleCode = model.UserRoleHr
		status = 0
		successCode, successMsg = 202, "企业账户注册成功，请等待管理员审核"
		if !strings.Contains(processedEmail, "@") {
			return 400, "企业注册必须提供完整邮箱", errors.New("email format error")
		}
	}

	// 3. 检查邮箱唯一性
	var count int64
	s.DB.Model(&model.User{}).Where("email = ?", processedEmail).Count(&count)
	if count > 0 {
		return 400, "邮箱已被注册", errors.New("duplicate email")
	}

	// 4. 密码加密
	hashedPwd, _ := utils.HashPassword(req.Password)

	// 5. 开启事务：创建用户及关联档案
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		user := model.User{
			Email:        processedEmail,
			PasswordHash: hashedPwd,
			Role:         roleCode,
			Status:       status, // 实际项目中 HR 可是 0 (PENDING)
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 根据角色创建 Student 或 Company 记录
		if roleCode == model.UserRoleStudent {
			student := model.Student{
				UserID:    user.ID,
				FullName:  req.Name,
				StudentID: "TEMP_" + strconv.Itoa(user.ID),
			}
			return tx.Create(&student).Error
		} else {
			company := model.Company{
				UserID:      user.ID,
				CompanyName: req.Name,
			}
			return tx.Create(&company).Error
		}
		return nil
	})
	if err != nil {
		return 500, "系统错误", err
	}
	return successCode, successMsg, nil
}

// Login 用户登录
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponseData, error) {
	// 1. 查找用户
	var user model.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在或密码错误")
	}

	// 2. 验证密码
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("用户不存在或密码错误") // 模糊报错是安全规范
	}

	// 3. 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, err
	}

	roleStr := "student"
	if user.Role == model.UserRoleHr {
		roleStr = "hr"
	}

	statusStr := "active"
	if user.Status == 0 {
		statusStr = "pending"
	}

	return &dto.LoginResponseData{
		Token: token,
		UserInfo: dto.UserInfo{
			ID:     user.ID,
			Email:  user.Email,
			Role:   roleStr,
			Status: statusStr,
		},
	}, nil
}
