package service

import (
	"fmt"
	"testing"
	"time"
	"tooladmin/server/scheduler/db"
	"tooladmin/server/scheduler/model"
)

func TestDB(t *testing.T) {
	var schedules []model.Schedule
	err := db.DB.Joins("inner join task on schedule.task_name=task.task_name").Model(&model.Schedule{}).Where("schedule.flag!='L' and task.status='on' and schedule.next_execute_time<=?", time.Now().UTC()).Find(&schedules).Error
	fmt.Println(err)
	fmt.Println(schedules)

	var task model.Task
	err = db.DB.Model(&model.Task{}).Where("task_name=?", "test").First(&task).Error
	fmt.Println(task)
	fmt.Println(err)
	err = db.DB.Model(&model.Task{}).Where("task_name=?", "test").Find(&task).Error
	fmt.Println(task)
	fmt.Println(err)
}

func TestTaskService_GetAllTasks(t *testing.T) {
	taskService := NewTaskService(db.DB)
	tasks, err := taskService.GetAllTasks()
	fmt.Println(tasks)
	fmt.Println(err)
}

func TestUpdate(t *testing.T) {
	var err error
	// newLog := model.Log{
	// 	TaskName: "test",
	// }
	// err := utils.DB.Create(&newLog).Error
	// fmt.Println("create error", err)

	// // 方法1: 使用Updates方法和结构体，会将所有非零值字段更新
	// err = utils.DB.Model(&model.Log{}).Where("id=1").Updates(&model.Log{
	// 	Result:  "",
	// 	TraceID: "",
	// 	Status:  "",
	// }).Error
	// fmt.Println("update with Update method error", err)

	// 方法2: 使用Select("*").Omit()更新除指定字段外的所有字段
	// err = utils.DB.Model(&model.Log{}).Where("id=1").Select("*").Omit("id", "task_name", "created_at").Updates(&model.Log{
	// 	Result:  "",
	// 	TraceID: "",
	// 	Status:  "",
	// }).Error
	// fmt.Println("update with Select(*).Omit() error", err)

	// 方法3: 创建一个新的结构体指针，只设置要更新的字段
	updateLog := &model.Log{
		Result:  "",
		TraceID: "",
		Status:  "",
	}
	err = db.DB.Model(&model.Log{}).Where("id<=2").Updates(updateLog).Error
	fmt.Println("update with struct pointer error", err)

	// // 注意：如果结构体字段有默认值标签，GORM可能不会将空字符串视为有效值，这时需要使用Select指定字段
	// err = utils.DB.Model(&model.Log{}).Where("id=1").Select("result", "trace_id", "status").Updates(map[string]interface{}{
	// 	"result":  "",
	// 	"trace_id": "",
	// 	"status":  "",
	// }).Error
	// fmt.Println("update with Select and map error", err)
}
