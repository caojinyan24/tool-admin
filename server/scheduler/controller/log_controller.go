package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/model"
	"tooladmin/server/scheduler/service"
)

// LogController 日志控制器
type LogController struct {
	LogService *service.LogService
}

// NewLogController 创建日志控制器实例
func NewLogController(logService *service.LogService) *LogController {
	return &LogController{
		LogService: logService,
	}
}

// GetLogs 获取日志列表
func (lc *LogController) GetLogs(ctx context.Context, c *app.RequestContext) {
	type GetLogsResponse struct {
		Data  []model.Log `json:"data"`
		Size  int         `json:"size"`
		Index int         `json:"index"`
		Count int         `json:"count"`
	}

	// 获取请求参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	taskName := c.DefaultQuery("taskName", "")
	status := c.DefaultQuery("status", "")
	executeType := c.DefaultQuery("executeType", "")
	startTime := c.DefaultQuery("startTime", "")
	endTime := c.DefaultQuery("endTime", "")

	// 转换参数类型
	pageInt := 1
	pageSizeInt := 10

	// 使用标准方式获取查询参数
	pageStr := c.Query("page")
	if pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &pageInt); err != nil {
			pageInt = 1
		}
	}

	pageSizeStr := c.Query("pageSize")
	if pageSizeStr != "" {
		if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSizeInt); err != nil {
			pageSizeInt = 10
		}
	}

	// 记录请求日志
	logger.Info(ctx, "GetLogs request params:",
		"page=", page,
		"pageSize=", pageSize,
		"taskName=", taskName,
		"status=", status,
		"executeType=", executeType,
		"startTime=", startTime,
		"endTime=", endTime)

	// 创建查询条件
	query := map[string]interface{}{
		"taskName":    taskName,
		"status":      status,
		"executeType": executeType,
		"startTime":   startTime,
		"endTime":     endTime,
		"page":        pageInt,
		"pageSize":    pageSizeInt,
	}

	// 查询日志列表
	logs, total, err := lc.LogService.GetLogs(query)
	if err != nil {
		logger.Error(ctx, "GetLogs error:", err)
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":  500,
			"msg":   "获取日志列表失败",
			"error": err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, utils.H{
		"code": 200,
		"msg":  "success",
		"data": GetLogsResponse{
			Data:  logs,
			Size:  pageSizeInt,
			Index: pageInt,
			Count: total,
		},
	})
}
