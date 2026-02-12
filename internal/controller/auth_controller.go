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
