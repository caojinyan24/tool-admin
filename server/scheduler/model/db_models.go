package model

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务表模型
type Task struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TaskName      string    `gorm:"size:100;not null;uniqueIndex" json:"task_name"`                  // 任务名称
	ApiAddr       string    `gorm:"size:255;not null" json:"api_addr"`                               // 接口地址
	TaskParam     string    `gorm:"column:task_param;type:text" json:"task_param"`                   // 接口参数
	ApiType       string    `gorm:"size:10;not null" json:"api_type"`                                // 接口类型(grpc, http)
	ApiName       string    `gorm:"size:255;not null" json:"api_name"`                               // 接口名称
	CronExpr      string    `gorm:"size:64;not null" json:"cron_expr"`                               // 调度时间表达式
	Status        string    `gorm:"size:5;not null;default:off" json:"status"`                       // 任务状态(on:启用, off:停用)
	RetryCount    int       `gorm:"column:retry_count;not null;default:0" json:"retry_count"`        // 重试次数
	RetryInterval int       `gorm:"column:retry_interval;not null;default:60" json:"retry_interval"` // 重试间隔(秒)
	Notifier      string    `gorm:"column:notifier;not null;default:256" json:"notifier"`
	Description   string    `gorm:"column:description;type:text" json:"description"` // 任务描述
	CreatedAt     time.Time `gorm:"created_at;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"updated_at;not null" json:"updated_at"`
}

// Log 日志表模型
//
//	type Log struct {
//		ID          uint      `gorm:"primaryKey" json:"id"`
//		TaskName    string    `gorm:"size:100;not null;index" json:"task_name"`              // 任务名称
//		ExecuteTime time.Time `gorm:"column:execute_time;not null" json:"execute_time"`      // 执行时间
//		Status      string    `gorm:"column:status;size:5;not null;default:I" json:"status"` // 执行状态(S:成功, F:失败, I:初始化)
//		Result      string    `gorm:"type:text" json:"result"`                               // 执行结果
//		TaskParam   string    `gorm:"type:text" json:"task_param"`                           // 执行结果
//		ExecuteType string    `gorm:"size:10;not null;default:''" json:"execute_type"`       // 执行类型(manual, scheduled)
//		TraceID     string    `gorm:"size:100" json:"trace_id"`                              // 追踪ID
//		CreatedAt   time.Time `gorm:"created_at;not null" json:"created_at"`
//		UpdatedAt   time.Time `gorm:"updated_at;not null" json:"updated_at"`
//	}
type Log struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TaskName    string    `gorm:"size:100;not null;index" json:"task_name"`              // 任务名称
	ExecuteTime time.Time `gorm:"column:execute_time;not null" json:"execute_time"`      // 执行时间
	Status      string    `gorm:"column:status;size:5;not null;default:I" json:"status"` // 执行状态(S:成功, F:失败, I:初始化)
	Result      string    `gorm:"type:text" json:"result"`                               // 执行结果
	TaskParam   string    `gorm:"type:text" json:"task_param"`                           // 执行结果
	ExecuteType string    `gorm:"size:10;not null;default:''" json:"execute_type"`       // 执行类型(manual, scheduled)
	TraceID     string    `gorm:"size:100" json:"trace_id"`                              // 追踪ID
	CreatedAt   time.Time `gorm:"created_at;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at;not null" json:"updated_at"`
}

func (s *Log) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now().UTC()
	return nil
}
func (s *Log) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now().UTC()
	return nil
}

// Schedule 调度表模型
//调度: 1. load到schedule, 2. 解锁, 3, 失败重试, 4, 按照schedule的下次执行时间做执行(log的插入和schedule的加锁放在一个事务里)

type Schedule struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TaskName        string    `gorm:"size:100;not null;uniqueIndex" json:"task_name"`                                           // 任务名称
	LastExecuteTime time.Time `gorm:"column:last_execute_time;not null;default:'1970-01-01 00:00:00'" json:"last_execute_time"` // 上次执行时间
	NextExecuteTime time.Time `gorm:"column:next_execute_time;not null;default:'1970-01-01 00:00:00'" json:"next_execute_time"` // 下次执行时间
	Flag            string    `gorm:"column:flag;size:5;not null;default:''" json:"flag"`                                       // 锁, 除失败重试外, 其他情况需要等待锁解开后再尝试; 重试条件为设置了重试次数, 且达到可重试间隔, 且log状态为失败
	RetryCount      int       `gorm:"column:retry_count;not null;default:0" json:"retry_count"`                                 // 重试次数--重试次数完成后且log达终态,或者重试成功 : 做解锁操作.
	LogID           uint      `gorm:"column:log_id;not null;default:0" json:"log_id"`                                           //根据logid的执行时间, 计算按照正常执行周期的下次执行时间
	UpdatedAt       time.Time `gorm:"updated_at;not null" json:"updated_at"`
	CreatedAt       time.Time `gorm:"created_at;not null" json:"created_at"`
}

// add updateAt自动更新
func (s *Schedule) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now().UTC()
	return nil
}
func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now().UTC()
	s.UpdatedAt = time.Now().UTC()
	return nil
}
