package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ResumeController 简历控制器
type ResumeController struct {
	ResumeService     *service.ResumeService
	ExperienceService *service.ExperienceService
}

// ==================== 简历草稿相关 ====================

// GetResumeDraft 获取简历草稿
func (c *ResumeController) GetResumeDraft(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	studentUserID := userID.(int)

	draft, err := c.ResumeService.GetResumeDraft(studentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取简历草稿失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    draft,
	})
}

// UpdateSkills 更新技能信息
func (c *ResumeController) UpdateSkills(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	var req dto.UpdateSkillsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ResumeService.UpdateSkills(studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新技能信息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// SetTemplate 设置简历模板
func (c *ResumeController) SetTemplate(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	var req dto.SetTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ResumeService.SetTemplate(studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "设置模板失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// ==================== 简历文件相关 ====================

// UploadResumeFile 上传简历文件
func (c *ResumeController) UploadResumeFile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未上传文件",
		})
		return
	}

	resp, err := c.ResumeService.UploadResumeFile(studentUserID, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// GetResumeFiles 获取简历文件列表
func (c *ResumeController) GetResumeFiles(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	files, err := c.ResumeService.GetResumeFiles(studentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取简历文件列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    files,
	})
}

// DeleteResumeFile 删除简历文件
func (c *ResumeController) DeleteResumeFile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的文件ID",
		})
		return
	}

	resp, err := c.ResumeService.DeleteResumeFile(id, studentUserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// ==================== 工作经历相关 ====================

// CreateWorkExperience 新增工作经历
func (c *ResumeController) CreateWorkExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	var req dto.CreateWorkExperienceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.CreateWorkExperience(studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// UpdateWorkExperience 修改工作经历
func (c *ResumeController) UpdateWorkExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req dto.UpdateWorkExperienceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.UpdateWorkExperience(id, studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// DeleteWorkExperience 删除工作经历
func (c *ResumeController) DeleteWorkExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := c.ExperienceService.DeleteWorkExperience(id, studentUserID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ==================== 项目经历相关 ====================

// CreateProjectExperience 新增项目经历
func (c *ResumeController) CreateProjectExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	var req dto.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.CreateProjectExperience(studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// UpdateProjectExperience 修改项目经历
func (c *ResumeController) UpdateProjectExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req dto.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.UpdateProjectExperience(id, studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// DeleteProjectExperience 删除项目经历
func (c *ResumeController) DeleteProjectExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := c.ExperienceService.DeleteProjectExperience(id, studentUserID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ==================== 组织经历相关 ====================

// CreateOrganizationExperience 新增组织经历
func (c *ResumeController) CreateOrganizationExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	var req dto.CreateOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.CreateOrganizationExperience(studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// UpdateOrganizationExperience 修改组织经历
func (c *ResumeController) UpdateOrganizationExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req dto.UpdateOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.UpdateOrganizationExperience(id, studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// DeleteOrganizationExperience 删除组织经历
func (c *ResumeController) DeleteOrganizationExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := c.ExperienceService.DeleteOrganizationExperience(id, studentUserID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ==================== 竞赛经历相关 ====================

// CreateCompetitionExperience 新增竞赛经历
func (c *ResumeController) CreateCompetitionExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	var req dto.CreateCompetitionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.CreateCompetitionExperience(studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// UpdateCompetitionExperience 修改竞赛经历
func (c *ResumeController) UpdateCompetitionExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req dto.UpdateCompetitionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.ExperienceService.UpdateCompetitionExperience(id, studentUserID, &req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// DeleteCompetitionExperience 删除竞赛经历
func (c *ResumeController) DeleteCompetitionExperience(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	studentUserID := userID.(int)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := c.ExperienceService.DeleteCompetitionExperience(id, studentUserID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
