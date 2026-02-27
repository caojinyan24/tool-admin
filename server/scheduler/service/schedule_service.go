package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"

	"tooladmin/server/common/config"

	"tooladmin/server/scheduler/db"
	"tooladmin/server/scheduler/model"

	"tooladmin/server/common/logger"
)

// ScheduleService 调度服务
type ScheduleService struct {
	DB          *gorm.DB
	TaskService *TaskService
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

type ManualService struct {
	DB          *gorm.DB
	TaskService *TaskService
}

// NewSchedulerService 创建调度服务实例
func NewSchedulerService(db *gorm.DB) *ScheduleService {
	return &ScheduleService{
		DB:          db,
		TaskService: NewTaskService(db),
		stopChan:    make(chan struct{}),
	}
}

// Start 启动调度服务
func (s *ScheduleService) Start() {
	s.wg.Add(3)
	go s.loadTask()
	go s.executeTask()
	go s.retryTask()
	// go s.registerGrpcCli()

}

// Stop 停止调度服务
func (s *ScheduleService) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

// startScheduler 启动调度循环
func (s *ScheduleService) loadTask() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(context.Background(), "loadTask failed", "error", r)
		}
	}()
	defer s.wg.Done()

	// 创建定时器，每秒检查一次
	ticker := time.NewTicker(time.Duration(config.GetConfig().Interval.Load) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndLoadTask(context.Background())
		case <-s.stopChan:
			return
		}
	}
}

func (s *ScheduleService) retryTask() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(context.Background(), "loadTask failed", "error", r)
		}
	}()
	defer s.wg.Done()

	// 创建定时器，每秒检查一次
	ticker := time.NewTicker(time.Duration(config.GetConfig().Interval.Retry) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndRetryTask(context.Background())
		case <-s.stopChan:
			return
		}
	}
}
func (s *ScheduleService) checkAndRetryTask(ctx context.Context) {
	var schedules []model.Schedule
	err := s.DB.Joins("inner join log on log.id=schedule.log_id").Model(&model.Schedule{}).Where("log.status =?", "F").Find(&schedules).Error
	if err != nil {
		logger.Error(ctx, "checkAndRetryTask failed", "error", err)
		return
	}
	for _, schedule := range schedules {
		//查下task
		var task model.Task
		err := s.DB.Model(&model.Task{}).Where("task_name = ?", schedule.TaskName).Find(&task).Error
		if err != nil {
			logger.Error(ctx, "checkAndRetryTask failed", "error", err)
			continue
		}
		if task.Status != "on" || task.RetryCount == 0 ||
			task.RetryInterval == 0 ||
			(task.RetryInterval > 0 && task.RetryCount > 0 && schedule.LastExecuteTime.Add(time.Duration(task.RetryInterval)*time.Second).After(time.Now().UTC())) { //不需要重试
			continue
		}
		var log model.Log
		err = s.DB.Model(&model.Log{}).Where("id=?", schedule.LogID).First(&log).Error
		if err != nil {
			logger.Error(ctx, "checkAndRetryTask failed", "error", err)
			continue
		}
		if schedule.RetryCount < task.RetryCount {
			var retryLog model.Log
			err = s.DB.Model(&model.Log{}).Where("task_name=?", task.TaskName).Where("execute_time>=?", log.ExecuteTime).Where("log.execute_type='retry'").Last(&retryLog).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					//没有重试记录, 直接重试
					s.TaskService.ExecuteAndLogTask(logger.AddTraceID(ctx), "retry", &task)
					s.DB.Model(&schedule).Where("task_name=?", schedule.TaskName).Update("retry_count", schedule.RetryCount+1)

					continue
				}
				logger.Error(ctx, "checkAndRetryTask failed", "error", err)
				continue
			}
			if retryLog.Status == "S" {
				//重试成功, 解锁
				schedule.Flag = "U"
				err = s.DB.Model(&schedule).Where("task_name=?", schedule.TaskName).Update("flag", "U").Error
				if err != nil {
					logger.Error(ctx, "checkAndRetryTask Update schedule failed", "error", err)
					continue
				}

			} else {
				if retryLog.ExecuteTime.Add(time.Duration(task.RetryInterval) * time.Second).Before(time.Now().UTC()) {
					//重试间隔超过, 重试
					logger.Info(ctx, "execute retry:", schedule)
					_, err = s.TaskService.ExecuteAndLogTask(logger.AddTraceID(ctx), "retry", &task)
					s.DB.Model(&schedule).Where("task_name=?", schedule.TaskName).Update("retry_count", schedule.RetryCount+1)
					logger.Info(ctx, "retry schedule:", schedule)
					if err != nil && schedule.RetryCount == task.RetryCount { //todo 最后一次重试失败, 发消息
						logger.Error(ctx, "Retry count exceeded, will send message, log:", log, task, schedule)

					}
				}
			}
		}

	}

}

