package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HRTalentpoolController struct {
	Service *service.HRTalentpoolService
}

// GetTalentPoolJobList 获取人才库岗位列表
// GET /api/hr/jobs
func (c *HRTalentpoolController) GetTalentPoolJobList(ctx *gin.Context) {
	// 1. 绑定query参数
	var req dto.TalentPoolJobListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, ctx)
		return
	}

	// 2. 从JWT获取companyID
	companyID, exists := ctx.Get("companyID")
	if !exists {
		utils.Result(401, "未授权", nil, ctx)
		return
	}

	// 3. 调用Service
	resp, err := c.Service.GetTalentPoolJobList(req, companyID.(int))
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 4. 返回
	utils.Result(200, "成功", resp, ctx)
}

// GetCandidateListByJob 获取岗位候选人列表
// GET /api/hr/talentpool/job/list/:job_id
func (c *HRTalentpoolController) GetCandidateListByJob(ctx *gin.Context) {
	// 1. 绑定参数
	var req dto.CandidateListRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, ctx)
		return
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, ctx)
		return
	}

	// 2. 从JWT获取companyID
	companyID, exists := ctx.Get("companyID")
	if !exists {
		utils.Result(401, "未授权", nil, ctx)
		return
	}

	// 3. 调用Service
	resp, err := c.Service.GetCandidateListByJob(req, companyID.(int))
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 4. 返回
	utils.Result(200, "成功", resp, ctx)
}

// GetApplicationDetail 获取候选人简历详情
// GET /api/hr/applications/:id
func (c *HRTalentpoolController) GetApplicationDetail(ctx *gin.Context) {
	// 1. 绑定path参数
	appIDStr := ctx.Param("id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		utils.Result(400, "参数格式错误", nil, ctx)
		return
	}

	// 2. 从JWT获取companyID
	companyID, exists := ctx.Get("companyID")
	if !exists {
		utils.Result(401, "未授权", nil, ctx)
		return
	}

	// 3. 调用Service
	resp, err := c.Service.GetApplicationDetail(appID, companyID.(int))
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 4. 返回
	utils.Result(200, "成功", resp, ctx)
}

// UpdateApplicationStatus 更新人才状态
// PUT /api/hr/applications/:id/status
func (c *HRTalentpoolController) UpdateApplicationStatus(ctx *gin.Context) {
	// 1. 绑定path参数
	appIDStr := ctx.Param("id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		utils.Result(400, "参数格式错误", nil, ctx)
		return
	}

	// 2. 绑定JSON body
	var req dto.UpdateApplicationStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, ctx)
		return
	}

	// 3. 从JWT获取companyID
	companyID, exists := ctx.Get("companyID")
	if !exists {
		utils.Result(401, "未授权", nil, ctx)
		return
	}

	// 4. 调用Service
	resp, err := c.Service.UpdateApplicationStatus(appID, req.Status, companyID.(int))
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 5. 返回
	utils.Result(200, "成功", resp, ctx)
}

// GetStudentResumePreview HR端获取学生简历预览
// GET /api/hr/resume/:studentUserId
func (c *HRTalentpoolController) GetStudentResumePreview(ctx *gin.Context) {
	// 1. 绑定参数
	studentUserIDStr := ctx.Param("studentUserId")
	studentUserID, err := strconv.Atoi(studentUserIDStr)
	if err != nil {
		utils.Result(400, "参数格式错误", nil, ctx)
		return
	}

	// 2. 调用Service
	resp, err := c.Service.GetStudentResumePreview(studentUserID)
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 3. 返回
	utils.Result(200, "成功", resp, ctx)
}
