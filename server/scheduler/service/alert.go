package service

import (
	"context"

	"tooladmin/server/scheduler/model"
)

// 增加告警配置表
type AlertConfig struct {
	ID         int    `gorm:"primaryKey"`
	TaskName   string `gorm:"column:task_name"`
	AlertType  string `gorm:"column:alert_type"` // 告警类型
	Threshold  int    `gorm:"column:threshold"`  // 阈值
	Recipients string `gorm:"column:recipients"` // 接收人
}

// 实现任务执行监控和告警发送
func (ls *LogService) checkAndSendAlert(ctx context.Context, log *model.Log) {
	// 根据执行结果检查是否需要发送告警
}
