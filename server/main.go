package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"

	"tooladmin/server/common/config"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			fmt.Println(string(debug.Stack()))
		}
	}()
	fmt.Println("Starting server...")

	// 加载配置
	configPath := filepath.Join("config", os.Getenv("StartEnv"), "conf.yml")
	fmt.Printf("Loading config from: %s\n", configPath)
	config.LoadConfig(configPath)
	fmt.Println("Config loaded successfully")

	// 初始化日志
	fmt.Println("Initializing logger...")
	logger.InitLog()
	fmt.Println("Logger initialized successfully")

	fmt.Println("Initializing scheduler...")
	scheduler.Init()
	fmt.Println("scheduler initialized successfully")

	// 创建Hertz服务器实例
	fmt.Printf("Creating server instance on port %s...\n", config.GetConfig().Server.Port)
	h := server.Default(server.WithHostPorts(fmt.Sprintf(":%s", config.GetConfig().Server.Port)))
	fmt.Println("Server instance created")

	// 添加CORS中间件
	fmt.Println("Configuring CORS middleware...")
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	fmt.Println("CORS middleware configured")

	// 初始化路由
	fmt.Println("Initializing router...")
	InitRouter(h)

	fmt.Println("Router initialized successfully")

	// devbot.RunBot(config.Env)

	// 优雅退出处理
	fmt.Println("Setting up graceful shutdown handler...")
	go func() {
		// 监听系统信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		fmt.Printf("Received signal: %v, shutting down server...\n", sig)

		logger.Info(context.Background(), "Shutting down server...")
		fmt.Println("Stopping scheduler...")
		scheduler.StopScheduler()
		fmt.Println("Scheduler stopped")

		// 创建5秒超时上下文
		fmt.Println("Creating shutdown context with 5s timeout...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 关闭服务器
		fmt.Println("Shutting down server...")
		if err := h.Shutdown(ctx); err != nil {
			fmt.Printf("Error during server shutdown: %v\n", err)
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		fmt.Println("Server shutdown complete")
		logger.Info(context.Background(), "Server exiting")
	}()

	// 启动服务器
	fmt.Printf("Starting server on port %s...\n", config.GetConfig().Server.Port)
	logger.Info(context.Background(), "Starting server on port %s...", config.GetConfig().Server.Port)
	h.Spin()
	logger.Info(context.Background(), "Server started")
}
