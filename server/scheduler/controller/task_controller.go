package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/model"
	"tooladmin/server/scheduler/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/robfig/cron/v3"
)

// TaskController 任务控制器
type TaskController struct {
	TaskService *service.TaskService
}

// NewTaskController 创建任务控制器实例
func NewTaskController(taskService *service.TaskService) *TaskController {
	return &TaskController{
		TaskService: taskService,
	}
}

// GetAllTasks 获取所有任务
func (tc *TaskController) GetAllTasks(ctx context.Context, c *app.RequestContext) {
	type GetAllTasksResponse struct {
		Data  []model.Task `json:"data"`
		Size  int          `json:"size"`
		Index int          `json:"index"`
		Count int          `json:"count"`
	}
	tasks, err := tc.TaskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":  500,
			"msg":   "获取任务列表失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": 200,
		"msg":  "success",
		"data": GetAllTasksResponse{
			Data:  tasks,
			Size:  10,
			Index: 1,
			Count: len(tasks),
		},
	})
}

// GetTaskByID 根据ID获取任务
func (tc *TaskController) GetTaskByTaskName(ctx context.Context, c *app.RequestContext) {
	taskName := c.Param("task_name")
	task, err := tc.TaskService.GetTaskByTaskName(taskName)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.H{
			"code":  404,
			"msg":   "任务不存在",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": 200,
		"msg":  "success",
		"data": task,
	})
}

// CreateTask 创建任务
func (tc *TaskController) CreateTask(ctx context.Context, c *app.RequestContext) {
	task := &model.Task{}
	if err := c.BindJSON(task); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":  400,
			"msg":   "请求参数错误",
			"error": err.Error(),
		})
		return
	}
	if task.ApiType != "http" {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":  400,
			"msg":   "暂时只支持http接口",
			"error": "api类型必须为http",
		})
		return
	}
	logger.Info(ctx, "CreateTask:%+v", task)
	if len(task.CronExpr) > 0 {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		_, err := parser.Parse(task.CronExpr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.H{
				"code":  400,
				"msg":   "cron表达式错误",
				"error": err.Error(),
			})
			return
		}

	}

	if len(task.Notifier) > 0 {
		emailList := make(map[string]model.Notifier, 0)
		json.Unmarshal([]byte(task.Notifier), &emailList)
		if len(emailList) > 0 {

		}

	}
	if err := tc.TaskService.CreateTask(ctx, task); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":  500,
			"msg":   "创建任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": 200,
		"msg":  "创建成功",
		"data": task,
	})
}

// UpdateTask 更新任务
func (tc *TaskController) UpdateTask(ctx context.Context, c *app.RequestContext) {
	task := &model.Task{}
	if err := c.BindJSON(task); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":  400,
			"msg":   "请求参数错误",
			"error": err.Error(),
		})
		return
	}
	if len(task.CronExpr) > 0 {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		_, err := parser.Parse(task.CronExpr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.H{
				"code":  400,
				"msg":   "cron表达式错误",
				"error": err.Error(),
			})
			return
		}

	}

	if err := tc.TaskService.UpdateTask(ctx, task); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":  500,
			"msg":   "更新任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": 200,
		"msg":  "更新成功",
		"data": task,
	})
}

// DeleteTask 删除任务
func (tc *TaskController) DeleteTask(ctx context.Context, c *app.RequestContext) {
	taskName := c.Param("task_name")

	if err := tc.TaskService.DeleteTask(ctx, taskName); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":  500,
			"msg":   "删除任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// TriggerTask 手动触发任务
func (tc *TaskController) TriggerTask(ctx context.Context, c *app.RequestContext) {
	ctx = logger.AddTraceID(ctx)
	logger.Info(ctx, "TriggerTask param:", string(c.Request.Body()))
	var taskParam = model.Task{}
	err := c.BindJSON(&taskParam)
	if err != nil {
		logger.Error(ctx, "parse param error:%v", err)
		c.JSON(http.StatusBadRequest, utils.H{
			"code":  400,
			"msg":   "请求参数错误",
			"error": err.Error(),
		})
		return
	}
	logger.Info(ctx, "taskParam:%+v", taskParam)
	result, err := tc.TaskService.ExecuteManualTask(ctx, &taskParam)
	if err != nil {
		logger.Error(ctx, "TriggerTask error:%v", err)
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":  500,
			"msg":   "触发任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": 200,
		"msg":  "触发成功",
		"data": result,
	})
}
