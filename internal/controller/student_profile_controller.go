package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StudentProfileController 学生档案控制器
type StudentProfileController struct {
	Service *service.StudentProfileService
}

// ==================== 获取档案信息 ====================

// GetMyProfile 获取学生档案信息（聚合）
func (c *StudentProfileController) GetMyProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	profile, err := c.Service.GetMyProfile(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取档案信息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取档案信息成功",
		"data":    profile,
	})
}

// ==================== 更新档案信息 ====================

// UpdateMyBaseProfile 更新学生档案基本信息
func (c *StudentProfileController) UpdateMyBaseProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := c.Service.UpdateMyBaseProfile(userID.(int), &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "档案更新失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "档案更新成功",
	})
}

// ==================== 头像上传 ====================

// UploadAvatar 上传学生头像
func (c *StudentProfileController) UploadAvatar(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	file, err := ctx.FormFile("avatar_file")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "未上传文件",
		})
		return
	}

	avatarURL, err := c.Service.UploadAvatar(userID.(int), file)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "头像上传成功",
		"data": gin.H{
			"new_avatar_url": avatarURL,
		},
	})
}

// ==================== 欢迎信息 ====================

// GetWelcomeInfo 获取欢迎信息
func (c *StudentProfileController) GetWelcomeInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	welcomeInfo, err := c.Service.GetWelcomeInfo(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取欢迎信息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取欢迎信息成功",
		"data":    welcomeInfo,
	})
}

// ==================== 修改密码 ====================

// ChangePassword 修改密码
func (c *StudentProfileController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 从query参数获取密码
	oldPassword := ctx.Query("old_password")
	newPassword := ctx.Query("new_password")

	if oldPassword == "" || newPassword == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "旧密码和新密码不能为空",
		})
		return
	}

	req := &dto.ChangePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	if err := c.Service.ChangePassword(userID.(int), req); err != nil {
		// 检查是否是旧密码错误
		if err.Error() == "原密码错误" {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "原密码错误",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码修改失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码修改成功",
	})
}

// ==================== 简历预览 ====================

// GetResumePreview 获取简历预览（含年龄计算）
func (c *StudentProfileController) GetResumePreview(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	preview, err := c.Service.GetResumePreview(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取简历预览失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取简历预览成功",
		"data":    preview,
	})
}
