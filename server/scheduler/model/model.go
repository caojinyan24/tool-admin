package model

type ScheduleFailMsg struct {
	TaskName   string   `json:"task_name"`
	LogUrl     string   `json:"log_url"`
	ReceiveIds []string `json:"receive_ids"`
}
