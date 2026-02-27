package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"tooladmin/server/common/logger"

	"strings"
	"time"
	"tooladmin/server/scheduler/model"
)

func NewRemoteRequest(task *model.Task) *RemoteRequest {
	return &RemoteRequest{Task: task}
}

type RemoteRequest struct {
	Task *model.Task
}

// SendRequest 发送请求
func (s *RemoteRequest) SendRequest(ctx context.Context) (string, error) {
	if s.Task == nil {
		return "", fmt.Errorf("task is nil")
	}

	if s.Task.ApiType == "http" {
		return s.doHttpRequest(ctx)
	} else {
		return "", fmt.Errorf("unsupported api type: %s", s.Task.ApiType)
	}

}

// 全局HTTP客户端实例，复用连接
var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
	Timeout: 30 * time.Second,
}

func (s *RemoteRequest) doHttpRequest(ctx context.Context) (string, error) {
	// 创建HTTP请求上下文，设置超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 判断请求方法（默认为GET）
	method := "POST"
	data := bytes.NewBuffer(nil)
	// 如果参数不为空，使用POST方法
	if s.Task.TaskParam != "" {
		data = bytes.NewBuffer([]byte(s.Task.TaskParam))
	}

	// 创建请求
	var url = ""
	var path = strings.TrimPrefix(s.Task.ApiName, "/")
	var addr = strings.TrimSuffix(s.Task.ApiAddr, "/")
	if strings.HasPrefix(s.Task.ApiAddr, "http") {
		url = fmt.Sprintf("%s/%s", addr, path)
	} else {
		url = fmt.Sprintf("http://%s/%s", addr, path)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, data)
	logger.Debug(ctx, "sendRequest req", url, s.Task.TaskParam)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trace_id", ctx.Value(logger.TraceIDKey{}).(string))

	// 发送请求 - 使用复用的HTTP客户端
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(ctx, "sendRequest error", err)
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, "sendRequest error", err, string(body))
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		logger.Error(ctx, "sendRequest error", "statusCode=", resp.StatusCode)
		return "", fmt.Errorf("request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
