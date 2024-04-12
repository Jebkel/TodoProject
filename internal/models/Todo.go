package models

import (
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	ID uint64 `json:"id" gorm:"primary_key"`

	TaskName        string     `json:"task_name"`
	TaskDescription *string    `json:"task_description,omitempty"`
	DueDate         *time.Time `json:"due_date" gorm:"type:TIMESTAMP NULL"`
	Completed       bool       `json:"completed" gorm:"default:0" json:"completed"`

	ParentTaskID *uint64 `json:"parent_task_id" gorm:"index"`
	ParentTask   *Todo   `json:"parent_task" gorm:"foreignKey:ParentTaskID"`

	ChildTasks []*Todo `json:"child_tasks" gorm:"foreignKey:ParentTaskID;references:ID"`

	UserId *uint64 `json:"user_id" gorm:"index"`
	User   *User   `json:"-" gorm:"foreignKey:UserId"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func TodoPreloadChilds(d *gorm.DB) *gorm.DB {
	return d.Preload("ChildTasks", TodoPreloadChilds)
}
