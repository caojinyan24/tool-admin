package service

import (
	"context"
	"fmt"
	"time"
	"tooladmin/server/common/logger"

	"github.com/robfig/cron/v3"
)

// ParseCronExpression 解析cron表达式并返回下次执行时间
func ParseCronExpression(expr string) (time.Time, error) {
	// 创建cron解析器
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cronSchedule, err := parser.Parse(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron expression: %w", err)
	}

	// 计算下次执行时间
	nextTime := cronSchedule.Next(time.Now().UTC())
	return nextTime, nil
}

// ShouldExecuteTask 判断任务是否应该执行
func ShouldExecuteTask(nextExecuteTime time.Time) bool {
	// 如果下次执行时间小于等于当前时间，则应该执行
	now := time.Now().UTC()
	result := !nextExecuteTime.After(now)
	logger.Debug(context.Background(), "ShouldExecuteTask: ", result, nextExecuteTime, now)
	return result
}

// CalculateNextExecuteTime 根据cron表达式计算下次执行时间
func CalculateNextExecuteTime(ctx context.Context, expr string, lastExecuteTime time.Time) (time.Time, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cronSchedule, err := parser.Parse(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron expression: %w", err)
	}

	// 使用上次执行时间作为基准计算下次执行时间
	nextTime := cronSchedule.Next(lastExecuteTime)
	logger.Debug(ctx, "cron表达式: ", expr, "上次执行时间:", lastExecuteTime, "下次执行时间: ", nextTime)

	return nextTime, nil
}
