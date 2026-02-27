package bot_test

import (
	"context"
	"os"
	"tooladmin/server/common/config"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/alert"
	"tooladmin/server/scheduler/db"
	"tooladmin/server/scheduler/model"

	"testing"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../../config/local/conf.yml")
	logger.InitLog()
	// 初始化数据库
	if err := db.InitDB(context.Background(), config.GetConfig().MysqlSchedule); err != nil {
		logger.Fatal(context.Background(), "Failed to initialize database: %v", err)
	}
	os.Exit(m.Run())
}
func TestSendScheduleFailMsg(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		msgData *model.ScheduleFailMsg
	}{
		{
			name: "test1",
			msgData: &model.ScheduleFailMsg{
				TaskName:   "test1",
				LogUrl:     "test1",
				ReceiveIds: []string{"ou_5415915c805338ea6a149399edb035fe"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert.SendScheduleFailMsg(context.Background(), tt.msgData)
		})
	}
}
