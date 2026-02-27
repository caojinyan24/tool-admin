package main

import (
	 "tooladmin/server/file_processor"

	"github.com/cloudwego/hertz/pkg/app/server"

	"tooladmin/server/scheduler"
)

// InitRouter 初始化路由
func InitRouter(h *server.Hertz) {

	// 定义API路由组
	api := h.Group("/api/v1")
	{
		// 任务相关路由
		tasks := api.Group("/tasks")
		{
			tasks.GET("", scheduler.TaskController.GetAllTasks)
			tasks.GET("/:task_name", scheduler.TaskController.GetTaskByTaskName)
			tasks.POST("/add", scheduler.TaskController.CreateTask)
			tasks.PUT("/edit", scheduler.TaskController.UpdateTask)
			tasks.DELETE("/:task_name", scheduler.TaskController.DeleteTask)
			tasks.POST("/trigger", scheduler.TaskController.TriggerTask)
		}

		// 日志相关路由
		logs := api.Group("/logs")
		{
			logs.GET("", scheduler.LogController.GetLogs)
		}

		// 文件处理相关路由
		fileProcessor := api.Group("/file-processor")
		{
			fileProcessor.POST("/merge-pdf", file_processor.FileProcessor.MergePDF)
		}

	}
}
