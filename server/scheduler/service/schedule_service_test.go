package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"tooladmin/server/scheduler/db"
	"tooladmin/server/scheduler/model"

	"gorm.io/gorm"
)

func TestScheduleService_checkAndExecuteTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		schedule *model.Schedule
		wantErr  bool
	}{
		{
			name: "Task is disabled,will skip,task_name:",
			db:   db.DB,
			schedule: &model.Schedule{
				TaskName: "test",
				Flag:     "L",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSchedulerService(tt.db)
			gotErr := s.checkAndExecuteTask(context.Background(), tt.schedule)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("checkAndExecuteTask() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("checkAndExecuteTask() succeeded unexpectedly")
			}
		})
	}
}

func TestScheduleService_RetryTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
	}{
		{
			name: "RetryTask",
			db:   db.DB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSchedulerService(tt.db)
			s.checkAndRetryTask(context.Background())
		})
	}
}
func Test(t *testing.T) {
	fmt.Println(strings.TrimPrefix("test", "/"))
	fmt.Println(strings.TrimPrefix("/test", "/"))
	fmt.Println(strings.TrimPrefix("/test/2/2", "/"))

}
