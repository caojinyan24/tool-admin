package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/model"
)

// LogService 日志服务
type LogService struct {
	DB *gorm.DB
}

// NewLogService 创建日志服务实例
func NewLogService(db *gorm.DB) *LogService {
	return &LogService{DB: db}
}

// GetLogs 获取日志列表，支持分页和过滤
func (s *LogService) GetLogs(query map[string]interface{}) ([]model.Log, int, error) {
	var logs []model.Log
	var total int64

	// 构建查询条件
	db := s.DB.Model(&model.Log{})

	// 任务名称过滤
	if taskName, ok := query["taskName"].(string); ok && taskName != "" {
		db = db.Where("task_name LIKE ?", fmt.Sprintf("%%%s%%", taskName))
	}

	// 状态过滤
	if status, ok := query["status"].(string); ok && status != "" {
		db = db.Where("status = ?", status)
	}

	// 执行类型过滤
	if executeType, ok := query["executeType"].(string); ok && executeType != "" {
		db = db.Where("execute_type = ?", executeType)
	}

	// 时间范围过滤
	// 处理时间范围
	if startTime, ok := query["startTime"].(string); ok && startTime != "" {
		// 尝试解析ISO格式时间
		start, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			// 如果失败，尝试解析原有格式
			start, err = time.Parse("2006-01-02 15:04:05", startTime)
		}
		if err == nil {
			db = db.Where("execute_time >= ?", start)
		}
	}

	if endTime, ok := query["endTime"].(string); ok && endTime != "" {
		// 尝试解析ISO格式时间
		end, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			// 如果失败，尝试解析原有格式
			end, err = time.Parse("2006-01-02 15:04:05", endTime)
		}
		if err == nil {
			db = db.Where("execute_time <= ?", end)
		}
	}

	// 获取总数
	db.Count(&total)

	// 排序
	db = db.Order("execute_time DESC")

	// 分页
	if page, ok := query["page"].(int); ok && page > 0 {
		if pageSize, ok := query["pageSize"].(int); ok && pageSize > 0 {
			offset := (page - 1) * pageSize
			db = db.Offset(offset).Limit(pageSize)
		}
	}

	// 查询数据
	err := db.Find(&logs).Error

	return logs, int(total), err
}

// GetLogByID 根据ID获取日志
func (s *LogService) GetLogByID(id uint) (*model.Log, error) {
	var log model.Log
	result := s.DB.First(&log, id)
	return &log, result.Error
}

// CreateLog 创建日志记录
func (s *LogService) CreateLog(ctx context.Context, log *model.Log) error {
	logger.Info(ctx, "CreateLog:", log)
	return s.DB.Create(log).Error
}

// DeleteLog 删除日志记录
func (s *LogService) DeleteLog(id uint) error {
	return s.DB.Delete(&model.Log{}, id).Error
}

// CleanLogs 清理指定时间之前的日志
func (s *LogService) CleanLogs(before time.Time) error {
	return s.DB.Where("execute_time < ?", before).Delete(&model.Log{}).Error
}
