package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"tooladmin/server/common/logger"

	"tooladmin/server/scheduler/model"
)

// TaskService 任务服务
type TaskService struct {
	DB *gorm.DB
}

// NewTaskService 创建任务服务实例
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{DB: db}
}

// GetAllTasks 获取所有任务
func (s *TaskService) GetAllTasks() ([]model.Task, error) {
	var tasks []model.Task
	result := s.DB.Model(&model.Task{}).Order("id DESC").Find(&tasks)
	return tasks, result.Error
}

// GetTaskByID 根据ID获取任务
func (s *TaskService) GetTaskByTaskName(taskName string) (*model.Task, error) {
	var task model.Task
	result := s.DB.First(&task, "task_name = ?", taskName)
	return &task, result.Error
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(ctx context.Context, task *model.Task) error {

	// 首先创建任务
	result := s.DB.Create(task)
	if result.Error != nil {
		return result.Error
	}

	if len(task.CronExpr) == 0 {
		return nil
	}
	// 然后创建调度记录
	schedule := &model.Schedule{
		TaskName:        task.TaskName,
		LastExecuteTime: time.Time{},
		NextExecuteTime: time.Time{}, // 将在调度器中更新
	}

	// 如果任务是启用状态，计算下次执行时间
	if task.Status == "on" {
		nextTime, err := ParseCronExpression(task.CronExpr)
		if err != nil {
			return err
		}
		schedule.NextExecuteTime = nextTime
	}

	return s.DB.Create(schedule).Error
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(ctx context.Context, task *model.Task) error {
	var currentTask model.Task
	if err := s.DB.Where("task_name = ?", task.TaskName).First(&currentTask).Error; err != nil {
		return err
	}
	if currentTask.Notifier != task.Notifier {
		if task.Notifier == "" {
			task.Notifier = "{}"
		} else {
			var emailMap = make(map[string]model.Notifier)
			err := json.Unmarshal([]byte(task.Notifier), &emailMap)
			if err != nil {
				return err
			}
			for _, notifier := range emailMap {
				if len(notifier.ReceiveId) == 0 {
					logger.Error(ctx, "ReceiveId is empty, will not send email")
				}
			}

			notifierBytes, err := json.Marshal(emailMap)
			if err != nil {
				return err
			}
			task.Notifier = string(notifierBytes)
		}
	}

	// 更新任务信息
	result := s.DB.Updates(task)
	if result.Error != nil {
		return result.Error
	}

	// 如果任务是启用状态，更新调度记录的下次执行时间
	if task.Status == "on" && len(task.CronExpr) > 0 {
		nextTime, err := ParseCronExpression(task.CronExpr)
		if err != nil {
			return err
		}

		schedule := &model.Schedule{}
		if err := s.DB.Where("task_name = ?", task.TaskName).First(schedule).Error; err != nil {
			// 如果调度记录不存在，创建新的
			schedule = &model.Schedule{
				TaskName:        task.TaskName,
				NextExecuteTime: nextTime,
			}
			return s.DB.Create(schedule).Error
		}

		// 更新现有调度记录
		schedule.NextExecuteTime = nextTime
		return s.DB.Save(schedule).Error
	}

	return nil
}

// DeleteTask 删除任务 - 支持通过任务名称或ID删除
func (s *TaskService) DeleteTask(ctx context.Context, taskName string) error {
	// 首先尝试通过任务名称查询
	task, err := s.GetTaskByTaskName(taskName)

	// 如果通过名称查询失败，尝试通过ID查询
	if err != nil {
		logger.Error(ctx, "GetTask error:", err)
		return err

	}

	// 删除任务
	if err := s.DB.Delete(&model.Task{}, "task_name = ?", task.TaskName).Error; err != nil {
		return err
	}

	return nil
}

// ExecuteTask 执行任务
// 需要确保schedule记录一定有
func (s *TaskService) ExecuteCroneTask(ctx context.Context, planedExecuteTime time.Time, task *model.Task) (string, error) {

	// 做下校验
	var currentSchedule = model.Schedule{}
	err := s.DB.Where("task_name = ?", task.TaskName).First(&currentSchedule).Error
	if err != nil {
		return "", err
	}

	// 创建日志记录
	tx := s.DB.Begin()
	err = tx.Set("gorm:query_option", "FOR UPDATE").Table("schedule").First(&currentSchedule, "task_name = ?", task.TaskName).Error
	defer tx.Rollback()
	if err != nil {
		tx.Rollback()
		return "", err
	}
	if currentSchedule.Flag == "L" || currentSchedule.NextExecuteTime != planedExecuteTime {
		tx.Rollback()
		return "", fmt.Errorf("task %s is locked", task.TaskName)
	}

	log := &model.Log{
		TaskName:    task.TaskName,
		ExecuteTime: time.Now().UTC(),
		ExecuteType: "scheduled",
		TaskParam:   task.TaskParam,
		TraceID:     ctx.Value(logger.TraceIDKey{}).(string),
	}
	tx.Model(&model.Log{}).Create(log)
	logger.Debug(ctx, "inserted log", log)
	tx.Model(&model.Schedule{}).Where("task_name = ?", task.TaskName).Updates(map[string]interface{}{
		"flag":        "L",
		"log_id":      log.ID,
		"retry_count": 0,
	})
	tx.Commit()

	// 执行HTTP请求
	result, err := NewRemoteRequest(task).SendRequest(ctx)
	if err != nil { //执行失败不解锁, 等load任务进行重试或解锁
		logger.Error(ctx, "ExecuteTask error:%v", err)
		log.Result = err.Error()
		log.Status = "F"
		s.DB.Model(&log).Where("id = ?", log.ID).Updates(map[string]interface{}{
			"status":     log.Status,
			"result":     log.Result,
			"updated_at": time.Now().UTC(),
		})
	} else {
		logger.Info(ctx, "ExecuteTask success:", result)
		// 更新日志状态为成功
		log.Status = "S"
		log.Result = result
		// 保存成功日志
		s.DB.Model(&log).Where("id = ?", log.ID).Updates(map[string]interface{}{
			"status":     log.Status,
			"result":     log.Result,
			"updated_at": time.Now().UTC(),
		})
	}

	// 更新调度记录的上次执行时间
	nextTime, err := CalculateNextExecuteTime(ctx, task.CronExpr, log.ExecuteTime)
	if err != nil {
		logger.Error(ctx, "CalculateNextExecuteTime error:%v", err)
		return "", err
	}
	var schedule = model.Schedule{}
	schedule.NextExecuteTime = nextTime
	schedule.LastExecuteTime = log.ExecuteTime
	schedule.TaskName = task.TaskName
	schedule.Flag = "U"
	s.DB.Model(&schedule).Where("task_name = ?", task.TaskName).Where("log_id=?", log.ID).Updates(map[string]interface{}{
		"next_execute_time": schedule.NextExecuteTime,
		"last_execute_time": schedule.LastExecuteTime,
		"flag":              schedule.Flag,
		"updated_at":        time.Now().UTC(),
	})
	return result, nil
}

// TriggerTask 手动触发任务
func (s *TaskService) ExecuteManualTask(ctx context.Context, task *model.Task) (string, error) {
	taskDb, err := s.GetTaskByTaskName(task.TaskName)
	if err != nil {
		return "", err
	}
	task.ApiAddr = taskDb.ApiAddr
	task.ApiName = taskDb.ApiName
	task.ApiType = taskDb.ApiType
	task.Status = taskDb.Status
	// 即使任务是停用状态，手动触发也应该执行
	// 创建日志记录

	log := &model.Log{
		TaskName:    task.TaskName,
		ExecuteTime: time.Now().UTC(),
		ExecuteType: "manual",
		TaskParam:   task.TaskParam,
		TraceID:     ctx.Value(logger.TraceIDKey{}).(string),
	}
	s.DB.Create(log)
	logger.Debug(ctx, "inserted log", log)
	// 执行HTTP请求
	result, err := NewRemoteRequest(task).SendRequest(ctx)
	if err != nil {
		logger.Error(ctx, "ExecuteTask error:%v", err)
		log.Result = err.Error()
		log.Status = "F"
		s.DB.Updates(log)
		// 保存失败日志
		return "", err
	}
	logger.Info(ctx, "ExecuteTask success:", result)
	// 更新日志状态为成功
	log.Status = "S"
	log.Result = result
	// 保存成功日志
	s.DB.Updates(log)
	return result, nil
}
func (s *TaskService) ExecuteAndLogTask(ctx context.Context, executeFlg string, task *model.Task) (string, error) {
	taskDb, err := s.GetTaskByTaskName(task.TaskName)
	if err != nil {
		logger.Error(ctx, "ExecuteAndLogTask error:", err)
		return "", err
	}
	task.ApiAddr = taskDb.ApiAddr
	task.ApiName = taskDb.ApiName
	task.ApiType = taskDb.ApiType
	task.Status = taskDb.Status
	// 即使任务是停用状态，手动触发也应该执行
	// 创建日志记录

	log := &model.Log{
		TaskName:    task.TaskName,
		ExecuteTime: time.Now().UTC(),
		ExecuteType: executeFlg,
		TaskParam:   task.TaskParam,
		TraceID:     ctx.Value(logger.TraceIDKey{}).(string),
	}
	s.DB.Create(log)
	logger.Debug(ctx, "inserted log", log)
	// 执行HTTP请求
	result, err := NewRemoteRequest(task).SendRequest(ctx)
	if err != nil {
		logger.Error(ctx, "ExecuteTask error:%v", err)
		log.Result = err.Error()
		log.Status = "F"
		s.DB.Updates(log)
		// 保存失败日志
		return "", err
	}
	logger.Info(ctx, "ExecuteTask success:", result)
	// 更新日志状态为成功
	log.Status = "S"
	log.Result = result
	// 保存成功日志
	s.DB.Updates(log)
	return result, nil
}
