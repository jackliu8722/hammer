package model

import "time"
const (
	TaskStatusUnstarted = iota + 1  // unstarted
	TaskStatusStarted				// Started
	TaskStatusDelivered				// Delivered
	TaskStatusRejected				// Rejected
	TaskStatusFinished				// Finished
	TaskStatusCancelled				// Cancelled
)
type Task struct {
	Id	int64 `json:"id" xorm:"pk autoincr unique"`
	Title string `json:"title"`
	PersonLimit int64 `json:"person_limit"`
	StartAt time.Time `json:"start_at"`
	EndAt time.Time `json:"end_at"`
	TaskCode string `json:"task_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_time"`
	DeletedAt time.Time `json:"deleted_time"`
	Amount float64 `json:"amount"`
}

func (t *Task) TableName() string {
	return "b_task"
}

type TaskLog struct {
	Id	int64 `json:"id" xorm:"pk autoincr unique"`
	TaskId int64 `json:"task_id"`
	Operator int64 `json:"operator"`
	FromStatus int16 `json:"from_status"`
	ToStatus int16 `json:"to_status"`
	CreatedAt time.Time `json:"Created_at"`
}

func (l *TaskLog) TableName() string {
	return "b_task_log"
}


type TaskUserRel struct {
	Id	int64 `json:"id" xorm:"pk autoincr unique"`
	Uid int64 `json:"uid"`
	TaskId int64 `json:"task_id"`
	Role int16 `json:"role"`
	Status int16 `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_time"`
	DeletedAt time.Time `json:"deleted_time"`
}

func (r *TaskUserRel) TableName() string {
	return "b_task_user_rel"
}