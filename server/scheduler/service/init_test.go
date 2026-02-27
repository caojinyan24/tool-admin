package service

import (
	"context"
	"os"
	"testing"

	"tooladmin/server/common/config"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/db"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../config/test/conf.yml")
	logger.InitLog()
	// 初始化数据库
	if err := db.InitDB(context.Background(), config.GetConfig().MysqlSchedule); err != nil {
		logger.Fatal(context.Background(), "Failed to initialize database: %v", err)
	}
	os.Exit(m.Run())
}
