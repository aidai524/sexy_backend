package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"sexy_backend/config"
	"sexy_backend/internal/handler"
	"sexy_backend/internal/middleware"
	"sexy_backend/internal/service"
	"sexy_backend/pkg/supabase"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// 初始化 Supabase 客户端
	supabaseClient, err := supabase.NewClient(
		config.AppConfig.Supabase.URL,
		config.AppConfig.Supabase.APIKey,
	)
	if err != nil {
		log.Fatal("Cannot initialize Supabase client:", err)
	}

	// 初始化服务
	userService := service.NewUserService(supabaseClient)

	// 创建 Gin 引擎
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())

	// 设置路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		userHandler := handler.NewUserHandler(userService)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// 启动服务器
	port := ":" + config.AppConfig.Server.Port
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
