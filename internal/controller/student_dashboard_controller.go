package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StudentDashboardController struct {
	Service *service.StudentDashboardService
}

// GetUserName 获取当前用户名
func (ctrl *StudentDashboardController) GetUserName(c *gin.Context) {
	// 从Context获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Result(401, "未登录", nil, c)
		return
	}

	currentUserID := userID.(int)

	// 调用Service
	response, err := ctrl.Service.GetUserName(currentUserID)
	if err != nil {
		utils.Result(500, "获取用户信息失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取用户信息成功", response, c)
}

// GetCalendar 获取求职日历数据
func (ctrl *StudentDashboardController) GetCalendar(c *gin.Context) {
	var req dto.CalendarRequest
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
	response, err := ctrl.Service.GetCalendar(req.Month, currentUserID)
	if err != nil {
		utils.Result(500, "获取日历数据失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取日历数据成功", response, c)
}

// GetUpcomingEvents 获取近期求职活动
func (ctrl *StudentDashboardController) GetUpcomingEvents(c *gin.Context) {
	// 解析limit参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	// 调用Service
	response, err := ctrl.Service.GetUpcomingEvents(limit)
	if err != nil {
		utils.Result(500, "获取近期活动失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取近期活动成功", response, c)
}

// GetRankedJobs 获取岗位热度排行榜
func (ctrl *StudentDashboardController) GetRankedJobs(c *gin.Context) {
	// 解析limit参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	// 调用Service
	response, err := ctrl.Service.GetRankedJobs(limit)
	if err != nil {
		utils.Result(500, "获取热度排行失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取热度排行成功", response, c)
}

// GetRecentJobs 获取近期招聘信息列表
func (ctrl *StudentDashboardController) GetRecentJobs(c *gin.Context) {
	var req dto.RecentJobsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Result(400, "参数错误", nil, c)
		return
	}

	// 调用Service
	response, err := ctrl.Service.GetRecentJobs(req.JobTypeFilter, req.Limit)
	if err != nil {
		utils.Result(500, "获取近期招聘列表失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "获取近期招聘列表成功", response, c)
}
