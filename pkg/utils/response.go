package utils

import (
	"github.com/gin-gonic/gin"
)

// ApiResponse 统一响应结构
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Result 统一返回方法
func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(200, ApiResponse{ // 注意：业务状态码放在 Body 里，HTTP 状态码通常设为 200
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
