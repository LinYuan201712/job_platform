package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 全局数据库对象 (为了 Day 1 演示简单直接放这里，之后会移到 global 或 internal 包)
var DB *gorm.DB

func main() {
	// 1. 初始化配置
	initConfig()

	// 2. 连接数据库
	initDB()

	// 3. 启动 Web 服务器
	r := gin.Default()

	// 添加一个测试路由，验证服务是否存活
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "Go Backend is running!",
		})
	})

	// 添加一个数据库测试路由，验证数据库是否连接成功
	r.GET("/db-check", func(c *gin.Context) {
		sqlDB, err := DB.DB()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get generic database object"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"error": "Database connection lost", "details": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "Database connection is healthy"})
	})

	port := viper.GetString("server.port")
	fmt.Printf("Starting server on port %s...\n", port)

	// 启动服务
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
