package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service *service.AuthService
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, c)
		return
	}

	code, msg, err := ctrl.Service.Register(req)
	if err != nil {
		utils.Result(code, msg, nil, c)
		return
	}
	utils.Result(code, msg, nil, c)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(400, "请输入邮箱和密码", nil, c)
		return
	}

	data, err := ctrl.Service.Login(req)
	if err != nil {
		utils.Result(400, err.Error(), nil, c)
		return
	}
	utils.Result(200, "登录成功", data, c)
}

// ChangePassword 修改密码
// PUT /auth/change-password
func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	// 1. 绑定参数
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, c)
		return
	}

	// 2. 从JWT获取userID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未授权", nil, c)
		return
	}

	// 3. 调用Service
	if err := ctrl.Service.ChangePassword(userID.(int), req.OldPassword, req.NewPassword); err != nil {
		utils.Result(400, err.Error(), nil, c)
		return
	}

	// 4. 返回
	utils.Result(200, "密码修改成功", nil, c)
}
