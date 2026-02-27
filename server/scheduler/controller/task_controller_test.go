package controller_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
	"tooladmin/server/scheduler/controller"
	"tooladmin/server/scheduler/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/robfig/cron/v3"
)

func TestTaskController_GetAllTasks(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		taskService *service.TaskService
		// Named input parameters for target function.
		c *app.RequestContext
	}{
		{
			name:        "GetAllTasks",
			taskService: testController.TaskService,
			c:           &app.RequestContext{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := controller.NewTaskController(tt.taskService)

			tc.GetAllTasks(context.Background(), tt.c)
			fmt.Println(string(tt.c.Response.Body()))
		})
	}
}

// 发起到http://localhost:8080/api/v1/tasks的get请求的单测
func TestTaskController_GetAllTasks_HTTP(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(body))
}
func TestCron_Parse(t *testing.T) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	scheduler, err := parser.Parse("0 0 22 ? * 2")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(scheduler.Next(time.Now().UTC()))
}
