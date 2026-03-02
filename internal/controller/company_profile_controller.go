package controller

import (
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/repository"
	"job-platform-go2/internal/service"
	"job-platform-go2/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CompanyProfileController struct {
	CompanyService *service.CompanyService
	DictRepo       *repository.DictionaryRepository
}

// GetCompanyProfile 获取企业信息
// GET /hr/company/profile
func (c *CompanyProfileController) GetCompanyProfile(ctx *gin.Context) {
	// 1. 从JWT获取companyID
	companyID, exists := ctx.Get("companyID")
	if !exists {
		utils.Result(401, "未授权", nil, ctx)
		return
	}

	// 2. 调用Service
	resp, err := c.CompanyService.GetCompanyProfile(companyID.(int))
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 3. 返回
	utils.Result(200, "成功", resp, ctx)
}

// UpdateCompanyProfile 更新企业信息
// PUT /hr/company/profile
func (c *CompanyProfileController) UpdateCompanyProfile(ctx *gin.Context) {
	// 1. 绑定JSON body
	var req dto.UpdateCompanyProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
	err := c.CompanyService.UpdateCompanyProfile(companyID.(int), req, c.DictRepo)
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 4. 返回
	utils.Result(200, "更新成功", nil, ctx)
}

// GetCompanyOptions 获取企业信息选项
// GET /hr/company/options
func (c *CompanyProfileController) GetCompanyOptions(ctx *gin.Context) {
	// 1. 调用Service
	resp, err := c.CompanyService.GetCompanyOptions(c.DictRepo)
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 2. 返回
	utils.Result(200, "成功", resp, ctx)
}

// UploadCompanyLogo 上传企业Logo
// POST /hr/upload/company-logo
func (c *CompanyProfileController) UploadCompanyLogo(ctx *gin.Context) {
	// 1. 从JWT获取companyID
	companyID, exists := ctx.Get("companyID")
	if !exists {
		utils.Result(401, "未授权", nil, ctx)
		return
	}

	// 2. 获取文件
	file, err := ctx.FormFile("logo_file")
	if err != nil {
		utils.Result(400, "未上传文件", nil, ctx)
		return
	}

	// 3. 调用Service
	logoURL, err := c.CompanyService.UploadCompanyLogo(companyID.(int), file)
	if err != nil {
		utils.Result(500, err.Error(), nil, ctx)
		return
	}

	// 4. 返回
	utils.Result(200, "上传成功", gin.H{"logo_url": logoURL}, ctx)
}
