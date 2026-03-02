package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	Service *service.EventService
}

// GetEventList 获取招聘活动列表（分页+搜索）
func (ctrl *EventController) GetEventList(c *gin.Context) {
	var req dto.EventListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Result(400, "参数错误", nil, c)
		return
	}

	// 调用Service
	response, err := ctrl.Service.GetEventList(req)
	if err != nil {
		utils.Result(500, "获取活动列表失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "success", response, c)
}

// GetEventDetail 获取招聘活动详情
func (ctrl *EventController) GetEventDetail(c *gin.Context) {
	eventIDStr := c.Param("event_id")

	// 处理前端传递的带前缀的ID (如 evt-uuid-002 或 uuid-002)
	if len(eventIDStr) > 0 {
		// 移除可能的 "evt-" 前缀
		if len(eventIDStr) >= 4 && eventIDStr[:4] == "evt-" {
			eventIDStr = eventIDStr[4:]
		}
		// 移除可能的 "uuid-" 前缀
		if len(eventIDStr) >= 5 && eventIDStr[:5] == "uuid-" {
			eventIDStr = eventIDStr[5:]
		}
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		utils.Result(400, "无效的活动ID", nil, c)
		return
	}

	// 调用Service
	response, err := ctrl.Service.GetEventDetail(eventID)
	if err != nil {
		utils.Result(500, "获取活动详情失败: "+err.Error(), nil, c)
		return
	}

	utils.Result(200, "success", response, c)
}
