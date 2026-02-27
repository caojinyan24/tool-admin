package scheduler

import (
	"context"
	"tooladmin/server/common/config"
	"tooladmin/server/scheduler/controller"
	"tooladmin/server/scheduler/db"

	"tooladmin/server/scheduler/service"
)

var TaskController *controller.TaskController
var LogController *controller.LogController
var schedulerService *service.ScheduleService

func Init() {
	// 初始化数据库
	if err := db.InitDB(context.Background(), config.GetConfig().MysqlSchedule); err != nil {
		// panic(err) todo add log
	}
	// 创建服务
	taskService := service.NewTaskService(db.DB)
	logService := service.NewLogService(db.DB)

	// 创建控制器
	TaskController = controller.NewTaskController(taskService)
	LogController = controller.NewLogController(logService)
	// 创建并启动调度服务
	schedulerService = service.NewSchedulerService(db.DB)
	schedulerService.Start()
}
func StopScheduler() {
	schedulerService.Stop()
}