// startScheduler 启动调度循环
func (s *ScheduleService) executeTask() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(context.Background(), "executeTask failed", "error", r)
		}
	}()
	defer s.wg.Done()

	ticker := time.NewTicker(time.Duration(config.GetConfig().Interval.Exec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndExecuteTasks(context.Background())
		case <-s.stopChan:
			return
		}
	}
}

// checkAndExecuteTasks 检查并执行需要执行的任务
func (s *ScheduleService) checkAndExecuteTasks(ctx context.Context) {
	// 查询所有调度记录，按下次执行时间排序
	//todo 做分段查询
	var schedules []model.Schedule
	err := db.DB.Joins("inner join task on schedule.task_name=task.task_name").Model(&model.Schedule{}).Where("schedule.flag!='L' and task.status='on' and schedule.next_execute_time<=?", time.Now().UTC()).Find(&schedules).Error
	if err != nil {
		logger.Error(ctx, "checkAndExecuteTasks Load schedule failed", "error", err)
		return
	}

	// 检查每个调度记录
	for _, schedule := range schedules {
		err := s.checkAndExecuteTask(ctx, &schedule)
		if err != nil {
			logger.Error(ctx, "checkAndExecuteTask failed", err, "schedule=", schedule)
		}
	}
}
func (s *ScheduleService) checkAndExecuteTask(ctx context.Context, schedule *model.Schedule) error {

	// 判断是否应该执行任务
	if ShouldExecuteTask(schedule.NextExecuteTime) {
		// 查询对应的任务信息
		task := &model.Task{}
		if err := s.DB.Where("task_name = ?", schedule.TaskName).First(task).Error; err != nil {
			// 任务不存在或已被删除，跳过
			return err
		}

		// 检查任务状态是否为启用
		if task.Status != "on" || len(task.CronExpr) == 0 || len(task.ApiAddr) == 0 || schedule.Flag == "L" {
			logger.Info(ctx, "checkAndExecuteTask Task is disabled,will skip,task_name:", task.TaskName)
			return nil
		}
		// 记录日志
		logger.Info(ctx, "checkAndExecuteTask Task execution scheduled", "task", task.TaskName, "nextExecuteTime", schedule.NextExecuteTime)

		// 异步执行任务
		s.executeTaskAsync(context.Background(), schedule.NextExecuteTime, task)
	}
	return nil
}

