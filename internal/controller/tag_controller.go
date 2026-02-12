package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	TagService *service.TagService
}

// GetAllTags GET /api/tags - 获取所有标签(按分类分组)
func (ctrl *TagController) GetAllTags(c *gin.Context) {
	resp, err := ctrl.TagService.GetAllTagsGrouped()
	if err != nil {
		utils.Result(500, err.Error(), nil, c)
		return
	}

	utils.Result(200, "success", resp, c)
}

// CreateTag POST /api/tags - 创建自定义标签
func (ctrl *TagController) CreateTag(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, c)
		return
	}

	// 从 JWT 中获取用户ID
	userID, _ := c.Get("userID")

	resp, err := ctrl.TagService.CreateTag(req, userID.(int))
	if err != nil {
		utils.Result(400, err.Error(), nil, c)
		return
	}

	utils.Result(200, "创建成功", resp, c)
}
