package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	JobService      *service.JobService
	JobAuditService *service.JobAuditService
	JobParseService *service.JobParseService
}

// CreateJob POST /api/hr/jobs - 创建岗位
func (ctrl *JobController) CreateJob(c *gin.Context) {
	var req dto.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 输出详细错误信息以便调试
		utils.Result(400, "参数格式错误: "+err.Error(), nil, c)
		return
	}

	// 从 JWT 中获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未找到用户信息", nil, c)
		return
	}
	// TODO: 从用户信息中获取 company_id,这里暂时硬编码
	companyID := 1 // 需要根据实际业务逻辑获取HR用户的company_id

	resp, err := ctrl.JobService.CreateJob(req, userID.(int), companyID)
	if err != nil {
		utils.Result(400, err.Error(), nil, c)
		return
	}

	utils.Result(200, "提交成功", resp, c)
}

// GetJobDetail GET /api/hr/jobs/:job_id - 获取岗位详情
func (ctrl *JobController) GetJobDetail(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		utils.Result(400, "岗位ID格式错误", nil, c)
		return
	}

	resp, err := ctrl.JobService.GetJobDetail(jobID)
	if err != nil {
		utils.Result(404, err.Error(), nil, c)
		return
	}

	utils.Result(200, "success", resp, c)
}

// UpdateJob PUT /api/hr/jobs/:job_id - 更新岗位
func (ctrl *JobController) UpdateJob(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		utils.Result(400, "岗位ID格式错误", nil, c)
		return
	}

	var req dto.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, c)
		return
	}

	userID, _ := c.Get("userID")
	resp, err := ctrl.JobService.UpdateJob(jobID, req, userID.(int))
	if err != nil {
		utils.Result(400, err.Error(), nil, c)
		return
	}

	utils.Result(200, "更新成功", resp, c)
}

// DeleteJobDraft DELETE /api/hr/jobs/:job_id - 删除岗位草稿
func (ctrl *JobController) DeleteJobDraft(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		utils.Result(400, "岗位ID格式错误", nil, c)
		return
	}

	userID, _ := c.Get("userID")
	if err := ctrl.JobService.DeleteJobDraft(jobID, userID.(int)); err != nil {
		utils.Result(400, err.Error(), nil, c)
		return
	}

	utils.Result(200, "删除成功", nil, c)
}

// OfflineJob PUT /api/hr/jobs/:job_id/status - 下线岗位
func (ctrl *JobController) OfflineJob(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		utils.Result(400, "岗位ID格式错误", nil, c)
		return
	}

	userID, _ := c.Get("userID")
	resp, err := ctrl.JobService.OfflineJob(jobID, userID.(int))
	if err != nil {
		utils.Result(400, err.Error(), nil, c)
		return
	}

	utils.Result(200, "下线成功", resp, c)
}

// GetJobAudit GET /api/hr/jobs/audit/:jobId - 获取岗位审核详情
func (ctrl *JobController) GetJobAudit(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("jobId"))
	if err != nil {
		utils.Result(400, "岗位ID格式错误", nil, c)
		return
	}

	resp, err := ctrl.JobAuditService.GetJobAudit(jobID)
	if err != nil {
		utils.Result(404, err.Error(), nil, c)
		return
	}

	utils.Result(200, "success", resp, c)
}

// ParseJob POST /api/hr/jobs/parse - 解析岗位信息(图片或文本)
func (ctrl *JobController) ParseJob(c *gin.Context) {
	inputType := c.PostForm("input_type")
	if inputType != "image" && inputType != "text" {
		utils.Result(400, "input_type 必须为 image 或 text", nil, c)
		return
	}

	var resp *dto.ParseJobResponse
	var err error

	if inputType == "image" {
		// 处理图片上传
		file, err := c.FormFile("image")
		if err != nil {
			utils.Result(400, "图片上传失败", nil, c)
			return
		}

		// TODO: 保存图片到临时目录,并使用真实路径调用AI服务
		// 暂时使用模拟数据
		_ = file // 避免未使用变量错误
		resp, err = ctrl.JobParseService.ParseJobFromImage("")
	} else {
		// 处理文本
		text := c.PostForm("text")
		if text == "" {
			utils.Result(400, "text 字段不能为空", nil, c)
			return
		}
		resp, err = ctrl.JobParseService.ParseJobFromText(text)
	}

	if err != nil {
		utils.Result(500, "解析失败", nil, c)
		return
	}

	utils.Result(200, "解析成功", resp, c)
}
