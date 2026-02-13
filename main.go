package main

import (
	"fmt"
	"job-platform-go2/internal/controller"
	"job-platform-go2/internal/middleware"
	"job-platform-go2/internal/repository"
	"job-platform-go2/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库对象 (为了 Day 1 演示简单直接放这里，之后会移到 global 或 internal 包)
var DB *gorm.DB

func main() {
	initConfig()
	initDB()

	// --- 依赖注入 (手动) ---
	// 1. 初始化仓库层 (Repository)
	jobRepo := &repository.JobRepository{DB: DB}
	tagRepo := &repository.TagRepository{DB: DB}
	jobAuditRepo := &repository.JobAuditRepository{DB: DB}
	studentRepo := &repository.StudentRepository{DB: DB}
	educationRepo := &repository.EducationExperienceRepository{DB: DB}
	workExpRepo := &repository.WorkExperienceRepository{DB: DB}
	projectExpRepo := &repository.ProjectExperienceRepository{DB: DB}
	orgExpRepo := &repository.OrganizationExperienceRepository{DB: DB}
	compExpRepo := &repository.CompetitionExperienceRepository{DB: DB}
	resumeRepo := &repository.ResumeRepository{DB: DB}

	// 2. 初始化 Service
	authService := &service.AuthService{DB: DB}
	jobService := &service.JobService{
		JobRepo: jobRepo,
		TagRepo: tagRepo,
		DB:      DB,
	}
	tagService := &service.TagService{TagRepo: tagRepo}
	jobAuditService := &service.JobAuditService{AuditRepo: jobAuditRepo}
	jobParseService := &service.JobParseService{}
	resumeService := &service.ResumeService{
		DB:               DB,
		StudentRepo:      studentRepo,
		EducationRepo:    educationRepo,
		WorkRepo:         workExpRepo,
		ProjectRepo:      projectExpRepo,
		OrganizationRepo: orgExpRepo,
		CompetitionRepo:  compExpRepo,
		ResumeRepo:       resumeRepo,
	}
	experienceService := &service.ExperienceService{
		WorkRepo:         workExpRepo,
		ProjectRepo:      projectExpRepo,
		OrganizationRepo: orgExpRepo,
		CompetitionRepo:  compExpRepo,
	}

	// 3. 初始化 Controller
	authController := &controller.AuthController{Service: authService}
	jobController := &controller.JobController{
		JobService:      jobService,
		JobAuditService: jobAuditService,
		JobParseService: jobParseService,
	}
	tagController := &controller.TagController{TagService: tagService}
	resumeController := &controller.ResumeController{
		ResumeService:     resumeService,
		ExperienceService: experienceService,
	}

	// --- 路由设置 ---
	r := gin.Default()
	// 注册 CORS 中间件 (必须放在所有路由之前)
	r.Use(middleware.CORS())

	// 开放接口 (无需登录)
	public := r.Group("/auth")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
	}

	// 受保护接口 (需要 Token)
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth()) // 挂载中间件
	{
		// 测试接口
		protected.GET("/profile", func(c *gin.Context) {
			// 从 Context 取出中间件存入的信息
			userID, _ := c.Get("userID")
			email, _ := c.Get("userEmail")
			c.JSON(200, gin.H{
				"message": "你已通过身份验证",
				"user_id": userID,
				"email":   email,
			})
		})

		// 标签管理接口 (所有用户都可访问)
		protected.GET("/tags", tagController.GetAllTags)
		protected.POST("/tags", tagController.CreateTag)

		// HR 岗位管理接口
		hr := protected.Group("/hr")
		{
			// 岗位基本操作
			hr.POST("/jobs", jobController.CreateJob)
			hr.GET("/jobs/:job_id", jobController.GetJobDetail)
			hr.PUT("/jobs/:job_id", jobController.UpdateJob)
			hr.DELETE("/jobs/:job_id", jobController.DeleteJobDraft)

			// 岗位状态管理
			hr.PUT("/jobs/:job_id/status", jobController.OfflineJob)

			// 岗位审核详情
			hr.GET("/jobs/audit/:jobId", jobController.GetJobAudit)

			// 岗位智能解析 (AI)
			hr.POST("/jobs/parse", jobController.ParseJob)
		}

		// 学生简历中心接口
		resumeCenter := protected.Group("/resume-center")
		{
			// 简历草稿
			resumeCenter.GET("/resume_draft", resumeController.GetResumeDraft)
			resumeCenter.PUT("/resume_draft/skills", resumeController.UpdateSkills)
			resumeCenter.PATCH("/resume_draft/template", resumeController.SetTemplate)

			// 简历文件
			resumeCenter.POST("/resume_files/upload", resumeController.UploadResumeFile)
			resumeCenter.GET("/resume_files", resumeController.GetResumeFiles)
			resumeCenter.DELETE("/resume_files/:id", resumeController.DeleteResumeFile)

			// 工作经历
			resumeCenter.POST("/resume_draft/work_experiences", resumeController.CreateWorkExperience)
			resumeCenter.PUT("/resume_draft/work_experiences/:id", resumeController.UpdateWorkExperience)
			resumeCenter.DELETE("/resume_draft/work_experiences/:id", resumeController.DeleteWorkExperience)

			// 项目经历
			resumeCenter.POST("/resume_draft/projects", resumeController.CreateProjectExperience)
			resumeCenter.PUT("/resume_draft/projects/:id", resumeController.UpdateProjectExperience)
			resumeCenter.DELETE("/resume_draft/projects/:id", resumeController.DeleteProjectExperience)

			// 组织经历
			resumeCenter.POST("/resume_draft/organizations", resumeController.CreateOrganizationExperience)
			resumeCenter.PUT("/resume_draft/organizations/:id", resumeController.UpdateOrganizationExperience)
			resumeCenter.DELETE("/resume_draft/organizations/:id", resumeController.DeleteOrganizationExperience)

			// 竞赛经历
			resumeCenter.POST("/resume_draft/competitions", resumeController.CreateCompetitionExperience)
			resumeCenter.PUT("/resume_draft/competitions/:id", resumeController.UpdateCompetitionExperience)
			resumeCenter.DELETE("/resume_draft/competitions/:id", resumeController.DeleteCompetitionExperience)
		}
	}

	port := viper.GetString("server.port")
	fmt.Printf("Starting server on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}

func initConfig() {
	viper.SetConfigName("config") // 配置文件名 (不带后缀)
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 搜索路径 (当前目录)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	fmt.Println("Configuration loaded successfully.")
}

func initDB() {
	dsn := viper.GetString("database.dsn")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 验证一下连接
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}

	// 设置连接池 (对应 Spring Boot 的 HikariCP)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	fmt.Println("Database connected successfully!")
}
