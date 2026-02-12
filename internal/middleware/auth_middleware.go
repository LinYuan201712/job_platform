package middleware

import (
	"job-platform-go2/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth 中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未携带 Token"})
			c.Abort()
			return
		}

		// 2. 格式通常是 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 格式错误"})
			c.Abort()
			return
		}

		// 3. 解析 Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 4. 将用户信息存入 Context，供后续接口使用
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}
