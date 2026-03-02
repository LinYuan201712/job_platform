package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JobCenterController struct {
	JobCenterService   *service.JobCenterService
	FavoriteService    *service.FavoriteService
	ApplicationService *service.ApplicationService
	CompanyService     *service.CompanyService
}

// ==================== 职位相关 ====================

// GetJobList 获取职位列表
func (ctrl *JobCenterController) GetJobList(c *gin.Context) {
	var req dto.JobListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Result(400, "参数错误", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	currentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.JobCenterService.GetJobList(req, currentUserID)
	if err != nil {
		utils.Result(500, "获取岗位列表失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取岗位列表成功", response, c)
}

// GetJobDetail 获取职位详情
func (ctrl *JobCenterController) GetJobDetail(c *gin.Context) {
	jobIDStr := c.Param("job_id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.Result(400, "无效的职位ID", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	currentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.JobCenterService.GetJobDetail(jobID, currentUserID)
	if err != nil {
		utils.Result(500, "获取岗位详情失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取岗位详情成功", response, c)
}

// ==================== 收藏相关 ====================

// AddFavorite 收藏职位
func (ctrl *JobCenterController) AddFavorite(c *gin.Context) {
	jobIDStr := c.Param("job_id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.Result(400, "无效的职位ID", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	err = ctrl.FavoriteService.AddFavorite(studentUserID, jobID)
	if err != nil {
		utils.Result(500, "收藏失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "收藏成功", nil, c)
}

// RemoveFavorite 取消收藏
func (ctrl *JobCenterController) RemoveFavorite(c *gin.Context) {
	jobIDStr := c.Param("job_id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.Result(400, "无效的职位ID", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	err = ctrl.FavoriteService.RemoveFavorite(studentUserID, jobID)
	if err != nil {
		utils.Result(500, "取消收藏失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "取消收藏成功", nil, c)
}

// GetFavoriteStatus 获取收藏状态
func (ctrl *JobCenterController) GetFavoriteStatus(c *gin.Context) {
	jobIDStr := c.Param("job_id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.Result(400, "无效的职位ID", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.FavoriteService.GetFavoriteStatus(studentUserID, jobID)
	if err != nil {
		utils.Result(500, "查询失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "查询成功", response, c)
}

// GetFavoriteList 获取收藏列表
func (ctrl *JobCenterController) GetFavoriteList(c *gin.Context) {
	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.FavoriteService.GetFavoriteList(studentUserID, page, pageSize)
	if err != nil {
		utils.Result(500, "获取收藏岗位失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取收藏岗位成功", response, c)
}

// SearchFavorites 搜索收藏列表
func (ctrl *JobCenterController) SearchFavorites(c *gin.Context) {
	var req dto.FavoriteSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Result(400, "参数错误", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.FavoriteService.SearchFavorites(req, studentUserID)
	if err != nil {
		utils.Result(500, "搜索失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取岗位列表成功", response, c)
}

// ==================== 投递相关 ====================

// ApplyJob 投递职位
func (ctrl *JobCenterController) ApplyJob(c *gin.Context) {
	var req dto.ApplyJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(400, "参数格式错误", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.ApplicationService.ApplyJob(req, studentUserID)
	if err != nil {
		utils.Result(500, err.Error(), nil, c)
		return
	}

	utils.Result(200, "岗位投递成功", response, c)
}

// ==================== 企业相关 ====================

// GetCompanyDetail 获取企业详情
func (ctrl *JobCenterController) GetCompanyDetail(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		utils.Result(400, "无效的企业ID", nil, c)
		return
	}

	// 调用Service
	response, err := ctrl.CompanyService.GetCompanyDetail(companyID)
	if err != nil {
		utils.Result(500, "获取企业详情失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取企业详情成功", response, c)
}

// ==================== 投递情况 ====================

// GetApplicationDetail 获取投递详情
func (ctrl *JobCenterController) GetApplicationDetail(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		utils.Result(400, "无效的投递ID", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.ApplicationService.GetApplicationDetail(appID, studentUserID)
	if err != nil {
		if err.Error() == "未找到该投递记录" {
			utils.Result(404, err.Error(), nil, c)
			return
		}
		utils.Result(500, "获取投递详情失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取投递详情成功", response, c)
}

// GetDeliveryList 获取已投递岗位列表
func (ctrl *JobCenterController) GetDeliveryList(c *gin.Context) {
	var req dto.DeliveryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Result(400, "参数错误", nil, c)
		return
	}

	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	studentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.ApplicationService.GetDeliveryList(req, studentUserID)
	if err != nil {
		utils.Result(500, "获取岗位列表失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取岗位列表成功", response, c)
}
