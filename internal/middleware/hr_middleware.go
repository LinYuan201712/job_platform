package middleware

import (
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"
	"job-platform-go2/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HRAuth HR权限中间件
// 1. 验证用户角色是否为HR
// 2. 将 companyID 注入 Context
func HRAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户角色
		role, exists := c.Get("userRole")
		if !exists {
			utils.Result(401, "未授权: 无法获取用户角色", nil, c)
			c.Abort()
			return
		}

		// 2. 验证角色
		if role.(int) != model.UserRoleHr {
			utils.Result(403, "权限不足: 仅限企业账号访问", nil, c)
			c.Abort()
			return
		}

		// 3. 获取 UserID
		userID, exists := c.Get("userID")
		if !exists {
			utils.Result(401, "未授权: 无法获取用户ID", nil, c)
			c.Abort()
			return
		}

		// 4. 查询 CompanyID
		repo := &repository.CompanyRepository{DB: db}
		company, err := repo.GetCompanyByUserID(userID.(int))
		if err != nil {
			utils.Result(401, "未找到关联的企业信息", nil, c)
			c.Abort()
			return
		}

		// 5. 注入 Context
		c.Set("companyID", company.CompanyID)
		c.Next()
	}
}
