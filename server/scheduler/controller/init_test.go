package controller_test

import (
	"context"
	"os"
	"testing"
	"tooladmin/server/common/config"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/controller"
	"tooladmin/server/scheduler/db"
	"tooladmin/server/scheduler/service"
)

var testController *controller.TaskController

func TestMain(m *testing.M) {
	config.LoadConfig("../../config/local/conf.yml")
	logger.InitLog()
	// 初始化数据库
	if err := db.InitDB(context.Background(), config.GetConfig().MysqlSchedule); err != nil {
		logger.Fatal(context.Background(), "Failed to initialize database: %v", err)
	}
	testController = controller.NewTaskController(service.NewTaskService(db.DB))

	os.Exit(m.Run())

}