func (s *ScheduleService) checkAndLoadTask(ctx context.Context) {
	var tasks []model.Task
	err := s.DB.Where("status =? and cron_expr!=?", "on", "").Find(&tasks).Error
	if err != nil {
		logger.Error(ctx, "Load task failed", "error", err)
		return
	}

	//First 找不到时报错, Find找不到时不报错
	// 检查每个调度记录
	for _, task := range tasks {
		//1. scheduler 是否存在
		var schedule model.Schedule
		err := s.DB.First(&schedule, "task_name = ?", task.TaskName).Error
		logger.Debug(ctx, "checkAndLoadTask, task:", task, "schedule:", schedule, "err=", err)
		if err != nil {
			logger.Error(ctx, "checkAndLoadTask, Load schedule failed", "error", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				//insert
				schedule.TaskName = task.TaskName
				schedule.NextExecuteTime = time.Now().UTC()
				schedule.LastExecuteTime = time.Now().UTC()
				schedule.Flag = "U"
				err = s.DB.Create(&schedule).Error
				if err != nil {
					logger.Error(ctx, "checkAndLoadTask Create schedule failed", "error", err)
				}
				continue
			}
			logger.Error(ctx, "checkAndLoadTask Load schedule failed", "error", err)
			continue
		}
		//2. 检查下schedule的各个字段
		if schedule.Flag == "U" && schedule.NextExecuteTime.After(time.Now().UTC()) {
			logger.Debug(ctx, "checkAndLoadTask task check done, schedule: ", schedule)
			continue
		}
		if schedule.Flag == "L" {
			err := s.processLockedTask(ctx, &task, &schedule)
			if err != nil {
				logger.Error(ctx, "checkAndLoadTask processLockedTask failed", err)
			}
		}
		//检查下次重试时间
		if schedule.NextExecuteTime.Add(time.Second * 60).Before(time.Now().UTC()) { //下次执行时间超过一分钟未更新
			logger.Error(ctx, "checkAndLoadTask NextExecuteTime not updated for a long time, task_name: ", task.TaskName, "scheduler:", schedule)
		}

	}

}
func (s *ScheduleService) processLockedTask(ctx context.Context, task *model.Task, schedule *model.Schedule) error {
	var log model.Log
	err := s.DB.Model(&log).Where("task_name=?", task.TaskName).Where("id=?", schedule.LogID).Last(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没有执行过, 直接解锁
			schedule.Flag = "U"
			err = s.DB.Update("flag", "U").Where("task_name=?", schedule.TaskName).Error
			if err != nil {
				logger.Error(ctx, "Update schedule failed", "error", err)
				return err
			}
		}
		logger.Error(ctx, "Load log failed", "error", err)
		return err
	}
	if log.Status == "S" { //执行成功, 更新下次执行时间, 解锁
		nextExecuteTime, err := CalculateNextExecuteTime(ctx, task.CronExpr, log.ExecuteTime)
		if err != nil {
			logger.Error(ctx, "Calculate next execute time failed", "error", err)
			return err
		}
		schedule.NextExecuteTime = nextExecuteTime
		schedule.Flag = "U"
		err = s.DB.Model(&schedule).Where("task_name=?", schedule.TaskName).Update("flag", "U").Update("next_execute_time", nextExecuteTime).Error
		if err != nil {
			logger.Error(ctx, "Update schedule failed", "error", err)
			return err
		}
		return nil
	}

	//重试次数超过最大重试次数, 解锁
	if task.RetryCount == 0 && log.ExecuteTime.Add(1*time.Minute).Before(time.Now().UTC()) || task.RetryCount > 0 && schedule.RetryCount >= task.RetryCount {
		logger.Debug(ctx, "Retry count exceeded, will unlock, log:", log, task, schedule)
		schedule.Flag = "U"
		err = s.DB.Model(&schedule).Where("task_name=?", schedule.TaskName).Updates(map[string]interface{}{
			"flag": "U",
		}).Error
		if err != nil {
			logger.Error(ctx, "Update schedule failed", "error", err)
			return err
		}
		if task.RetryCount > 0 && len(task.Notifier) > 0 { //todo 发消息
			logger.Info(ctx, "Retry count exceeded, will send message, log:", log, task, schedule)
		}
	}
	return nil

}

// executeTaskAsync 异步执行任务
func (s *ScheduleService) executeTaskAsync(ctx context.Context, planedExecuteTime time.Time, task *model.Task) {
	go func() {
		// 执行任务
		result, err := s.TaskService.ExecuteCroneTask(logger.AddTraceID(ctx), planedExecuteTime, task)
		if err != nil {
			// 记录错误日志
			// 注意：这里简化处理，实际项目中应该使用更完善的日志系统
			logger.Error(ctx, "Task execution failed", "task", task.TaskName, "error", err)
		} else {
			// 记录成功日志
			logger.Info(ctx, "Task execution succeeded", "task", task.TaskName, "result", result)
		}
	}()
}
func (m *ManualService) TriggerTask(ctx context.Context, task *model.Task) (string, error) {
	return m.TaskService.ExecuteAndLogTask(ctx, "manual", task)
}
